package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type KafkaStore[A event_store.Aggregate] struct {
	producer sarama.SyncProducer
	client   sarama.Client
	brokers  []string
	topic    string
	new      func(uuid.UUID) A
	eventFactory map[string]func() interface{}
	// Local subscribers
	subscribers []chan event_store.Record
	mu          sync.RWMutex
}

var _ event_store.Store[event_store.Aggregate] = &KafkaStore[event_store.Aggregate]{}

func NewKafkaStore[A event_store.Aggregate](
	brokers []string, 
	new func(uuid.UUID) A,
	eventFactory map[string]func() interface{},
) (*KafkaStore[A], error) {
	fmt.Printf("NewKafkaStore: connecting to brokers %v\n", brokers)
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 5

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}
	fmt.Println("NewKafkaStore: client created")

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}
	fmt.Println("NewKafkaStore: producer created")


	// Determine topic name from Aggregate Type
	aggregate := new(uuid.New())
	aggregateType := reflect.TypeOf(aggregate)
	if aggregateType.Kind() == reflect.Ptr {
		aggregateType = aggregateType.Elem()
	}
	parts := strings.Split(aggregateType.String(), ".")
	topic := parts[len(parts)-1]

	return &KafkaStore[A]{
		producer:     producer,
		client:       client,
		brokers:      brokers,
		topic:        topic,
		new:          new,
		eventFactory: eventFactory,
		subscribers:  make([]chan event_store.Record, 0),
	},
	nil
}

func (s *KafkaStore[A]) Close() error {
	if err := s.producer.Close(); err != nil {
		return err
	}
	return s.client.Close()
}

func (s *KafkaStore[A]) Subscribe(ctx context.Context) <-chan event_store.Record {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan event_store.Record, 100)
	s.subscribers = append(s.subscribers, ch)
	
	go s.consume(ctx, ch)

	return ch
}

func (s *KafkaStore[A]) consume(ctx context.Context, ch chan event_store.Record) {
	consumer, err := sarama.NewConsumerFromClient(s.client)
	if err != nil {
		fmt.Printf("consume error: %v\n", err)
		close(ch)
		return
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(s.topic)
	if err != nil {
		fmt.Printf("consume partitions error: %v\n", err)
		close(ch)
		return
	}

	var wg sync.WaitGroup
	for _, p := range partitions {
		wg.Add(1)
		pc, err := consumer.ConsumePartition(s.topic, p, sarama.OffsetNewest)
		if err != nil {
			wg.Done()
			continue
		}
		
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			defer pc.Close()
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-pc.Messages():
					record, err := s.deserialize(msg.Value)
					if err == nil {
						ch <- record
					} else {
						fmt.Printf("deserialize error: %v\n", err)
					}
				}
			}
		}(pc)
	}
	wg.Wait()
	close(ch)
}

func (s *KafkaStore[A]) Append(ctx context.Context, record event_store.Record) error {
	val, err := s.serialize(record)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: s.topic,
		Key:   sarama.StringEncoder(record.AggregateID.String()),
		Value: sarama.ByteEncoder(val),
	}

	_, _, err = s.producer.SendMessage(msg)
	return err
}

func (s *KafkaStore[A]) GetAggregate(ctx context.Context, id uuid.UUID) (A, error) {
	var zero A
	
	consumer, err := sarama.NewConsumerFromClient(s.client)
	if err != nil {
		return zero, err
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(s.topic)
	if err != nil {
		return zero, err
	}

	records := []event_store.Record{}
	
	for _, p := range partitions {
		limit, err := s.client.GetOffset(s.topic, p, sarama.OffsetNewest)
		if err != nil {
			return zero, err
		}
		
		if limit == 0 {
			continue
		}

		pc, err := consumer.ConsumePartition(s.topic, p, sarama.OffsetOldest)
		if err != nil {
			return zero, err
		}
		
		keepReading := true
		for keepReading {
			select {
			case <-ctx.Done():
				pc.Close()
				return zero, ctx.Err()
			case msg := <-pc.Messages():
				if string(msg.Key) == id.String() {
					rec, err := s.deserialize(msg.Value)
					if err == nil {
						records = append(records, rec)
					}
				}
				if msg.Offset >= limit-1 {
					keepReading = false
				}
			case <-time.After(5 * time.Second):
				keepReading = false
			}
		}
		pc.Close()
	}

	aggregate := s.new(id)
	if err := aggregate.Hydrate(records); err != nil {
		return zero, fmt.Errorf("failed to hydrate aggregate: %w", err)
	}

	return aggregate, nil
}

// Serialization helpers

type jsonRecord struct {
	AggregateID uuid.UUID       `json:"aggregate_id"`
	Version     uint64          `json:"version"`
	EventType   string          `json:"event_type"`
	EventData   json.RawMessage `json:"event_data"`
}

func (s *KafkaStore[A]) serialize(record event_store.Record) ([]byte, error) {
	data, err := json.Marshal(record.Event.Content())
	if err != nil {
		return nil, err
	}

	jr := jsonRecord{
		AggregateID: record.AggregateID,
		Version:     record.Version,
		EventType:   record.Event.Type(),
		EventData:   json.RawMessage(data),
	}
	return json.Marshal(jr)
}

func (s *KafkaStore[A]) deserialize(data []byte) (event_store.Record, error) {
	var jr jsonRecord
	if err := json.Unmarshal(data, &jr); err != nil {
		return event_store.Record{}, err
	}

	factory, ok := s.eventFactory[jr.EventType]
	if !ok {
		return event_store.Record{}, fmt.Errorf("unknown event type: %s", jr.EventType)
	}

	eventPtr := factory()
	if err := json.Unmarshal(jr.EventData, eventPtr); err != nil {
		return event_store.Record{}, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	// Dereference to get the value, since our events use value receivers and factory returns pointer
	contentVal := reflect.ValueOf(eventPtr).Elem().Interface()

	ge := &genericEvent{
		eventType: jr.EventType,
		content:   contentVal,
	}

	return event_store.Record{
		AggregateID: jr.AggregateID,
		Version:     jr.Version,
		Event:       ge,
	},
	nil
}

type genericEvent struct {
	eventType string
	content   interface{}
}

func (e *genericEvent) Type() string {
	return e.eventType
}

func (e *genericEvent) Content() any {
	return e.content
}
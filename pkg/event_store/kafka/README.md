# Kafka Event Store

This package implements an `event_store.Store` backed by Apache Kafka.

## Overview

The Kafka Event Store persists domain events to Kafka topics, enabling an event-sourced architecture with a distributed log as the source of truth.

### Key Features

*   **Topic Mapping:** Each Aggregate type is mapped to a specific Kafka topic. The topic name is derived from the Go type name of the Aggregate (e.g., `Account`, `Transaction`).
*   **Partitioning:** Events are produced with the Aggregate ID as the message Key. This ensures that all events for a specific aggregate are routed to the same partition, preserving causal ordering.
*   **Serialization:** Events are serialized as JSON before being stored in Kafka.
*   **Event Replay (Hydration):** Aggregates are hydrated by replaying the event stream from the beginning of the topic. The store scans partitions to find events matching the Aggregate ID.
*   **Subscriptions:** The store supports subscribing to new events, allowing Projections (Read Models) to update in real-time as events are appended.

## Configuration

The store is initialized with a list of Kafka brokers.

```go
brokers := []string{"localhost:9092"}
store, err := kafka.NewKafkaStore(brokers, MyAggregateFactory)
```

## Implementation Details

### Append

When `Append` is called:
1.  The event record is serialized to JSON.
2.  A Kafka message is created with:
    *   `Topic`: Derived from Aggregate Type.
    *   `Key`: Aggregate ID.
    *   `Value`: Serialized JSON event.
3.  The message is sent synchronously to Kafka with `WaitForAll` acks to ensure durability.

### GetAggregate

When `GetAggregate` is called:
1.  The store creates a consumer.
2.  It queries the current High Water Mark (latest offset) for all partitions of the topic.
3.  It consumes each partition from `OffsetOldest` up to the determined High Water Mark.
4.  It filters messages where the Key matches the requested Aggregate ID.
5.  Matching events are deserialized and applied to the Aggregate using its `Hydrate` method.

**Note:** This approach ("Event Sourcing on Log") can be slow for aggregates with long histories or topics with massive throughput, as it involves scanning. In production, this would typically be optimized with Snapshotting or by using a separate view database optimized for lookups (CQRS).

### Subscribe

When `Subscribe` is called:
1.  A background goroutine is started.
2.  It consumes all partitions of the topic starting from `OffsetNewest`.
3.  New events are deserialized and sent to the returned channel.
4.  This allows projections to react to new events immediately.

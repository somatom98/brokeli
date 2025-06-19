package transaction

import "context"

type Dispatcher struct{}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) CreateExpense(ctx context.Context, cmd CreateExpense) error {
	_, err := HandleCreateExpense(cmd)
	if err != nil {
		return err
	}

	// TODO: store event

	return nil
}

func (d *Dispatcher) CreateIncome(ctx context.Context, cmd CreateIncome) error {
	_, err := HandleCreateIncome(cmd)
	if err != nil {
		return err
	}

	// TODO: store event

	return nil
}

func (d *Dispatcher) CreateTransfer(ctx context.Context, cmd CreateTransfer) error {
	_, err := HandleCreateTransfer(cmd)
	if err != nil {
		return err
	}

	// TODO: store event

	return nil
}

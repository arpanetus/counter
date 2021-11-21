package service

type CounterServicer interface {
	Parse() error
	Count() uint64
	Add()
	Write() error
	Close() error
}

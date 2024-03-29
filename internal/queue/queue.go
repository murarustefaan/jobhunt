package queue

type WorkScheduler[E interface{}] interface {
	Push(item *E) error
	Pop() (*E, error)
}

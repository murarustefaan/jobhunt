package data

type Identifiable struct {
	Id string
}

type Storer[T Identifiable] interface {
	Create(item *T) error
	Find(id string) (*T, error)
}

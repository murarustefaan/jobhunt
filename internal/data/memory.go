package data

type MemoryStore struct {
	data map[string][]Identifiable
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string][]Identifiable),
	}
}

// todo: will this Store need to be refactored due to storing pointers in slices???

func (s *MemoryStore) Create(item *Identifiable) error {
	if items, _ := s.data[item.Id]; len(items) != 0 {
		items = append(items, *item)
	} else {
		s.data[item.Id] = []Identifiable{*item}
	}

	return nil
}

func (s *MemoryStore) Find(id string) ([]Identifiable, error) {
	if items, ok := s.data[id]; ok {
		return items, nil
	}
	return nil, nil
}

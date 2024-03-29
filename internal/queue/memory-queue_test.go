package queue

import (
	"slices"
	"testing"
)

func TestMemoryQueue_Pop(t *testing.T) {
	//t.Run("pop from empty data blocks", func(t *testing.T) {
	//	q := &MemoryQueue[int]{}
	//
	//	ctx, cancel := context.WithCancel(context.Background())
	//	defer cancel()
	//
	//	go func() {
	//		_, _ = q.Pop()
	//		cancel()
	//	}()
	//
	//	select {
	//	case <-ctx.Done():
	//	case <-time.After(100 * time.Millisecond):
	//		t.Error("expected pop to block")
	//	}
	//})
	//
	//t.Run("pop from data with one item returns the item", func(t *testing.T) {
	//	q := &MemoryQueue[int]{}
	//
	//	item := 42
	//	_ = q.Push(&item)
	//
	//	val, err := q.Pop()
	//	if err != nil {
	//		t.Errorf("unexpected error: %v", err)
	//	}
	//
	//	if *val != item {
	//		t.Errorf("expected %d, got %d", item, *val)
	//	}
	//})

	t.Run("pop from data with multiple items returns the items in LIFO order", func(t *testing.T) {
		q := NewMemoryQueue[int]()

		items := []int{1, 2, 3}
		for _, item := range items {
			_ = q.Push(&item)
		}

		slices.Reverse(items)
		for _, item := range items {
			val, err := q.Pop()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if *val != item {
				t.Errorf("expected %d, got %d", item, *val)
			}
		}

		if q.data != nil {
			t.Errorf("expected empty data, got %v", q.data)
		}
	})
}

//func TestMemoryQueue_Push(t *testing.T) {
//	t.Run("pushing an item to the data", func(t *testing.T) {
//		q := &MemoryQueue[int]{}
//
//		item := 42
//		err := q.Push(&item)
//		if err != nil {
//			t.Errorf("unexpected error: %v", err)
//		}
//
//		val, err := q.Pop()
//		if err != nil {
//			t.Errorf("unexpected error: %v", err)
//		}
//
//		if *val != item {
//			t.Errorf("expected %d, got %d", item, *val)
//		}
//	})
//}

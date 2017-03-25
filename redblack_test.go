package tree_test

import (
	"testing"

	"github.com/ryanfowler/tree"
)

func TestDeleteMax(t *testing.T) {
	const size = 1000

	var rb tree.RedBlackTree
	for i := size - 1; i >= 0; i-- {
		item := rb.Upsert(tree.Int(i))
		if item != nil {
			t.Fatalf("Unexpected replacement from insert: %+v", item)
		}
	}

	if rb.Size() != size {
		t.Fatalf("Unexpected size: %d", rb.Size())
	}

	for i := size - 1; i >= 0; i-- {
		item := int(rb.DeleteMax().(tree.Int))
		if item != i {
			t.Fatalf("Unexpected min value: %d", item)
		}
		if rb.Size() != i {
			t.Fatalf("Unexpceted size: %d", rb.Size())
		}
	}
	item := rb.DeleteMax()
	if item != nil {
		t.Fatalf("Unexpected non-nil value deleted: %+v", item)
	}
}

func TestDeleteMin(t *testing.T) {
	const size = 1000

	var rb tree.RedBlackTree
	for i := 0; i < size; i++ {
		item := rb.Upsert(tree.Int(i))
		if item != nil {
			t.Fatalf("Unexpected replacement from insert: %+v", item)
		}
	}

	if rb.Size() != size {
		t.Fatalf("Unexpected size: %d", rb.Size())
	}

	for i := 0; i < size; i++ {
		item := int(rb.DeleteMin().(tree.Int))
		if item != i {
			t.Fatalf("Unexpected min value: %d", item)
		}
		if rb.Size() != size-i-1 {
			t.Fatalf("Unexpceted size: %d", rb.Size())
		}
	}
	item := rb.DeleteMin()
	if item != nil {
		t.Fatalf("Unexpected non-nil value deleted: %+v", item)
	}
}

func TestAscendDescend(t *testing.T) {
	const size = 1000

	var rb tree.RedBlackTree
	rb.Ascend(nil)
	rb.Descend(nil)

	for i := 0; i < size; i++ {
		item := rb.Upsert(tree.Int(i))
		if item != nil {
			t.Fatalf("Unexpected replacement from insert: %+v", item)
		}
	}

	var i int
	rb.Ascend(func(item tree.Item) bool {
		val := int(item.(tree.Int))
		if val != i {
			t.Fatalf("Unexpected value in ascend: %d - %d", val, i)
		}
		i++
		return true
	})

	i = size - 1
	rb.Descend(func(item tree.Item) bool {
		val := int(item.(tree.Int))
		if val != i {
			t.Fatalf("Unexpected value in descend: %d - %d", val, i)
		}
		i--
		return true
	})
}

func TestDelete(t *testing.T) {
	var rb tree.RedBlackTree
	for i := 100; i > 0; i-- {
		rb.Upsert(tree.Int(i))
	}
	for i := 101; i <= 200; i++ {
		rb.Upsert(tree.Int(i))
	}

	for i, j := 100, 101; i > 0 && j <= 200; {
		it := rb.Delete(tree.Int(i))
		if it == nil || int(it.(tree.Int)) != i {
			t.Fatalf("Unexpected item deleted: %v", it)
		}
		it = rb.Delete(tree.Int(j))
		if it == nil || int(it.(tree.Int)) != j {
			t.Fatalf("Unexpected item deleted: %v", it)
		}
		i--
		j++
	}

	it := rb.Delete(tree.Int(400))
	if it != nil {
		t.Fatalf("Unexpected item deleted: %v", it)
	}
}

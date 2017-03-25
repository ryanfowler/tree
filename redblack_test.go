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

	it := rb.Delete(tree.Int(400))
	if it != nil {
		t.Fatalf("Unexpected item deleted: %v", it)
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

	it = rb.Delete(tree.Int(400))
	if it != nil {
		t.Fatalf("Unexpected item deleted: %v", it)
	}
}

func TestGet(t *testing.T) {
	var rb tree.RedBlackTree

	// Check Get, Max, Min on empty tree.
	if it := rb.Get(tree.Int(1)); it != nil {
		t.Fatalf("Unexpected item from Get: %v", it)
	}
	if it := rb.Max(); it != nil {
		t.Fatalf("Unexpected item from Max: %v", it)
	}
	if it := rb.Min(); it != nil {
		t.Fatalf("Unexpected item from Min: %v", it)
	}

	// Insert items.
	for i := 0; i < 500; i += 5 {
		rb.Upsert(tree.Int(i))
	}

	// Verify Get, Min, Max calls.
	if it := rb.Get(tree.Int(600)); it != nil {
		t.Fatalf("Unexpected non-nil item: %v", it)
	}
	if it := rb.Get(tree.Int(400)); int(it.(tree.Int)) != 400 {
		t.Fatalf("Unexpected item for value 400: %v", it)
	}
	if it := rb.Max(); int(it.(tree.Int)) != 495 {
		t.Fatalf("Unexpected max item: %v", it)
	}
	if it := rb.Min(); int(it.(tree.Int)) != 0 {
		t.Fatalf("Unexpected min item: %v", it)
	}

}

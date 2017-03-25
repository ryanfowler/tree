# tree [![Build Status](https://travis-ci.org/ryanfowler/tree.svg?branch=master)](https://travis-ci.org/ryanfowler/tree) [![Go Report Card](https://goreportcard.com/badge/github.com/ryanfowler/tree)](https://goreportcard.com/report/github.com/ryanfowler/tree) [![GoDoc](https://godoc.org/github.com/ryanfowler/tree?status.svg)](https://godoc.org/github.com/ryanfowler/tree)

In-memory Red-Black Tree implementation in Go.

Red-black trees are a self-balancing binary search tree that allow for O(log(n))
search, insertion, and deletion. To learn more, you can check out the
[Wikipedia page](https://en.wikipedia.org/wiki/Red%E2%80%93black_tree).

This library purposely maintains a similiar interface to the excellent btree 
implementation by Google:

[https://github.com/google/btree](https://github.com/google/btree)

Which, in turn, was inspired by the great:

[https://github.com/petar/GoLLRB](https://github.com/petar/GoLLRB)

## Usage

The zero value of a RedBlackTree is a ready to use empty tree.

```go
import "github.com/ryanfowler/tree"

var rb tree.RedBlackTree
```

### Items

Items put into the tree must implement the Item interface, which consists on a
single method, `Less(Item) bool`. Two Items 'a' & 'b' are considered equal if:

```go
equal := !a.Less(b) && !b.Less(a)
```

Convenience types that implement the Item interface are:

  - `Int`
  - `String`
  - `Bytes`

However, most of the time you'll want to use a custom type.

### Inserting

Items can be inserted or replaced in a RedBlackTree using `Upsert`. Upsert will
return the Item that was replaced by the Upsert. If there was no matching Item
in the tree, nil is returned.

```go
var rb tree.RedBlackTree

item := rb.Upsert(tree.Int(10))
fmt.Println(item)
// Will print: <nil>

item = rb.Upsert(tree.Int(10))
fmt.Println(item)
// Will print: 10
```

### Retrieving

An Item can be retrieved from a RedBlackTree using the following methods:

  - `Get(item Item) Item` - returns an Item in the tree _equal_ to the provided Item, or nil if there is no matching Item.
  - `Max() Item` - returns the largest Item in the tree, or nil if the tree is empty.
  - `Min() Item` - returns the smallest Item in the tree, or nil if the tree is empty.

```go
var rb tree.RedBlackTree

rb.Upsert(tree.Int(4))
rb.Upsert(tree.Int(5))
rb.Upsert(tree.Int(6))

item := rb.Get(tree.Int(5))
fmt.Println(item)
// Will print: 5

item = rb.Get(tree.Int(8))
fmt.Println(item)
// Will print: <nil>

item = rb.Max()
fmt.Println(item)
// Will print: 6

item = rb.Min()
fmt.Println(item)
// Will print: 4
```

### Deleting

An Item can be removed from a RedBlackTree using the following methods:

  - `Delete(item Item) Item` - removes an Item in the tree _equal_ to the provided Item, returning it. If there is no matching Item, nil is returned.
  - `DeleteMax() Item` - deletes and returns the largest Item in the tree, or nil if the tree is empty.
  - `DeleteMin() Item` - deletes and returns the smallest Item in the tree, or nil if the tree is empty.

```go
var rb tree.RedBlackTree

rb.Upsert(tree.Int(4))
rb.Upsert(tree.Int(5))
rb.Upsert(tree.Int(6))

fmt.Println(rb.Size())
// Will print: 3

item := rb.Delete(tree.Int(5))
fmt.Println(item)
// Will print: 5
fmt.Println(rb.Size())
// Will print: 2

item = rb.Delete(tree.Int(8))
fmt.Println(item)
// Will print: <nil>
fmt.Println(rb.Size())
// Will print: 2

item = rb.DeleteMax()
fmt.Println(item)
// Will print: 6
fmt.Println(rb.Size())
// Will print: 1

item = rb.Min()
fmt.Println(item)
// Will print: 4
fmt.Println(rb.Size())
// Will print: 0
```

### Example

Currently, a RedBlackTree does _not_ allow for the storing of equal items. This
can be worked around by including a second, unique field in the type's `Less` 
method.

For example, suppose you were designing a game and had a collection of users 
with their highest scores. If you wanted to retrieve the user(s) with the 
highest score, it could look something like this:

```go
import (
	"fmt"

	"github.com/ryanfowler/tree"
)

type User struct {
	ID    int
	Score int
}

func (u *User) Less(than tree.Item) bool {
	thanUser := than.(*User)
	return u.Score < thanUser.Score || (u.Score == thanUser.Score && u.ID < thanUser.ID)
}

func main() {
	var rb tree.RedBlackTree

	// Create users.
	u1 := &User{ID: 1, Score: 10}
	u2 := &User{ID: 2, Score: 7}
	u3 := &User{ID: 3, Score: 12}
	u4 := &User{ID: 4, Score: 10}

	// Insert into the RedBlackTree.
	rb.Upsert(u1)
	rb.Upsert(u2)
	rb.Upsert(u3)
	rb.Upsert(u4)

	// Get the user with the highest score.
	item := rb.Max()
	if item == nil {
		// The tree is empty.
		return
	}
	maxItem := item.(*User)
	fmt.Printf("%+v\n", maxItem)
	// Will print:
	// &{ID: 3, Score: 12}

	// Insert another user with a matching high score.
	u5 := &User{ID: 5, Score: 12}
	rb.Upsert(u5)

	// Get all users with the highest score.
	var maxScore int
	var maxUsers []*User
	rb.Descend(func(item tree.Item) bool {
		user := item.(*User)
		if len(maxUsers) == 0 {
			maxScore = user.Score
			maxUsers = append(maxUsers, user)
			return true
		}
		if user.Score < maxScore {
			return false
		}
		maxUsers = append(maxUsers, user)
		return true
	})
	for _, user := range maxUsers {
		fmt.Printf("%+v\n", user)
	}
	// Will print:
	// &{ID: 5, Score: 12}
	// &{ID: 3, Score: 12}
}
```

## License

MIT License

Copyright (c) 2017 Ryan Fowler

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

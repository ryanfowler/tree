// MIT License
//
// Copyright (c) 2017 Ryan Fowler
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package tree provides an implementation of a red-black tree.
package tree

// Int represents an integer that implements the Item interface.
type Int int

// Less returns true if the Int is less than the provided Int. If the provided
// Item is not an Int, Less will panic.
func (i Int) Less(than Item) bool {
	return i < than.(Int)
}

// Item is the interface that wraps the Less method.
//
// Less should return 'true' if the instance is "less than" the provided Item.
// Items are considered equal if neither are less than each other.
// E.g. Items 'a' & 'b' are considered equal if: (!a.Less(b) && !b.Less(a))
type Item interface {
	Less(Item) bool
}

// RedBlackTree is an in-memory implementation of a red-black tree.
//
// The internal data structure will automatically re-balance, and therefore
// allow for O(log(n)) retrieval, insertion, and deletion.
//
// Note: While read-only operations may occur concurrently, any write operation
// must be serially executed (typically protected with a mutex).
type RedBlackTree struct {
	root *node
	size int
}

// Ascend (O(n)) starts at the first Item and calls 'fn' for each Item until no
// Items remain or fn returns 'false'.
func (t *RedBlackTree) Ascend(fn func(Item) bool) {
	if t.root == nil {
		return
	}
	n := t.root.min()
	for n != nil && fn(n.item) {
		n = n.next()
	}
}

// Descend (O(n)) starts at the last Item and calls 'fn' for each Item until no
// Items remain or fn returns 'false'.
func (t *RedBlackTree) Descend(fn func(Item) bool) {
	if t.root == nil {
		return
	}
	n := t.root.max()
	for n != nil && fn(n.item) {
		n = n.prev()
	}
}

// Delete (O(log(n))) deletes an item in the RedBlackTree equal to the provided
// item. If an item was deleted, it is returned. Otherwise, nil is returned.
//
// Note: equality for items a & b is: (!a.Less(b) && !b.Less(a)).
func (t *RedBlackTree) Delete(item Item) Item {
	if t.root == nil {
		return nil
	}
	return t.root.deleteItem(t, item)
}

// DeleteMax (O(log(n))) deletes the maximum item in the RedBlackTree, returning
// it. If the tree is empty, nil is returned.
func (t *RedBlackTree) DeleteMax() Item {
	if t.root == nil {
		return nil
	}
	return t.root.deleteMax(t)
}

// DeleteMin (O(log(n))) deletes the minimum item in the RedBlackTree, returning
// it. If the tree is empty, nil is returned.
func (t *RedBlackTree) DeleteMin() Item {
	if t.root == nil {
		return nil
	}
	return t.root.deleteMin(t)
}

// Get (O(log(n))) retrieves an item in the RedBlackTree equal to the provided
// item. If an item was found, it is returned. Otherwise, nil is returned.
//
// Note: equality for items a & b is: (!a.Less(b) && !b.Less(a)).
func (t *RedBlackTree) Get(item Item) Item {
	n := t.root.find(item)
	if n == nil {
		return nil
	}
	return n.item
}

// Insert (O(log(n))) inserts (or replaces) an item into the RedBlackTree. If an
// item was replaced, it is returned. Otherwise, nil is returned.
//
// Note: equality for items a & b is: (!a.Less(b) && !b.Less(a)).
func (t *RedBlackTree) Insert(item Item) Item {
	if t.root == nil {
		t.root = newNode(nil, item)
		t.root.colour = colourBlack
		t.size++
		return nil
	}
	n, oldItem := t.root.insert(item)
	if oldItem == nil {
		t.size++
		n.rebalanceInsert(t)
	}
	return oldItem
}

// Exists (O(log(n))) returns 'true' if an item equal to the provided item
// exists in the RedBlackTree.
//
// Note: equality for items a & b is: (!a.Less(b) && !b.Less(a)).
func (t *RedBlackTree) Exists(item Item) bool {
	return t.Get(item) != nil
}

// Min (O(log(n))) returns the minimum item in the RedBlackTree. If the tree is
// empty, nil is returned.
func (t *RedBlackTree) Min() Item {
	if t.root == nil {
		return nil
	}
	n := t.root
	for n.left != nil {
		n = n.left
	}
	return n.item
}

// Max (O(log(n))) returns the maximum item in the RedBlackTree. If the tree is
// empty, nil is returned.
func (t *RedBlackTree) Max() Item {
	if t.root == nil {
		return nil
	}
	n := t.root
	for n.right != nil {
		n = n.right
	}
	return n.item
}

// Size (O(1)) returns the number of items in the RedBlackTree.
func (t *RedBlackTree) Size() int {
	return t.size
}

type colour uint8

const (
	colourRed   colour = 0
	colourBlack colour = 1
)

type node struct {
	colour      colour
	parent      *node
	left, right *node
	item        Item
}

func newNode(parent *node, item Item) *node {
	return &node{
		colour: colourRed,
		parent: parent,
		item:   item,
	}
}

func (n *node) find(item Item) *node {
	for n != nil {
		switch {
		case item.Less(n.item):
			n = n.left
		case n.item.Less(item):
			n = n.right
		default:
			return n
		}
	}
	return nil
}

func (n *node) deleteMax(t *RedBlackTree) Item {
	return n.max().deleteNode(t)
}

func (n *node) deleteMin(t *RedBlackTree) Item {
	return n.min().deleteNode(t)
}

func (n *node) deleteItem(t *RedBlackTree, item Item) Item {
	n = n.find(item)
	if n == nil {
		return nil
	}
	return n.deleteNode(t)
}

func (n *node) deleteNode(t *RedBlackTree) Item {
	t.size--
	delItem := n.item

	var child, parent *node
	for {
		if n.left == nil {
			child = n.right
			parent = n.parent
			n.replaceNode(t, n.right)
			break
		}
		if n.right == nil {
			child = n.left
			parent = n.parent
			n.replaceNode(t, n.left)
			break
		}
		// replace minimum value in right subtree with node to delete.
		min := n.right.min()
		n.item = min.item
		n = min
	}

	if n.isRed() {
		return delItem
	}
	if child.isRed() {
		child.colour = colourBlack
		return delItem
	}
	child.rebalanceDelete(t, parent)
	return delItem
}

func (n *node) rebalanceDelete(t *RedBlackTree, parent *node) {
	var s *node
	for {
		// Case 1.
		if n == t.root {
			return
		}
		if n != nil {
			parent = n.parent
		}
		// Case 2.
		s = n.sibling(parent)
		if s.isRed() {
			parent.colour = colourRed
			s.colour = colourBlack
			if n == parent.left {
				parent.rotateLeft(t)
			} else {
				parent.rotateRight(t)
			}
		}
		// Case 3.
		s = n.sibling(parent)
		if parent.isBlack() && s.isBlack() && s != nil && s.left.isBlack() && s.right.isBlack() {
			s.colour = colourRed
			n = parent
			if n != nil {
				parent = n.parent
			} else {
				parent = nil
			}
			continue
		}
		break
	}
	// Case 4.
	if parent.isRed() &&
		s.isBlack() &&
		s != nil &&
		s.left.isBlack() &&
		s.right.isBlack() {
		s.colour = colourRed
		parent.colour = colourBlack
		return
	}
	// Case 5.
	if s.isBlack() && s != nil {
		if n == parent.left && s.right.isBlack() && s.left.isRed() {
			s.colour = colourRed
			s.left.colour = colourBlack
			s.rotateRight(t)
		} else if n == parent.right && s.left.isBlack() && s.right.isRed() {
			s.colour = colourRed
			s.right.colour = colourBlack
			s.rotateLeft(t)
		}
	}
	// Case 6.
	s = n.sibling(parent)
	if s != nil {
		s.colour = parent.colour
		parent.colour = colourBlack
		if n == parent.left {
			s.right.colour = colourBlack
			parent.rotateLeft(t)
		} else {
			s.left.colour = colourBlack
			parent.rotateRight(t)
		}
	}
}

func (n *node) isRed() bool {
	return n != nil && n.colour == colourRed
}

func (n *node) isBlack() bool {
	return n == nil || n.colour == colourBlack
}

func (n *node) sibling(parent *node) *node {
	if n == parent.left {
		return parent.right
	}
	return parent.left
}

func (n *node) replaceNode(t *RedBlackTree, child *node) {
	switch {
	case n.parent == nil:
		t.root = child
	case n == n.parent.left:
		n.parent.left = child
	default:
		n.parent.right = child
	}
	if child != nil {
		child.parent = n.parent
	}
}

func (n *node) min() *node {
	for n.left != nil {
		n = n.left
	}
	return n
}

func (n *node) max() *node {
	for n.right != nil {
		n = n.right
	}
	return n
}

func (n *node) next() *node {
	if n.right != nil {
		return n.right.min()
	}
	parent := n.parent
	for parent != nil && parent.right == n {
		n = parent
		parent = n.parent
	}
	return parent
}

func (n *node) prev() *node {
	if n.left != nil {
		return n.left.max()
	}
	parent := n.parent
	for parent != nil && parent.left == n {
		n = parent
		parent = n.parent
	}
	return parent
}

func (n *node) insert(item Item) (*node, Item) {
	for {
		switch {
		case item.Less(n.item):
			if n.left == nil {
				n.left = newNode(n, item)
				return n.left, nil
			}
			n = n.left
		case n.item.Less(item):
			if n.right == nil {
				n.right = newNode(n, item)
				return n.right, nil
			}
			n = n.right
		default:
			oldItem := n.item
			n.item = item
			return n, oldItem
		}
	}
}

func (n *node) rebalanceInsert(t *RedBlackTree) {
	var g *node
	for {
		// Case 1.
		if n.parent == nil {
			n.colour = colourBlack
			return
		}
		// Case 2.
		if n.parent.colour == colourBlack {
			return
		}
		// Case 3.
		g = n.grandparent()
		var ps *node
		if g != nil {
			if n.parent == g.left {
				ps = g.right
			} else {
				ps = g.left
			}
		}
		if ps == nil || ps.colour == colourBlack {
			break
		}
		n.parent.colour = colourBlack
		ps.colour = colourBlack
		g.colour = colourRed
		n = g
	}
	// Case 4.
	if n == n.parent.right && n.parent == g.left {
		n.parent.rotateLeft(t)
		n = n.left
		g = n.grandparent()
	} else if n == n.parent.left && n.parent == g.right {
		n.parent.rotateRight(t)
		n = n.right
		g = n.grandparent()
	}
	// Case 5.
	n.parent.colour = colourBlack
	g.colour = colourRed
	if n == n.parent.left {
		g.rotateRight(t)
	} else {
		g.rotateLeft(t)
	}
}

func (n *node) rotateLeft(t *RedBlackTree) {
	right := n.right
	n.right = right.left
	if right.left != nil {
		right.left.parent = n
	}
	right.parent = n.parent
	switch {
	case n.parent == nil:
		t.root = right
	case n == n.parent.left:
		n.parent.left = right
	default:
		n.parent.right = right
	}
	right.left = n
	n.parent = right
}

func (n *node) rotateRight(t *RedBlackTree) {
	left := n.left
	n.left = left.right
	if left.right != nil {
		left.right.parent = n
	}
	left.parent = n.parent
	switch {
	case n.parent == nil:
		t.root = left
	case n == n.parent.right:
		n.parent.right = left
	default:
		n.parent.left = left
	}
	left.right = n
	n.parent = left
}

func (n *node) grandparent() *node {
	if n == nil || n.parent == nil {
		return nil
	}
	return n.parent.parent
}

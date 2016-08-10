/*
<!--
Copyright (c) 2016 Christoph Berger. Some rights reserved.
Use of this text is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Balancing a binary search tree"
description = "This article describes a basic tree balancing technique, coded in Go, and applied to the binary search tree from last week's article."
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-08-11"
publishdate = "2016-08-11"
domains = ["Algorithms And Data Strucutures"]
tags = ["Tree", "Balanced Tree", "Binary Tree", "Search Tree"]
categories = ["Tutorial"]
+++

Only a well-balanced search tree can provide optimal search performance. This article adds automatic balancing to the binary search tree from the previous article.

<!--more-->

## How a tree can get out of balance

As we have seen in last week's article, search performance is best if the tree's height is small. Unfortunately, without any further measure, our simple binary search tree can quickly get out of shape - or never reach a good shape in the first place.

The picture below shows a balanced tree on the left and an extreme case of an unbalanced tree at the right. In the balanced tree, element #6 can be reached in three steps, whereas in the extremely unbalanced case, it takes six steps to find element #6.

![Tree Shapes](BinTreeShapes.png)

Unfortunately, the extreme case can occur quite easily: Just create the tree from a sorted list.

```go
tree.Insert(1)
tree.Insert(2)
tree.Insert(3)
tree.Insert(4)
tree.Insert(5)
tree.Insert(6)
```

According to `Insert`'s logic, each new element is added as the right child of the rightmost node, because it is larger than any of the elements that were already inserted.

We need a way to avoid this.


## A Definition Of "Balanced"

For our purposes, a good working definition of "balanced" is:

> The heights of the two child subtrees of any node differ by at most one.
>
> (Wikipedia: [AVL-Tree](https://en.wikipedia.org/wiki/AVL_tree))

Why "at most one"? Shouldn't we demand *zero* difference for perfect balance? Actually, no, as we can see on this very simple two-node tree:

![Two-node tree](TwoNodeTree.png)

The left subtree is a single node, hence the height is 1, and the right "subtree" is empty, hence the height is zero. There is no way to make both subtrees exactly the same height, except perhaps by adding a third "fake" node that has no other purpose of providing perfect balance. But we would gain nothing from this, so a height difference of 1 is perfectly acceptable.

Note that our definition of *balanced* does not include the *size* of the left and right subtrees of a node. That is, the following tree is completely fine:

![No Weight Balance](BinTreeNoWeightBalance.png)

The left subtree is considerably larger than the right one; yet for either of the two subtrees, any node can be reached with at most four search steps. And the heights of both subtrees differs only by one.


## How to keep a tree in balance

Now that we know what balance means, we need to take care of always keeping the tree in balance. This task consists of two parts: First, we need to be able to detect when a (sub-)tree goes out of balance. And second, we need a way to rearrange the nodes so that the tree is in balance again.


### Step 1. Detecting an imbalance

Balance is related to subtree heights, so we might think of writing a "height" method that descends a given subtree to calculate its height. But this can be come quite costly in terms of CPU time, as these calculations would need to be done repeatedly as we try to determine the balance of each subtee and each subtree's subree, and so on.

Instead, we store a "balance factor" in each node. This factor is an integer that tells the height difference between the node's left and right subtrees. Based on our definition of "balanced", the balance factor of a balanced tree can be -1, 0, or +1. If the balance factor is outside that range (that is, either smaller than -1 or larger than +1), the tree is out of balance and needs to be rebalanced.

The balance factor is maintained by the `Insert` and `Delete` operations.

*For brevity, this article only handles the `Insert` case.*

Here is how `Insert` maintains the balance factors:

1. First, `Insert` descends recursively down the tree until it finds a node `n` to append the new value. `n` is either a leaf or a half-leaf.
2. If `n` is a leaf, adding a new child node increases the height of the subtree `n` by 1. If the child node is added to the left, the balance of `n` changes from 0 to -1. If the child is added to the right, the balance changes from 0 to 1.
2. `Insert` now adds a new child node to node `n`.
3. The height increase is passed back to `n`'s parent node.
4. Depending on whether `n` is the left or the right child, the parent node adjusts its balance accordingly.

**An imbalance is detected if the balance factor of a node changes to +2 or -2, respectively.** At this point, the affected node must start the rebalancing.

HYPE[Balance Factors](BalanceFactors.html)

### Removing the imbalance

Let's assume the unbalanced node has a balance factor of -2. This means that its left subtree is too high. Two situations can occur here.

#### 1. The left child node has a balance of 0 or -1.

In other words, the left child node's left subtree is higher than its right subtree. This is an easy case. All we have to do is to "rotate" the tree:

1. Make the left child node the root node.
2. If the former left child node has a right subtree, add this subtree to the former root node as the left child.
3. Make the former root node the new root node's right child.

This may sound a bit complicated, so here is a visualization:

HYPE[Rotation](Rotation.html)

Clicking on "Without right child" shows the simples form of rotation. This is only step 1 and 3 of the above sequence. Step 2 occurs when the left child node also has a right child - click "With right child" to watch this scenario.


#### 2. The left child node has a balance of 1.

This means that the left child's *right* subtree is higher than its left subtree. We can try to apply the same type of rotation to this scenario. Click the button "Single Rotation" and see what happens:

HYPE[Double Rotation](DoubleRotation.html)

The tree is again unbalanced; the root node's balance factor changed from -2 to +2. Obviously, a simple rotation as in case 1 does not work here.

Now try the second button, "Double Rotation". Here, the unbalanced node's left subtree is rotated first, and now the situation is similar to case 1. Now rotating the tree to the right rebalances the tree.


#### Two more cases and a summary

The two cases above assumed that the unbalanced node's balance factor is -2. If the balance factor is +2, the same cases apply in an analgous way, except that everything is mirror-reversed.


To summarize, here is a scenario where all of the above is included - double rotation as well as reassigning a child node/tree to a rotated node.

HYPE[Re-balance](Rebalance.html)


## The Code

Now, after all this theory, let's see how to add the balancing into the code from the previous article.

First, we set up two helper functions, `min` and `max`, that we will need later.

*/

// ### Imports, helper functions, and globals
package main

import (
	"fmt"
	"strings"
)

// `min` is like math.Min but for int.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// `max` is math.Max for int.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// `Node` gets a new field, `bal`, to store the height difference between the node's subtrees.
type Node struct {
	Value string
	Data  string
	Left  *Node
	Right *Node
	bal   int // height(n.Right) - height(n.Left)
}

// ### The modified `Insert` function

// `Insert` returns:
//
// * `true` if the height of the tree has increased.
// * `false` otherwise.
func (n *Node) Insert(value, data string) bool {

	switch {
	case value == n.Value:
		return false // Node already exists, nothing changes
	case value < n.Value:
		// If there is no left child, create a new one.
		if n.Left == nil {
			// A new left child reduces the balance of this node by one, making it either 0 or -1.
			n.bal--
			// Create a new node.
			n.Left = &Node{Value: value, Data: data}
			// If there is no right child, the new child node has increased the height of this subtree.
			if n.Right == nil {
				return true
			}
			return false
		}
		// The left child is not nil. Continue in the left subtree.
		if n.Left.Insert(value, data) {
			// The left subtree has grown by one: Decrease the balance by one.
			n.bal--
		}
	case value > n.Value:
		if n.Right == nil {
			// A new right child increases the balance of this node by one, making it either 0 or +1.
			n.bal++
			n.Right = &Node{Value: value, Data: data}
			// If there is no left child, the new child node has increased the height of this subtree.
			if n.Left == nil {
				return true
			}
			return false
		}
		if n.Right.Insert(value, data) {
			// The right subtree has grown by one. Increase the balance by one.
			n.bal++
		}
	}
	// If rebalancing is required, the method `rebalance()` takes care of all the different rebalancing scenarios.
	if n.bal < -1 || n.bal > 1 {
		n.rebalance()
	}
	if n.bal != 0 {
		return true
	}
	// No more adjustments to the ancestor nodes required.
	return false
}

// ### The new `rebalance()` method and its helpers `rotateLeft()`, `rotateRight()`, `rotateLeftRight()`, and `rotateRightLeft`.

// `rotateLeft` takes a parent node and rotates the current node's subtree to the left.
func (n *Node) rotateLeft(p *Node) *Node {
	// Save `n`'s right child.
	r := n.Right
	// `r`'s right subtree gets reassigned to `n`.
	n.Right = r.Left
	// `n` becomes the left child of `r`.
	r.Left = n
	// Make the parent node point to the new root node.
	if p != nil {
		if n == p.Left {
			p.Left = r
		} else {
			p.Right = r
		}
	}
	// Finally, adjust the balances.
	if r.bal == 0 { // This case does not apply to inserts, only to deletes.
		n.bal = 1
		r.bal = -1
	} else {
		n.bal = 0
		r.bal = 0
	}
	return r
}

// `rotateRight` is the mirrored version of `rotateLeft`.
func (n *Node) rotateRight(p *Node) *Node {
	l := n.Left
	n.Left = l.Right
	l.Right = n
	if p != nil {
		if n == p.Left {
			p.Left = l
		} else {
			p.Right = l
		}
	}
	return l
}

// `rotateRightLeft` first rotates the right child to the right, then the current node to the left.
func (n *Node) rotateRightLeft(p *Node) *Node {
	// TODO
}

// `rotateLeftRight` first rotates the right child to the left, then the current node to the right.
func (n *Node) rotateRightLeft(p *Node) *Node {
	// TODO
}

// `rebalance` brings the tree back into a balanced state.
func (n *Node) rebalance() {
	fmt.Println("rebalance " + n.Value)
	// TODO
}

// `Find` stays the same.
func (n *Node) Find(s string) (string, bool) {

	if n == nil {
		return "", false
	}

	switch {
	case s == n.Value:
		return n.Data, true
	case s < n.Value:
		return n.Left.Find(s)
	default:
		return n.Right.Find(s)
	}
}

// `Dump` dumps the structure of the subtree starting at node `n`, including node search values and balance factors.
func (n *Node) Dump(i int) {
	if n == nil {
		return
	}
	indent := ""
	if i > 0 {
		indent = strings.Repeat(" ", (i-1)*4) + "+" + strings.Repeat("-", 3)
	}
	fmt.Printf("%s%s[%d]\n", indent, n.Value, n.bal)
	n.Left.Dump(i + 1)
	n.Right.Dump(i + 1)
}

/*
## Tree

The Tree type is largely unchanged, except that `Delete` is gone and a new method, `Dump`, exist for invoking `Node.Dump`.

*/

//
type Tree struct {
	Root *Node
}

func (t *Tree) Insert(value, data string) {
	if t.Root == nil {
		t.Root = &Node{Value: value, Data: data}
		return
	}
	t.Root.Insert(value, data)
}

func (t *Tree) Find(s string) (string, bool) {
	if t.Root == nil {
		return "", false
	}
	return t.Root.Find(s)
}

func (t *Tree) Traverse(n *Node, f func(*Node)) {
	if n == nil {
		return
	}
	t.Traverse(n.Left, f)
	f(n)
	t.Traverse(n.Right, f)
}

// `Dump` dumps the tree structure.
func (t *Tree) Dump() {
	t.Root.Dump(0)
}

func main() {
	values := []string{"d", "b", "g", "c", "e", "a", "h", "f", "i", "j", "k", "l"}
	data := []string{"delta", "bravo", "golf", "charlie", "echo", "alpha", "hotel", "foxtrot", "india", "juliett", "kilo", "lima"}

	tree := &Tree{}
	for i := 0; i < len(values); i++ {
		tree.Insert(values[i], data[i])
		tree.Dump()
		fmt.Println()
	}

	fmt.Print("Sorted values: | ")
	tree.Traverse(tree.Root, func(n *Node) { fmt.Print(n.Value, ": ", n.Data, " | ") })
	fmt.Println()

	tree.Dump()
	tree.Root.Right.Right.Right.Right.rotateLeft(tree.Root.Right.Right.Right)
	tree.Dump()
}

/*
As always, the code is available on GitHub. Using `-d` on `go get` avoids installing the binary into $GOPATH/bin.

```sh
go get -d github.com/appliedgo/balancedtree
cd $GOPATH/src/github.com/appliedgo/balancedtree
go build
./balancedtree
```

## Conclusion



*/

package main

import (
	"fmt"
	"math"
	"testing"
)

type tree struct {
	name  string
	value []string
	data  []string
}

var (
	trees = []tree{
		{
			name:  "empty",
			value: []string{},
			data:  []string{},
		},
		{
			name:  "onenode",
			value: []string{"0"},
			data:  []string{"zero"},
		},
		{
			name:  "twonodes",
			value: []string{"0", "1"},
			data:  []string{"zero", "one"},
		},
		{
			name:  "random",
			value: []string{"d", "b", "g", "g", "c", "e", "a", "h", "f", "i", "j", "l", "k"},
			data:  []string{"delta", "bravo", "golang", "golf", "charlie", "echo", "alpha", "hotel", "foxtrot", "india", "juliett", "lima", "kilo"},
		},
		{
			name:  "ascending",
			value: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m"},
			data:  []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel", "india", "juliett", "kilo", "lima", "mike"},
		},
		{
			name:  "descending",
			value: []string{"m", "l", "k", "j", "i", "h", "g", "f", "e", "d", "c", "b", "a"},
			data:  []string{"mike", "lima", "kilo", "juliett", "india", "hotel", "golf", "foxtrot", "echo", "delta", "charlie", "bravo", "alpha"},
		},
		{
			name:  "issue2",
			value: []string{"3", "5", "1", "0", "2", "4", "6", "7", "8"},
			data:  []string{"3", "5", "1", "0", "2", "4", "6", "7", "8"},
		},
		{
			name:  "balancedfromthestart",
			value: []string{"4", "2", "6", "1", "7", "3", "5"},
			data:  []string{"4", "2", "6", "1", "7", "3", "5"},
		},
	}
)

func newTree(t tree) *Tree {
	tree := &Tree{}
	for i := 0; i < len(t.value); i++ {
		tree.Insert(t.value[i], t.data[i])
	}
	return tree
}

// calculate the height recursively, without relying on n.height
func (n *Node) recHeight() int {
	if n == nil {
		return 0
	}
	return 1 + max(n.Left.recHeight(), n.Right.recHeight())
}

func (n *Node) checkHeight() (*Node, bool) {
	if n == nil {
		return nil, true
	}

	if n.height != n.recHeight() {
		return n, false
	}

	if node, ok := n.Left.checkHeight(); !ok {
		return node, false
	}

	if node, ok := n.Right.checkHeight(); !ok {
		return node, false
	}
	return nil, true
}

// A (sub-)tree is balanced if the heights of the two child subtrees of any node differ by at most one.
func (n *Node) isBalanced() bool {
	return n == nil || n.Right.recHeight()-n.Left.recHeight() <= 1
}

func (n *Node) checkBalances() (problem string) {
	if n == nil {
		return ""
	}
	rh, lh := n.Right.recHeight(), n.Left.recHeight()
	if n.Bal() != rh-lh {
		problem = fmt.Sprintf("Node %s has balance %d but right height %d and left height %d\n", n.Value, n.Bal(), rh, lh)
	}
	return problem + n.Right.checkBalances() + n.Left.checkBalances()
}

func (t *Tree) containsAllElements(source tree) (string, bool) {
	for _, v := range source.value {
		_, found := t.Find(v)
		if !found {
			return v, false
		}
	}
	return "", true
}

func (t *Tree) isSorted() bool {
	var sorted func(*Node) bool
	sorted = func(n *Node) bool {
		if n == nil {
			return true
		}
		if (n.Left != nil && n.Value < n.Left.Value) ||
			(n.Right != nil && n.Value > n.Right.Value) {
			return false
		}
		return sorted(n.Left) && sorted(n.Right)
	}
	return sorted(t.Root)
}

func TestTree_rebalance(t *testing.T) {
	for _, tree := range trees {
		t.Run(tree.name, func(t *testing.T) {
			fmt.Println("Creating tree ", tree.name)
			tt := newTree(tree)
			tt.Dump()
			h := tt.Root.recHeight()
			lh, rh := 0, 0
			if tt.Root != nil {
				lh = tt.Root.Left.recHeight()
				rh = tt.Root.Right.recHeight()
			}
			exh := 2.0*math.Log2(float64(len(tree.value))+1.44) - 0.328

			heightImbalance := ""
			if float64(h) > exh {
				heightImbalance = fmt.Sprintf("Height: %d - expected: %0f\nLeft.Height(): %d, Right.Height(): %d\n", h, exh, lh, rh)
			}
			wrongBalanceFactors := tt.Root.checkBalances()
			problem := heightImbalance + wrongBalanceFactors

			if v, ok := tt.containsAllElements(tree); !ok {
				problem += fmt.Sprintf("Some data in the tree is missing or wrong: %s\n", v)
			}

			if !tt.isSorted() {
				problem += fmt.Sprintf("Tree %s is not balanced\n", tree.name)
			}

			if n, ok := tt.Root.checkHeight(); !ok {
				problem += fmt.Sprintf("Actual height %d differs from recorded height %d in node %s\n", n.recHeight(), n.height, n.Value)
			}

			if len(problem) > 0 {
				t.Error(problem)
			}
		})
	}
}

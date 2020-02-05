package main

import (
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	recursiveWalk(t, ch)
	close(ch)
}

func recursiveWalk(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		recursiveWalk(t.Left, ch)
	}

	ch <- t.Value

	if t.Right != nil {
		recursiveWalk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for c1 := range ch1 {
		c2, ch2ok := <-ch2
		if !ch2ok || c1 != c2 {
			return false
		}
	}

	_, ok := <-ch2
	if ok {
		return false
	}

	return true
}

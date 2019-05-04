// yourtool
package main

import (
	"fmt"
	"strings"
)

//trie
type Vertex struct {
	next [127]int
	leaf bool
	p int
	pch byte
	link int	
	gonow [127]int
}


	t []Vertex
	size int

func add_string(const s &string, t []Vertex) {
	v := 0
	for i :=0; i < len(s); ++i {
		c := s[i] - 'a'
		if t[v].next[c] == -1 {
			
			t[v].next[c] = size++
		}
		v = t[v].next[c]
	}
	t[v].leaf = true
}

func init() {
	t[0].p = t[0].link
}

func main() {

	//Initialisation of trie
	size = 1
	
}

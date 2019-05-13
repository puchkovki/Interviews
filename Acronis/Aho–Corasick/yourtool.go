// yourtool
package main

import (
	"fmt"
	//"strings"
)

//trie
type Vertex struct {
	//edges
	next [255]int

	//end of the template
	leaf bool

	//parent
	p int

	//symbol that parent enter by
	pch byte
	//suffix
	link int

	gonow [255]int
}

var (
	//NMAX is total length of the templates
	t [NMAX + 1]Vertex

	//length of builth trie
	size int
)

//init initiates a trie
func init() {
	t[0].p, t[0].link = -1, -1

	//At the beginning there aren't any transitions
	for i, _ := range t[0].next {
		t[0].next[i] = -1
		t[0].gonow[i] = -1
	}

	//size of root
	size = 1
}

//add inserts a string to the trie
func add(s string, t []Vertex) {
	//Start from the root
	v := 0

	//Transition on every symbol of template
	for i := 0; i < len(s); i++ {
		var c byte = s[i]

		//Non-existence — creating a state
		if t[v].next[c] == -1 {
			for i := range t[size].next {
				//переход в дереве
				t[size].next[i] = -1

				//переход в автомате
				t[size].gonow[i] = -1
			}
			t[size].link = -1
			t[size].p = v
			t[size].pch = c

			//amount of allocated vertices
			t[v].next[c] = size
			size++
		}

		//Entering
		v = t[v].next[c]
	}
	t[v].leaf = true
}

func gonow(v int, c byte) int {
	if t[v].gonow[c] == -1 {
		if t[v].next[c] != -1 {
			t[v].gonow[c] = t[v].next[c]
		} else if v == 0 {
			t[v].gonow[c] = 0
		} else {
			gonow(getlink(v), c)
		}
	}
	return t[v].gonow[c]
}

func getlink(v int) int {
	if t[v].link == -1 {
		if v == 0 || t[v].p == 0 {
			t[v].link = 0
		} else {
			t[v].link = gonow(getlink(t[v].p), t[v].pch)
		}
	}
	return t[v].link
}

func main() {

	size = 1
	fmt.Println("Hello")
}

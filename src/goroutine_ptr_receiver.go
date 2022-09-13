package main

import (
	"fmt"
	"sync"
)

// START OMIT
type Foo int

func (u *Foo) Print(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(*u)
}

func main() {
	const n = 4
	foos := make([]Foo, n)
	for i := 0; i < n; i++ {
		foos[i] = Foo(i)
	}
	fmt.Println(foos)
	var wg sync.WaitGroup
	wg.Add(n)
	for _, foo := range foos {
		go foo.Print(&wg)
	}
	wg.Wait()
}

// END OMIT

package main

import (
	"fmt"
	"unsafe"
)

const uintSizeInBytes = 4 << ((^uint(0)) >> 32 & 1)

type User struct {
	admin       bool
	age         uint64
	active      bool
	nbFollowers uint64
	retired     bool
}

func main() {
	fmt.Println(uintSizeInBytes)
	fmt.Println(unsafe.Sizeof(User{}))
}

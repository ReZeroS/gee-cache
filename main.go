package main

import (
	"sync"
)

func test(a int, res sync.Map) {
	sum := -1
	key := a
	for a > 0 {
		sum *= a
		a--
	}

	res.Store(key, sum)
}
func main() {
	var res sync.Map
	a := 1
	res.Store(99999999999990, 109)
	for a <= 10 {
		go func(a int) {
			test(a, res)
		}(a)
		a++
	}
	//if v, ok := res.Load(5); ok {
	//	fmt.Println(v)
	//}
	//fmt . Println(res)
}

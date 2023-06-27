package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x []int
	fmt.Println(reflect.ValueOf(x).Kind())

	var y = make([]string, 40, 50)
	fmt.Println("y\tLen:", len(y), "\tCap:", cap(y))

	var z = []int{10, 20, 30, 40, 50}
	fmt.Println("z\tLen:", len(z), "\tCap:", cap(z))
	fmt.Println(z)

	var a = new([50]int)[0:5]
	fmt.Println("a\tLen:", len(a), "\tCap:", cap(a))
	fmt.Println(a)

	var b = make([]int, 1, 10)
	fmt.Println(b)
	b = append(b, 20)
	fmt.Println(b)

	var c = []int{20, 30, 40, 50}
	fmt.Println(c[0])
	fmt.Println(c[0:2])

	var d = []int{10, 20, 30, 40}
	fmt.Println(d)
	d[1] = 35
	fmt.Println(d)
	var e = []string{"krishna", "sai", "pavan"}
	fmt.Println(e)
	var f = []string{"vedantham"}
	e = append(e, f...)
	fmt.Println(e)

	var g = []int{44}
	var h = []int{22}
	copy(g, h)
	fmt.Println("G:", g)

---------------------------------------------------------------------------
PS D:\go\sample1> go run slice.go
slice
y       Len: 40         Cap: 50
z       Len: 5  Cap: 5
[10 20 30 40 50]
a       Len: 5  Cap: 50
[0 0 0 0 0]
[0]
[0 20]
20
[20 30]
[10 20 30 40]
[10 35 30 40]
[krishna sai pavan]
[krishna sai pavan vedantham]
G: [22]

	

}

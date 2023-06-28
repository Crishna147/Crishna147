package main

import (
	"fmt"
	"math"
)

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.radius
}

type Rectangle struct {
	width, height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func main() {
	circle := Circle{radius: 5}
	fmt.Println("Circle Area:", circle.Area())
	fmt.Println("Circle Circumference:", circle.Circumference())

	rectangle := Rectangle{width: 10, height: 5}
	fmt.Println("Rectangle Area:", rectangle.Area())
	fmt.Println("Rectangle Perimeter:", rectangle.Perimeter())
}
---------------------------------------------------------------------------
PS D:\go\sample1> go run structmethod.go
Circle Area: 78.53981633974483
Circle Circumference: 31.41592653589793
Rectangle Area: 50
Rectangle Perimeter: 30











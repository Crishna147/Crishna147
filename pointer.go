package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	person := Person{Name: "Crishna", Age: 22}
	personPtr := &person

	fmt.Println("Person:", person)
	fmt.Println("Memory address of person:", personPtr)
	fmt.Println("Name:", personPtr.Name)
	fmt.Println("Age:", personPtr.Age)
}
-------------------------------------------------------------
PS D:\go\sample1> go run pointers.go    
Person: {Crishna 22}
Memory address of person: &{Crishna 22}
Name: Crishna
Age: 22

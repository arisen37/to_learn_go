package main

import (
	"errors"
	"fmt"
)

type Greeter interface{ Greet() string }
type Person struct{ Name string }

func (p Person) Greet() string { return "Hello, " + p.Name }
func (p *Person) Rename(name string) error {
	if name == "" {
		return errors.New("empty name")
	}
	p.Name = name
	return nil
}
func greet(g Greeter) { fmt.Println(g.Greet()) }
func main() {
	var zero int
	fmt.Println("zero value:", zero)
	nums := []int{1, 2, 3}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	fmt.Println("sum:", sum)
	p := Person{Name: "Go learner"}
	greet(p)
	if err := p.Rename("Gopher"); err != nil {
		fmt.Println(err)
		return
	}
	greet(&p)
	labels := map[string]int{"easy": 7, "medium": 7, "hard": 7}
	fmt.Println("easy assignments:", labels["easy"])
}

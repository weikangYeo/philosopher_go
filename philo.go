package main

import "fmt"

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
	state     string
}

func (p *Philosopher) ShowCurrentState() string {
	return fmt.Sprintf("%s is %s", p.name, p.state)
}

func (p *Philosopher) Eat() {
	fmt.Println(p.name, " is eating")
}

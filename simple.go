package main

import (
	"fmt"
)

type Env map[string]Reducible

// Types
type Number struct {
	value int
}

type Add struct {
	left Reducible
	right Reducible
}

type Multiply struct {
	left Reducible
	right Reducible
}

type Boolean struct {
	value bool
}

type Reducible interface {
	Evaluate(env Env) (Reducible, Env)
	Reduce(env Env) (Reducible, Env)
	IsReducible() bool
}

// Evaluate is Implementations of Reducible
func (n Number) Evaluate(env Env) (Reducible, Env) {
	return n, env
}

func (b Boolean) Evaluate(env Env) (Reducible, Env) {
	return b, env
}

func (a Add) Evaluate(env Env) (Reducible, Env) {
	var rLeft, rRight Reducible
	rLeft, env = a.left.Evaluate(env)
	rRight, env = a.right.Evaluate(env)
	numLeft := rLeft.(Number)
	numRight := rRight.(Number)
	return Number{numLeft.value + numRight.value}, env
}

func (m Multiply) Evaluate(env Env) (Reducible, Env) {
	var rLeft, rRight Reducible
	rLeft, env = m.left.Evaluate(env)
	rRight, env = m.right.Evaluate(env)
	numLeft := rLeft.(Number)
	numRight := rRight.(Number)
	return Number{numLeft.value * numRight.value}, env
}

// Reduce is Implementations of Reducible
func (n Number) Reduce(env Env) (Reducible, Env) {
	return n, env
}

func (b Boolean) Reduce(env Env) (Reducible, Env) {
	return b, env
}

func (a Add) Reduce(env Env) (Reducible, Env) {
	if a.left.IsReducible() {
		var newleft Reducible
		newleft, env = a.left.Evaluate(env)
		return Add{newleft, a.right}, env

	} else if a.right.IsReducible() {
		var newright Reducible
		newright, env = a.right.Evaluate(env)
		return Add{a.left, newright}, env

	} else {
		return a.Evaluate(env)
	}
}

func (m Multiply) Reduce(env Env) (Reducible, Env) {
	if m.left.IsReducible() {
		var newleft Reducible
		newleft, env = m.left.Evaluate(env)
		return Multiply{newleft, m.right}, env

	} else if m.right.IsReducible() {
		var newright Reducible
		newright, env = m.right.Evaluate(env)
		return Multiply{m.left, newright}, env

	} else {
		return m.Evaluate(env)
	}
}

// IsReducible is function to determine if the value is reducible
func (n Number) IsReducible() bool {
	return false
}

func (b Boolean) IsReducible() bool {
	return false
}

func (a Add) IsReducible() bool {
	return true
}

func (m Multiply) IsReducible() bool {
	return true
}

// String returns a formula
func (n Number) String() string {
	return fmt.Sprintf("%d", n.value)
}

func (a Add) String() string {
	return fmt.Sprintf("%v + %v", a.left, a.right)
}

func (m Multiply) String() string {
	return fmt.Sprintf("%v * %v", m.left, m.right)
}

func (b Boolean) String() string {
	return fmt.Sprintf("%v", b.value)
}

// Macine
type Machine struct {
	expression Reducible
}

func (m *Machine) Step(env Env) Env {
	m.expression, env = m.expression.Reduce(env)
	return env
}

func (m *Machine) Run(env Env) Env {
	for m.expression.IsReducible() {
		fmt.Printf("%v\n", m.expression)
		env = m.Step(env)
	}
	fmt.Printf("%v\n", m.expression)
	return env
}

func main() {
	env := Env{}
	n1 := Number{1}
	n2 := Number{2}
	n3 := Number{3}
	n4 := Number{4}
	m1 := Multiply{n1, n2}
	m2 := Multiply{n3, n4}

	a := Add{m1, m2}

	machine := Machine{a}
	machine.Run(env)
}
// Result:
// 1 * 2 + 3 * 4
// 2 + 3 * 4
// 2 + 12
// 14
// false

// ref: https://github.com/quux00/Understanding.Computation.in.Go/blob/master/2meaning/src/simple/simple.go
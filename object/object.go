/*
We're going to create an object system to represent the values from the AST in memory!
an easy way to evaluate the parsed code! : JIT (Just In Time) Evaluation! not the best but many programming langauges started slow but got extented!

let a = 9;
...
a = a + 5; <-- We need to bind the int a to value "9" and when we come cross a we've to fetch its value from the memory, but we've to get the value 9 as a is represented as *ast.IntegarLiteral.

The point is this: there are a lot of diﬀerent ways to represent values of the interpreted languages in the host language.

Maybe the simplest way is to represent each value as Object.

*/

package object

import "fmt"

type ObjectType string

// why ? as each value needs different internal representations
type Object interface {
	Type() ObjectType
	Inspect() string
}

// datatypes
type Int struct {
	Value int64
}
type Bool struct {
	Value bool
}
type Null struct{}

type ReturnValue struct {
	Value Object
}

const (
	INT_OBJ    = "INT"
	BOOL_OBJ   = "BOOL"
	NULL_OBJ   = "NULL"
	RETURN_OBJ = "RETURN" // wrap the return value into an object
)

// --- Integar
func (i *Int) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Int) Type() ObjectType {
	return INT_OBJ
}

// --- Boolean
func (b *Bool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Bool) Type() ObjectType {
	return BOOL_OBJ
}

// --- Null
// Tony Hoare's “billion-dollar mistake”.
func (n *Null) Inspect() string {
	return "Null"
}
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}
func (r *ReturnValue) Type() ObjectType {
	return RETURN_OBJ
}

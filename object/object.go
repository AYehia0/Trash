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

import (
	"bytes"
	"fmt"
	"strings"
	"trash/ast"
)

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
type String struct {
	Value string
}
type Bool struct {
	Value bool
}
type Null struct{}

type ReturnValue struct {
	Value Object
}
type Error struct {
	Message string
}

const (
	INT_OBJ    = "INT"
	STR_OBJ    = "STRING"
	BOOL_OBJ   = "BOOL"
	NULL_OBJ   = "NULL"
	RETURN_OBJ = "RETURN" // wrap the return value into an object
	ERROR_OBJ  = "ERROR"
	FUNC_OBJ   = "FUNCTION"
)

// --- Integar
func (i *Int) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Int) Type() ObjectType {
	return INT_OBJ
}

// --- String
func (st *String) Inspect() string {
	return fmt.Sprintf("%s", st.Value)
}
func (st *String) Type() ObjectType {
	return STR_OBJ
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

// --- Return
func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}
func (r *ReturnValue) Type() ObjectType {
	return RETURN_OBJ
}

// --- Error
func (e *Error) Inspect() string {
	return "Error: " + e.Message
}
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

// the reason I am putting Env here to allow direct access of the Environment where function is defined in, this is useful for adding closures
type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Env
}

func (f *Function) Type() ObjectType {
	return FUNC_OBJ
}
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

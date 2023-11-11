/*
We're going to create an object system to represent the values from the AST in memory!
an easy way to evaluate the parsed code! : JIT (Just In Time) Evaluation! not the best but many programming langauges started slow but got extented!

let a = 9;
...
a = a + 5; <-- We need to bind the int a to value "9" and when we come cross a we've to fetch its value from the memory, but we've to get the value 9 as a is represented as *ast.IntegerLiteral.

The point is this: there are a lot of diﬀerent ways to represent values of the interpreted languages in the host language.

Maybe the simplest way is to represent each value as Object.

*/

package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
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
type BuiltinFuncs func(args ...Object) Object

type Builtin struct {
	Func BuiltinFuncs
}
type List struct {
	Values []Object // list can contain anything lol, JS are you happy now ?
}

// why hashing the key you may ask ?
/*
	key1 := &object.String{Value: "name"}
	value := &object.String{Value: "key"}

	pairs := map[object.Object]object.Object{}

	pairs[key1] = value
	fmt.Printf("pairs[key1]=%+v\n", pairs[key1])
	// => pairs[key1]=&{Value:value}

	key2 := &object.String{Value: "name"}
	fmt.Printf("pairs[key2]=%+v\n", pairs[key2])
	// => pairs[key2]=<nil>

	fmt.Printf("(key1 == key2)=%t\n", key1 == key2)
	// => (key1 == key2)=false

Hasing keys : Strings, Booleans, Integers
requirements: hashing the key must be unique, no two keys should have equal hashes.
*/

type HashKey struct {
	Type  ObjectType
	Value uint64
}

// used in the eval to check if the object is a usable as Hashkey when evaluating hash literals or indexing keys in hashmap.
type Hashable interface {
	HashKey() HashKey // TODO: Performance - Cache the return values
}

func (b *Bool) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (s *String) HashKey() HashKey {
	hash := fnv.New64a()
	hash.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: hash.Sum64()}
}

func (i *Int) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hashmap struct {
	Store map[HashKey]HashPair // to be able to print the key and value when Inspect(), also if later implementing something like iter (range)
}

const (
	INT_OBJ     = "INT"
	STR_OBJ     = "STRING"
	BOOL_OBJ    = "BOOL"
	NULL_OBJ    = "NULL"
	RETURN_OBJ  = "RETURN" // wrap the return value into an object
	ERROR_OBJ   = "ERROR"
	FUNC_OBJ    = "FUNCTION"
	BUILTIN_OBJ = "BUILTIN"
	LIST_OBJ    = "LIST"
	HASHMAP_OBJ = "HASH"
)

// --- Hashmap
func (hm *Hashmap) Type() ObjectType {
	return HASHMAP_OBJ
}
func (hm *Hashmap) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range hm.Store {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// --- List
func (ls *List) Type() ObjectType {
	return LIST_OBJ
}
func (ls *List) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ls.Values {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// --- Builtin functions
func (fn *Builtin) Inspect() string {
	return "Built-in Function"
}
func (fn *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

// --- Integer
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

/*
A solution for binded variables through expanding the enviroment :-

	let x = 10

	fn printNumber(x) {
		print(x)
	}

	printNumber(20) // 20
	print(x) // 10

args x will have a ref in the Env (GlobalEnv) of 10, so evaluating will print 10 instead of 20, and in case we overwritten x with 20, then 20 will be printed instead of 10!

for this case we must have a separate Env for each function (I know not the best solution, but it's pretty basic).

we create a new instance of object.Env with a pointer to the environment it should extend. By doing that we enclose a fresh and empty environment with an existing one.

When the new environment’s Get method is called and it itself doesn’t have a value associated
with the given name, it calls the Get of the enclosing environment. That’s the environment it’s
extending. And if that enclosing environment can’t find the value, it calls its own enclosing
environment and so on until there is no enclosing environment anymore and we can safely say
that we have an “ERROR: unknown identifier: foobar”.

aka : Scopes

	     Outer Env
	 -----------------
	|      x: 10      |
	|    Inner Env    |
	|    ---------    |
	|   |         |   |
	|   |  x: 20  |   |
	|   |         |   |
	|    ---------    |
	 -----------------
*/
package object

// --- Environment : used to keep track of assigned objects (basically a hashmap)
type Env struct {
	store map[string]Object
	outer *Env
}

func NewEnv() *Env {
	env := make(map[string]Object)
	return &Env{
		store: env,
		outer: nil,
	}
}

// getters and setters for our store
// get from the nearest env
func (env *Env) Get(key string) (Object, bool) {
	obj, ok := env.store[key]
	if !ok && env.outer != nil {
		obj, ok = env.outer.Get(key)
	}
	return obj, ok
}

func (env *Env) Set(key string, val Object) {
	env.store[key] = val
}

func NewEnclosedEnv(outerEnv *Env) *Env {
	env := NewEnv()
	env.outer = outerEnv

	return env
}

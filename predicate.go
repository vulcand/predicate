package predicate

type Def struct {
	Operators Operators
	Functions map[string]interface{}
}

type Operators struct {
	EQ  interface{}
	NEQ interface{}

	LT interface{}
	GT interface{}

	OR  interface{}
	AND interface{}
}

type Parser interface {
	Parse(string) (interface{}, error)
}

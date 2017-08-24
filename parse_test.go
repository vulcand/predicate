package predicate

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type PredicateSuite struct {
}

var _ = Suite(&PredicateSuite{})

func (s *PredicateSuite) getParser(c *C) Parser {
	return s.getParserWithOpts(c, nil, nil)
}

func (s *PredicateSuite) getParserWithOpts(c *C, getID GetIdentifierFn, getProperty GetPropertyFn) Parser {
	p, err := NewParser(Def{
		Operators: Operators{
			AND: numberAND,
			OR:  numberOR,
			GT:  numberGT,
			LT:  numberLT,
			EQ:  numberEQ,
			NEQ: numberNEQ,
			LE:  numberLE,
			GE:  numberGE,
		},
		Functions: map[string]interface{}{
			"DivisibleBy":        divisibleBy,
			"Remainder":          numberRemainder,
			"Len":                stringLength,
			"number.DivisibleBy": divisibleBy,
			"Equals":             equals,
		},
		GetIdentifier: getID,
		GetProperty:   getProperty,
	})
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)
	return p
}

func (s *PredicateSuite) TestSinglePredicate(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("DivisibleBy(2)")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(2))
	fn := pr.(numberPredicate)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, false)
}

func (s *PredicateSuite) TestModulePredicate(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("number.DivisibleBy(2)")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(2))
	fn := pr.(numberPredicate)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, false)
}

func (s *PredicateSuite) TestJoinAND(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("DivisibleBy(2) && DivisibleBy(3)")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(2), Equals, false)
	c.Assert(fn(3), Equals, false)
	c.Assert(fn(6), Equals, true)
}

func (s *PredicateSuite) TestJoinOR(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("DivisibleBy(2) || DivisibleBy(3)")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, true)
	c.Assert(fn(5), Equals, false)
}

func (s *PredicateSuite) TestGT(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("Remainder(3) > 1")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(1), Equals, false)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, false)
	c.Assert(fn(4), Equals, false)
	c.Assert(fn(5), Equals, true)
}

func (s *PredicateSuite) TestGTE(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("Remainder(3) >= 1")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(1), Equals, true)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, false)
	c.Assert(fn(4), Equals, true)
	c.Assert(fn(5), Equals, true)
}

func (s *PredicateSuite) TestLT(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("Remainder(3) < 2")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(1), Equals, true)
	c.Assert(fn(2), Equals, false)
	c.Assert(fn(3), Equals, true)
	c.Assert(fn(4), Equals, true)
	c.Assert(fn(5), Equals, false)
}

func (s *PredicateSuite) TestLE(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("Remainder(3) <= 2")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(1), Equals, true)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, true)
	c.Assert(fn(4), Equals, true)
	c.Assert(fn(5), Equals, true)
}

func (s *PredicateSuite) TestEQ(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("Remainder(3) == 2")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(1), Equals, false)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, false)
	c.Assert(fn(4), Equals, false)
	c.Assert(fn(5), Equals, true)
}

func (s *PredicateSuite) TestNEQ(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("Remainder(3) != 2")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(1), Equals, true)
	c.Assert(fn(2), Equals, false)
	c.Assert(fn(3), Equals, true)
	c.Assert(fn(4), Equals, true)
	c.Assert(fn(5), Equals, false)
}

func (s *PredicateSuite) TestParen(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("(Remainder(3) != 1) && (Remainder(3) != 0)")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(0), Equals, false)
	c.Assert(fn(1), Equals, false)
	c.Assert(fn(2), Equals, true)
}

func (s *PredicateSuite) TestStrings(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse(`Remainder(3) == Len("hi")`)
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(0), Equals, false)
	c.Assert(fn(1), Equals, false)
	c.Assert(fn(2), Equals, true)
}

func (s *PredicateSuite) TestGTFloat64(c *C) {
	p := s.getParser(c)

	pr, err := p.Parse("Remainder(3) > 1.2")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(1))
	fn := pr.(numberPredicate)
	c.Assert(fn(1), Equals, false)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, false)
	c.Assert(fn(4), Equals, false)
	c.Assert(fn(5), Equals, true)
}

func (s *PredicateSuite) TestIdentifier(c *C) {
	getID := func(fields []string) (interface{}, error) {
		c.Assert(fields, DeepEquals, []string{"first", "second", "third"})
		return 2, nil
	}
	p := s.getParserWithOpts(c, getID, nil)

	pr, err := p.Parse("DivisibleBy(first.second.third)")
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(2))
	fn := pr.(numberPredicate)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, false)
}

func (s *PredicateSuite) TestMap(c *C) {
	getID := func(fields []string) (interface{}, error) {
		c.Assert(fields, DeepEquals, []string{"first", "second"})
		return map[string]int{"key": 2}, nil
	}
	getProperty := func(mapVal, keyVal interface{}) (interface{}, error) {
		m := mapVal.(map[string]int)
		k := keyVal.(string)
		return m[k], nil
	}

	p := s.getParserWithOpts(c, getID, getProperty)

	pr, err := p.Parse(`DivisibleBy(first.second["key"])`)
	c.Assert(err, IsNil)
	c.Assert(pr, FitsTypeOf, divisibleBy(2))
	fn := pr.(numberPredicate)
	c.Assert(fn(2), Equals, true)
	c.Assert(fn(3), Equals, false)
}

func (s *PredicateSuite) TestIdentifierAndFunction(c *C) {
	getID := func(fields []string) (interface{}, error) {
		switch fields[0] {
		case "firstSlice":
			return []string{"a"}, nil
		case "secondSlice":
			return []string{"b"}, nil
		case "a":
			return "a", nil
		case "b":
			return "b", nil
		}
		return nil, nil
	}
	p := s.getParserWithOpts(c, getID, nil)

	pr, err := p.Parse("Equals(firstSlice, firstSlice)")
	c.Assert(err, IsNil)
	fn := pr.(boolPredicate)
	c.Assert(fn(), Equals, true)

	pr, err = p.Parse("Equals(a, a)")
	c.Assert(err, IsNil)
	fn = pr.(boolPredicate)
	c.Assert(fn(), Equals, true)

	pr, err = p.Parse("Equals(firstSlice, secondSlice)")
	c.Assert(err, IsNil)
	fn = pr.(boolPredicate)
	c.Assert(fn(), Equals, false)
}

func (s *PredicateSuite) TestUnhappyCases(c *C) {
	cases := []string{
		")(",                      // invalid expression
		"SomeFunc",                // unsupported id
		"Remainder(banana)",       // unsupported argument
		"Remainder(1, 2)",         // unsupported arguments count
		"Remainder(Len)",          // unsupported argument
		`Remainder(Len("Ho"))`,    // unsupported argument
		"Bla(1)",                  // unknown method call
		"0.2 && Remainder(1)",     // unsupported value
		`Len("Ho") && 0.2`,        // unsupported value
		"func(){}()",              // function call
		"Remainder(3) >> 3",       // unsupported operator
		`Remainder(3) > "banana"`, // unsupported comparison type
	}
	p := s.getParser(c)
	for _, expr := range cases {
		pr, err := p.Parse(expr)
		c.Assert(err, NotNil)
		c.Assert(pr, IsNil)
	}
}

type numberPredicate func(v int) bool
type numberMapper func(v int) int
type boolPredicate func() bool

// equals can compare complex objects, e.g. arrays of strings
// and strings together
func equals(a interface{}, b interface{}) boolPredicate {
	return func() bool {
		switch aval := a.(type) {
		case string:
			bval, ok := b.(string)
			return ok && aval == bval
		case []string:
			bval, ok := b.([]string)
			if !ok {
				return false
			}
			if len(aval) != len(bval) {
				return false
			}
			for i := range aval {
				if aval[i] != bval[i] {
					return false
				}
			}
			return true
		default:
			return false
		}
	}
}

func divisibleBy(divisor int) numberPredicate {
	return func(v int) bool {
		return v%divisor == 0
	}
}

func numberAND(a, b numberPredicate) numberPredicate {
	return func(v int) bool {
		return a(v) && b(v)
	}
}

func numberOR(a, b numberPredicate) numberPredicate {
	return func(v int) bool {
		return a(v) || b(v)
	}
}

func numberRemainder(divideBy int) numberMapper {
	return func(v int) int {
		return v % divideBy
	}
}

func numberGT(m numberMapper, value interface{}) (numberPredicate, error) {
	switch value.(type) {
	case int:
	case float64:
	default:
		return nil, fmt.Errorf("GT: unsupported argument type: %T", value)
	}
	return func(v int) bool {
		switch val := value.(type) {
		case int:
			return m(v) > val
		case float64:
			return m(v) > int(val)
		default:
			return true
		}
	}, nil
}

func numberGE(m numberMapper, value int) (numberPredicate, error) {
	return func(v int) bool {
		return m(v) >= value
	}, nil
}

func numberLE(m numberMapper, value int) (numberPredicate, error) {
	return func(v int) bool {
		return m(v) <= value
	}, nil
}

func numberLT(m numberMapper, value int) numberPredicate {
	return func(v int) bool {
		return m(v) < value
	}
}

func numberEQ(m numberMapper, value int) numberPredicate {
	return func(v int) bool {
		return m(v) == value
	}
}

func numberNEQ(m numberMapper, value int) numberPredicate {
	return func(v int) bool {
		return m(v) != value
	}
}

func stringLength(v string) int {
	return len(v)
}

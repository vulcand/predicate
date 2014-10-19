package predicate

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type PredicateSuite struct {
}

var _ = Suite(&PredicateSuite{})

func (s *PredicateSuite) getParser(c *C) Parser {
	p, err := NewParser(Def{
		Operators: Operators{
			AND: numberAND,
			OR:  numberOR,
			GT:  numberGT,
			LT:  numberLT,
		},
		Functions: map[string]interface{}{
			"DivisibleBy": divisibleBy,
			"Remainder":   numberRemainder,
		},
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

type numberPredicate func(v int) bool
type numberMapper func(v int) int

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

func numberGT(m numberMapper, value int) numberPredicate {
	return func(v int) bool {
		return m(v) > value
	}
}

func numberLT(m numberMapper, value int) numberPredicate {
	return func(v int) bool {
		return m(v) < value
	}
}

package predicate

import (
	"fmt"
	"testing"

	"github.com/gravitational/trace"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	suite.Run(t, new(PredicateSuite))
}

type PredicateSuite struct {
	suite.Suite
}

func (s *PredicateSuite) getParser() Parser {
	return s.getParserWithOpts(nil, nil)
}

func (s *PredicateSuite) getParserWithOpts(getID GetIdentifierFn, getProperty GetPropertyFn) Parser {
	s.T().Helper()

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
			NOT: numberNOT,
		},
		Functions: map[string]interface{}{
			"DivisibleBy":        divisibleBy,
			"Remainder":          numberRemainder,
			"Len":                stringLength,
			"number.DivisibleBy": divisibleBy,
			"Equals":             Equals,
			"Contains":           Contains,
			"fnreturn": func(arg interface{}) (interface{}, error) {
				return arg, nil
			},
			"fnerr": func(arg interface{}) (interface{}, error) {
				return nil, trace.BadParameter("don't like this parameter")
			},
		},
		GetIdentifier: getID,
		GetProperty:   getProperty,
	})

	s.NoError(err)
	s.NotNil(p)

	return p
}

func (s *PredicateSuite) TestSinglePredicate() {
	p := s.getParser()

	pr, err := p.Parse("DivisibleBy(2)")
	s.NoError(err)

	s.IsType(divisibleBy(2), pr)

	fn := pr.(numberPredicate)
	s.True(fn(2))
	s.False(fn(3))
}

func (s *PredicateSuite) TestSinglePredicateNot() {
	p := s.getParser()

	pr, err := p.Parse("!DivisibleBy(2)")
	s.NoError(err)

	s.IsType(divisibleBy(2), pr)

	fn := pr.(numberPredicate)
	s.False(fn(2))
	s.True(fn(3))
}

func (s *PredicateSuite) TestSinglePredicateWithFunc() {
	p := s.getParser()

	pr, err := p.Parse("DivisibleBy(fnreturn(2))")
	s.Require().NoError(err)

	s.IsType(divisibleBy(2), pr)

	fn := pr.(numberPredicate)
	s.True(fn(2))
	s.False(fn(3))
}

func (s *PredicateSuite) TestSinglePredicateWithFuncErr() {
	p := s.getParser()

	_, err := p.Parse("DivisibleBy(fnerr(2))")
	s.Error(err)
}

func (s *PredicateSuite) TestModulePredicate() {
	p := s.getParser()

	pr, err := p.Parse("number.DivisibleBy(2)")
	s.NoError(err)

	s.IsType(divisibleBy(2), pr)

	fn := pr.(numberPredicate)
	s.True(fn(2))
	s.False(fn(3))
}

func (s *PredicateSuite) TestJoinAND() {
	p := s.getParser()

	pr, err := p.Parse("DivisibleBy(2) && DivisibleBy(3)")
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.False(fn(2))
	s.False(fn(3))
	s.True(fn(6))
}

func (s *PredicateSuite) TestJoinOR() {
	p := s.getParser()

	pr, err := p.Parse("DivisibleBy(2) || DivisibleBy(3)")
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.True(fn(2))
	s.True(fn(3))
	s.False(fn(5))
}

func (s *PredicateSuite) TestGT() {
	p := s.getParser()

	pr, err := p.Parse("Remainder(3) > 1")
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.False(fn(1))
	s.True(fn(2))
	s.False(fn(3))
	s.False(fn(4))
	s.True(fn(5))
}

func (s *PredicateSuite) TestGTE() {
	p := s.getParser()

	pr, err := p.Parse("Remainder(3) >= 1")
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.True(fn(1))
	s.True(fn(2))
	s.False(fn(3))
	s.True(fn(4))
	s.True(fn(5))
}

func (s *PredicateSuite) TestLT() {
	p := s.getParser()

	pr, err := p.Parse("Remainder(3) < 2")
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.True(fn(1))
	s.False(fn(2))
	s.True(fn(3))
	s.True(fn(4))
	s.False(fn(5))
}

func (s *PredicateSuite) TestLE() {
	p := s.getParser()

	pr, err := p.Parse("Remainder(3) <= 2")
	s.NoError(err)
	s.IsType(divisibleBy(1), pr)
	fn := pr.(numberPredicate)
	s.True(fn(1))
	s.True(fn(2))
	s.True(fn(3))
	s.True(fn(4))
	s.True(fn(5))
}

func (s *PredicateSuite) TestEQ() {
	p := s.getParser()

	pr, err := p.Parse("Remainder(3) == 2")
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.False(fn(1))
	s.True(fn(2))
	s.False(fn(3))
	s.False(fn(4))
	s.True(fn(5))
}

func (s *PredicateSuite) TestNEQ() {
	p := s.getParser()

	pr, err := p.Parse("Remainder(3) != 2")
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.True(fn(1))
	s.False(fn(2))
	s.True(fn(3))
	s.True(fn(4))
	s.False(fn(5))
}

func (s *PredicateSuite) TestParen() {
	p := s.getParser()

	pr, err := p.Parse("(Remainder(3) != 1) && (Remainder(3) != 0)")
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.False(fn(0))
	s.False(fn(1))
	s.True(fn(2))
}

func (s *PredicateSuite) TestStrings() {
	p := s.getParser()

	pr, err := p.Parse(`Remainder(3) == Len("hi")`)
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.False(fn(0))
	s.False(fn(1))
	s.True(fn(2))
}

func (s *PredicateSuite) TestGTFloat64() {
	p := s.getParser()

	pr, err := p.Parse("Remainder(3) > 1.2")
	s.NoError(err)

	s.IsType(divisibleBy(1), pr)

	fn := pr.(numberPredicate)
	s.False(fn(1))
	s.True(fn(2))
	s.False(fn(3))
	s.False(fn(4))
	s.True(fn(5))
}

func (s *PredicateSuite) TestSelectExpr() {
	getID := func(fields []string) (interface{}, error) {
		s.Equal([]string{"first", "second", "third"}, fields)
		return 2, nil
	}
	p := s.getParserWithOpts(getID, nil)

	// Test selector expression.
	pr, err := p.Parse("Remainder(4) <= first.second.third")
	s.NoError(err)

	s.IsType(divisibleBy(2), pr)

	fn := pr.(numberPredicate)
	s.True(fn(2))
	s.False(fn(3))

	// Test selector expression inside call expression.
	pr, err = p.Parse("DivisibleBy(first.second.third)")
	s.NoError(err)

	s.IsType(divisibleBy(2), pr)

	fn = pr.(numberPredicate)
	s.True(fn(2))
	s.False(fn(3))
}

func (s *PredicateSuite) TestIndexExpr() {
	getID := func(fields []string) (interface{}, error) {
		s.Equal([]string{"first", "second"}, fields)
		return map[string]int{"key": 2}, nil
	}
	getProperty := func(mapVal, keyVal interface{}) (interface{}, error) {
		m := mapVal.(map[string]int)
		k := keyVal.(string)
		return m[k], nil
	}

	p := s.getParserWithOpts(getID, getProperty)

	// Test index expression.
	pr, err := p.Parse(`Remainder(4) <= first.second["key"]`)
	s.NoError(err)

	s.IsType(divisibleBy(2), pr)

	fn := pr.(numberPredicate)
	s.True(fn(2))
	s.False(fn(3))

	// Test index expression inside call expression.
	pr, err = p.Parse(`DivisibleBy(first.second["key"])`)
	s.NoError(err)

	s.IsType(divisibleBy(2), pr)

	fn = pr.(numberPredicate)
	s.True(fn(2))
	s.False(fn(3))
}

func (s *PredicateSuite) TestIdentifierExpr() {
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
		case "num":
			return 2, nil
		}
		return nil, nil
	}
	p := s.getParserWithOpts(getID, nil)

	pr, err := p.Parse("Equals(firstSlice, firstSlice)")
	s.NoError(err)
	fn := pr.(BoolPredicate)
	s.True(fn())

	pr, err = p.Parse("Equals(a, a)")
	s.NoError(err)
	fn = pr.(BoolPredicate)
	s.True(fn())

	pr, err = p.Parse("Equals(firstSlice, secondSlice)")
	s.NoError(err)

	fn = pr.(BoolPredicate)
	s.False(fn())

	pr, err = p.Parse("Remainder(4) <= num")
	s.NoError(err)
	fn2 := pr.(numberPredicate)
	s.True(fn2(2))
	s.False(fn2(3))
}

func (s *PredicateSuite) TestContains() {
	val := TestStruct{}
	val.Param.Key1 = map[string][]string{"key": {"a", "b", "c"}}

	getID := func(fields []string) (interface{}, error) {
		return GetFieldByTag(val, "json", fields[1:])
	}
	p := s.getParserWithOpts(getID, GetStringMapValue)

	pr, err := p.Parse(`Contains(val.param.key1["key"], "a")`)
	s.NoError(err)
	s.True(pr.(BoolPredicate)())

	pr, err = p.Parse(`Contains(val.param.key1["key"], "z")`)
	s.NoError(err)
	s.False(pr.(BoolPredicate)())

	pr, err = p.Parse(`Contains(val.param.key1["missing"], "a")`)
	s.NoError(err)
	s.False(pr.(BoolPredicate)())
}

func (s *PredicateSuite) TestEquals() {
	val := TestStruct{}
	val.Param.Key2 = map[string]string{"key": "a"}

	getID := func(fields []string) (interface{}, error) {
		return GetFieldByTag(val, "json", fields[1:])
	}
	p := s.getParserWithOpts(getID, GetStringMapValue)

	pr, err := p.Parse(`Equals(val.param.key2["key"], "a")`)
	s.NoError(err)
	s.True(pr.(BoolPredicate)())

	pr, err = p.Parse(`Equals(val.param.key2["key"], "b")`)
	s.NoError(err)
	s.False(pr.(BoolPredicate)())

	pr, err = p.Parse(`Contains(val.param.key2["missing"], "z")`)
	s.NoError(err)
	s.False(pr.(BoolPredicate)())

	pr, err = p.Parse(`Contains(val.param.key1["missing"], "z")`)
	s.NoError(err)
	s.False(pr.(BoolPredicate)())
}

// TestStruct is a test structure with json tags.
type TestStruct struct {
	Param struct {
		Key1 map[string][]string `json:"key1,omitempty"`
		Key2 map[string]string   `json:"key2,omitempty"`
	} `json:"param,omitempty"`
}

func (s *PredicateSuite) TestGetTagField() {
	val := TestStruct{}
	val.Param.Key1 = map[string][]string{"key": {"val"}}

	type testCase struct {
		tag    string
		fields []string
		val    interface{}
		expect interface{}
		err    error
	}
	testCases := []testCase{
		// nested field
		{tag: "json", val: val, fields: []string{"param", "key1"}, expect: val.Param.Key1},
		// pointer to struct
		{tag: "json", val: &val, fields: []string{"param", "key1"}, expect: val.Param.Key1},
		// not found field
		{tag: "json", val: &val, fields: []string{"param", "key3"}, err: trace.NotFound("not found")},
		// nil pointer
		{tag: "json", val: nil, fields: []string{"param", "key1"}, err: trace.BadParameter("bad param")},
	}

	for i, tc := range testCases {
		comment := fmt.Sprintf("test case %d", i)

		out, err := GetFieldByTag(tc.val, tc.tag, tc.fields)
		if tc.err != nil {
			s.IsType(tc.err, err)
		} else {
			s.NoError(err, comment)
			s.Equal(tc.expect, out, comment)
		}
	}
}

func (s *PredicateSuite) TestUnhappyCases() {
	cases := []string{
		")(",                      // invalid expression
		"SomeFunc",                // unsupported id
		"Remainder(banana)",       // unsupported argument
		"Remainder(1, 2)",         // unsupported arguments count
		"Remainder(Len)",          // unsupported argument
		"Bla(1)",                  // unknown method call
		"0.2 && Remainder(1)",     // unsupported value
		`Len("Ho") && 0.2`,        // unsupported value
		"func(){}()",              // function call
		"Remainder(3) >> 3",       // unsupported operator
		`Remainder(3) > "banana"`, // unsupported comparison type
	}
	p := s.getParser()
	for _, expr := range cases {
		pr, err := p.Parse(expr)
		s.Error(err)
		s.Nil(pr)
	}
}

type (
	numberPredicate func(v int) bool
	numberMapper    func(v int) int
)

func divisibleBy(divisor int) numberPredicate {
	return func(v int) bool {
		return v%divisor == 0
	}
}

func numberNOT(a numberPredicate) numberPredicate {
	return func(v int) bool {
		return !a(v)
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

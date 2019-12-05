package sexp

import (
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestParse(t *testing.T) {
	testData := []struct {
		Name     string
		Pattern  string
		Expected *Sexp
	}{
		{
			Name:    "pattern 1 - empty",
			Pattern: `()`,
			Expected: &Sexp{
				Atom:     nil,
				Children: []*Sexp{},
			},
		},
		{
			Name:    "pattern 2 - only a symbol",
			Pattern: `(a)`,
			Expected: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "a"}},
				},
			},
		},
		{
			Name:    "pattern 3 - list",
			Pattern: `(a b c)`,
			Expected: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "a"}},
					{Atom: &Token{Type: TokenTypeSymbol, Value: "b"}},
					{Atom: &Token{Type: TokenTypeSymbol, Value: "c"}},
				},
			},
		},
		{
			Name:    "pattern 4 - nest",
			Pattern: `(a (b c))`,
			Expected: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "a"}},
					{
						Atom: nil,
						Children: []*Sexp{
							{Atom: &Token{Type: TokenTypeSymbol, Value: "b"}},
							{Atom: &Token{Type: TokenTypeSymbol, Value: "c"}},
						},
					},
				},
			},
		},
		{
			Name:     "pattern 5 - no parens",
			Pattern:  "",
			Expected: nil,
		},
		{
			Name:     "pattern 6 - incomplete",
			Pattern:  "(a",
			Expected: nil,
		},
		{
			Name:     "pattern 7 - incomplete, nested",
			Pattern:  "(a (b) (",
			Expected: nil,
		},
		{
			Name:    "pattern - wast 1",
			Pattern: `(assert_return (invoke "add" (i32.const 1) (i32.const 1)) (i32.const 2))`,
			Expected: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "assert_return"}},
					{
						Atom: nil,
						Children: []*Sexp{
							{Atom: &Token{Type: TokenTypeSymbol, Value: "invoke"}},
							{Atom: &Token{Type: TokenTypeString, Value: `"add"`}},
							{
								Atom: nil,
								Children: []*Sexp{
									{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
									{Atom: &Token{Type: TokenTypeNumber, Value: "1"}},
								},
							},
							{
								Atom: nil,
								Children: []*Sexp{
									{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
									{Atom: &Token{Type: TokenTypeNumber, Value: "1"}},
								},
							},
						},
					},
					{
						Atom: nil,
						Children: []*Sexp{
							{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
							{Atom: &Token{Type: TokenTypeNumber, Value: "2"}},
						},
					},
				},
			},
		},
		{
			Name:    "pattern - wast 2",
			Pattern: `(assert_trap (invoke "div_s" (i32.const 1) (i32.const 0)) "integer divide by zero")`,
			Expected: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "assert_trap"}},
					{
						Atom: nil,
						Children: []*Sexp{
							{Atom: &Token{Type: TokenTypeSymbol, Value: "invoke"}},
							{Atom: &Token{Type: TokenTypeString, Value: `"div_s"`}},
							{
								Atom: nil,
								Children: []*Sexp{
									{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
									{Atom: &Token{Type: TokenTypeNumber, Value: "1"}},
								},
							},
							{
								Atom: nil,
								Children: []*Sexp{
									{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
									{Atom: &Token{Type: TokenTypeNumber, Value: "0"}},
								},
							},
						},
					},
					{
						Atom: &Token{Type: TokenTypeString, Value: `"integer divide by zero"`},
					},
				},
			},
		},
	}

	for _, data := range testData {
		data := data // capture
		t.Run(data.Name, func(t *testing.T) {
			//t.Parallel()

			a, err := Parse(data.Pattern)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(data.Expected, a) {
				t.Fatalf("\n%s", pretty.Compare(data.Expected, a))
			}
		})
	}
}

func TestSexpString(t *testing.T) {
	testData := []struct {
		Name     string
		Sexp     *Sexp
		Expected string
	}{
		{
			Name: "pattern 1 - empty",
			Sexp: &Sexp{
				Atom:     nil,
				Children: []*Sexp{},
			},
			Expected: `()`,
		},
		{
			Name: "pattern 2 - only a symbol",
			Sexp: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "a"}},
				},
			},
			Expected: `(a)`,
		},
		{
			Name: "pattern 3 - list",
			Sexp: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "a"}},
					{Atom: &Token{Type: TokenTypeSymbol, Value: "b"}},
					{Atom: &Token{Type: TokenTypeSymbol, Value: "c"}},
				},
			},
			Expected: `(a b c)`,
		},
		{
			Name: "pattern 4 - nest",
			Sexp: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "a"}},
					{
						Atom: nil,
						Children: []*Sexp{
							{Atom: &Token{Type: TokenTypeSymbol, Value: "b"}},
							{Atom: &Token{Type: TokenTypeSymbol, Value: "c"}},
						},
					},
				},
			},
			Expected: `(a (b c))`,
		},
		{
			Name: "pattern - wast 1",
			Sexp: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "assert_return"}},
					{
						Atom: nil,
						Children: []*Sexp{
							{Atom: &Token{Type: TokenTypeSymbol, Value: "invoke"}},
							{Atom: &Token{Type: TokenTypeString, Value: `"add"`}},
							{
								Atom: nil,
								Children: []*Sexp{
									{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
									{Atom: &Token{Type: TokenTypeNumber, Value: "1"}},
								},
							},
							{
								Atom: nil,
								Children: []*Sexp{
									{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
									{Atom: &Token{Type: TokenTypeNumber, Value: "1"}},
								},
							},
						},
					},
					{
						Atom: nil,
						Children: []*Sexp{
							{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
							{Atom: &Token{Type: TokenTypeNumber, Value: "2"}},
						},
					},
				},
			},
			Expected: `(assert_return (invoke "add" (i32.const 1) (i32.const 1)) (i32.const 2))`,
		},
		{
			Name: "pattern - wast 2",
			Sexp: &Sexp{
				Atom: nil,
				Children: []*Sexp{
					{Atom: &Token{Type: TokenTypeSymbol, Value: "assert_trap"}},
					{
						Atom: nil,
						Children: []*Sexp{
							{Atom: &Token{Type: TokenTypeSymbol, Value: "invoke"}},
							{Atom: &Token{Type: TokenTypeString, Value: `"div_s"`}},
							{
								Atom: nil,
								Children: []*Sexp{
									{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
									{Atom: &Token{Type: TokenTypeNumber, Value: "1"}},
								},
							},
							{
								Atom: nil,
								Children: []*Sexp{
									{Atom: &Token{Type: TokenTypeSymbol, Value: "i32.const"}},
									{Atom: &Token{Type: TokenTypeNumber, Value: "0"}},
								},
							},
						},
					},
					{
						Atom: &Token{Type: TokenTypeString, Value: `"integer divide by zero"`},
					},
				},
			},
			Expected: `(assert_trap (invoke "div_s" (i32.const 1) (i32.const 0)) "integer divide by zero")`,
		},
	}

	for _, data := range testData {
		data := data // capture
		t.Run(data.Name, func(t *testing.T) {
			//t.Parallel()

			a := data.Sexp.String()
			if data.Expected != a {
				t.Fatalf("\n%s", pretty.Compare(data.Expected, a))
			}
		})
	}
}

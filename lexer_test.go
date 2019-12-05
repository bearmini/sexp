package sexp

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestNextToken(t *testing.T) {
	testData := []struct {
		Name     string
		Pattern  string
		Expected []*Token
	}{
		{
			Name:    "pattern 1",
			Pattern: "a b c",
			Expected: []*Token{
				{Type: TokenTypeSymbol, Value: "a"},
				{Type: TokenTypeSymbol, Value: "b"},
				{Type: TokenTypeSymbol, Value: "c"},
			},
		},
		{
			Name:    "pattern 2",
			Pattern: `(assert_return (invoke "add" (i32.const 1) (i32.const 1)) (i32.const 2))`,
			Expected: []*Token{
				{Type: TokenTypeOpenParen, Value: "("},
				{Type: TokenTypeSymbol, Value: "assert_return"},
				{Type: TokenTypeOpenParen, Value: "("},
				{Type: TokenTypeSymbol, Value: "invoke"},
				{Type: TokenTypeString, Value: `"add"`},
				{Type: TokenTypeOpenParen, Value: "("},
				{Type: TokenTypeSymbol, Value: "i32.const"},
				{Type: TokenTypeNumber, Value: "1"},
				{Type: TokenTypeCloseParen, Value: ")"},
				{Type: TokenTypeOpenParen, Value: "("},
				{Type: TokenTypeSymbol, Value: "i32.const"},
				{Type: TokenTypeNumber, Value: "1"},
				{Type: TokenTypeCloseParen, Value: ")"},
				{Type: TokenTypeCloseParen, Value: ")"},
				{Type: TokenTypeOpenParen, Value: "("},
				{Type: TokenTypeSymbol, Value: "i32.const"},
				{Type: TokenTypeNumber, Value: "2"},
				{Type: TokenTypeCloseParen, Value: ")"},
				{Type: TokenTypeCloseParen, Value: ")"},
			},
		},
	}

	for _, data := range testData {
		data := data // capture
		t.Run(data.Name, func(t *testing.T) {
			//t.Parallel()

			lex := NewLexer(strings.NewReader(data.Pattern))
			as := []*Token{}
			for {
				a := lex.NextToken()
				if a == nil {
					break
				}
				as = append(as, a)
			}
			if !reflect.DeepEqual(data.Expected, as) {
				t.Fatalf("%s", pretty.Compare(data.Expected, as))
			}
		})
	}
}

func TestReadSymbol(t *testing.T) {
	testData := []struct {
		Name     string
		Pattern  string
		Expected string
	}{
		{
			Name:     "pattern 1 - 1 character",
			Pattern:  "a b c",
			Expected: "a",
		},
		{
			Name:     "pattern 2 - symbol using _",
			Pattern:  `assert_return (invoke "add"`,
			Expected: "assert_return",
		},
	}

	for _, data := range testData {
		data := data // capture
		t.Run(data.Name, func(t *testing.T) {
			//t.Parallel()

			a, err := readSymbol(bufio.NewReader(strings.NewReader(data.Pattern)))
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if data.Expected != a {
				t.Fatalf("\nExpected: %s\nActual:   %s", data.Expected, a)
			}
		})
	}
}

func TestReadNumber(t *testing.T) {
	testData := []struct {
		Name     string
		Pattern  string
		Expected string
	}{
		{
			Name:     "pattern 1 - normal",
			Pattern:  `"123 b c"`,
			Expected: `"123 b c"`,
		},
		{
			Name:     "pattern 2 - escaped",
			Pattern:  `"he said \"hello\""`,
			Expected: `"he said \"hello\""`,
		},
	}

	for _, data := range testData {
		data := data // capture
		t.Run(data.Name, func(t *testing.T) {
			//t.Parallel()

			a, err := readString(bufio.NewReader(strings.NewReader(data.Pattern)))
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if data.Expected != a {
				t.Fatalf("\nExpected: %s\nActual:   %s", data.Expected, a)
			}
		})
	}
}

func TestUnread(t *testing.T) {
	lex := NewLexer(strings.NewReader("(a b c)"))

	// read the first token '('
	token := lex.NextToken()
	if token == nil {
		t.Fatalf("expected an open paren")
	}
	if token.Type != TokenTypeOpenParen {
		t.Fatalf("expected an open paren, but got %s", token.Type.String())
	}
	if token.Value != "(" {
		t.Fatalf("expected '(', but got %s", token.Value)
	}

	// unread it
	err := lex.Unread()
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	// read the first token '(' again
	token = lex.NextToken()
	if token == nil {
		t.Fatalf("expected an open paren")
	}
	if token.Type != TokenTypeOpenParen {
		t.Fatalf("expected an open paren, but got %s", token.Type.String())
	}
	if token.Value != "(" {
		t.Fatalf("expected '(', but got %s", token.Value)
	}

	// read the second token 'a' again
	token = lex.NextToken()
	if token == nil {
		t.Fatalf("expected a symbol")
	}
	if token.Type != TokenTypeSymbol {
		t.Fatalf("expected a symbol, but got %s", token.Type.String())
	}
	if token.Value != "a" {
		t.Fatalf("expected 'a', but got %s", token.Value)
	}

	err = lex.Unread()
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	err = lex.Unread()
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	token = lex.NextToken()
	if token == nil {
		t.Fatalf("expected an open paren")
	}
	if token.Type != TokenTypeOpenParen {
		t.Fatalf("expected an open paren, but got %s", token.Type.String())
	}
	if token.Value != "(" {
		t.Fatalf("expected '(', but got %s", token.Value)
	}

	token = lex.NextToken()
	if token == nil {
		t.Fatalf("expected a symbol")
	}
	if token.Type != TokenTypeSymbol {
		t.Fatalf("expected a symbol, but got %s", token.Type.String())
	}
	if token.Value != "a" {
		t.Fatalf("expected 'a', but got %s", token.Value)
	}

}


func TestReadString(t *testing.T) {
	testData := []struct {
		Name     string
		Pattern  string
		Expected string
	}{
		{
			Name:     "pattern 1 - integer",
			Pattern:  "123 b c",
			Expected: "123",
		},
		{
			Name:     "pattern 2 - negative integer",
			Pattern:  "-123 234 56",
			Expected: "-123",
		},
		{
			Name:     "pattern 3 - hex",
			Pattern:  "0x00000001 0x00000002 0x00000003",
			Expected: "0x00000001",
		},
		{
			Name:     "pattern 4 - floating number",
			Pattern:  "-0.0 6.28318",
			Expected: "-0.0",
		},
		{
			Name:     "pattern 5 - floating number with exp",
			Pattern:  "6.023e23 -0.0 6.28318",
			Expected: "6.023e23",
		},
	}

	for _, data := range testData {
		data := data // capture
		t.Run(data.Name, func(t *testing.T) {
			//t.Parallel()

			a, err := readSymbol(bufio.NewReader(strings.NewReader(data.Pattern)))
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if data.Expected != a {
				t.Fatalf("\nExpected: %s\nActual:   %s", data.Expected, a)
			}
		})
	}
}

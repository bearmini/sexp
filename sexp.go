package sexp

import (
	"strings"

	"github.com/pkg/errors"
)

type Sexp struct {
	Atom *Token
	Children []*Sexp
}

func Parse(str string) (*Sexp, error) {
	l := NewLexer(strings.NewReader(str))
	return parse(l)
}

func parse(l *Lexer) (*Sexp, error) {
	token := l.NextToken()
	switch token.Type {
	case TokenTypeSymbol, TokenTypeString, TokenTypeNumber:
		return &Sexp{Atom: token}, nil
	}

	if token.Type != TokenTypeOpenParen {
		return nil, errors.Errorf("expected open paren, but found %s", token.Type.String())
	}

	children := []*Sexp{}

loop:
	for {
		token = l.NextToken()
		if token == nil {
			break
		}

		switch token.Type {
		case TokenTypeOpenParen:
			err := l.Unread()
			if err != nil {
				return nil, err
			}
			s, err := parse(l)
			if err != nil {
				return nil, err
			}
			children = append(children, s)
		case TokenTypeSymbol, TokenTypeString, TokenTypeNumber:
			children = append(children, &Sexp{Atom: token})
		case TokenTypeCloseParen:
			break loop
		}
	}

	return &Sexp{Children: children}, nil
}

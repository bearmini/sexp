package sexp

import (
	"bufio"
	"errors"
	"io"
	"unicode"
)

type Lexer struct {
	br *bufio.Reader
	history []*Token
	unread []*Token
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		br: bufio.NewReader(r),
		history: []*Token{},
		unread: []*Token{},
	}
}

func (lex *Lexer) Unread() error {
	if len(lex.history) == 0 {
		return errors.New("unable to unread")
	}

	last := len(lex.history) - 1
	lex.unread = append(lex.unread, lex.history[last])
	lex.history = lex.history[:last]
	return nil
}

func (lex *Lexer) Peek() *Token {
	token := lex.NextToken()
	if token == nil {
		return nil
	}

	err := lex.Unread()
	if err != nil {
		return nil
	}

	return token
}

func (lex *Lexer) NextToken() *Token {
	if len(lex.unread) > 0 {
		last := len(lex.unread) - 1
		token := lex.unread[last]
		lex.unread = lex.unread[:last]
		lex.history = append(lex.history, token)
		return token
	}

	var token *Token
loop:
	for {
		r, _, err := lex.br.ReadRune()
		if err != nil {
			return nil
		}

		if unicode.IsSpace(r) {
			continue
		}

		switch {
		case r == '(':
			token = &Token{
				Type:  TokenTypeOpenParen,
				Value: string([]rune{r}),
			}
			break loop
		case r == ')':
			token = &Token{
				Type:  TokenTypeCloseParen,
				Value: string([]rune{r}),
			}
			break loop
		case isNumberStartRune(r):
			err = lex.br.UnreadRune()
			if err != nil {
				return nil
			}
			s, err := readNumber(lex.br)
			if err != nil {
				return nil
			}
			token = &Token{
				Type:  TokenTypeNumber,
				Value: s,
			}
			break loop
		case r == '"':
			err = lex.br.UnreadRune()
			if err != nil {
				return nil
			}
			s, err := readString(lex.br)
			if err != nil {
				return nil
			}
			token = &Token{
				Type:  TokenTypeString,
				Value: s,
			}
			break loop
		case isSymbolRune(r):
			err = lex.br.UnreadRune()
			if err != nil {
				return nil
			}
			s, err := readSymbol(lex.br)
			if err != nil {
				return nil
			}
			token = &Token{
				Type:  TokenTypeSymbol,
				Value: s,
			}
			break loop
		}
	}
	lex.history = append(lex.history, token)
	return token
}

func isSymbolRune(r rune) bool {
	if r == '(' || r == ')' {
		return false
	}
	if r == '"' {
		return false
	}
	if unicode.IsSpace(r) {
		return false
	}
	return true
}

func isNumberStartRune(r rune) bool {
	if unicode.IsNumber(r) {
		return true
	}

	if r == '-' {
		return true
	}

	return false
}

func isNumberRune(r rune) bool {
	if unicode.IsNumber(r) {
		return true
	}

	if r == '-' || r == '.' || r == 'e' {
		return true
	}

	return false
}

func readSymbol(br *bufio.Reader) (string, error) {
	buf := []rune{}
	for {
		r, _, err := br.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if !isSymbolRune(r) {
			err = br.UnreadRune()
			if err != nil {
				return "", err
			}
			break
		}
		buf = append(buf, r)
	}
	return string(buf), nil
}

func readNumber(br *bufio.Reader) (string, error) {
	buf := []rune{}
	for {
		r, _, err := br.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if !isNumberRune(r) {
			err = br.UnreadRune()
			if err != nil {
				return "", err
			}
			break
		}

		buf = append(buf, r)
	}
	return string(buf), nil
}

func readString(br *bufio.Reader) (string, error) {
	buf := []rune{}
	pr, _, err := br.ReadRune()
	if err != nil {
		return "", err
	}
	if pr != '"' {
		return "", errors.New(`string should start with "`)
	}
	buf = append(buf, pr)

	for {
		r, _, err := br.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		buf = append(buf, r)

		if r == '"' && pr != '\\' {
			break
		}

		pr = r
	}
	return string(buf), nil
}

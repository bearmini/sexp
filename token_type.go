//go:generate stringer -type=TokenType -output=token_type_string.go
package sexp

type TokenType int

const (
	TokenTypeOpenParen TokenType = iota
	TokenTypeCloseParen
	TokenTypeSymbol
	TokenTypeNumber
	TokenTypeString
)

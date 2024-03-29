// Code generated by "stringer -type=TokenType -output=token_type_string.go"; DO NOT EDIT.

package sexp

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokenTypeOpenParen-0]
	_ = x[TokenTypeCloseParen-1]
	_ = x[TokenTypeSymbol-2]
	_ = x[TokenTypeNumber-3]
	_ = x[TokenTypeString-4]
}

const _TokenType_name = "TokenTypeOpenParenTokenTypeCloseParenTokenTypeSymbolTokenTypeNumberTokenTypeString"

var _TokenType_index = [...]uint8{0, 18, 37, 52, 67, 82}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}

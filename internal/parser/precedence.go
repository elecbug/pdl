package parser

import "github.com/elecbug/pdl/internal/token"

const (
	PRECEDENCE_LOWEST = iota
	PRECEDENCE_OR
	PRECEDENCE_AND
	PRECEDENCE_COMPARE
	PRECEDENCE_SUM
	PRECEDENCE_PRODUCT
)

// precedence returns the precedence level of the given token type. It is used to determine the order
// of operations when parsing expressions. Higher precedence values indicate higher precedence.
func precedence(t token.TokenType) int {
	switch t {
	case token.OR_OR_SIGN:
		return PRECEDENCE_OR
	case token.AND_AND_SIGN:
		return PRECEDENCE_AND
	case token.LANGLE_SIGN, token.RANGLE_SIGN,
		token.LESS_EQUAL_SIGN, token.GREATER_EQUAL_SIGN,
		token.EQUAL_EQUAL_SIGN, token.NOT_EQUAL_SIGN:
		return PRECEDENCE_COMPARE
	case token.PLUS_SIGN, token.MINUS_SIGN:
		return PRECEDENCE_SUM
	case token.STAR_SIGN, token.SLASH_SIGN:
		return PRECEDENCE_PRODUCT
	default:
		return PRECEDENCE_LOWEST
	}
}

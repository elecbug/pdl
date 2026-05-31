package parser

import "github.com/elecbug/pdl/internal/token"

const (
	PRECEDENCE_LOWEST = iota
	PRECEDENCE_SUM
	PRECEDENCE_PRODUCT
)

func precedence(t token.TokenType) int {
	switch t {
	case token.PLUS_SIGN, token.MINUS_SIGN:
		return PRECEDENCE_SUM
	case token.STAR_SIGN, token.SLASH_SIGN:
		return PRECEDENCE_PRODUCT
	default:
		return PRECEDENCE_LOWEST
	}
}

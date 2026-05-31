package parser

import "github.com/elecbug/pdl/internal/token"

const (
	precLowest = iota
	precSum
	precProduct
)

func precedence(t token.TokenType) int {
	switch t {
	case token.TokenPlus, token.TokenMinus:
		return precSum
	case token.TokenStar, token.TokenSlash:
		return precProduct
	default:
		return precLowest
	}
}

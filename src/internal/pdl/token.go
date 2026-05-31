package pdl

type TokenType string

const (
	TokenEOF     TokenType = "EOF"
	TokenIllegal TokenType = "ILLEGAL"

	TokenIdent  TokenType = "IDENT"
	TokenNumber TokenType = "NUMBER"

	TokenPacket TokenType = "packet"
	TokenSet    TokenType = "set"
	TokenMode   TokenType = "mode"
	TokenDef    TokenType = "def"
	TokenOut    TokenType = "out"
	TokenJSON   TokenType = "json"

	TokenFrom   TokenType = "from"
	TokenTo     TokenType = "to"
	TokenLength TokenType = "length"
	TokenEnd    TokenType = "end"

	TokenLBrace TokenType = "{"
	TokenRBrace TokenType = "}"
	TokenLParen TokenType = "("
	TokenRParen TokenType = ")"
	TokenLAngle TokenType = "<"
	TokenRAngle TokenType = ">"

	TokenStar  TokenType = "*"
	TokenPlus  TokenType = "+"
	TokenMinus TokenType = "-"
	TokenSlash TokenType = "/"

	TokenColon TokenType = ":"
)

type Token struct {
	Type TokenType
	Lit  string
	Line int
	Col  int
}

var keywords = map[string]TokenType{
	"packet": TokenPacket,
	"set":    TokenSet,
	"mode":   TokenMode,
	"def":    TokenDef,
	"out":    TokenOut,
	"json":   TokenJSON,
	"from":   TokenFrom,
	"to":     TokenTo,
	"length": TokenLength,
	"end":    TokenEnd,
}

func LookupIdent(s string) TokenType {
	if tok, ok := keywords[s]; ok {
		return tok
	}
	return TokenIdent
}

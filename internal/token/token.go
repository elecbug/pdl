package token

type TokenType string

const (
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"
	IDENT   TokenType = "IDENT"
	NUMBER  TokenType = "NUMBER"
	STRING  TokenType = "STRING"

	PACKET_KEYWORD TokenType = "packet"
	SET_KEYWORD    TokenType = "set"
	MODE_KEYWORD   TokenType = "mode"
	DEF_KEYWORD    TokenType = "def"
	OUT_KEYWORD    TokenType = "out"
	JSON_KEYWORD   TokenType = "json"
	VAR_KEYWORD    TokenType = "var"
	FROM_KEYWORD   TokenType = "from"
	TO_KEYWORD     TokenType = "to"
	LENGTH_KEYWORD TokenType = "length"
	END_KEYWORD    TokenType = "end"

	EQUAL_SIGN  TokenType = "="
	LBRACE_SIGN TokenType = "{"
	RBRACE_SIGN TokenType = "}"
	LPAREN_SIGN TokenType = "("
	RPAREN_SIGN TokenType = ")"
	LANGLE_SIGN TokenType = "<"
	RANGLE_SIGN TokenType = ">"
	STAR_SIGN   TokenType = "*"
	PLUS_SIGN   TokenType = "+"
	MINUS_SIGN  TokenType = "-"
	SLASH_SIGN  TokenType = "/"
	COLON_SIGN  TokenType = ":"
	QUOTE_SIGN  TokenType = "\""
)

type Token struct {
	Type TokenType
	Lit  string
	Line int
	Col  int
}

var keywords = map[string]TokenType{
	"packet": PACKET_KEYWORD,
	"set":    SET_KEYWORD,
	"mode":   MODE_KEYWORD,
	"var":    VAR_KEYWORD,
	"def":    DEF_KEYWORD,
	"out":    OUT_KEYWORD,
	"json":   JSON_KEYWORD,
	"from":   FROM_KEYWORD,
	"to":     TO_KEYWORD,
	"length": LENGTH_KEYWORD,
	"end":    END_KEYWORD,
}

func LookupIdent(s string) TokenType {
	if tok, ok := keywords[s]; ok {
		return tok
	}
	return IDENT
}

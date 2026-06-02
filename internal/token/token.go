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

// Token represents a lexical token with its type, literal value, and position (line and column) in the source code. It is used by the lexer to produce a stream of tokens that the parser can consume to build the abstract syntax tree (AST) of the PDL document.
type Token struct {
	// Type is the type of the token, such as IDENT, NUMBER, STRING, or a specific keyword or symbol.
	Type TokenType
	// Lit is the literal string value of the token as it appears in the source code. For example, for an IDENT token, Lit would be the identifier name; for a NUMBER token, it would be the numeric literal as a string; and for a STRING token, it would be the contents of the string literal.
	Lit string
	// Line is the line number in the source code where the token was found. It is used for error reporting and debugging purposes.
	Line int
	// Col is the column number in the source code where the token starts. It is used for error reporting and debugging purposes.
	Col int
}

// keywords is a map of reserved keywords in the PDL language to their corresponding token types. It is used by the lexer to identify when an identifier matches a reserved keyword and should be tokenized as that keyword instead of a generic IDENT token.
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

// LookupIdent checks if the given string is a reserved keyword in the PDL language. If it is, it returns the corresponding token type for that keyword. If it is not a reserved keyword, it returns the IDENT token type, indicating that the string should be treated as a generic identifier. This function is used by the lexer when tokenizing identifiers to determine if they are keywords or not.
func LookupIdent(s string) TokenType {
	if tok, ok := keywords[s]; ok {
		return tok
	}
	return IDENT
}

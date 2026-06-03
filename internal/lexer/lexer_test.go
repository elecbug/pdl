package lexer_test

import (
	"testing"

	"github.com/elecbug/pdl/internal/lexer"
	"github.com/elecbug/pdl/internal/token"
)

func TestLexer(t *testing.T) {
	input := `packet TCP
# full line comment should be ignored
set mode BIG_ENDIAN # inline comment too
var {
	alpha_1 = 123
	bin = 0b1010
	hex = 0xAF
	text = "hello"
	path = header.length_words
	idx = flags[7]
	expr = (1+2) * 3 / 4 - 5
	map<4> : end
}
@`

	tests := []struct {
		wantType token.TokenType
		wantLit  string
		line     int
		col      int
	}{
		{token.PACKET_KEYWORD, "packet", 1, 1},
		{token.IDENT, "TCP", 1, 8},
		{token.SET_KEYWORD, "set", 3, 1},
		{token.MODE_KEYWORD, "mode", 3, 5},
		{token.IDENT, "BIG_ENDIAN", 3, 10},
		{token.VAR_KEYWORD, "var", 4, 1},
		{token.LBRACE_SIGN, "{", 4, 5},

		{token.IDENT, "alpha_1", 5, 2},
		{token.EQUAL_SIGN, "=", 5, 10},
		{token.NUMBER, "123", 5, 12},

		{token.IDENT, "bin", 6, 2},
		{token.EQUAL_SIGN, "=", 6, 6},
		{token.NUMBER, "0b1010", 6, 8},

		{token.IDENT, "hex", 7, 2},
		{token.EQUAL_SIGN, "=", 7, 6},
		{token.NUMBER, "0xAF", 7, 8},

		{token.IDENT, "text", 8, 2},
		{token.EQUAL_SIGN, "=", 8, 7},
		{token.STRING, "hello", 8, 9},

		{token.IDENT, "path", 9, 2},
		{token.EQUAL_SIGN, "=", 9, 7},
		{token.IDENT, "header.length_words", 9, 9},

		{token.IDENT, "idx", 10, 2},
		{token.EQUAL_SIGN, "=", 10, 6},
		{token.IDENT, "flags[7]", 10, 8},

		{token.IDENT, "expr", 11, 2},
		{token.EQUAL_SIGN, "=", 11, 7},
		{token.LPAREN_SIGN, "(", 11, 9},
		{token.NUMBER, "1", 11, 10},
		{token.PLUS_SIGN, "+", 11, 11},
		{token.NUMBER, "2", 11, 12},
		{token.RPAREN_SIGN, ")", 11, 13},
		{token.STAR_SIGN, "*", 11, 15},
		{token.NUMBER, "3", 11, 17},
		{token.SLASH_SIGN, "/", 11, 19},
		{token.NUMBER, "4", 11, 21},
		{token.MINUS_SIGN, "-", 11, 23},
		{token.NUMBER, "5", 11, 25},

		{token.IDENT, "map", 12, 2},
		{token.LANGLE_SIGN, "<", 12, 5},
		{token.NUMBER, "4", 12, 6},
		{token.RANGLE_SIGN, ">", 12, 7},
		{token.COLON_SIGN, ":", 12, 9},
		{token.END_KEYWORD, "end", 12, 11},

		{token.RBRACE_SIGN, "}", 13, 1},
		{token.ILLEGAL, "@", 14, 1},
		{token.EOF, "", 14, 2},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.wantType || tok.Lit != tt.wantLit || tok.Line != tt.line || tok.Col != tt.col {
			t.Fatalf(
				"test[%d] token mismatch: got {Type:%s Lit:%q Line:%d Col:%d}, want {Type:%s Lit:%q Line:%d Col:%d}",
				i,
				tok.Type,
				tok.Lit,
				tok.Line,
				tok.Col,
				tt.wantType,
				tt.wantLit,
				tt.line,
				tt.col,
			)
		}
	}

	// Ensure lexer stays at EOF after input exhaustion.
	eofAgain := l.NextToken()
	if eofAgain.Type != token.EOF {
		t.Fatalf("expected EOF on repeated read, got %s (%q)", eofAgain.Type, eofAgain.Lit)
	}
}

func TestLexer_UnterminatedString(t *testing.T) {
	l := lexer.New("text = \"unterminated")

	tok1 := l.NextToken()
	if tok1.Type != token.IDENT || tok1.Lit != "text" {
		t.Fatalf("first token = {%s %q}, want IDENT text", tok1.Type, tok1.Lit)
	}

	tok2 := l.NextToken()
	if tok2.Type != token.EQUAL_SIGN || tok2.Lit != "=" {
		t.Fatalf("second token = {%s %q}, want '='", tok2.Type, tok2.Lit)
	}

	tok3 := l.NextToken()
	if tok3.Type != token.ILLEGAL {
		t.Fatalf("third token type = %s, want ILLEGAL", tok3.Type)
	}

	if tok3.Lit != "unterminated string" {
		t.Fatalf("third token lit = %q, want %q", tok3.Lit, "unterminated string")
	}

}

package lexer

import (
	"fmt"

	"unicode"

	"github.com/elecbug/pdl/internal/token"
)

type Lexer struct {
	input []rune
	pos   int

	line int
	col  int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		line:  1,
		col:   1,
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaceAndComments()

	startLine := l.line
	startCol := l.col

	ch := l.peek()
	if ch == 0 {
		return token.Token{Type: token.EOF, Lit: "", Line: startLine, Col: startCol}
	}

	switch ch {
	case '{':
		l.advance()
		return token.Token{Type: token.LBRACE_SIGN, Lit: "{", Line: startLine, Col: startCol}
	case '}':
		l.advance()
		return token.Token{Type: token.RBRACE_SIGN, Lit: "}", Line: startLine, Col: startCol}
	case '(':
		l.advance()
		return token.Token{Type: token.LPAREN_SIGN, Lit: "(", Line: startLine, Col: startCol}
	case ')':
		l.advance()
		return token.Token{Type: token.RPAREN_SIGN, Lit: ")", Line: startLine, Col: startCol}
	case '<':
		l.advance()
		return token.Token{Type: token.LANGLE_SIGN, Lit: "<", Line: startLine, Col: startCol}
	case '>':
		l.advance()
		return token.Token{Type: token.RANGLE_SIGN, Lit: ">", Line: startLine, Col: startCol}
	case '*':
		l.advance()
		return token.Token{Type: token.STAR_SIGN, Lit: "*", Line: startLine, Col: startCol}
	case '+':
		l.advance()
		return token.Token{Type: token.PLUS_SIGN, Lit: "+", Line: startLine, Col: startCol}
	case '-':
		l.advance()
		return token.Token{Type: token.MINUS_SIGN, Lit: "-", Line: startLine, Col: startCol}
	case '/':
		l.advance()
		return token.Token{Type: token.SLASH_SIGN, Lit: "/", Line: startLine, Col: startCol}
	case ':':
		l.advance()
		return token.Token{Type: token.COLON_SIGN, Lit: ":", Line: startLine, Col: startCol}
	case '=':
		l.advance()
		return token.Token{Type: token.EQUAL_SIGN, Lit: "=", Line: startLine, Col: startCol}
	case '"':
		lit, err := l.readString()
		if err != nil {
			return token.Token{Type: token.ILLEGAL, Lit: err.Error(), Line: startLine, Col: startCol}
		}
		return token.Token{Type: token.STRING, Lit: lit, Line: startLine, Col: startCol}
	}

	if isIdentStart(ch) {
		lit := l.readIdent()
		return token.Token{
			Type: token.LookupIdent(lit),
			Lit:  lit,
			Line: startLine,
			Col:  startCol,
		}
	}

	if unicode.IsDigit(ch) {
		lit := l.readNumber()
		return token.Token{
			Type: token.NUMBER,
			Lit:  lit,
			Line: startLine,
			Col:  startCol,
		}
	}

	l.advance()
	return token.Token{
		Type: token.ILLEGAL,
		Lit:  fmt.Sprintf("%c", ch),
		Line: startLine,
		Col:  startCol,
	}
}

func (l *Lexer) readString() (string, error) {
	l.advance() // skip opening "

	start := l.pos

	for {
		ch := l.peek()
		if ch == 0 || ch == '\n' {
			return "", fmt.Errorf("unterminated string")
		}

		if ch == '"' {
			lit := string(l.input[start:l.pos])
			l.advance() // skip closing "
			return lit, nil
		}

		l.advance()
	}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		ch := l.peek()

		for unicode.IsSpace(ch) {
			l.advance()
			ch = l.peek()
		}

		if ch == '#' {
			for ch != '\n' && ch != 0 {
				l.advance()
				ch = l.peek()
			}
			continue
		}

		return
	}
}

func (l *Lexer) readIdent() string {
	start := l.pos

	for isIdentPart(l.peek()) {
		l.advance()
	}

	return string(l.input[start:l.pos])
}

func (l *Lexer) readNumber() string {
	start := l.pos

	// Support:
	// 123
	// 0b1010
	// 0xFF
	if l.peek() == '0' {
		next := l.peekN(1)
		if next == 'b' || next == 'B' || next == 'x' || next == 'X' {
			l.advance()
			l.advance()

			for isNumberPart(l.peek()) {
				l.advance()
			}

			return string(l.input[start:l.pos])
		}
	}

	for unicode.IsDigit(l.peek()) {
		l.advance()
	}

	return string(l.input[start:l.pos])
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekN(n int) rune {
	idx := l.pos + n
	if idx >= len(l.input) {
		return 0
	}
	return l.input[idx]
}

func (l *Lexer) advance() rune {
	ch := l.peek()
	if ch == 0 {
		return 0
	}

	l.pos++

	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}

	return ch
}

func isIdentStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isIdentPart(ch rune) bool {
	return unicode.IsLetter(ch) ||
		unicode.IsDigit(ch) ||
		ch == '_' ||
		ch == '.' ||
		ch == '[' ||
		ch == ']'
}

func isNumberPart(ch rune) bool {
	return unicode.IsDigit(ch) ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

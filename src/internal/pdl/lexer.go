package pdl

import (
	"fmt"
	"unicode"
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

func (l *Lexer) NextToken() Token {
	l.skipWhitespaceAndComments()

	startLine := l.line
	startCol := l.col

	ch := l.peek()
	if ch == 0 {
		return Token{Type: TokenEOF, Lit: "", Line: startLine, Col: startCol}
	}

	switch ch {
	case '{':
		l.advance()
		return Token{Type: TokenLBrace, Lit: "{", Line: startLine, Col: startCol}
	case '}':
		l.advance()
		return Token{Type: TokenRBrace, Lit: "}", Line: startLine, Col: startCol}
	case '(':
		l.advance()
		return Token{Type: TokenLParen, Lit: "(", Line: startLine, Col: startCol}
	case ')':
		l.advance()
		return Token{Type: TokenRParen, Lit: ")", Line: startLine, Col: startCol}
	case '<':
		l.advance()
		return Token{Type: TokenLAngle, Lit: "<", Line: startLine, Col: startCol}
	case '>':
		l.advance()
		return Token{Type: TokenRAngle, Lit: ">", Line: startLine, Col: startCol}
	case '*':
		l.advance()
		return Token{Type: TokenStar, Lit: "*", Line: startLine, Col: startCol}
	case '+':
		l.advance()
		return Token{Type: TokenPlus, Lit: "+", Line: startLine, Col: startCol}
	case '-':
		l.advance()
		return Token{Type: TokenMinus, Lit: "-", Line: startLine, Col: startCol}
	case '/':
		l.advance()
		return Token{Type: TokenSlash, Lit: "/", Line: startLine, Col: startCol}
	case ':':
		l.advance()
		return Token{Type: TokenColon, Lit: ":", Line: startLine, Col: startCol}
	}

	if isIdentStart(ch) {
		lit := l.readIdent()
		return Token{
			Type: LookupIdent(lit),
			Lit:  lit,
			Line: startLine,
			Col:  startCol,
		}
	}

	if unicode.IsDigit(ch) {
		lit := l.readNumber()
		return Token{
			Type: TokenNumber,
			Lit:  lit,
			Line: startLine,
			Col:  startCol,
		}
	}

	l.advance()
	return Token{
		Type: TokenIllegal,
		Lit:  fmt.Sprintf("%c", ch),
		Line: startLine,
		Col:  startCol,
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
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}

func isNumberPart(ch rune) bool {
	return unicode.IsDigit(ch) ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

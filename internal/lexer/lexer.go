package lexer

import (
	"fmt"

	"unicode"

	"github.com/elecbug/pdl/internal/token"
)

// Lexer is responsible for tokenizing the input PDL source code. It reads the input string and produces
// a stream of tokens that can be consumed by the parser.
type Lexer struct {
	// The input source code as a slice of runes, allowing for proper handling of Unicode characters.
	input []rune
	// The current position in the input (index of the next character to read).
	pos int

	// The current line number, starting at 1.
	line int
	// The current column number, starting at 1.
	col int
}

// New creates a Lexer and initializes line and column counters to 1.
func New(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		line:  1,
		col:   1,
	}
}

// NextToken reads the next token from the input and returns it. It handles whitespace, comments, identifiers,
// numbers, strings, and punctuation.
// If it encounters an unrecognized character, it returns an ILLEGAL token.
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
		if l.peekN(1) == '=' {
			l.advance()
			l.advance()
			return token.Token{Type: token.LESS_EQUAL_SIGN, Lit: "<=", Line: startLine, Col: startCol}
		}
		l.advance()
		return token.Token{Type: token.LANGLE_SIGN, Lit: "<", Line: startLine, Col: startCol}
	case '>':
		if l.peekN(1) == '=' {
			l.advance()
			l.advance()
			return token.Token{Type: token.GREATER_EQUAL_SIGN, Lit: ">=", Line: startLine, Col: startCol}
		}
		l.advance()
		return token.Token{Type: token.RANGLE_SIGN, Lit: ">", Line: startLine, Col: startCol}
	case '=':
		if l.peekN(1) == '=' {
			l.advance()
			l.advance()
			return token.Token{Type: token.EQUAL_EQUAL_SIGN, Lit: "==", Line: startLine, Col: startCol}
		}
		l.advance()
		return token.Token{Type: token.EQUAL_SIGN, Lit: "=", Line: startLine, Col: startCol}
	case '!':
		if l.peekN(1) == '=' {
			l.advance()
			l.advance()
			return token.Token{Type: token.NOT_EQUAL_SIGN, Lit: "!=", Line: startLine, Col: startCol}
		}
	case '&':
		if l.peekN(1) == '&' {
			l.advance()
			l.advance()
			return token.Token{Type: token.AND_AND_SIGN, Lit: "&&", Line: startLine, Col: startCol}
		}
	case '|':
		if l.peekN(1) == '|' {
			l.advance()
			l.advance()
			return token.Token{Type: token.OR_OR_SIGN, Lit: "||", Line: startLine, Col: startCol}
		}
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

// readString reads a string literal from the input, starting with the opening quote. It continues
// until it finds a closing quote or reaches the end of the line/input. If it encounters an
// unterminated string, it returns an error.
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

// skipWhitespaceAndComments advances the lexer's position past any whitespace characters and comments.
// It keeps advancing until it reaches a non-whitespace, non-comment
// character or the end of input.
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

// readIdent reads an identifier from the input, starting with a valid identifier character. It continues
// until it encounters a character that is not valid in an identifier.
// It returns the identifier text.
func (l *Lexer) readIdent() string {
	start := l.pos

	for isIdentPart(l.peek()) {
		l.advance()
	}

	return string(l.input[start:l.pos])
}

// readNumber reads a number literal.
// It supports decimal, binary (0b prefix), and hexadecimal (0x prefix)
// forms, and returns the raw literal text.
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

// peek returns the next character without advancing.
// It returns 0 at end of input.
func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

// peekN returns the character at the given offset without advancing.
// If it reaches the end of the input, it returns 0.
func (l *Lexer) peekN(n int) rune {
	idx := l.pos + n
	if idx >= len(l.input) {
		return 0
	}
	return l.input[idx]
}

// advance moves one character forward and updates line/column counters.
// It returns the consumed character, or 0 at end of input.
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

// isIdentStart reports whether ch can start an identifier.
func isIdentStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// isIdentPart reports whether ch is valid in an identifier.
func isIdentPart(ch rune) bool {
	return unicode.IsLetter(ch) ||
		unicode.IsDigit(ch) ||
		ch == '_' ||
		ch == '.' ||
		ch == '[' ||
		ch == ']'
}

// isNumberPart reports whether ch is valid in a numeric literal.
func isNumberPart(ch rune) bool {
	return unicode.IsDigit(ch) ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

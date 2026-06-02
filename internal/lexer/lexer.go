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

// New creates a new Lexer instance with the given input string. It initializes the line and column numbers to 1.
func New(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		line:  1,
		col:   1,
	}
}

// NextToken reads the next token from the input and returns it. It handles whitespace, comments, identifiers,
// numbers, strings, and various punctuation. If it encounters an unrecognized character, it returns an ILLEGAL token.
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
// It continues to skip until it encounters a non-whitespace, non-comment character or reaches the end of the input.
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
// until it encounters a character that is not valid in an identifier. It returns the string representation of the identifier.
func (l *Lexer) readIdent() string {
	start := l.pos

	for isIdentPart(l.peek()) {
		l.advance()
	}

	return string(l.input[start:l.pos])
}

// readNumber reads a number from the input, supporting decimal, binary (0b prefix), and hexadecimal (0x prefix) formats. It continues
// until it encounters a character that is not valid in a number. It returns the string representation of the number.
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

// peek returns the next character in the input without advancing the position. If it reaches the end of the input, it returns 0.
func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

// peekN returns the character at the specified offset from the current position without advancing the position.
// If it reaches the end of the input, it returns 0.
func (l *Lexer) peekN(n int) rune {
	idx := l.pos + n
	if idx >= len(l.input) {
		return 0
	}
	return l.input[idx]
}

// advance moves the lexer's position forward by one character and updates the line and column numbers accordingly.
// It returns the character that was advanced over. If it reaches the end of the input, it returns 0.
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

// isIdentStart checks if the given character is a valid starting character for an identifier, which can be a letter or an underscore.
func isIdentStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// isIdentPart checks if the given character is a valid part of an identifier, which can be a letter, digit, underscore, dot, or square brackets.
func isIdentPart(ch rune) bool {
	return unicode.IsLetter(ch) ||
		unicode.IsDigit(ch) ||
		ch == '_' ||
		ch == '.' ||
		ch == '[' ||
		ch == ']'
}

// isNumberPart checks if the given character is a valid part of a number, which can be a digit or a hexadecimal character (a-f, A-F).
func isNumberPart(ch rune) bool {
	return unicode.IsDigit(ch) ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

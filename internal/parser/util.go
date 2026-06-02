package parser

import (
	"fmt"

	"github.com/elecbug/pdl/internal/token"
)

// next advances the parser to the next token. It updates the current token (cur) to the peek token and
// then fetches the next token from the lexer to update the peek token. This allows the parser to look ahead one token while parsing.
func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

// expect checks if the current token matches the expected token type. If it does, it advances to the next token.
func (p *Parser) expect(t token.TokenType) error {
	if p.cur.Type != t {
		return p.errf("expected %s, got %s %q", t, p.cur.Type, p.cur.Lit)
	}
	p.next()
	return nil
}

// errf formats an error message with the current token's line and column information. It is used to
// provide detailed error messages when the parser encounters unexpected tokens or syntax errors.
func (p *Parser) errf(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("parse error at %d:%d: %s", p.cur.Line, p.cur.Col, msg)
}

package parser

import (
	"fmt"

	"github.com/elecbug/pdl/internal/token"
)

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) expect(t token.TokenType) error {
	if p.cur.Type != t {
		return p.errf("expected %s, got %s %q", t, p.cur.Type, p.cur.Lit)
	}
	p.next()
	return nil
}

func (p *Parser) errf(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("parse error at %d:%d: %s", p.cur.Line, p.cur.Col, msg)
}

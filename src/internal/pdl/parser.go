package pdl

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	l *Lexer

	cur  Token
	peek Token
}

func NewParser(input string) *Parser {
	l := NewLexer(input)

	p := &Parser{l: l}
	p.next()
	p.next()

	return p
}

func ParseString(input string) (*Document, error) {
	return NewParser(input).Parse()
}

func (p *Parser) Parse() (*Document, error) {
	doc := &Document{}

	for p.cur.Type != TokenEOF {
		switch p.cur.Type {
		case TokenPacket:
			if err := p.parsePacket(doc); err != nil {
				return nil, err
			}

		case TokenSet:
			if err := p.parseSetMode(doc); err != nil {
				return nil, err
			}

		case TokenDef:
			defs, err := p.parseDefBlock()
			if err != nil {
				return nil, err
			}
			doc.Defs = append(doc.Defs, defs...)

		case TokenOut:
			outs, err := p.parseOutBlock()
			if err != nil {
				return nil, err
			}
			doc.Outputs = append(doc.Outputs, outs...)

		default:
			return nil, p.errf("unexpected token %s %q", p.cur.Type, p.cur.Lit)
		}
	}

	return doc, nil
}

func (p *Parser) parsePacket(doc *Document) error {
	if err := p.expect(TokenPacket); err != nil {
		return err
	}

	if p.cur.Type != TokenIdent {
		return p.errf("expected packet name, got %s %q", p.cur.Type, p.cur.Lit)
	}

	doc.PacketName = p.cur.Lit
	p.next()
	return nil
}

func (p *Parser) parseSetMode(doc *Document) error {
	if err := p.expect(TokenSet); err != nil {
		return err
	}

	if err := p.expect(TokenMode); err != nil {
		return err
	}

	if p.cur.Type != TokenIdent {
		return p.errf("expected mode name, got %s %q", p.cur.Type, p.cur.Lit)
	}

	doc.Mode = p.cur.Lit
	p.next()
	return nil
}

func (p *Parser) parseDefBlock() ([]Def, error) {
	if err := p.expect(TokenDef); err != nil {
		return nil, err
	}

	if err := p.expect(TokenLBrace); err != nil {
		return nil, err
	}

	var defs []Def

	for p.cur.Type != TokenRBrace && p.cur.Type != TokenEOF {
		def, err := p.parseDefLine()
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}

	if err := p.expect(TokenRBrace); err != nil {
		return nil, err
	}

	return defs, nil
}

func (p *Parser) parseDefLine() (Def, error) {
	if p.cur.Type != TokenIdent {
		return Def{}, p.errf("expected field name in def block, got %s %q", p.cur.Type, p.cur.Lit)
	}

	def := Def{Name: p.cur.Lit}
	p.next()

	if err := p.expect(TokenFrom); err != nil {
		return Def{}, err
	}

	from, err := p.parseExpr(0)
	if err != nil {
		return Def{}, err
	}
	def.From = from

	switch p.cur.Type {
	case TokenLength:
		p.next()

		length, err := p.parseExpr(0)
		if err != nil {
			return Def{}, err
		}

		def.Length = length
		def.UseLength = true

	case TokenTo:
		p.next()

		to, err := p.parseExpr(0)
		if err != nil {
			return Def{}, err
		}

		def.To = to
		def.UseTo = true

	default:
		return Def{}, p.errf("expected length or to in def line, got %s %q", p.cur.Type, p.cur.Lit)
	}

	return def, nil
}

func (p *Parser) parseOutBlock() ([]Output, error) {
	if err := p.expect(TokenOut); err != nil {
		return nil, err
	}

	if err := p.expect(TokenJSON); err != nil {
		return nil, err
	}

	if err := p.expect(TokenLBrace); err != nil {
		return nil, err
	}

	var outs []Output

	for p.cur.Type != TokenRBrace && p.cur.Type != TokenEOF {
		out, err := p.parseOutLine()
		if err != nil {
			return nil, err
		}
		outs = append(outs, out)
	}

	if err := p.expect(TokenRBrace); err != nil {
		return nil, err
	}

	return outs, nil
}

func (p *Parser) parseOutLine() (Output, error) {
	if p.cur.Type != TokenIdent {
		return Output{}, p.errf("expected field name in out block, got %s %q", p.cur.Type, p.cur.Lit)
	}

	out := Output{
		Field: p.cur.Lit,
	}

	p.next()

	if p.cur.Type == TokenLAngle {
		p.next()

		if p.cur.Type != TokenNumber {
			return Output{}, p.errf("expected bit index number")
		}

		v, err := parseNumber(p.cur.Lit)
		if err != nil {
			return Output{}, p.errf("invalid bit index %q", p.cur.Lit)
		}

		idx := int(v)
		out.BitIndex = &idx
		p.next()

		if err := p.expect(TokenRAngle); err != nil {
			return Output{}, err
		}
	}

	if p.cur.Type != TokenIdent {
		return Output{}, p.errf("expected output path, got %s %q", p.cur.Type, p.cur.Lit)
	}

	out.Path = p.cur.Lit
	p.next()

	// Normal output:
	// src_port source_port DEC
	if p.cur.Type == TokenIdent {
		out.Format = p.cur.Lit
		p.next()
		return out, nil
	}

	// Bracket output:
	// flags<6> syn { 0 : false 1 : true }
	if p.cur.Type == TokenLBrace {
		m, err := p.parseMapBlock()
		if err != nil {
			return Output{}, err
		}
		out.Map = m
		return out, nil
	}

	return Output{}, p.errf("expected format or map block in out line, got %s %q", p.cur.Type, p.cur.Lit)
}

func (p *Parser) parseMapBlock() (map[string]string, error) {
	if err := p.expect(TokenLBrace); err != nil {
		return nil, err
	}

	m := make(map[string]string)

	for p.cur.Type != TokenRBrace && p.cur.Type != TokenEOF {
		key := p.cur.Lit

		if p.cur.Type != TokenNumber && p.cur.Type != TokenIdent {
			return nil, p.errf("expected map key, got %s %q", p.cur.Type, p.cur.Lit)
		}
		p.next()

		if err := p.expect(TokenColon); err != nil {
			return nil, err
		}

		if p.cur.Type != TokenIdent && p.cur.Type != TokenNumber {
			return nil, p.errf("expected map value, got %s %q", p.cur.Type, p.cur.Lit)
		}

		val := p.cur.Lit
		m[key] = val
		p.next()
	}

	if err := p.expect(TokenRBrace); err != nil {
		return nil, err
	}

	return m, nil
}

const (
	precLowest = iota
	precSum
	precProduct
)

func precedence(t TokenType) int {
	switch t {
	case TokenPlus, TokenMinus:
		return precSum
	case TokenStar, TokenSlash:
		return precProduct
	default:
		return precLowest
	}
}

func (p *Parser) parseExpr(minPrec int) (Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		prec := precedence(p.cur.Type)
		if prec <= minPrec {
			break
		}

		op := p.cur.Lit
		p.next()

		right, err := p.parseExpr(prec)
		if err != nil {
			return nil, err
		}

		left = BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) parsePrimary() (Expr, error) {
	switch p.cur.Type {
	case TokenNumber:
		raw := p.cur.Lit
		v, err := parseNumber(raw)
		if err != nil {
			return nil, p.errf("invalid number %q", raw)
		}
		p.next()
		return NumberExpr{Raw: raw, Value: v}, nil

	case TokenIdent:
		name := p.cur.Lit
		p.next()
		return IdentExpr{Name: name}, nil

	case TokenEnd:
		p.next()
		return EndExpr{}, nil

	case TokenStar:
		p.next()

		if p.cur.Type != TokenIdent {
			return nil, p.errf("expected field name after *")
		}

		name := p.cur.Lit
		p.next()

		return FieldValueExpr{Name: name}, nil

	case TokenLParen:
		p.next()

		expr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}

		if err := p.expect(TokenRParen); err != nil {
			return nil, err
		}

		return expr, nil

	default:
		return nil, p.errf("expected expression, got %s %q", p.cur.Type, p.cur.Lit)
	}
}

func parseNumber(raw string) (int64, error) {
	switch {
	case strings.HasPrefix(raw, "0b") || strings.HasPrefix(raw, "0B"):
		return strconv.ParseInt(raw[2:], 2, 64)
	case strings.HasPrefix(raw, "0x") || strings.HasPrefix(raw, "0X"):
		return strconv.ParseInt(raw[2:], 16, 64)
	default:
		return strconv.ParseInt(raw, 10, 64)
	}
}

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) expect(t TokenType) error {
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

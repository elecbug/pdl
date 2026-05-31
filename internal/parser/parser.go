package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/elecbug/pdl/internal/ast"
	"github.com/elecbug/pdl/internal/lexer"
	"github.com/elecbug/pdl/internal/token"
)

type Parser struct {
	l *lexer.Lexer

	cur  token.Token
	peek token.Token
}

func NewParser(input string) *Parser {
	l := lexer.NewLexer(input)

	p := &Parser{l: l}
	p.next()
	p.next()

	return p
}

func ParseString(input string) (*ast.Document, error) {
	return NewParser(input).Parse()
}

func (p *Parser) Parse() (*ast.Document, error) {
	doc := &ast.Document{}

	for p.cur.Type != token.TokenEOF {
		switch p.cur.Type {
		case token.TokenPacket:
			if err := p.parsePacket(doc); err != nil {
				return nil, err
			}

		case token.TokenSet:
			if err := p.parseSetMode(doc); err != nil {
				return nil, err
			}

		case token.TokenDef:
			defs, err := p.parseDefBlock()
			if err != nil {
				return nil, err
			}
			doc.Defs = append(doc.Defs, defs...)

		case token.TokenOut:
			outs, err := p.parseOutBlock()
			if err != nil {
				return nil, err
			}
			doc.Outputs = append(doc.Outputs, outs...)
		case token.TokenVar:
			vars, err := p.parseVarBlock()
			if err != nil {
				return nil, err
			}
			doc.Vars = append(doc.Vars, vars...)
		default:
			return nil, p.errf("unexpected token %s %q", p.cur.Type, p.cur.Lit)
		}
	}

	return doc, nil
}

func (p *Parser) parsePacket(doc *ast.Document) error {
	if err := p.expect(token.TokenPacket); err != nil {
		return err
	}

	if p.cur.Type != token.TokenIdent {
		return p.errf("expected packet name, got %s %q", p.cur.Type, p.cur.Lit)
	}

	doc.PacketName = p.cur.Lit
	p.next()
	return nil
}

func (p *Parser) parseSetMode(doc *ast.Document) error {
	if err := p.expect(token.TokenSet); err != nil {
		return err
	}

	if err := p.expect(token.TokenMode); err != nil {
		return err
	}

	if p.cur.Type != token.TokenIdent {
		return p.errf("expected byte order")
	}

	switch p.cur.Lit {
	case "BIG_ENDIAN":
		doc.ByteOrder = ast.BIG_ENDIAN

	case "LITTLE_ENDIAN":
		doc.ByteOrder = ast.LITTLE_ENDIAN

	default:
		return p.errf("unknown byte order %q", p.cur.Lit)
	}

	p.next()

	if p.cur.Type != token.TokenIdent {
		return p.errf("expected bit order")
	}

	switch p.cur.Lit {
	case "MSB_FIRST":
		doc.BitOrder = ast.MSB_FIRST

	case "LSB_FIRST":
		doc.BitOrder = ast.LSB_FIRST

	default:
		return p.errf("unknown bit order %q", p.cur.Lit)
	}

	p.next()

	return nil
}

func (p *Parser) parseDefBlock() ([]ast.Def, error) {
	if err := p.expect(token.TokenDef); err != nil {
		return nil, err
	}

	if err := p.expect(token.TokenLBrace); err != nil {
		return nil, err
	}

	var defs []ast.Def

	for p.cur.Type != token.TokenRBrace && p.cur.Type != token.TokenEOF {
		def, err := p.parseDefLine()
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}

	if err := p.expect(token.TokenRBrace); err != nil {
		return nil, err
	}

	return defs, nil
}

func (p *Parser) parseDefLine() (ast.Def, error) {
	if p.cur.Type != token.TokenIdent {
		return ast.Def{}, p.errf("expected field name in def block, got %s %q", p.cur.Type, p.cur.Lit)
	}

	def := ast.Def{Name: p.cur.Lit}
	p.next()

	if err := p.expect(token.TokenFrom); err != nil {
		return ast.Def{}, err
	}

	from, err := p.parseExpr(0)
	if err != nil {
		return ast.Def{}, err
	}
	def.From = from

	switch p.cur.Type {
	case token.TokenLength:
		p.next()

		length, err := p.parseExpr(0)
		if err != nil {
			return ast.Def{}, err
		}

		def.Length = length
		def.UseLength = true

	case token.TokenTo:
		p.next()

		to, err := p.parseExpr(0)
		if err != nil {
			return ast.Def{}, err
		}

		def.To = to
		def.UseTo = true

	default:
		return ast.Def{}, p.errf("expected length or to in def line, got %s %q", p.cur.Type, p.cur.Lit)
	}

	return def, nil
}

func (p *Parser) parseOutBlock() ([]ast.Output, error) {
	if err := p.expect(token.TokenOut); err != nil {
		return nil, err
	}

	if err := p.expect(token.TokenJSON); err != nil {
		return nil, err
	}

	if err := p.expect(token.TokenLBrace); err != nil {
		return nil, err
	}

	var outs []ast.Output

	for p.cur.Type != token.TokenRBrace && p.cur.Type != token.TokenEOF {
		out, err := p.parseOutLine()
		if err != nil {
			return nil, err
		}
		outs = append(outs, out)
	}

	if err := p.expect(token.TokenRBrace); err != nil {
		return nil, err
	}

	return outs, nil
}

func (p *Parser) parseOutLine() (ast.Output, error) {
	if p.cur.Type != token.TokenIdent {
		return ast.Output{}, p.errf("expected field name in out block, got %s %q", p.cur.Type, p.cur.Lit)
	}

	out := ast.Output{
		Field: p.cur.Lit,
	}

	p.next()

	if p.cur.Type == token.TokenLAngle {
		p.next()

		if p.cur.Type != token.TokenNumber {
			return ast.Output{}, p.errf("expected bit index number")
		}

		v, err := parseNumber(p.cur.Lit)
		if err != nil {
			return ast.Output{}, p.errf("invalid bit index %q", p.cur.Lit)
		}

		idx := int(v)
		out.BitIndex = &idx
		p.next()

		if err := p.expect(token.TokenRAngle); err != nil {
			return ast.Output{}, err
		}
	}

	if p.cur.Type != token.TokenIdent {
		return ast.Output{}, p.errf("expected output path, got %s %q", p.cur.Type, p.cur.Lit)
	}

	out.Path = p.cur.Lit
	p.next()

	// Normal output:
	// src_port source_port DEC
	if p.cur.Type == token.TokenIdent {
		out.Format = p.cur.Lit
		p.next()
		return out, nil
	}

	// Bracket output:
	// flags<6> syn { 0 : false 1 : true }
	if p.cur.Type == token.TokenLBrace {
		m, err := p.parseMapBlock()
		if err != nil {
			return ast.Output{}, err
		}
		out.Map = m
		return out, nil
	}

	return ast.Output{}, p.errf("expected format or map block in out line, got %s %q", p.cur.Type, p.cur.Lit)
}

func (p *Parser) parseMapBlock() (map[string]string, error) {
	if err := p.expect(token.TokenLBrace); err != nil {
		return nil, err
	}

	m := make(map[string]string)

	for p.cur.Type != token.TokenRBrace && p.cur.Type != token.TokenEOF {
		key := p.cur.Lit

		if p.cur.Type != token.TokenNumber && p.cur.Type != token.TokenIdent {
			return nil, p.errf("expected map key, got %s %q", p.cur.Type, p.cur.Lit)
		}
		p.next()

		if err := p.expect(token.TokenColon); err != nil {
			return nil, err
		}

		if p.cur.Type != token.TokenIdent && p.cur.Type != token.TokenNumber {
			return nil, p.errf("expected map value, got %s %q", p.cur.Type, p.cur.Lit)
		}

		val := p.cur.Lit
		m[key] = val
		p.next()
	}

	if err := p.expect(token.TokenRBrace); err != nil {
		return nil, err
	}

	return m, nil
}

const (
	precLowest = iota
	precSum
	precProduct
)

func precedence(t token.TokenType) int {
	switch t {
	case token.TokenPlus, token.TokenMinus:
		return precSum
	case token.TokenStar, token.TokenSlash:
		return precProduct
	default:
		return precLowest
	}
}

func (p *Parser) parseExpr(minPrec int) (ast.Expr, error) {
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

		left = ast.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	switch p.cur.Type {
	case token.TokenNumber:
		raw := p.cur.Lit
		v, err := parseNumber(raw)
		if err != nil {
			return nil, p.errf("invalid number %q", raw)
		}
		p.next()
		return ast.NumberExpr{Raw: raw, Value: v}, nil

	case token.TokenIdent:
		name := p.cur.Lit
		p.next()
		return ast.IdentExpr{Name: name}, nil

	case token.TokenEnd:
		p.next()
		return ast.EndExpr{}, nil

	case token.TokenStar:
		p.next()

		if p.cur.Type != token.TokenIdent {
			return nil, p.errf("expected field name after *")
		}

		name := p.cur.Lit
		p.next()

		return ast.FieldValueExpr{Name: name}, nil

	case token.TokenLParen:
		p.next()

		expr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}

		if err := p.expect(token.TokenRParen); err != nil {
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

func (p *Parser) parseVarBlock() ([]ast.Var, error) {
	if err := p.expect(token.TokenVar); err != nil {
		return nil, err
	}

	if err := p.expect(token.TokenLBrace); err != nil {
		return nil, err
	}

	var vars []ast.Var

	for p.cur.Type != token.TokenRBrace && p.cur.Type != token.TokenEOF {
		v, err := p.parseVarLine()
		if err != nil {
			return nil, err
		}
		vars = append(vars, v)
	}

	if err := p.expect(token.TokenRBrace); err != nil {
		return nil, err
	}

	return vars, nil
}

func (p *Parser) parseVarLine() (ast.Var, error) {
	if p.cur.Type != token.TokenIdent {
		return ast.Var{}, p.errf("expected variable name, got %s %q", p.cur.Type, p.cur.Lit)
	}

	name := p.cur.Lit
	p.next()

	if err := p.expect(token.TokenEqual); err != nil {
		return ast.Var{}, err
	}

	expr, err := p.parseExpr(0)
	if err != nil {
		return ast.Var{}, err
	}

	return ast.Var{
		Name: name,
		Expr: expr,
	}, nil
}

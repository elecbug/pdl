package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/elecbug/pdl/internal/document"
	"github.com/elecbug/pdl/internal/document/order"
	"github.com/elecbug/pdl/internal/lexer"
	"github.com/elecbug/pdl/internal/token"
)

// Parser consumes lexer tokens and constructs a Document.
// It keeps current and lookahead tokens to support predictive parsing.
type Parser struct {
	// The lexer instance used to tokenize the input source code.
	l *lexer.Lexer

	// The current token being processed by the parser.
	cur token.Token
	// The next token in the stream, used for lookahead to make parsing decisions.
	peek token.Token
}

// New creates a parser for the given input and primes current/peek tokens.
func New(input string) *Parser {
	l := lexer.New(input)

	p := &Parser{l: l}
	p.next()
	p.next()

	return p
}

// ParseString parses input into a Document.
func ParseString(input string) (*document.Document, error) {
	return New(input).Parse()
}

// ParseWithMultiSources parses multiple input strings into a DocumentSet, ensuring that each document has a
// unique packet name and that the first document is set as the root for decoding.
func ParseWithMultiSources(root string, inputs ...string) (*document.DocumentSet, error) {
	set := &document.DocumentSet{
		Documents: make(map[string]*document.Document),
	}

	for i, input := range inputs {
		doc, err := ParseString(input)
		if err != nil {
			return nil, fmt.Errorf("parse input[%d]: %w", i, err)
		}

		if doc.PacketName == "" {
			return nil, fmt.Errorf("parse input[%d]: packet name is empty", i)
		}

		if _, exists := set.Documents[doc.PacketName]; exists {
			return nil, fmt.Errorf("duplicate packet %q", doc.PacketName)
		}

		set.Documents[doc.PacketName] = doc

		if set.Root == nil && doc.PacketName == root {
			set.Root = doc
		}
	}

	return set, nil
}

// Parse processes the token stream and constructs a Document structure based on the PDL source code.
// It handles the various sections of the document, such as packet definition, mode settings, variable
// definitions, field definitions, and output specifications. If it encounters any syntax errors or
// unexpected tokens, it returns an error with a descriptive message.
func (p *Parser) Parse() (*document.Document, error) {
	doc := &document.Document{}

	for p.cur.Type != token.EOF {
		switch p.cur.Type {
		case token.PACKET_KEYWORD:
			if err := p.parsePacket(doc); err != nil {
				return nil, err
			}

		case token.SET_KEYWORD:
			if err := p.parseSetMode(doc); err != nil {
				return nil, err
			}

		case token.DEF_KEYWORD:
			defs, err := p.parseDefBlock()
			if err != nil {
				return nil, err
			}
			doc.Defs = append(doc.Defs, defs...)

		case token.OUT_KEYWORD:
			outs, err := p.parseOutBlock()
			if err != nil {
				return nil, err
			}
			doc.Outs = append(doc.Outs, outs...)
		case token.VAR_KEYWORD:
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

// parsePacket parses the packet definition section of the document, expecting the "packet" keyword
// followed by an identifier for the packet name.
func (p *Parser) parsePacket(doc *document.Document) error {
	if err := p.expect(token.PACKET_KEYWORD); err != nil {
		return err
	}

	if p.cur.Type != token.IDENT {
		return p.errf("expected packet name, got %s %q", p.cur.Type, p.cur.Lit)
	}

	doc.PacketName = p.cur.Lit
	p.next()
	return nil
}

// parseSetMode parses "set mode" with byte order and bit order values.
// It updates doc.ByteOrder and doc.BitOrder.
func (p *Parser) parseSetMode(doc *document.Document) error {
	if err := p.expect(token.SET_KEYWORD); err != nil {
		return err
	}

	if err := p.expect(token.MODE_KEYWORD); err != nil {
		return err
	}

	if p.cur.Type != token.IDENT {
		return p.errf("expected byte order")
	}

	switch p.cur.Lit {
	case "BIG_ENDIAN":
		doc.ByteOrder = order.BIG_ENDIAN

	case "LITTLE_ENDIAN":
		doc.ByteOrder = order.LITTLE_ENDIAN

	default:
		return p.errf("unknown byte order %q", p.cur.Lit)
	}

	p.next()

	if p.cur.Type != token.IDENT {
		return p.errf("expected bit order")
	}

	switch p.cur.Lit {
	case "MSB_FIRST":
		doc.BitOrder = order.MSB_FIRST

	case "LSB_FIRST":
		doc.BitOrder = order.LSB_FIRST

	default:
		return p.errf("unknown bit order %q", p.cur.Lit)
	}

	p.next()

	return nil
}

// parseVarBlock parses a var block enclosed in braces.
// It returns parsed variable definitions.
func (p *Parser) parseVarBlock() ([]document.Var, error) {
	if err := p.expect(token.VAR_KEYWORD); err != nil {
		return nil, err
	}

	if err := p.expect(token.LBRACE_SIGN); err != nil {
		return nil, err
	}

	var vars []document.Var

	for p.cur.Type != token.RBRACE_SIGN && p.cur.Type != token.EOF {
		v, err := p.parseVarLine()
		if err != nil {
			return nil, err
		}
		vars = append(vars, v)
	}

	if err := p.expect(token.RBRACE_SIGN); err != nil {
		return nil, err
	}

	return vars, nil
}

// parseVarLine parses one variable definition: name = expression.
func (p *Parser) parseVarLine() (document.Var, error) {
	if p.cur.Type != token.IDENT {
		return document.Var{}, p.errf("expected variable name, got %s %q", p.cur.Type, p.cur.Lit)
	}

	name := p.cur.Lit
	p.next()

	if err := p.expect(token.EQUAL_SIGN); err != nil {
		return document.Var{}, err
	}

	expr, err := p.parseExpr(0)
	if err != nil {
		return document.Var{}, err
	}

	return document.Var{
		Name: name,
		Expr: expr,
	}, nil
}

// parseDefBlock parses a def block enclosed in braces.
// It returns parsed field definitions.
func (p *Parser) parseDefBlock() ([]document.Def, error) {
	if err := p.expect(token.DEF_KEYWORD); err != nil {
		return nil, err
	}

	if err := p.expect(token.LBRACE_SIGN); err != nil {
		return nil, err
	}

	var defs []document.Def

	for p.cur.Type != token.RBRACE_SIGN && p.cur.Type != token.EOF {
		def, err := p.parseDefLine()
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}

	if err := p.expect(token.RBRACE_SIGN); err != nil {
		return nil, err
	}

	return defs, nil
}

// parseDefLine parses one field definition line.
// Supported forms are "from X length Y" and "from X to Y".
func (p *Parser) parseDefLine() (document.Def, error) {
	if p.cur.Type != token.IDENT {
		return document.Def{}, p.errf("expected field name in def block, got %s %q", p.cur.Type, p.cur.Lit)
	}

	def := document.Def{Name: p.cur.Lit}
	p.next()

	if err := p.expect(token.FROM_KEYWORD); err != nil {
		return document.Def{}, err
	}

	from, err := p.parseExpr(0)
	if err != nil {
		return document.Def{}, err
	}
	def.From = from

	switch p.cur.Type {
	case token.LENGTH_KEYWORD:
		p.next()

		length, err := p.parseExpr(0)
		if err != nil {
			return document.Def{}, err
		}

		def.Length = length
		def.UseLength = true

	case token.TO_KEYWORD:
		p.next()

		to, err := p.parseExpr(0)
		if err != nil {
			return document.Def{}, err
		}

		def.To = to
		def.UseTo = true

	default:
		return document.Def{}, p.errf("expected length or to in def line, got %s %q", p.cur.Type, p.cur.Lit)
	}

	return def, nil
}

// parseOutBlock parses an out json block enclosed in braces.
// It returns parsed output rules.
func (p *Parser) parseOutBlock() ([]document.Out, error) {
	if err := p.expect(token.OUT_KEYWORD); err != nil {
		return nil, err
	}

	if err := p.expect(token.JSON_KEYWORD); err != nil {
		return nil, err
	}

	if err := p.expect(token.LBRACE_SIGN); err != nil {
		return nil, err
	}

	var outs []document.Out

	for p.cur.Type != token.RBRACE_SIGN && p.cur.Type != token.EOF {
		out, err := p.parseOutLine()
		if err != nil {
			return nil, err
		}
		outs = append(outs, out)
	}

	if err := p.expect(token.RBRACE_SIGN); err != nil {
		return nil, err
	}

	return outs, nil
}

// parseOutLine parses one output rule.
// It expects a field name, optional bit index, output path, and either format
// or map block.
func (p *Parser) parseOutLine() (document.Out, error) {
	if p.cur.Type != token.IDENT {
		return document.Out{}, p.errf("expected field name in out block, got %s %q", p.cur.Type, p.cur.Lit)
	}

	out := document.Out{
		Field: p.cur.Lit,
	}

	p.next()

	if p.cur.Type == token.LANGLE_SIGN {
		p.next()

		if p.cur.Type != token.NUMBER {
			return document.Out{}, p.errf("expected bit index number")
		}

		v, err := parseNumber(p.cur.Lit)
		if err != nil {
			return document.Out{}, p.errf("invalid bit index %q", p.cur.Lit)
		}

		out.HasBitIndex = true
		out.BitIndex = int(v)

		p.next()

		if err := p.expect(token.RANGLE_SIGN); err != nil {
			return document.Out{}, err
		}
	}

	if p.cur.Type != token.IDENT {
		return document.Out{}, p.errf("expected output path, got %s %q", p.cur.Type, p.cur.Lit)
	}

	out.Path = p.cur.Lit
	p.next()

	if p.cur.Type == token.AS_KEYWORD {
		p.next()

		if p.cur.Type != token.IDENT {
			return document.Out{}, p.errf("expected packet name after as")
		}

		out.AsPacket = p.cur.Lit
		p.next()
		return out, nil
	}

	if p.cur.Type == token.IDENT {
		out.Format = p.cur.Lit
		p.next()
		return out, nil
	}

	if p.cur.Type == token.LBRACE_SIGN {
		m, err := p.parseMapBlock()
		if err != nil {
			return document.Out{}, err
		}

		out.Map = m
		return out, nil
	}

	return document.Out{}, p.errf("expected format, as packet, or map block in out line, got %s %q", p.cur.Type, p.cur.Lit)
}

// parseMapBlock parses a brace-enclosed key-value mapping.
func (p *Parser) parseMapBlock() (map[string]string, error) {
	if err := p.expect(token.LBRACE_SIGN); err != nil {
		return nil, err
	}

	m := make(map[string]string)

	for p.cur.Type != token.RBRACE_SIGN && p.cur.Type != token.EOF {
		var key string

		switch p.cur.Type {
		case token.NUMBER:
			v, err := parseNumber(p.cur.Lit)
			if err != nil {
				return nil, p.errf("invalid map key %q", p.cur.Lit)
			}
			key = strconv.FormatInt(v, 10)

		case token.IDENT:
			key = p.cur.Lit

		default:
			return nil, p.errf("expected map key, got %s %q", p.cur.Type, p.cur.Lit)
		}

		p.next()

		if err := p.expect(token.COLON_SIGN); err != nil {
			return nil, err
		}

		var val string

		switch p.cur.Type {
		case token.STRING, token.IDENT, token.NUMBER:
			val = p.cur.Lit
		default:
			return nil, p.errf("expected map value, got %s %q", p.cur.Type, p.cur.Lit)
		}

		m[key] = val
		p.next()
	}

	if err := p.expect(token.RBRACE_SIGN); err != nil {
		return nil, err
	}

	return m, nil
}

// parseExpr parses an expression using precedence climbing.
// minPrec is the minimum precedence required to continue parsing.
func (p *Parser) parseExpr(minPrec int) (document.Expr, error) {
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

		left = document.BinaryExpr{
			Op:    op,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

// parsePrimary parses a primary expression:
// number, identifier, end, field reference (*name), or parenthesized expression.
func (p *Parser) parsePrimary() (document.Expr, error) {
	switch p.cur.Type {
	case token.NUMBER:
		raw := p.cur.Lit
		v, err := parseNumber(raw)
		if err != nil {
			return nil, p.errf("invalid number %q", raw)
		}
		p.next()
		return document.NumberExpr{Raw: raw, Value: v}, nil

	case token.IDENT:
		name := p.cur.Lit
		p.next()
		return document.IdentExpr{Name: name}, nil

	case token.END_KEYWORD:
		p.next()
		return document.EndExpr{}, nil

	case token.STAR_SIGN:
		p.next()

		if p.cur.Type != token.IDENT {
			return nil, p.errf("expected field name after *")
		}

		name := p.cur.Lit
		p.next()

		return document.FieldValueExpr{Name: name}, nil

	case token.LPAREN_SIGN:
		p.next()

		expr, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}

		if err := p.expect(token.RPAREN_SIGN); err != nil {
			return nil, err
		}

		return expr, nil

	default:
		return nil, p.errf("expected expression, got %s %q", p.cur.Type, p.cur.Lit)
	}
}

// parseNumber parses decimal, binary (0b), or hexadecimal (0x) literals.
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

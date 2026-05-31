package ast

type Expr interface {
	exprNode()
}

type NumberExpr struct {
	Raw   string
	Value int64
}

func (NumberExpr) exprNode() {}

type IdentExpr struct {
	Name string
}

func (IdentExpr) exprNode() {}

type FieldValueExpr struct {
	Name string
}

func (FieldValueExpr) exprNode() {}

type EndExpr struct{}

func (EndExpr) exprNode() {}

type BinaryExpr struct {
	Op    string
	Left  Expr
	Right Expr
}

func (BinaryExpr) exprNode() {}

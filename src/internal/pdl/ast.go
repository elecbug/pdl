package pdl

type Document struct {
	PacketName string
	Mode       string

	Vars    []Var
	Defs    []Def
	Outputs []Output
}

type Var struct {
	Name string
	Expr Expr
}

type Def struct {
	Name string

	From Expr

	Length Expr
	To     Expr

	UseLength bool
	UseTo     bool
}

type Output struct {
	Field  string
	Path   string
	Format string

	BitIndex *int
	Map      map[string]string
}

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

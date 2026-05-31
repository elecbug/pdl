package ast

type Document struct {
	PacketName string

	ByteOrder ByteOrder
	BitOrder  BitOrder

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

package document

// Expr represents an expression in the document, which can be a number, an identifier, a field value,
// an end marker, or a binary expression.
type Expr interface {
	// exprNode is a marker method to indicate that a type implements the Expr interface. It does not
	// perform any operations and is used solely for type identification.
	exprNode()
}

// NumberExpr represents a numeric literal expression, containing the raw string representation of
// the number and its evaluated integer value.
type NumberExpr struct {
	// Raw is the original string representation of the number as it appears in the document.
	Raw string
	// Value is the evaluated integer value of the number, derived from the Raw string.
	Value int64
}

// exprNode is a marker method to indicate that NumberExpr implements the Expr interface. It does not
// perform any operations and is used solely for type identification.
func (NumberExpr) exprNode() {}

// IdentExpr represents an identifier expression, which refers to a variable defined in the document.
// It contains the name of the variable.
type IdentExpr struct {
	// Name is the name of the variable being referenced by this identifier expression. It should correspond
	// to a variable defined in the document's Vars section.
	Name string
}

// exprNode is a marker method to indicate that IdentExpr implements the Expr interface. It does not
// perform any operations and is used solely for type identification.
func (IdentExpr) exprNode() {}

// FieldValueExpr represents an expression that references the value of a field defined in the document.
// It contains the name of the field being referenced.
type FieldValueExpr struct {
	// Name is the name of the field whose value is being referenced by this expression. It should correspond
	// to a field defined in the document's Defs section.
	Name string
}

// exprNode is a marker method to indicate that FieldValueExpr implements the Expr interface. It does not
// perform any operations and is used solely for type identification.
func (FieldValueExpr) exprNode() {}

// EndExpr represents an expression that indicates the end of the input data, which can be used to calculate
// the length of a field based on the remaining bits in the input. It does not contain any fields.
type EndExpr struct{}

// exprNode is a marker method to indicate that EndExpr implements the Expr interface. It does not
// perform any operations and is used solely for type identification.
func (EndExpr) exprNode() {}

// BinaryExpr combines two sub-expressions with an operator (for example,
// +, -, *, or /).
// It contains the operator and the left and right sub-expressions.
type BinaryExpr struct {
	// Op is the operator used in the binary expression, such as "+", "-", "*", or "/".
	Op string
	// Left is the left sub-expression.
	Left Expr
	// Right is the right sub-expression.
	Right Expr
}

// exprNode is a marker method to indicate that BinaryExpr implements the Expr interface. It does not
// perform any operations and is used solely for type identification.
func (BinaryExpr) exprNode() {}

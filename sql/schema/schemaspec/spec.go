package schemaspec

import (
	"fmt"
	"strconv"
)

type (
	// Spec holds a specification for a schema resource (such as a Table, Column or Index).
	Spec interface {
		spec()
	}

	// Encoder is the interface that wraps the Encode method.
	//
	// Encoder takes a Spec and returns a byte slice representing that Spec in some configuration
	// format (for instance, HCL).
	Encoder interface {
		Encode(Spec) ([]byte, error)
	}

	// Decoder is the interface that wraps the Decode method.
	//
	// Decoder takes a byte slice representing a Spec and decodes it into a Spec.
	Decoder interface {
		Decode([]byte, Spec) error
	}

	// Codec wraps Encoder and Decoder.
	Codec interface {
		Encoder
		Decoder
	}

	// Resource is a generic container for resources described in configurations.
	Resource struct {
		Name     string
		Type     string
		Attrs    []*Attr
		Children []*Resource
	}

	// Schema holds a specification for a Schema.
	Schema struct {
		Name   string
		Tables []*Table
	}

	// Table holds a specification for an SQL table.
	Table struct {
		Name        string
		SchemaName  string
		Columns     []*Column
		PrimaryKey  *PrimaryKey
		ForeignKeys []*ForeignKey
		Indexes     []*Index
		Attrs       []*Attr
		Children    []*Resource
	}

	// Column holds a specification for a column in an SQL table.
	Column struct {
		Name     string
		Type     string
		Default  *LiteralValue
		Null     bool
		Attrs    []*Attr
		Children []*Resource
	}

	// PrimaryKey holds a specification for the primary key of a table.
	PrimaryKey struct {
		Columns  []*ColumnRef
		Attrs    []*Attr
		Children []*Resource
	}

	// ForeignKey holds a specification for a foreign key of a table.
	ForeignKey struct {
		Symbol     string
		Columns    []*ColumnRef
		RefColumns []*ColumnRef
		OnUpdate   string
		OnDelete   string
		Attrs      []*Attr
		Children   []*Resource
	}

	// Index holds a specification for an index of a table.
	Index struct {
		Name     string
		Columns  []*ColumnRef
		Unique   bool
		Attrs    []*Attr
		Children []*Resource
	}

	// ColumnRef is a reference to a Column described in another spec.
	ColumnRef struct {
		Name  string
		Table string
	}

	// TableRef is a reference to a Table described in another spec.
	TableRef struct {
		Name   string
		Schema string
	}

	// Element is an object that can be encoded into bytes to be written to a configuration file representing
	// Schema resources.
	Element interface {
		elem()
	}

	// Attr is an attribute of a Spec.
	Attr struct {
		K string
		V Value
	}

	// Value represents the value of a Attr.
	Value interface {
		val()
	}

	// LiteralValue implements Value and represents a literal value (string, number, etc.)
	LiteralValue struct {
		V string
	}

	// ListValue implements Value and represents a list of literal value (string, number, etc.)
	ListValue struct {
		V []string
	}
)

// Attr returns the value of the Column attribute named `name` and reports whether such an attribute exists.
func (c *Column) Attr(name string) (*Attr, bool) {
	return getAttrVal(c.Attrs, name)
}

// Attr returns the value of the Table attribute named `name` and reports whether such an attribute exists.
func (t *Table) Attr(name string) (*Attr, bool) {
	return getAttrVal(t.Attrs, name)
}

func getAttrVal(attrs []*Attr, name string) (*Attr, bool) {
	for _, attr := range attrs {
		if attr.K == name {
			return attr, true
		}
	}
	return nil, false
}

// Int returns an integer from the Value of the Attr. If The value is not a LiteralValue or the value
// cannot be converted to an integer an error is returned.
func (a *Attr) Int() (int, error) {
	lit, ok := a.V.(*LiteralValue)
	if !ok {
		return 0, fmt.Errorf("schema: cannot read attribute %q as literal", a.K)
	}
	s, err := strconv.Atoi(lit.V)
	if err != nil {
		return 0, fmt.Errorf("schema: cannot read attribute %q as integer", a.K)
	}
	return s, nil
}

// Strings returns a slice of strings from the Value of the Attr. If The value is not a ListValue or the its
// values cannot be converted to strings an error is returned.
func (a *Attr) Strings() ([]string, error) {
	lst, ok := a.V.(*ListValue)
	if !ok {
		return nil, fmt.Errorf("schema: attribute %q is not a list", a.K)
	}
	out := make([]string, 0, len(lst.V))
	for _, item := range lst.V {
		unquote, err := strconv.Unquote(item)
		if err != nil {
			return nil, fmt.Errorf("schema: failed parsing item %q to string: %w", item, err)
		}
		out = append(out, unquote)
	}
	return out, nil
}

func (*LiteralValue) val() {}
func (*ListValue) val()    {}

func (*Resource) elem()   {}
func (*Attr) elem()       {}
func (*Column) elem()     {}
func (*Table) elem()      {}
func (*Schema) elem()     {}
func (*PrimaryKey) elem() {}
func (*ForeignKey) elem() {}
func (*Index) elem()      {}

func (*Column) spec()     {}
func (*Table) spec()      {}
func (*Schema) spec()     {}
func (*Resource) spec()   {}
func (*PrimaryKey) spec() {}
func (*ForeignKey) spec() {}
func (*Index) spec()      {}

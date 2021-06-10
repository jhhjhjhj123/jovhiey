package mysql

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"ariga.io/atlas/sql/schema"

	"golang.org/x/mod/semver"
)

// A Diff provides diff capabilities for schema elements.
type Diff struct{ *Driver }

// SchemaDiff implements the schema.Differ interface and returns a list of
// changes that need to be applied in order to move from one state to the other.
func (d *Diff) SchemaDiff(from, to *schema.Schema) ([]schema.Change, error) {
	var changes []schema.Change
	// Charset change.
	if change := d.charsetChange(from.Attrs, from.Realm.Attrs, to.Attrs); change != noChange {
		changes = append(changes, change)
	}
	// Collation change.
	if change := d.collationChange(from.Attrs, from.Realm.Attrs, to.Attrs); change != noChange {
		changes = append(changes, change)
	}

	// Drop or modify tables.
	for _, t1 := range from.Tables {
		t2, ok := to.Table(t1.Name)
		if !ok {
			changes = append(changes, &schema.DropTable{T: t1})
			continue
		}
		change, err := d.TableDiff(t1, t2)
		if err != nil {
			return nil, err
		}
		if len(change) > 0 {
			changes = append(changes, &schema.ModifyTable{
				T:       t1,
				Changes: change,
			})
		}
	}
	// Add tables.
	for _, t1 := range to.Tables {
		if _, ok := from.Table(t1.Name); !ok {
			changes = append(changes, &schema.AddTable{T: t1})
		}
	}
	return changes, nil
}

// TableDiff implements the schema.TableDiffer interface and returns a list of
// changes that need to be applied in order to move from one state to the other.
func (d *Diff) TableDiff(from, to *schema.Table) ([]schema.Change, error) {
	var changes []schema.Change
	// PK modification is not support.
	if pk1, pk2 := from.PrimaryKey, to.PrimaryKey; (pk1 != nil) != (pk2 != nil) || (pk1 != nil) && d.indexChange(pk1, pk2) != schema.NoChange {
		return nil, fmt.Errorf("changing %q table primary key is not supported", to.Name)
	}

	// Charset change.
	if change := d.charsetChange(from.Attrs, from.Schema.Attrs, to.Attrs); change != noChange {
		changes = append(changes, change)
	}
	// Collation change.
	if change := d.collationChange(from.Attrs, from.Schema.Attrs, to.Attrs); change != noChange {
		changes = append(changes, change)
	}

	// Drop or modify checks.
	for _, c1 := range checks(from.Attrs) {
		switch c2, ok := checkByName(to.Attrs, c1.Name); {
		case !ok:
			changes = append(changes, &schema.DropAttr{
				A: c1,
			})
		case c1.Clause != c2.Clause || c1.Enforced != c2.Enforced:
			changes = append(changes, &schema.ModifyAttr{
				From: c1,
				To:   c2,
			})
		}
	}
	// Add checks.
	for _, c1 := range checks(to.Attrs) {
		if _, ok := checkByName(from.Attrs, c1.Name); !ok {
			changes = append(changes, &schema.AddAttr{
				A: c1,
			})
		}
	}

	// Drop or modify columns.
	for _, c1 := range from.Columns {
		c2, ok := to.Column(c1.Name)
		if !ok {
			changes = append(changes, &schema.DropColumn{C: c1})
			continue
		}
		change, err := d.columnChange(c1, c2)
		if err != nil {
			return nil, err
		}
		if change != schema.NoChange {
			changes = append(changes, &schema.ModifyColumn{
				From:   c1,
				To:     c2,
				Change: change,
			})
		}
	}
	// Add columns.
	for _, c1 := range to.Columns {
		if _, ok := from.Column(c1.Name); !ok {
			changes = append(changes, &schema.AddColumn{C: c1})
		}
	}

	// Drop or modify indexes.
	for _, idx1 := range from.Indexes {
		idx2, ok := to.Index(idx1.Name)
		if !ok {
			changes = append(changes, &schema.DropIndex{I: idx1})
			continue
		}
		if change := d.indexChange(idx1, idx2); change != schema.NoChange {
			changes = append(changes, &schema.ModifyIndex{
				From:   idx1,
				To:     idx2,
				Change: change,
			})
		}
	}
	// Add indexes.
	for _, idx1 := range to.Indexes {
		if _, ok := from.Index(idx1.Name); !ok {
			changes = append(changes, &schema.AddIndex{I: idx1})
		}
	}

	// Drop or modify foreign-keys.
	for _, fk1 := range from.ForeignKeys {
		fk2, ok := to.ForeignKey(fk1.Symbol)
		if !ok {
			changes = append(changes, &schema.DropForeignKey{F: fk1})
			continue
		}
		if change := d.fkChange(fk1, fk2); change != schema.NoChange {
			changes = append(changes, &schema.ModifyForeignKey{
				From:   fk1,
				To:     fk2,
				Change: change,
			})
		}
	}
	// Add foreign-keys.
	for _, fk1 := range to.ForeignKeys {
		if _, ok := from.ForeignKey(fk1.Symbol); !ok {
			changes = append(changes, &schema.AddForeignKey{F: fk1})
		}
	}
	return changes, nil
}

// columnChange returns the schema changes (if any) for migrating one column to the other.
func (d *Diff) columnChange(from, to *schema.Column) (schema.ChangeKind, error) {
	var change schema.ChangeKind
	if from.Type.Null != to.Type.Null {
		change |= schema.ChangeNull
	}
	change |= commentChange(from.Attrs, to.Attrs)
	var c1, c2 schema.Collation
	if has(from.Attrs, &c1) != has(to.Attrs, &c2) || c1.V != c2.V {
		change |= schema.ChangeCollation
	}
	var cr1, cr2 schema.Charset
	if has(from.Attrs, &cr1) != has(to.Attrs, &cr2) || cr1.V != cr2.V {
		change |= schema.ChangeCharset
	}
	typeChanged, err := d.typeChanged(from, to)
	if err != nil {
		return schema.NoChange, err
	}
	if typeChanged {
		change |= schema.ChangeType
	}
	d1, d2 := from.Default, to.Default
	if (d1 != nil) != (d2 != nil) || (d1 != nil && d1.(*schema.RawExpr).X != d2.(*schema.RawExpr).X) {
		change |= schema.ChangeDefault
	}
	return change, nil
}

func (d *Diff) typeChanged(from, to *schema.Column) (bool, error) {
	fromT, toT := from.Type.Type, to.Type.Type
	if fromT == nil || toT == nil {
		return false, fmt.Errorf("missing type infromation for column %q", from.Name)
	}
	if reflect.TypeOf(fromT) != reflect.TypeOf(toT) {
		// TODO: Add type conversion rating (need check, can fail, etc).
		return true, nil
	}
	var changed bool
	switch fromT := fromT.(type) {
	case *schema.BinaryType:
		toT := toT.(*schema.BinaryType)
		changed = fromT.T != toT.T || fromT.Size != toT.Size
	case *schema.BoolType:
		toT := toT.(*schema.BinaryType)
		changed = fromT.T != toT.T
	case *schema.DecimalType:
		toT := toT.(*schema.DecimalType)
		changed = fromT.T != toT.T || fromT.Scale != toT.Scale || fromT.Precision != toT.Precision
	case *schema.EnumType:
		toT := toT.(*schema.EnumType)
		changed = !valuesEqual(fromT.Values, toT.Values)
	case *schema.FloatType:
		toT := toT.(*schema.FloatType)
		changed = fromT.T != toT.T || fromT.Precision != toT.Precision
	case *schema.IntegerType:
		toT := toT.(*schema.IntegerType)
		// MySQL v8.0.19 dropped the display width
		// information from the information schema.
		if semver.Compare("v"+d.version, "v8.0.19") == -1 {
			ft, _, _, err := parseColumn(fromT.T)
			if err != nil {
				return false, err
			}
			tt, _, _, err := parseColumn(toT.T)
			if err != nil {
				return false, err
			}
			fromT.T, toT.T = ft[0], tt[0]
		}
		fromW, toW := displayWidth(fromT.Attrs), displayWidth(toT.Attrs)
		changed = fromT.T != toT.T || fromT.Unsigned != toT.Unsigned ||
			(fromW != nil) != (toW != nil) || (fromW != nil && fromW.N != toW.N)
	case *schema.JSONType:
		toT := toT.(*schema.JSONType)
		changed = fromT.T != toT.T
	case *schema.StringType:
		toT := toT.(*schema.StringType)
		changed = fromT.T != toT.T || fromT.Size != toT.Size
	case *schema.SpatialType:
		toT := toT.(*schema.SpatialType)
		changed = fromT.T != toT.T
	case *schema.TimeType:
		toT := toT.(*schema.TimeType)
		changed = fromT.T != toT.T
	case *BitType:
		toT := toT.(*BitType)
		changed = fromT.T != toT.T
	case *SetType:
		toT := toT.(*SetType)
		changed = !valuesEqual(fromT.Values, toT.Values)
	default:
		return false, fmt.Errorf("unsupported type %T", fromT)
	}
	return changed, nil
}

// indexChange returns the schema changes (if any) for migrating one index to the other.
func (d *Diff) indexChange(from, to *schema.Index) schema.ChangeKind {
	var change schema.ChangeKind
	change |= partsChange(from.Parts, to.Parts)
	change |= commentChange(from.Attrs, to.Attrs)
	change |= indexTypeChange(from.Attrs, to.Attrs)
	if from.Unique != to.Unique {
		change |= schema.ChangeUnique
	}
	return change
}

// collationChange returns the schema change for migrating the collation if
// it was changed and its not the default attribute inherited from its parent.
func (d *Diff) collationChange(from, top, to []schema.Attr) schema.Change {
	var fromC, topC, toC schema.Collation
	switch fromHas, topHas, toHas := has(from, &fromC), has(top, &topC), has(to, &toC); {
	case !fromHas && !toHas:
	case !fromHas:
		return &schema.AddAttr{
			A: &toC,
		}
	case !toHas:
		if !topHas || fromC.V != topC.V {
			return &schema.DropAttr{
				A: &fromC,
			}
		}
	case fromC.V != toC.V:
		return &schema.ModifyAttr{
			From: &fromC,
			To:   &toC,
		}
	}
	return noChange
}

// charsetChange returns the schema change for migrating the collation if
// it was changed and its not the default attribute inherited from its parent.
func (d *Diff) charsetChange(from, top, to []schema.Attr) schema.Change {
	var fromC, topC, toC schema.Charset
	switch fromHas, topHas, toHas := has(from, &fromC), has(top, &topC), has(to, &toC); {
	case !fromHas && !toHas:
	case !fromHas:
		return &schema.AddAttr{
			A: &toC,
		}
	case !toHas:
		if !topHas || fromC.V != topC.V {
			return &schema.DropAttr{
				A: &fromC,
			}
		}
	case fromC.V != toC.V:
		return &schema.ModifyAttr{
			From: &fromC,
			To:   &toC,
		}
	}
	return noChange
}

// fkChange returns the schema changes (if any) for migrating one index to the other.
func (d *Diff) fkChange(from, to *schema.ForeignKey) schema.ChangeKind {
	var change schema.ChangeKind
	switch {
	case from.Table.Name != to.Table.Name:
		change |= schema.ChangeRefTable | schema.ChangeRefColumn
	case len(from.RefColumns) != len(to.RefColumns):
		change |= schema.ChangeRefColumn
	default:
		for i := range from.RefColumns {
			if from.RefColumns[i].Name != to.RefColumns[i].Name {
				change |= schema.ChangeRefColumn
			}
		}
	}
	switch {
	case len(from.Columns) != len(to.Columns):
		change |= schema.ChangeColumn
	default:
		for i := range from.Columns {
			if from.Columns[i].Name != to.Columns[i].Name {
				change |= schema.ChangeColumn
			}
		}
	}
	if actionChanged(from.OnUpdate, to.OnUpdate) {
		change |= schema.ChangeUpdateAction
	}
	if actionChanged(from.OnDelete, to.OnDelete) {
		change |= schema.ChangeDeleteAction
	}
	return change
}

func actionChanged(from, to schema.ReferenceOption) bool {
	// According to MySQL docs, specifying RESTRICT (or NO ACTION)
	// is the same as omitting the ON DELETE or ON UPDATE clause.
	if from == "" || from == schema.Restrict {
		from = schema.NoAction
	}
	if to == "" || to == schema.Restrict {
		to = schema.NoAction
	}
	return from != to
}

func commentChange(from, to []schema.Attr) schema.ChangeKind {
	var c1, c2 schema.Comment
	if has(from, &c1) != has(to, &c2) || c1.Text != c2.Text {
		return schema.ChangeComment
	}
	return schema.NoChange
}

func indexTypeChange(from, to []schema.Attr) schema.ChangeKind {
	var c1, c2 IndexType
	if has(from, &c1) != has(to, &c2) || c1.T != c2.T {
		return schema.ChangeAttr
	}
	return schema.NoChange
}

// indexCollation returns the index collation from its attribute.
// The default collation is ascending if no order was specified.
func indexCollation(attr []schema.Attr) *schema.Collation {
	c := &schema.Collation{V: "A"}
	if has(attr, c) {
		c.V = strings.ToUpper(c.V)
	}
	return c
}

func partsChange(from, to []*schema.IndexPart) schema.ChangeKind {
	if len(from) != len(to) {
		return schema.ChangeParts
	}
	// MySQL starts counting the sequence number from 1, but internal tools start counting
	// from 0. Therefore, we care only about the parts order and not their seqno attribute.
	sort.Slice(to, func(i, j int) bool { return to[i].SeqNo < to[j].SeqNo })
	sort.Slice(from, func(i, j int) bool { return from[i].SeqNo < from[j].SeqNo })
	for i := range from {
		switch {
		case indexCollation(from[i].Attrs).V != indexCollation(to[i].Attrs).V:
			return schema.ChangeParts
		case from[i].C != nil && to[i].C != nil:
			if from[i].C.Name != to[i].C.Name {
				return schema.ChangeParts
			}
			var s1, s2 SubPart
			if has(from[i].Attrs, &s1) != has(to[i].Attrs, &s2) || s1.Len != s2.Len {
				return schema.ChangeParts
			}
		case from[i].X != nil && to[i].X != nil:
			if from[i].X.(*schema.RawExpr).X != to[i].X.(*schema.RawExpr).X {
				return schema.ChangeParts
			}
		default: // (C1 != nil) != (C2 != nil) || (X1 != nil) != (X2 != nil).
			return schema.ChangeParts
		}
	}
	return schema.NoChange
}

// noChange describes a zero change.
var noChange struct{ schema.Change }

func checks(attr []schema.Attr) (checks []*Check) {
	for i := range attr {
		if c, ok := attr[i].(*Check); ok {
			checks = append(checks, c)
		}
	}
	return checks
}

func checkByName(attr []schema.Attr, name string) (*Check, bool) {
	for i := range attr {
		if c, ok := attr[i].(*Check); ok && c.Name == name {
			return c, true
		}
	}
	return nil, false
}

func has(attrs []schema.Attr, target interface{}) bool {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("mysql: target must be a non-nil pointer")
	}
	for _, attr := range attrs {
		if reflect.TypeOf(attr).AssignableTo(rv.Type()) {
			rv.Elem().Set(reflect.ValueOf(attr).Elem())
			return true
		}
	}
	return false
}

func displayWidth(attr []schema.Attr) *DisplayWidth {
	var (
		z *ZeroFill
		d *DisplayWidth
	)
	for i := range attr {
		switch at := attr[i].(type) {
		case *ZeroFill:
			z = at
		case *DisplayWidth:
			d = at
		}
	}
	// Accept the display width only if
	// the zerofill attribute is defined.
	if z == nil || d == nil {
		return nil
	}
	return d
}

func valuesEqual(v1, v2 []string) bool {
top:
	for i := range v1 {
		for j := range v2 {
			if v1[i] == v2[j] {
				continue top
			}
		}
		return false
	}
	return true
}

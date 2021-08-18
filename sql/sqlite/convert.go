// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package sqlite

import (
	"fmt"
	"strings"

	"ariga.io/atlas/sql/internal/schemautil"
	"ariga.io/atlas/sql/internal/sqlx"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/schema/schemaspec"
)

// ConvertSchema converts a schemaspec.Schema into a schema.Schema.
func (d *Driver) ConvertSchema(spec *schemaspec.Schema, tables []*schemaspec.Table) (*schema.Schema, error) {
	return schemautil.ConvertSchema(spec, tables, d.ConvertTable)
}

// ConvertTable converts a schemaspec.Table to a schema.Table. Table conversion is done without converting
// ForeignKeySpecs into ForeignKeys, as the target tables do not necessarily exist in the schema
// at this point. Instead, the linking is done by the ConvertSchema function.
func (d *Driver) ConvertTable(spec *schemaspec.Table, parent *schema.Schema) (*schema.Table, error) {
	return schemautil.ConvertTable(spec, parent, d.ConvertColumn, d.ConvertPrimaryKey, d.ConvertIndex)
}

// ConvertPrimaryKey converts a schemaspec.PrimaryKey to a schema.Index.
func (d *Driver) ConvertPrimaryKey(spec *schemaspec.PrimaryKey, parent *schema.Table) (*schema.Index, error) {
	return schemautil.ConvertPrimaryKey(spec, parent)
}

// ConvertIndex converts an schemaspec.Index to a schema.Index.
func (d *Driver) ConvertIndex(spec *schemaspec.Index, parent *schema.Table) (*schema.Index, error) {
	return schemautil.ConvertIndex(spec, parent)
}

// ConvertColumn converts a schemaspec.Column into a schema.Column.
func (d *Driver) ConvertColumn(spec *schemaspec.Column, parent *schema.Table) (*schema.Column, error) {
	if override := spec.Override(sqlx.VersionPermutations(Name, d.version)...); override != nil {
		if err := schemautil.Override(spec, override); err != nil {
			return nil, err
		}
	}
	return schemautil.ConvertColumn(spec, ConvertColumnType)
}

// ConvertColumnType converts a schemaspec.Column into a concrete sqlite schema.Type.
func ConvertColumnType(spec *schemaspec.Column) (schema.Type, error) {
	switch spec.Type {
	case "int", "int8", "int16", "int64", "uint", "uint8", "uint16", "uint64":
		return convertInteger(spec)
	case "string":
		return convertString(spec)
	case "binary":
		return convertBinary(spec)
	case "enum":
		return convertEnum(spec)
	case "boolean":
		return convertBoolean(spec)
	case "decimal":
		return convertDecimal(spec)
	case "float":
		return convertFloat(spec)
	case "time":
		return convertTime(spec)
	}
	return parseRawType(spec.Type)
}

func convertInteger(spec *schemaspec.Column) (schema.Type, error) {
	if strings.HasPrefix(spec.Type, "u") {
		// todo(rotemtam): support his once we can express CHECK(col >= 0)
		return nil, fmt.Errorf("sqlite: unsigned integers currently not supported")
	}
	typ := &schema.IntegerType{
		T: "integer",
	}
	return typ, nil
}

func convertBinary(spec *schemaspec.Column) (schema.Type, error) {
	bt := &schema.BinaryType{
		T: "blob",
	}
	return bt, nil
}

func convertString(spec *schemaspec.Column) (schema.Type, error) {
	st := &schema.StringType{
		T: "text",
	}
	if attr, ok := spec.Attr("size"); ok {
		s, err := attr.Int()
		if err != nil {
			return nil, err
		}
		st.Size = s
	}
	return st, nil
}

func convertEnum(spec *schemaspec.Column) (schema.Type, error) {
	// sqlite does not have a enum column type
	return &schema.StringType{T: "text"}, nil
}

func convertBoolean(spec *schemaspec.Column) (schema.Type, error) {
	return &schema.BoolType{T: "boolean"}, nil
}

func convertTime(spec *schemaspec.Column) (schema.Type, error) {
	return &schema.TimeType{T: "datetime"}, nil
}

func convertDecimal(spec *schemaspec.Column) (schema.Type, error) {
	dt := &schema.DecimalType{
		T: "decimal",
	}
	if precision, ok := spec.Attr("precision"); ok {
		p, err := precision.Int()
		if err != nil {
			return nil, err
		}
		dt.Precision = p
	}
	if scale, ok := spec.Attr("scale"); ok {
		s, err := scale.Int()
		if err != nil {
			return nil, err
		}
		dt.Scale = s
	}
	return dt, nil
}

func convertFloat(spec *schemaspec.Column) (schema.Type, error) {
	ft := &schema.FloatType{
		T: "real",
	}
	if precision, ok := spec.Attr("precision"); ok {
		p, err := precision.Int()
		if err != nil {
			return nil, err
		}
		ft.Precision = p
	}
	return ft, nil
}

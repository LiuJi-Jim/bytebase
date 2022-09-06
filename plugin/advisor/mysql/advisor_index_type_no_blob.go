package mysql

// Framework code is generated by the generator.

import (
	"fmt"
	"strings"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/parser/types"
	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/plugin/advisor"
	"github.com/bytebase/bytebase/plugin/advisor/catalog"
	"github.com/bytebase/bytebase/plugin/advisor/db"
)

var (
	_ advisor.Advisor = (*IndexTypeNoBlobAdvisor)(nil)
	_ ast.Visitor     = (*indexTypeNoBlobChecker)(nil)
)

func init() {
	advisor.Register(db.MySQL, advisor.MySQLIndexTypeNoBlob, &IndexTypeNoBlobAdvisor{})
	advisor.Register(db.TiDB, advisor.MySQLIndexTypeNoBlob, &IndexTypeNoBlobAdvisor{})
}

// IndexTypeNoBlobAdvisor is the advisor checking for index type no blob.
type IndexTypeNoBlobAdvisor struct {
}

// Check checks for index type no blob.
func (*IndexTypeNoBlobAdvisor) Check(ctx advisor.Context, statement string) ([]advisor.Advice, error) {
	stmtList, errAdvice := parseStatement(statement, ctx.Charset, ctx.Collation)
	if errAdvice != nil {
		return errAdvice, nil
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &indexTypeNoBlobChecker{
		level:            level,
		title:            string(ctx.Rule.Type),
		database:         ctx.Database,
		tablesNewColumns: make(map[string]columnNameToColumnDef),
	}

	for _, stmt := range stmtList {
		checker.text = stmt.Text()
		checker.line = stmt.OriginTextPosition()
		(stmt).Accept(checker)
	}

	if len(checker.adviceList) == 0 {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return checker.adviceList, nil
}

type indexTypeNoBlobChecker struct {
	adviceList       []advisor.Advice
	level            advisor.Status
	title            string
	text             string
	line             int
	database         *catalog.Database
	tablesNewColumns tableNewColumn
}

// Enter implements the ast.Visitor interface.
func (v *indexTypeNoBlobChecker) Enter(in ast.Node) (ast.Node, bool) {
	var pkDataList []pkData
	switch node := in.(type) {
	case *ast.CreateTableStmt:
		tableName := node.Table.Name.String()
		for _, column := range node.Cols {
			pds := v.addNewColumn(tableName, column.OriginTextPosition(), column)
			pkDataList = append(pkDataList, pds...)
		}
		for _, constraint := range node.Constraints {
			pds := v.addConstraint(tableName, constraint.OriginTextPosition(), constraint)
			pkDataList = append(pkDataList, pds...)
		}
	case *ast.AlterTableStmt:
		tableName := node.Table.Name.String()
		for _, spec := range node.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns:
				for _, column := range spec.NewColumns {
					pds := v.addNewColumn(tableName, node.OriginTextPosition(), column)
					pkDataList = append(pkDataList, pds...)
				}
			case ast.AlterTableAddConstraint:
				pds := v.addConstraint(tableName, node.OriginTextPosition(), spec.Constraint)
				pkDataList = append(pkDataList, pds...)
			case ast.AlterTableChangeColumn, ast.AlterTableModifyColumn:
				newColumnDef := spec.NewColumns[0]
				oldColumnName := newColumnDef.Name.Name.String()
				if spec.OldColumnName != nil {
					oldColumnName = spec.OldColumnName.Name.String()
				}
				pds := v.changeColumn(tableName, oldColumnName, node.OriginTextPosition(), newColumnDef)
				pkDataList = append(pkDataList, pds...)
			}
		}
	case *ast.CreateIndexStmt:
		tableName := node.Table.Name.String()
		for _, indexSpec := range node.IndexPartSpecifications {
			pds := v.addIndex(tableName, node.OriginTextPosition(), indexSpec)
			pkDataList = append(pkDataList, pds...)
		}
	}
	for _, pd := range pkDataList {
		v.adviceList = append(v.adviceList, advisor.Advice{
			Status:  v.level,
			Code:    advisor.IndexTypeNoBlob,
			Title:   v.title,
			Content: fmt.Sprintf("Columns in index must not be BLOB but `%s`.`%s` is %s", pd.table, pd.column, pd.columnType),
			Line:    pd.line,
		})
	}
	return in, false
}

// Leave implements the ast.Visitor interface.
func (*indexTypeNoBlobChecker) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func (v *indexTypeNoBlobChecker) addNewColumn(tableName string, line int, colDef *ast.ColumnDef) []pkData {
	var pkDataList []pkData
	for _, option := range colDef.Options {
		if option.Tp == ast.ColumnOptionUniqKey {
			tp := v.getBlobStr(colDef.Tp)
			if v.isBlob(tp) {
				pkDataList = append(pkDataList, pkData{
					table:      tableName,
					column:     colDef.Name.String(),
					columnType: tp,
					line:       line,
				})
			}
		}
	}
	v.tablesNewColumns.set(tableName, colDef.Name.String(), colDef)
	return pkDataList
}

func (v *indexTypeNoBlobChecker) addIndex(tableName string, line int, indexSpec *ast.IndexPartSpecification) []pkData {
	var pkDataList []pkData
	columnName := indexSpec.Column.Name.String()
	columnType, err := v.getColumnType(tableName, columnName)
	if err != nil {
		return nil
	}
	if v.isBlob(columnType) {
		pkDataList = append(pkDataList, pkData{
			table:      tableName,
			column:     columnName,
			columnType: columnType,
			line:       line,
		})
	}
	return pkDataList
}

func (v *indexTypeNoBlobChecker) changeColumn(tableName, oldColumnName string, line int, newColumnDef *ast.ColumnDef) []pkData {
	var pkDataList []pkData
	v.tablesNewColumns.delete(tableName, oldColumnName)
	for _, option := range newColumnDef.Options {
		if option.Tp == ast.ColumnOptionPrimaryKey || option.Tp == ast.ColumnOptionUniqKey {
			tp := v.getBlobStr(newColumnDef.Tp)
			if v.isBlob(tp) {
				pkDataList = append(pkDataList, pkData{
					table:      tableName,
					column:     newColumnDef.Name.String(),
					columnType: tp,
					line:       line,
				})
			}
		}
	}
	v.tablesNewColumns.set(tableName, newColumnDef.Name.String(), newColumnDef)
	return pkDataList
}

func (v *indexTypeNoBlobChecker) addConstraint(tableName string, line int, constraint *ast.Constraint) []pkData {
	var pkDataList []pkData
	if constraint.Tp == ast.ConstraintPrimaryKey || constraint.Tp == ast.ConstraintUniqKey || constraint.Tp == ast.ConstraintKey ||
		constraint.Tp == ast.ConstraintIndex || constraint.Tp == ast.ConstraintUniqIndex || constraint.Tp == ast.ConstraintUniq {
		for _, key := range constraint.Keys {
			columnName := key.Column.Name.String()
			columnType, err := v.getColumnType(tableName, columnName)
			if err != nil {
				continue
			}
			if v.isBlob(columnType) {
				pkDataList = append(pkDataList, pkData{
					table:      tableName,
					column:     columnName,
					columnType: columnType,
					line:       line,
				})
			}
		}
	}
	return pkDataList
}

// getPKColumnType gets the column type string from v.tablesNewColumns or catalog, returns empty string and non-nil error if cannot find the column in given table.
func (v *indexTypeNoBlobChecker) getColumnType(tableName string, columnName string) (string, error) {
	if colDef, ok := v.tablesNewColumns.get(tableName, columnName); ok {
		return v.getBlobStr(colDef.Tp), nil
	}
	column := v.database.FindColumn(&catalog.ColumnFind{
		TableName:  tableName,
		ColumnName: columnName,
	})
	if column != nil {
		return column.Type, nil
	}
	return "", errors.Errorf("cannot find the type of `%s`.`%s`", tableName, columnName)
}

// getIntOrBigIntStr returns the type string of tp.
func (*indexTypeNoBlobChecker) getBlobStr(tp *types.FieldType) string {
	switch tp.GetType() {
	case mysql.TypeTinyBlob:
		return "tinyblob"
	case mysql.TypeBlob:
		return "blob"
	case mysql.TypeMediumBlob:
		return "mediumblob"
	case mysql.TypeLongBlob:
		return "longblob"
	}
	return tp.String()
}

func (*indexTypeNoBlobChecker) isBlob(tp string) bool {
	up := strings.ToUpper(tp)
	return up == "TINYBLOB" || up == "BLOB" || up == "MEDIUMBLOB" || up == "LONGBLOB"
}

package catalog

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	mysql "github.com/bytebase/mysql-parser"

	mysqlparser "github.com/bytebase/bytebase/backend/plugin/parser/mysql"
)

func (d *DatabaseState) mysqlV2WalkThrough(stmt string) error {
	// We define the Catalog as Database -> Schema -> Table. The Schema is only for PostgreSQL.
	// So we use a Schema whose name is empty for other engines, such as MySQL.
	// If there is no empty-string-name schema, create it to avoid corner cases.
	if _, exists := d.schemaSet[""]; !exists {
		d.createSchema("")
	}

	nodeList, err := mysqlparser.ParseMySQL(stmt + ";")
	if err != nil {
		return NewParseError(err.Error())
	}
	for _, node := range nodeList {
		if err := d.mysqlV2ChangeState(node); err != nil {
			return err
		}
	}

	return nil
}

type mysqlV2Listener struct {
	*mysql.BaseMySQLParserListener

	baseLine      int
	text          string
	databaseState *DatabaseState
	err           *WalkThroughError
}

func (l *mysqlV2Listener) EnterQuery(ctx *mysql.QueryContext) {
	l.text = ctx.GetParser().GetTokenStream().GetTextFromRuleContext(ctx)
}

func (d *DatabaseState) mysqlV2ChangeState(in *mysqlparser.ParseResult) (err *WalkThroughError) {
	defer func() {
		if err == nil {
			return
		}
		if err.Line == 0 {
			err.Line = in.BaseLine
		}
	}()

	if d.deleted {
		return &WalkThroughError{
			Type:    ErrorTypeDatabaseIsDeleted,
			Content: fmt.Sprintf("Database `%s` is deleted", d.name),
		}
	}

	listener := &mysqlV2Listener{
		baseLine:      in.BaseLine,
		databaseState: d,
	}
	antlr.ParseTreeWalkerDefault.Walk(listener, in.Tree)
	if listener.err != nil {
		return listener.err
	}
	return nil
}

// EnterCreateTable is called when production createTable is entered.
func (l *mysqlV2Listener) EnterCreateTable(ctx *mysql.CreateTableContext) {
	if ctx.TableName() == nil {
		return
	}
	databaseName, tableName := mysqlparser.NormalizeMySQLTableName(ctx.TableName())
	if databaseName != "" && !l.databaseState.isCurrentDatabase(databaseName) {
		l.err = &WalkThroughError{
			Type:    ErrorTypeAccessOtherDatabase,
			Content: fmt.Sprintf("Database `%s` is not the current database `%s`", databaseName, l.databaseState.name),
		}
		return
	}

	schema, exists := l.databaseState.schemaSet[""]
	if !exists {
		schema = l.databaseState.createSchema("")
	}
	if _, exists = schema.getTable(tableName); exists {
		if ctx.IfNotExists() != nil {
			return
		}
		l.err = &WalkThroughError{
			Type:    ErrorTypeTableExists,
			Content: fmt.Sprintf("Table `%s` already exists", tableName),
		}
		return
	}

	if ctx.DuplicateAsQueryExpression() != nil {
		l.err = &WalkThroughError{
			Type:    ErrorTypeUseCreateTableAs,
			Content: fmt.Sprintf("Disallow the CREATE TABLE AS statement but \"%s\" uses", l.text),
		}
		return
	}

	if ctx.LIKE_SYMBOL() != nil {
		_, referTable := mysqlparser.NormalizeMySQLTableRef(ctx.TableRef())
		l.err = l.databaseState.mysqlV2CopyTable(databaseName, tableName, referTable)
		return
	}

	table := &TableState{
		name:      tableName,
		engine:    newEmptyStringPointer(),
		collation: newEmptyStringPointer(),
		comment:   newEmptyStringPointer(),
		columnSet: make(columnStateMap),
		indexSet:  make(IndexStateMap),
	}
	schema.tableSet[table.name] = table

	if ctx.TableElementList() == nil {
		return
	}

	hasAutoIncrement := false
	for _, tableElement := range ctx.TableElementList().AllTableElement() {
		switch {
		// handle column
		case tableElement.ColumnDefinition() != nil:
			if tableElement.ColumnDefinition().FieldDefinition() == nil {
				continue
			}
			if mysqlparser.IsAutoIncrement(tableElement.ColumnDefinition().FieldDefinition()) {
				if hasAutoIncrement {
					l.err = &WalkThroughError{
						Type: ErrorTypeAutoIncrementExists,
						// The content comes from MySQL error content.
						Content: fmt.Sprintf("There can be only one auto column for table `%s`", table.name),
					}
				}
				hasAutoIncrement = true
			}
			if err := table.mysqlV2CreateColumn(l.databaseState.ctx, tableElement.ColumnDefinition()); err != nil {
				err.Line = l.baseLine + tableElement.GetStart().GetLine()
				l.err = err
				return
			}
		case tableElement.TableConstraintDef() != nil:
			if err := table.mysqlV2CreateConstraint(l.databaseState.ctx, tableElement.TableConstraintDef()); err != nil {
				err.Line = tableElement.GetStart().GetLine()
				l.err = err
				return
			}
		}
	}
}

// EnterDropTable is called when production dropTable is entered.
func (l *mysqlV2Listener) EnterDropTable(ctx *mysql.DropTableContext) {
	if ctx.TableRefList() == nil {
		return
	}

	for _, tableRef := range ctx.TableRefList().AllTableRef() {
		databaseName, tableName := mysqlparser.NormalizeMySQLTableRef(tableRef)
		if databaseName != "" && !l.databaseState.isCurrentDatabase(databaseName) {
			l.err = &WalkThroughError{
				Type:    ErrorTypeAccessOtherDatabase,
				Content: fmt.Sprintf("Database `%s` is not the current database `%s`", databaseName, tableName),
			}
		}

		schema, exists := l.databaseState.schemaSet[""]
		if !exists {
			schema = l.databaseState.createSchema("")
		}

		table, exists := schema.getTable(tableName)
		if !exists {
			if ctx.IfExists() != nil || !l.databaseState.ctx.CheckIntegrity {
				return
			}
			l.err = &WalkThroughError{
				Type:    ErrorTypeTableNotExists,
				Content: fmt.Sprintf("Table `%s` does not exist", tableName),
			}
			return
		}

		delete(schema.tableSet, table.name)
	}
}

func (d *DatabaseState) mysqlV2CopyTable(databaseName, tableName, referTable string) *WalkThroughError {
	targetTable, err := d.mysqlV2FindTableState(databaseName, referTable, true /* createIncompleteTable */)
	if err != nil {
		return err
	}

	schema := d.schemaSet[""]
	table := targetTable.copy()
	table.name = tableName
	schema.tableSet[table.name] = table
	return nil
}

func (d *DatabaseState) mysqlV2FindTableState(databaseName, tableName string, createIncompleteTable bool) (*TableState, *WalkThroughError) {
	if databaseName != "" && !d.isCurrentDatabase(databaseName) {
		return nil, NewAccessOtherDatabaseError(d.name, databaseName)
	}

	schema, exists := d.schemaSet[""]
	if !exists {
		schema = d.createSchema("")
	}

	table, exists := schema.getTable(tableName)
	if !exists {
		if schema.ctx.CheckIntegrity {
			return nil, NewTableNotExistsError(tableName)
		}
		if createIncompleteTable {
			table = schema.createIncompleteTable(tableName)
		} else {
			return nil, nil
		}
	}

	return table, nil
}

func (t *TableState) mysqlV2CreateConstraint(ctx *FinderContext, constraintDef mysql.ITableConstraintDefContext) *WalkThroughError {
	if constraintDef.GetType_() != nil {
		switch constraintDef.GetType_().GetTokenType() {
		// PRIMARY KEY.
		case mysql.MySQLParserPRIMARY_SYMBOL:
			if constraintDef.KeyListVariants() == nil {
				// never reach here.
				return nil
			}
			keyList := mysqlparser.NormalizeKeyListVariants(constraintDef.KeyListVariants())
			if err := t.mysqlV2ValidateKeyStringList(ctx, keyList, true /* primary */, false /* isSpatial*/); err != nil {
				return err
			}
			if err := t.mysqlV2CreatePrimaryKey(keyList, mysqlV2GetIndexType(constraintDef)); err != nil {
				return err
			}
		// normal KEY/INDEX.
		case mysql.MySQLParserKEY_SYMBOL, mysql.MySQLParserINDEX_SYMBOL:
			if constraintDef.KeyListVariants() == nil {
				// never reach here.
				return nil
			}
			keyList := mysqlparser.NormalizeKeyListVariants(constraintDef.KeyListVariants())
			if err := t.mysqlV2ValidateKeyStringList(ctx, keyList, false /* primary */, false /* isSpatial */); err != nil {
				return err
			}

			indexName := ""
			if constraintDef.IndexNameAndType() != nil && constraintDef.IndexNameAndType().IndexName() != nil {
				indexName = mysqlparser.NormalizeIndexName(constraintDef.IndexNameAndType().IndexName())
			}
			if err := t.mysqlV2CreateIndex(indexName, keyList, false /* unique */, mysqlV2GetIndexType(constraintDef), constraintDef); err != nil {
				return err
			}
		// UNIQUE KEY.
		case mysql.MySQLParserUNIQUE_SYMBOL:
			if constraintDef.KeyListVariants() == nil {
				// never reach here.
				return nil
			}
			keyList := mysqlparser.NormalizeKeyListVariants(constraintDef.KeyListVariants())
			if err := t.mysqlV2ValidateKeyStringList(ctx, keyList, false /* primary */, false /* isSpatial*/); err != nil {
				return err
			}

			indexName := ""
			if constraintDef.ConstraintName() != nil {
				indexName = mysqlparser.NormalizeConstraintName(constraintDef.ConstraintName())
			}
			if constraintDef.IndexNameAndType() != nil && constraintDef.IndexNameAndType().IndexName() != nil {
				indexName = mysqlparser.NormalizeIndexName(constraintDef.IndexNameAndType().IndexName())
			}
			if err := t.mysqlV2CreateIndex(indexName, keyList, true /* unique */, mysqlV2GetIndexType(constraintDef), constraintDef); err != nil {
				return err
			}
		// FULLTEXT KEY.
		case mysql.MySQLParserFULLTEXT_SYMBOL:
			if constraintDef.KeyListVariants() == nil {
				// never reach here.
				return nil
			}
			keyList := mysqlparser.NormalizeKeyListVariants(constraintDef.KeyListVariants())
			if err := t.mysqlV2ValidateKeyStringList(ctx, keyList, false /* primary */, false /* isSpatial*/); err != nil {
				return err
			}
			indexName := ""
			if constraintDef.IndexName() != nil {
				indexName = mysqlparser.NormalizeIndexName(constraintDef.IndexName())
			}
			if err := t.mysqlV2CreateIndex(indexName, keyList, false /* unique */, mysqlV2GetIndexType(constraintDef), constraintDef); err != nil {
				return err
			}
		case mysql.MySQLParserFOREIGN_SYMBOL:
			// we do not deal with FOREIGN KEY constraints.
		}
	}

	// we do not deal with check constraints.
	// if constraintDef.CheckConstraint() != nil {}
	return nil
}

func (t *TableState) mysqlV2ValidateKeyStringList(ctx *FinderContext, keyList []string, primary bool, isSpatial bool) *WalkThroughError {
	for _, columnName := range keyList {
		column, exists := t.columnSet[columnName]
		if !exists {
			if ctx.CheckIntegrity {
				return NewColumnNotExistsError(t.name, columnName)
			}
		} else {
			if primary {
				column.nullable = newFalsePointer()
			}
			if isSpatial && column.nullable != nil && *column.nullable {
				return &WalkThroughError{
					Type: ErrorTypeSpatialIndexKeyNullable,
					// The error content comes from MySQL.
					Content: fmt.Sprintf("All parts of a SPATIAL index must be NOT NULL, but `%s` is nullable", column.name),
				}
			}
		}
	}
	return nil
}

func mysqlV2GetIndexType(tableConstraint mysql.ITableConstraintDefContext) string {
	if tableConstraint.GetType_() == nil {
		return "BTREE"
	}

	// I still need to handle IndexNameAndType to get index type(algorithm).
	switch tableConstraint.GetType_().GetTokenType() {
	case mysql.MySQLParserPRIMARY_SYMBOL,
		mysql.MySQLParserKEY_SYMBOL,
		mysql.MySQLParserINDEX_SYMBOL,
		mysql.MySQLParserUNIQUE_SYMBOL:

		if tableConstraint.IndexNameAndType() != nil {
			if tableConstraint.IndexNameAndType().IndexType() != nil {
				indexType := tableConstraint.IndexNameAndType().IndexType().GetText()
				return strings.ToUpper(indexType)
			}
		}

		for _, option := range tableConstraint.AllIndexOption() {
			if option == nil || option.IndexTypeClause() == nil {
				continue
			}

			indexType := option.IndexTypeClause().IndexType().GetText()
			return strings.ToUpper(indexType)
		}
	case mysql.MySQLParserFULLTEXT_SYMBOL:
		return "FULLTEXT"
	case mysql.MySQLParserFOREIGN_SYMBOL:
	}
	// for mysql, we use BTREE as default index type.
	return "BTREE"
}

func (t *TableState) mysqlV2CreateColumn(_ *FinderContext, columnDef mysql.IColumnDefinitionContext) *WalkThroughError {
	if columnDef.ColumnName() == nil || columnDef.FieldDefinition() == nil {
		// todo: add more error info
		return nil
	}
	_, _, columnName := mysqlparser.NormalizeMySQLColumnName(columnDef.ColumnName())
	if _, exists := t.columnSet[columnName]; exists {
		return &WalkThroughError{
			Type:    ErrorTypeColumnExists,
			Content: fmt.Sprintf("Column `%s` already exists in table `%s`", columnName, t.name),
		}
	}

	// todo: handle position.
	pos := len(t.columnSet) + 1
	columnType := ""
	characterSet := ""
	collation := ""
	if columnDef.FieldDefinition() == nil || columnDef.FieldDefinition().DataType() == nil {
		// todo: add more error detail.
		return nil
	}
	columnType = mysqlparser.NormalizeMySQLDataType(columnDef.FieldDefinition().DataType(), true /* compact */)
	characterSet = mysqlparser.GetCharSetName(columnDef.FieldDefinition().DataType())
	collation = mysqlparser.GetCollationName(columnDef.FieldDefinition())

	col := &ColumnState{
		name:         columnName,
		position:     &pos,
		defaultValue: nil,
		nullable:     newTruePointer(),
		columnType:   newStringPointer(columnType),
		characterSet: newStringPointer(characterSet),
		collation:    newStringPointer(collation),
		comment:      newEmptyStringPointer(),
	}
	setNullDefault := false

	for _, attribute := range columnDef.FieldDefinition().AllColumnAttribute() {
		if attribute == nil {
			continue
		}
		if attribute.CheckConstraint() != nil {
			// we do not deal with CHECK constraint.
			continue
		}
		// not null.
		if attribute.NullLiteral() != nil && attribute.NOT_SYMBOL() != nil {
			col.nullable = newFalsePointer()
		}
		if attribute.GetValue() != nil {
			switch attribute.GetValue().GetTokenType() {
			// default value.
			case mysql.MySQLParserDEFAULT_SYMBOL:
				if err := mysqlV2CheckDefault(columnName, columnDef.FieldDefinition()); err != nil {
					return err
				}
				if attribute.SignedLiteral() == nil {
					continue
				}
				// handle default null.
				if attribute.SignedLiteral().Literal() != nil && attribute.SignedLiteral().Literal().NullLiteral() != nil {
					setNullDefault = true
					continue
				}
				// handle default 'null' etc.
				defaultValue := mysqlparser.NormalizeMySQLSignedLiteral(attribute.SignedLiteral())
				col.defaultValue = &defaultValue
			// comment.
			case mysql.MySQLParserCOMMENT_SYMBOL:
				if attribute.TextLiteral() == nil {
					continue
				}
				comment := mysqlparser.NormalizeMySQLTextLiteral(attribute.TextLiteral())
				col.comment = &comment
			// on update now().
			case mysql.MySQLParserON_SYMBOL:
				if attribute.UPDATE_SYMBOL() == nil || attribute.NOW_SYMBOL() == nil {
					continue
				}
				if !mysqlparser.IsTimeType(columnDef.FieldDefinition().DataType()) {
					return &WalkThroughError{
						Type:    ErrorTypeOnUpdateColumnNotDatetimeOrTimestamp,
						Content: fmt.Sprintf("Column `%s` use ON UPDATE but is not DATETIME or TIMESTAMP", col.name),
					}
				}
			// primary key.
			case mysql.MySQLParserKEY_SYMBOL:
				// the key attribute for in a column meaning primary key.
				col.nullable = newFalsePointer()
				// we need to check the key type which generated by tidb parser.
				if err := t.mysqlV2CreatePrimaryKey([]string{col.name}, "BTREE"); err != nil {
					return err
				}
			// unique key.
			case mysql.MySQLParserUNIQUE_SYMBOL:
				// unique index.
				if err := t.mysqlV2CreateIndex("", []string{col.name}, true /* unique */, "BTREE", mysql.NewEmptyTableConstraintDefContext()); err != nil {
					return err
				}
			// auto_increment.
			case mysql.MySQLParserAUTO_INCREMENT_SYMBOL:
				// we do not deal with AUTO_INCREMENT.
			// column_format.
			case mysql.MySQLParserCOLUMN_FORMAT_SYMBOL:
				// we do not deal with COLUMN_FORMAT.
			// storage.
			case mysql.MySQLParserSTORAGE_SYMBOL:
				// we do not deal with STORAGE.
			}
		}
	}

	if col.nullable != nil && !*col.nullable && setNullDefault {
		return &WalkThroughError{
			Type: ErrorTypeSetNullDefaultForNotNullColumn,
			// Content comes from MySQL Error content.
			Content: fmt.Sprintf("Invalid default value for column `%s`", col.name),
		}
	}

	t.columnSet[col.name] = col
	return nil
}

func (t *TableState) mysqlV2CreateIndex(name string, keyList []string, unique bool, tp string, tableConstraint mysql.ITableConstraintDefContext) *WalkThroughError {
	if len(keyList) == 0 {
		return &WalkThroughError{
			Type:    ErrorTypeIndexEmptyKeys,
			Content: fmt.Sprintf("Index `%s` in table `%s` has empty key", name, t.name),
		}
	}
	// construct a index name if name is empty.
	if name != "" {
		if _, exists := t.indexSet[name]; exists {
			return NewIndexExistsError(t.name, name)
		}
	} else {
		suffix := 1
		for {
			name = keyList[0]
			if suffix > 1 {
				name = fmt.Sprintf("%s_%d", keyList[0], suffix)
			}
			if _, exists := t.indexSet[name]; !exists {
				break
			}
			suffix++
		}
	}

	index := &IndexState{
		name:           name,
		expressionList: keyList,
		indexType:      &tp,
		unique:         &unique,
		primary:        newFalsePointer(),
		visible:        newTruePointer(),
		comment:        newEmptyStringPointer(),
	}

	// need to check the visibility of index.
	// we need a for-loop to determined the visibility of index.

	// NORMAL KEY/INDEX.
	// PRIMARY KEY.
	// UNIQUE KEY.
	for _, attribute := range tableConstraint.AllIndexOption() {
		if attribute == nil || attribute.CommonIndexOption() == nil {
			continue
		}
		if attribute.CommonIndexOption().Visibility() != nil && attribute.CommonIndexOption().Visibility().INVISIBLE_SYMBOL() != nil {
			index.visible = newFalsePointer()
		}
	}

	// FULLTEXT INDEX.
	for _, attribute := range tableConstraint.AllFulltextIndexOption() {
		if attribute == nil || attribute.CommonIndexOption() == nil {
			continue
		}
		if attribute.CommonIndexOption().Visibility() != nil && attribute.CommonIndexOption().Visibility().INVISIBLE_SYMBOL() != nil {
			index.visible = newFalsePointer()
		}
	}

	// SPATIAL INDEX.
	for _, attribute := range tableConstraint.AllSpatialIndexOption() {
		if attribute == nil || attribute.CommonIndexOption() == nil {
			continue
		}
		if attribute.CommonIndexOption().Visibility() != nil && attribute.CommonIndexOption().Visibility().INVISIBLE_SYMBOL() != nil {
			index.visible = newFalsePointer()
		}
	}

	t.indexSet[name] = index
	return nil
}

func (t *TableState) mysqlV2CreatePrimaryKey(keys []string, tp string) *WalkThroughError {
	if _, exists := t.indexSet[PrimaryKeyName]; exists {
		return &WalkThroughError{
			Type:    ErrorTypePrimaryKeyExists,
			Content: fmt.Sprintf("Primary key exists in table `%s`", t.name),
		}
	}

	pk := &IndexState{
		name:           PrimaryKeyName,
		expressionList: keys,
		indexType:      &tp,
		unique:         newTruePointer(),
		primary:        newTruePointer(),
		visible:        newTruePointer(),
		comment:        newEmptyStringPointer(),
	}
	t.indexSet[pk.name] = pk
	return nil
}

func mysqlV2CheckDefault(columnName string, fieldDefinition mysql.IFieldDefinitionContext) *WalkThroughError {
	if fieldDefinition.DataType() == nil || fieldDefinition.DataType().GetType_() == nil {
		return nil
	}

	switch fieldDefinition.DataType().GetType_().GetTokenType() {
	case mysql.MySQLParserTEXT_SYMBOL,
		mysql.MySQLParserTINYTEXT_SYMBOL,
		mysql.MySQLParserMEDIUMTEXT_SYMBOL,
		mysql.MySQLParserLONGTEXT_SYMBOL,
		mysql.MySQLParserBLOB_SYMBOL,
		mysql.MySQLParserTINYBLOB_SYMBOL,
		mysql.MySQLParserLONGBLOB_SYMBOL,
		mysql.MySQLParserJSON_SYMBOL,
		mysql.MySQLParserGEOMETRY_SYMBOL,
		mysql.MySQLParserGEOMETRYCOLLECTION_SYMBOL,
		mysql.MySQLParserPOINT_SYMBOL,
		mysql.MySQLParserMULTIPOINT_SYMBOL,
		mysql.MySQLParserLINESTRING_SYMBOL,
		mysql.MySQLParserMULTILINESTRING_SYMBOL,
		mysql.MySQLParserPOLYGON_SYMBOL,
		mysql.MySQLParserMULTIPOLYGON_SYMBOL:
		return &WalkThroughError{
			Type: ErrorTypeInvalidColumnTypeForDefaultValue,
			// Content comes from MySQL Error content.
			Content: fmt.Sprintf("BLOB, TEXT, GEOMETRY or JSON column `%s` can't have a default value", columnName),
		}
	}

	// todo: check column type with default value.
	return nil
}

package ast

// DropTableStmt is the struct for drop table or view statement.
type DropTableStmt struct {
	node

	TableList []*TableDef
}

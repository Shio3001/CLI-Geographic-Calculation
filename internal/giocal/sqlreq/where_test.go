package sqlreq

import (
	"testing"
)

// isWhereExprInで使うテスト
func TestParseWhereExprIn1(t *testing.T)	{
	query := "SELECT * FROM rail WHERE company IN ('東日本旅客鉄道', '東海旅客鉄道')"
	parsed := ParseSQLQuery(query)
	firstStmt := GetFirstStmt(parsed)
	whereClause := GetWhereClause(firstStmt)
	t.Log(whereClause)

	// whereClauseからA_Exprを取り出す
	aExpr := whereClause.GetAExpr()
	if aExpr == nil {
		t.Fatal("whereClause is not A_Expr")
	}

	// lexprとrexprを取り出す
	lexpr := aExpr.GetLexpr()
	rexpr := aExpr.GetRexpr()

	column, values, ok := ParseWhereExprIn(lexpr, rexpr)
	if !ok {
		t.Fatal("ParseWhereExprIn failed")
	}

	/**
	where_test.go:30: Column: company
    where_test.go:31: Values: [東日本旅客鉄道 東海旅客鉄道]
	*/
	t.Logf("Column: %s", column)
	t.Logf("Values: %v", values)
}

// 単一の時
func TestParseWhereExprIn2(t *testing.T)	{	
	query := "SELECT * FROM rail WHERE company IN ('東日本旅客鉄道')"
	parsed := ParseSQLQuery(query)
	firstStmt := GetFirstStmt(parsed)
	whereClause := GetWhereClause(firstStmt)
	t.Log(whereClause)

	// whereClauseからA_Exprを取り出す
	aExpr := whereClause.GetAExpr()
	if aExpr == nil {
		t.Fatal("whereClause is not A_Expr")
	}

	// lexprとrexprを取り出す
	lexpr := aExpr.GetLexpr()
	rexpr := aExpr.GetRexpr()

	column, values, ok := ParseWhereExprIn(lexpr, rexpr)
	if !ok {
		t.Fatal("ParseWhereExprIn failed")
	}

	/**
	where_test.go:62: Column: company
	where_test.go:63: Values: [東日本旅客鉄道]
	*/
	t.Logf("Column: %s", column)
	t.Logf("Values: %v", values)
}	
package sqlreq

import (
	"CLI-Geographic-Calculation/internal/giocal"
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"
	"CLI-Geographic-Calculation/internal/giocal/sqlreq/testutil"
	"testing"
)

// isWhereExprInで使うテスト
func TestParseWhereExprIn1(t *testing.T) {
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
func TestParseWhereExprIn2(t *testing.T) {
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

// ParseWhereClauseのテスト
func TestParseWhereClause(t *testing.T) {
	query := "SELECT * FROM rail WHERE company = '東日本旅客鉄道' AND line IN ('山手線', '中央線')"
	parsed := ParseSQLQuery(query)
	firstStmt := GetFirstStmt(parsed)
	whereClause := GetWhereClause(firstStmt)
	t.Log(whereClause)

	// ダミーのDatasetResourceを作成

	drp := giocaltype.DatasetResourcePath{
		Rail: testutil.ProjectRootPath(
			"internal/giodata_public/N02-23_RailroadSection.json",
		),
		Station: testutil.ProjectRootPath(
			"internal/giodata_public/N02-23_Station.json",
		),
	}
	//LoadDatasetResourceで作成する
	dummyDRS, error := giocal.LoadDatasetResource(drp)
	if error != nil {
		t.Fatalf("Failed to load DatasetResource: %v", error)
	}

	resultIndices := ParseWhereClause(dummyDRS, whereClause, []int{})

	// 結果のインデックスを表示
	t.Logf("Result Indices: %v", resultIndices)
}

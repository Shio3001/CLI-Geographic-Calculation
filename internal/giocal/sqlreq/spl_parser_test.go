package sqlreq

import "testing"

func TestParseSQLQuery1(t *testing.T) {
	query := "SELECT * FROM rail WHERE year = 2023"
	result := ParseSQLQuery(query)
	//表示するだけ
	t.Log(result)
}

// 会社名でフィルタリングする場合
func TestParseSQLQuery2(t *testing.T) {
	query := "SELECT * FROM rail WHERE company = '東日本旅客鉄道'"
	result := ParseSQLQuery(query)
	//表示するだけ
	t.Log(result)
}

// 会社数が複数ある場合
func TestParseSQLQuery3(t *testing.T) {
	query := "SELECT * FROM rail WHERE company IN ('東日本旅客鉄道', '東海旅客鉄道')"
	result := ParseSQLQuery(query)
	//表示するだけ
	t.Log(result)
}

// 会社と路線名でフィルタリングする場合
func TestParseSQLQuery4(t *testing.T) {
	query := "SELECT * FROM rail WHERE company = '東日本旅客鉄道' AND line IN ('山手線', '中央線')"
	result := ParseSQLQuery(query)
	//表示するだけ
	t.Log(result)
}


func TestGetFirstStmt(t *testing.T) {
	query := "SELECT * FROM rail WHERE company = '東日本旅客鉄道' AND line IN ('山手線', '中央線')"
	parsed := ParseSQLQuery(query)
	firstStmt := GetFirstStmt(parsed)
	t.Log(firstStmt)
}

func TestGetFromClause(t *testing.T) {
	query := "SELECT * FROM rail WHERE company = '東日本旅客鉄道' AND line IN ('山手線', '中央線')"
	parsed := ParseSQLQuery(query)
	firstStmt := GetFirstStmt(parsed)
	fromClause := GetFromClause(firstStmt)
	t.Log(fromClause)
}

func TestGetWhereClause(t *testing.T) {
	query := "SELECT * FROM rail WHERE company = '東日本旅客鉄道' AND line IN ('山手線', '中央線')"
	parsed := ParseSQLQuery(query)
	firstStmt := GetFirstStmt(parsed)
	whereClause := GetWhereClause(firstStmt)
	t.Log(whereClause)
}
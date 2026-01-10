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
	query := "SELECT * FROM station WHERE company = '東日本旅客鉄道'"
	result := ParseSQLQuery(query)
	//表示するだけ
	t.Log(result)
}

// 会社数が複数ある場合
func TestParseSQLQuery3(t *testing.T) {
	query := "SELECT * FROM station WHERE company IN ('東日本旅客鉄道', '東海旅客鉄道')"
	result := ParseSQLQuery(query)
	//表示するだけ
	t.Log(result)
}

package sqlreq

import (
	"CLI-Geographic-Calculation/pkg/giocal"
	"CLI-Geographic-Calculation/pkg/giocal/giocaltype"
	"CLI-Geographic-Calculation/pkg/giocal/linefilter"
	"CLI-Geographic-Calculation/pkg/giocal/sqlreq/testutil"
	"testing"
)

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

// whereが奥行きのある場合
/***
  spl_parser_test.go:65: bool_expr:{boolop:OR_EXPR args:{bool_expr:{boolop:AND_EXPR args:{a_expr:{kind:AEXPR_OP name:{string:{sval:"="}} lexpr:{column_ref:{fields:{string:{sval:"company"}} location:26}} rexpr:{a_const:{sval:{sval:"東日本旅客鉄道"} location:36}} location:34}} args:{a_expr:{kind:AEXPR_IN name:{string:{sval:"="}} lexpr:{column_ref:{fields:{string:{sval:"line"}} location:64}} rexpr:{list:{items:{a_const:{sval:{sval:"山手線"} location:73}} items:{a_const:{sval:{sval:"中央線"} location:86}}}} location:69}} location:60}} args:{a_expr:{kind:AEXPR_OP name:{string:{sval:"="}} lexpr:{column_ref:{fields:{string:{sval:"year"}} location:103}} rexpr:{a_const:{ival:{ival:2023} location:110}} location:108}} location:100}
*/
/**
bool_expr  (boolop=OR_EXPR)
├─ args[0]: bool_expr (boolop=AND_EXPR)
│  ├─ args[0]: a_expr (kind=AEXPR_OP, op="=")
│  │  ├─ lexpr: column_ref
│  │  │  └─ fields[0]: string "company"  (location=26)
│  │  └─ rexpr: a_const (string)
│  │     └─ sval "東日本旅客鉄道"  (location=36)
│  └─ args[1]: a_expr (kind=AEXPR_IN)
│     ├─ lexpr: column_ref
│     │  └─ fields[0]: string "line"  (location=64)
│     └─ rexpr: list  (location=69)
│        ├─ items[0]: a_const (string) sval "山手線"  (location=73)
│        └─ items[1]: a_const (string) sval "中央線"  (location=86)
└─ args[1]: a_expr (kind=AEXPR_OP, op="=")  (location=100)
   ├─ lexpr: column_ref
   │  └─ fields[0]: string "year"  (location=103)
   └─ rexpr: a_const (int)
      └─ ival 2023  (location=110)

*/
func TestGetWhereClauseDeep(t *testing.T) {
	query := "SELECT * FROM rail WHERE (company = '東日本旅客鉄道' AND line IN ('山手線', '中央線')) OR year = 2023"
	parsed := ParseSQLQuery(query)
	firstStmt := GetFirstStmt(parsed)
	whereClause := GetWhereClause(firstStmt)
	t.Log(whereClause)
}

// groupBy句の解析
/**
  spl_parser_test.go:66: select_stmt:{target_list:{res_target:{val:{column_ref:{fields:{string:{sval:"line"}}  location:7}}  location:7}}  target_list:{res_target:{val:{func_call:{funcname:{string:{sval:"count"}}  agg_star:true  funcformat:COERCE_EXPLICIT_CALL  location:13}}  location:13}}  from_clause:{range_var:{relname:"rail"  inh:true  relpersistence:"p"  location:27}}  where_clause:{a_expr:{kind:AEXPR_OP  name:{string:{sval:"="}}  lexpr:{column_ref:{fields:{string:{sval:"company"}}  location:38}}  rexpr:{a_const:{sval:{sval:"東日本旅客鉄道"}  location:48}}  location:46}}  group_clause:{column_ref:{fields:{string:{sval:"line"}}  location:81}}  limit_option:LIMIT_OPTION_DEFAULT  op:SETOP_NONE}
*/
func TestGetGroupByClause(t *testing.T) {
	query := "SELECT line, COUNT(*) FROM rail WHERE company = '東日本旅客鉄道' GROUP BY line"
	parsed := ParseSQLQuery(query)
	firstStmt := GetFirstStmt(parsed)
	groupByClauses := GetGroupByClauses(firstStmt)
	for _, clause := range groupByClauses {
		t.Log(clause)
	}
}

// groupByが複数
func TestGetGroupByClauseMultiple(t *testing.T) {
	query := "SELECT company, line, COUNT(*) FROM rail WHERE year = 2023 GROUP BY company, line"
	parsed := ParseSQLQuery(query)

	firstStmt := GetFirstStmt(parsed)
	groupByClauses := GetGroupByClauses(firstStmt)
	for _, clause := range groupByClauses {
		t.Log(clause)
	}
}

/**
version: 170004
stmts:
  stmt:
    select_stmt:
      target_list:
        - res_target:
            location: 7
            val:
              column_ref:
                location: 7
                fields:
                  - string: { sval: "company" }

        - res_target:
            location: 16
            val:
              column_ref:
                location: 16
                fields:
                  - string: { sval: "line" }

        - res_target:
            location: 22
            val:
              func_call:
                location: 22
                funcname:
                  - string: { sval: "count" }
                args:
                  - column_ref:
                      location: 28
                      fields:
                        - string: { sval: "node" }
                funcformat: COERCE_EXPLICIT_CALL

      from_clause:
        - range_var:
            relname: "rail"
            inh: true
            relpersistence: "p"
            location: 39

      where_clause:
        a_expr:
          kind: AEXPR_OP
          location: 55
          name:
            - string: { sval: "=" }
          lexpr:
            column_ref:
              location: 50
              fields:
                - string: { sval: "year" }
          rexpr:
            a_const:
              location: 57
              ival: { ival: 2023 }

      group_clause:
        - column_ref:
            location: 71
            fields:
              - string: { sval: "company" }
        - column_ref:
            location: 80
            fields:
              - string: { sval: "line" }

      limit_option: LIMIT_OPTION_DEFAULT
      op: SETOP_NONE
*/
// groupByが複数
func TestGetGroupByClauseMultipleCount(t *testing.T) {
	query := "SELECT company, line, COUNT(node) FROM rail WHERE year = 2023 GROUP BY company, line"
	parsed := ParseSQLQuery(query)
	t.Log(parsed)
	firstStmt := GetFirstStmt(parsed)
	groupByClauses := GetGroupByClauses(firstStmt)
	for _, clause := range groupByClauses {
		t.Log(clause)
	}
}

// SQLToGraphのテスト
func TestSQLToGraph(t *testing.T) {
	query := "SELECT * FROM rail WHERE company = '東日本旅客鉄道' AND line IN ('山手線', '中央線')"
	parsed := ParseSQLQuery(query)
	drp := giocaltype.DatasetResourcePath{
		Rail: testutil.ProjectRootPath(
			"pkg/giodata_public/N02-23_RailroadSection.json",
		),
		Station: testutil.ProjectRootPath(
			"pkg/giodata_public/N02-23_Station.json",
		),
	}

	drs, _ := giocal.LoadDatasetResource(drp)
	graph := SQLToGraph(linefilter.FilterRailroadSectionByProperties, parsed, drs)
	t.Log("Graph: ")
	testutil.OutputGraph(graph, "test_output_graph_east.json")

}

// SQLToGraphのテスト
func TestSQLToGraphTokai(t *testing.T) {
	query := "SELECT * FROM rail WHERE company IN ('東日本旅客鉄道' , '東海旅客鉄道') AND line IN ('山手線', '中央線')"
	parsed := ParseSQLQuery(query)
	drp := giocaltype.DatasetResourcePath{
		Rail: testutil.ProjectRootPath(
			"pkg/giodata_public/N02-23_RailroadSection.json",
		),
		Station: testutil.ProjectRootPath(
			"pkg/giodata_public/N02-23_Station.json",
		),
	}

	drs, _ := giocal.LoadDatasetResource(drp)
	graph := SQLToGraph(linefilter.FilterRailroadSectionByProperties, parsed, drs)
	t.Log("Graph: ")
	testutil.OutputGraph(graph, "test_output_graph_tokai.json")
}

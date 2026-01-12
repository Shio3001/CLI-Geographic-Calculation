package sqlreq

import (
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"

	pg_query "github.com/pganalyze/pg_query_go/v6"
)

/*
*

	spl_parser_test.go:57: bool_expr:{boolop:AND_EXPR args:{a_expr:{kind:AEXPR_OP name:{string:{sval:"="}} lexpr:{column_ref:{fields:{string:{sval:"company"}} location:25}} rexpr:{a_const:{sval:{sval:"東日本旅客鉄道"} location:35}} location:33}} args:{a_expr:{kind:AEXPR_IN name:{string:{sval:"="}} lexpr:{column_ref:{fields:{string:{sval:"line"}} location:63}} rexpr:{list:{items:{a_const:{sval:{sval:"山手線"} location:72}} items:{a_const:{sval:{sval:"中央線"} location:85}}}} location:68}} location:59}
*/
func ParseWhereClause(drs *giocaltype.DatasetResource, whereClause *pg_query.Node) {
	// AND条件
	if whereClause.GetBoolExpr() != nil && whereClause.GetBoolExpr().GetBoolop() == pg_query.BoolExprType_AND_EXPR {
		args := whereClause.GetBoolExpr().GetArgs()
		if len(args) != 2 {
			return
		}

		left := args[0]
		right := args[1]

		// 左辺: company = '東日本旅客鉄道'
		if isAexpr(left) {
		}

		// 右辺: line IN ('山手線', '中央線')
		if isAexpr(right) {
		}
	}
}

/**
  ├─ where_clause: BoolExpr
  │  ├─ boolop: AND_EXPR
  │  ├─ args[0]: A_Expr
  │  │  ├─ kind: AEXPR_OP           (通常の二項演算)
  │  │  ├─ op: "="
  │  │  ├─ left : ColumnRef("company")  (location: 25)
  │  │  ├─ right: ConstString("東日本旅客鉄道") (location: 35)
  │  │  └─ location: 33
  │  ├─ args[1]: A_Expr
  │  │  ├─ kind: AEXPR_IN           (IN 条件)
  │  │  ├─ left : ColumnRef("line") (location: 63)
  │  │  ├─ right: List
  │  │  │  ├─ ConstString("山手線") (location: 72)
  │  │  └  └─ ConstString("中央線") (location: 85)
  │  │  └─ location: 68
  │  └─ location: 59
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

func isAexpr(node *pg_query.Node) bool {
	if node.GetAExpr() != nil {
		return true
	}
	return false
}

// tree
type ParserNode func(drs *giocaltype.DatasetResource, required1 []int, required2 []int) []int

// 必要な行数を返す
func WhereAnd(drs *giocaltype.DatasetResource, required1 []int, required2 []int) []int {

	// required1とrequired2の共通部分を返す
	result := make([]int, 0)
	m := make(map[int]bool)

	for _, v := range required1 {
		m[v] = true
	}

	for _, v := range required2 {
		if m[v] {
			result = append(result, v)
		}
	}

	return result
}

// 必要な行数を返す
func WhereOr(drs *giocaltype.DatasetResource, required1 []int, required2 []int) []int {
	// required1とrequired2の和集合を返す
	result := make([]int, 0)
	m := make(map[int]bool)

	for _, v := range required1 {
		m[v] = true
	}

	for _, v := range required2 {
		m[v] = true
	}

	for k := range m {
		result = append(result, k)
	}

	return result
}

/**
  │  │  ├─ left : ColumnRef("line") (location: 63)
  │  │  ├─ right: List
  │  │  │  ├─ ConstString("山手線") (location: 72)
  │  │  └  └─ ConstString("中央線") (location: 85)
*/

func ParseWhereExprIn(
	left *pg_query.Node,
	right *pg_query.Node,
) (string, []string, bool) {

	if left == nil || right == nil {
		return "", nil, false
	}

	// 左辺: ColumnRef
	colRef := left.GetColumnRef()
	if colRef == nil {
		return "", nil, false
	}

	// column 名取得（単一カラム想定）
	if len(colRef.Fields) != 1 {
		return "", nil, false
	}
	column := colRef.Fields[0].GetString_().GetSval()

	// 右辺: IN (...) の List
	list := right.GetList()
	if list == nil {
		return "", nil, false
	}

	values := make([]string, 0, len(list.Items))

	for _, item := range list.Items {
		c := item.GetAConst()
		if c == nil {
			continue
		}

		if sval := c.GetSval(); sval != nil {
			values = append(values, sval.Sval)
		}
	}

	if len(values) == 0 {
		return "", nil, false
	}

	return column, values, true
}

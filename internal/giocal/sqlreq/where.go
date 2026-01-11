package sqlreq

import pg_query "github.com/pganalyze/pg_query_go/v6"

/**
  spl_parser_test.go:57: bool_expr:{boolop:AND_EXPR args:{a_expr:{kind:AEXPR_OP name:{string:{sval:"="}} lexpr:{column_ref:{fields:{string:{sval:"company"}} location:25}} rexpr:{a_const:{sval:{sval:"東日本旅客鉄道"} location:35}} location:33}} args:{a_expr:{kind:AEXPR_IN name:{string:{sval:"="}} lexpr:{column_ref:{fields:{string:{sval:"line"}} location:63}} rexpr:{list:{items:{a_const:{sval:{sval:"山手線"} location:72}} items:{a_const:{sval:{sval:"中央線"} location:85}}}} location:68}} location:59}
*/
func ParseWhereClause(whereClause  *pg_query.Node)   {
}

func WhereAnd(left *pg_query.Node, right *pg_query.Node) *[]int {
  return &[]int{}
}

func WhereOr(left *pg_query.Node, right *pg_query.Node) *[]int {
  return &[]int{}
}

// カラムと値の比較
func WhereEqual(column string, value *[]string) *[]int {
  return &[]int{}
}

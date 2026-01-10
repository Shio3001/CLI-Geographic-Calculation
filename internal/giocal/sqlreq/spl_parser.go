package sqlreq

//github.com/pganalyze/pg_query_go/v6を使う
import (
	pg_query "github.com/pganalyze/pg_query_go/v6"
)

func ParseSQLQuery(query string)  {
	parsed, err := pg_query.ParseToJSON(query)
	if err != nil {
		panic(err)
	}
	println(parsed)
}
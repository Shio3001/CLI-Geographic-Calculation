//get requestを受信SQL文を受け取る
//SQL Likeな文でgeojsonのデータを整理できるようにする
//url routing は

package native

import (
	"io"
	"net/http"
)


func habdleSQLRequest(w http.ResponseWriter, r *http.Request) {
	sqlQuery := r.URL.Query().Get("query")
	if sqlQuery == "" {
		http.Error(w, "Missing 'query' parameter", http.StatusBadRequest)
		return
	}

	// SQL文を解析し、対応するgeojsonデータを取得・処理するロジックをここに実装
	// 例: resultGeoJSON := executeSQLQuery(sqlQuery)

	// 仮のレスポンスとして空のGeoJSONを返す
	resultGeoJSON := `{"type": "FeatureCollection", "features": []}`

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, resultGeoJSON)
}


func main() {
	http.HandleFunc("/query", habdleSQLRequest)
	http.ListenAndServe(":8080", nil)
}
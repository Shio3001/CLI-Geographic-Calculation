//get requestを受信SQL文を受け取る
//SQL Likeな文でgeojsonのデータを整理できるようにする
//url routing は

package main

import (
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"
	"CLI-Geographic-Calculation/internal/giocal/sqlreq"
	"net/http"
	"strconv"
	"strings"
)

// 連想配列でルーティングを定義するが、下記のようなパスを想定
// 2023/rail/クエリパラメータ
// 2024/station/クエリパラメータ
// クエリパラメータはSQ文でgeojsonのデータを整理できるようにする
// 例: /2023/rail?query=SELECT * FROM rails WHERE length > 1000
// 例: /2024/station?query=SELECT * FROM stations WHERE city = 'Tokyo'
// 解析するSQL Likeな文は通常のSQL文に近い形で実装する

// ルーティング用のキー
// routeKeyの年は、データセット内における年次データを用いるとき、どの年度を使うかを指定するためのもの
type routeKey struct {
	year    int
	Resource string
}

type Dataset struct {
	Handler  datasetHandler
	Resources giocaltype.DatasetResourcePath
}


var datasets = map[routeKey]Dataset{
	{year: 2023, Resource: "rail"}: {
		Handler: handleRail,
		Resources: giocaltype.DatasetResourcePath{
			Rail: "internal/giodata_public/N02-23_RailroadSection.json",
			Station: "internal/giodata_public/N02-23Station.json",
		},
	},

}

func main(){
	println("[HTTP] GIOCAL")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	println("[HANDLER] : ", r.URL.Path)
	// "/rail/2023/" → ["rail", "2023"]
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.Error(w, "path must be /{year}/{resource}", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}

	key := routeKey{
		Resource: parts[0],
		year:    year,
	}

	// パラメータでSQLクエリを受け取る
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	// パラメータ（Options）で駅から駅までのルート探索なども指定できるようにする
	routeSearch := r.URL.Query().Get("rs")
	if routeSearch != "" {
		// ルート探索のパラメータ処理（省略）
		println("[ROUTE SEARCH] : ", routeSearch)
	}

	// ルーティングマップからハンドラを取得
	handler, ok := datasets[key]
	if !ok {
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	sqlreq.ParseSQLQuery(query)
	
	// ハンドラを呼び出す（ここでは単純にレスポンスを書き込む例）
	_ = handler // 実際のハンドラ呼び出しは省略

	w.Write([]byte(
		"year=" + strconv.Itoa(year) + ", resource=" + key.Resource,
	))
}

type datasetHandler func(datasetResource giocaltype.DatasetResourcePath , year int)

func handleRail(datasetResource giocaltype.DatasetResourcePath, year int) {
	println("[HANDLE RAIL] Year:", year, "Rail Resource:", datasetResource.Rail)
}
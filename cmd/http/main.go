//get requestを受信SQL文を受け取る
//SQL Likeな文でgeojsonのデータを整理できるようにする
//url routing は

package main

import (
	"net/http"
	"strconv"
	"strings"
)

// 連想配列でルーティングを定義するが、下記のようなパスを想定
// 2023/rail/クエリパラメータ
// 2024/station/クエリパラメータ
// クエリパラメータはSQL Likeな文でgeojsonのデータを整理できるようにする
// 例: /2023/rail?select=company,line,station_from,station_to&where=company='東日本旅客鉄道' and line='山手線'
// 例: /2024/station?select=station_name,passengers_2022&where=passengers_2022>1000000

type routeKey struct {
	Year     int
	Resource string
}

datasets := map[routeKey]datasetHandler{
	{Year: 2023, Resource: "rail"}:    "dataset_2023_rail.geojson",
	{Year: 2024, Resource: "station"}: "dataset_2024_station.geojson",
}

func main(){
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// "/2023/rail" → ["2023", "rail"]
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.Error(w, "path must be /{year}/{resource}", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "invalid year", http.StatusBadRequest)
		return
	}

	key := routeKey{
		Year:     year,
		Resource: parts[1],
	}

	// ルーティングマップからハンドラを取得
	handler, ok := datasets[key]
	if !ok {
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}
	
	// ハンドラを呼び出す（ここでは単純にレスポンスを書き込む例）
	_ = handler // 実際のハンドラ呼び出しは省略

	w.Write([]byte(
		"year=" + strconv.Itoa(key.Year) + ", resource=" + key.Resource,
	))
}

type datasetHandler func hello(filename string) {
	// ここにデータセットを処理するロジックを実装
}
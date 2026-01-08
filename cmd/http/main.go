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
// クエリパラメータはSQ文でgeojsonのデータを整理できるようにする
// 例: /2023/rail?query=SELECT * FROM rails WHERE length > 1000
// 例: /2024/station?query=SELECT * FROM stations WHERE city = 'Tokyo'


type routeKey struct {
	Year     int
	Resource string
}

type Dataset struct {
	Handler  datasetHandler
	Filename string
}

var datasets = map[routeKey]Dataset{
	{Year: 2023, Resource: "rail"}:    { handleRail, "rails_2023.geojson" },
	{Year: 2024, Resource: "station"}: { handleStation, "stations_2024.geojson" },
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

type datasetHandler func(filename string)

func handleRail(filename string) {
	
}
func handleStation(filename string) {

}
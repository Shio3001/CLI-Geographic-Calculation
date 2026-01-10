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
// 解析するSQL Likeな文は通常のSQL文に近い形で実装する

// ルーティング用のキー
// routeKeyの年は、データセット内における年次データを用いるとき、どの年度を使うかを指定するためのもの
type routeKey struct {
	Resource string
}

type Dataset struct {
	Handler  datasetHandler
	Resources DatasetResource
}

type DatasetResource struct {
	rail string
	station string
	history string
	passengers string 
}


var datasets = map[routeKey]Dataset{
	{Resource: "rail"}: {
		Handler: handleRail,
		Resources: DatasetResource{
			rail: "internal/giodata/N02-23_RailroadSection.json",
			station: "internal/giodata/N02-23Station.json",
			history: "internal/giodata/N05-24_RailroadHistory.json",
			passengers: "internal/giodata/S12-24_Passengers.json",
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
	}

	// ルーティングマップからハンドラを取得
	handler, ok := datasets[key]
	if !ok {
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	sqlreq.parseSQLQuery("SELECT * FROM " + key.Resource + " WHERE year = " + strconv.Itoa(year));
	
	// ハンドラを呼び出す（ここでは単純にレスポンスを書き込む例）
	_ = handler // 実際のハンドラ呼び出しは省略

	w.Write([]byte(
		"year=" + strconv.Itoa(year) + ", resource=" + key.Resource,
	))
}

type datasetHandler func(datasetResource DatasetResource , year int)

func handleRail(datasetResource DatasetResource, year int) {
	println("[HANDLE RAIL] Year:", year, "Rail Resource:", datasetResource.rail)
}
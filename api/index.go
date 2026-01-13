package api

import (
	"CLI-Geographic-Calculation/internal/dataResolve"
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"
	"CLI-Geographic-Calculation/internal/giocal/sqlreq"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type routeKey struct {
	year     int
	Resource string
}

type Dataset struct {
	Handler   datasetHandler
	Resources giocaltype.DatasetResourcePath
}

var datasets = map[routeKey]Dataset{
	{year: 2023, Resource: "rail"}: {
		Handler: handleRail,
		// DevResources: giocaltype.DatasetResourcePath{
		// 	Rail:    "internal/giodata_public/N02-23_RailroadSection.json",
		// 	Station: "internal/giodata_public/N02-23Station.json",
		// },
		// ProdResources: giocaltype.DatasetResourcePath{
		// 	Rail:    "https://github.com/Shio3001/giojson/blob/main/N02-23_RailroadSection.json",
		// 	Station: "https://github.com/Shio3001/giojson/blob/main/N02-23_Station.json",
		// },

		Resources: giocaltype.DatasetResourcePath{
			Rail:    "N02-23_RailroadSection.json",
			Station: "N02-23_Station.json",
		},
	},
}

func Handler(w http.ResponseWriter, r *http.Request) {
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
		year:     year,
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
	ds, ok := datasets[key]
	if !ok {
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	resolved, err := resolveResources(ds.Resources)
	if err != nil {
		http.Error(w, "failed to resolve dataset resources: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sqlreq.ParseSQLQuery(query)

	// ハンドラを呼び出す（ここでは単純にレスポンスを書き込む例）
	ds.Handler(resolved, year)

	w.Write([]byte(
		"year=" + strconv.Itoa(year) + ", resource=" + key.Resource,
	))
}

func resolveResources(r giocaltype.DatasetResourcePath) (giocaltype.DatasetResourcePath, error) {
	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "" {
		env = "dev"
	}

	if env == "prod" {
		cacheDir := filepath.Join(os.TempDir(), "gio-cache")

		p := dataResolve.BlobURLProvider{
			CacheDir: cacheDir,
			Client:   &http.Client{Timeout: 20 * time.Second},
		}
		return p.Resolve(r)
	}

	// dev: ローカル baseDir を付けるだけ
	base := os.Getenv("GIO_LOCAL_BASE")
	if base == "" {
		base = "internal/giodata_public"
	}
	return giocaltype.DatasetResourcePath{
		Rail:    filepath.Join(base, r.Rail),
		Station: filepath.Join(base, r.Station),
	}, nil
}

type datasetHandler func(datasetResource giocaltype.DatasetResourcePath, year int)

func handleRail(datasetResource giocaltype.DatasetResourcePath, year int) {
	println("[HANDLE RAIL] Year:", year, "Rail Resource:", datasetResource.Rail)

	// ここで、datasetResourceを使ってデータを処理するロジックを実装

}

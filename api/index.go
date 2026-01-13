package main

import (
	"CLI-Geographic-Calculation/pkg/dataResolve"
	"CLI-Geographic-Calculation/pkg/giocal/giocaltype"
	"CLI-Geographic-Calculation/pkg/giocal/sqlreq"
	"encoding/json"
	"errors"
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
		// 	Rail:    "pkg/giodata_public/N02-23_RailroadSection.json",
		// 	Station: "pkg/giodata_public/N02-23Station.json",
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

func parseYearResource(path string) (int, string, error) {
	p := strings.Trim(path, "/")
	parts := strings.Split(p, "/")

	// 許可パターン:
	// - /api/{year}/{resource}
	// - /{year}/{resource}
	switch len(parts) {
	case 2:
		year, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, "", err
		}
		return year, parts[1], nil
	case 3:
		if parts[0] != "api" {
			return 0, "", errBadPath
		}
		year, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, "", err
		}
		return year, parts[2], nil
	default:
		return 0, "", errBadPath
	}
}

var errBadPath = errors.New("path must be /api/{year}/{resource} or /{year}/{resource}")

func Handler(w http.ResponseWriter, r *http.Request) {
	println("[HANDLER] : ", r.URL.Path)
	// "/rail/2023/" → ["rail", "2023"]
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.Error(w, "path must be /{year}/{resource}", http.StatusBadRequest)
		return
	}

	year, resource, err := parseYearResource(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := routeKey{
		Resource: resource,
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
	ds.Handler(w, year, resolved, nil, query)

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
		base = "pkg/giodata_public"
	}
	return giocaltype.DatasetResourcePath{
		Rail:    filepath.Join(base, r.Rail),
		Station: filepath.Join(base, r.Station),
	}, nil
}

type datasetHandler func(w http.ResponseWriter, year int, res giocaltype.DatasetResourcePath, parsed any, rawSQL string)

func handleRail(w http.ResponseWriter, year int, res giocaltype.DatasetResourcePath, parsed any, rawSQL string) {
	println("[HANDLE RAIL] Year:", year, "Rail Resource:", res.Rail)

	out := map[string]any{
		"ok":       true,
		"year":     year,
		"resource": "rail",
		"sql":      rawSQL,
		"resolved_paths": map[string]string{
			"rail":    res.Rail,
			"station": res.Station,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)

}

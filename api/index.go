package handler

import (
	"CLI-Geographic-Calculation/pkg/dataResolve"
	"CLI-Geographic-Calculation/pkg/giocal"
	"CLI-Geographic-Calculation/pkg/giocal/giocaltype"
	"CLI-Geographic-Calculation/pkg/giocal/linefilter"
	"CLI-Geographic-Calculation/pkg/giocal/sqlreq"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pg_query "github.com/pganalyze/pg_query_go/v6"
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
			Rail:    "https://t4rttlmu64pagxmi.public.blob.vercel-storage.com/N02-23_RailroadSection.json",
			Station: "https://t4rttlmu64pagxmi.public.blob.vercel-storage.com/N02-23_Station.json",
		},
	},
}

func parseYearResource(path string) (int, string, error) {
	p := strings.Trim(path, "/")
	parts := strings.Split(p, "/")
	if len(parts) != 2 {
		return 0, "", errBadPath
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", err
	}
	return year, parts[1], nil
}

var errBadPath = errors.New("path must be {year}/{resource}")

func Handler(w http.ResponseWriter, r *http.Request) {
	println("[HANDLER] urlPath:", r.URL.Path, "pathQuery:", r.URL.Query().Get("path"))

	p := r.URL.Query().Get("path")
	if p == "" {
		p = strings.TrimPrefix(r.URL.Path, "/api/")
	}

	year, resource, err := parseYearResource(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := routeKey{
		Resource: resource,
		year:     year,
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query parameter: query", http.StatusBadRequest)
		return
	}

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

	parsed := sqlreq.ParseSQLQuery(query)

	ds.Handler(w, year, resolved, parsed, query)
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

type datasetHandler func(w http.ResponseWriter, year int, res giocaltype.DatasetResourcePath, parsed *pg_query.ParseResult, rawSQL string)

func handleRail(
	w http.ResponseWriter,
	year int,
	res giocaltype.DatasetResourcePath,
	parsed *pg_query.ParseResult,
	rawSQL string,
) {
	// 1) データセット読み込み
	drs, err := giocal.LoadDatasetResource(res)
	if err != nil {
		http.Error(w, "Failed to load DatasetResource: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2) SQL -> Graph
	graph := sqlreq.SQLToGraph(
		linefilter.FilterRailroadSectionByProperties,
		parsed,
		drs,
	)

	// 3) 返却
	out := map[string]any{
		"ok":       true,
		"year":     year,
		"resource": "rail",
		"sql":      rawSQL,
		"resolved_paths": map[string]string{
			"rail":    res.Rail,
			"station": res.Station,
		},
		"graph": graph,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}

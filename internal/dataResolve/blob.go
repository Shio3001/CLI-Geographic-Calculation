package dataResolve

import (
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"
	"errors"
	"net/http"
	"os"
)

type BlobURLProvider struct {
	CacheDir string
	Client   *http.Client
}

func (p BlobURLProvider) Resolve(r giocaltype.DatasetResourcePath) (giocaltype.DatasetResourcePath, error) {
	railURL := os.Getenv("BLOB_RAIL_URL")
	stationURL := os.Getenv("BLOB_STATION_URL")
	if railURL == "" || stationURL == "" {
		return giocaltype.DatasetResourcePath{}, errors.New("missing env: BLOB_RAIL_URL / BLOB_STATION_URL")
	}

	crp := CachedRemoteProvider{
		CacheDir: p.CacheDir,
		Client:   p.Client,
	}

	railPath, err := crp.fetchToCache("rail", railURL)
	if err != nil {
		return giocaltype.DatasetResourcePath{}, err
	}
	stationPath, err := crp.fetchToCache("station", stationURL)
	if err != nil {
		return giocaltype.DatasetResourcePath{}, err
	}

	return giocaltype.DatasetResourcePath{Rail: railPath, Station: stationPath}, nil
}

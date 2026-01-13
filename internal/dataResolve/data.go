package dataResolve

import (
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type CachedRemoteProvider struct {
	RemoteBase string // raw base
	CacheDir   string // /tmp/app-cache
	Client     *http.Client
}

func (p CachedRemoteProvider) Resolve(r giocaltype.DatasetResourcePath) (giocaltype.DatasetResourcePath, error) {
	railURL := strings.TrimRight(p.RemoteBase, "/") + "/" + r.Rail
	stationURL := strings.TrimRight(p.RemoteBase, "/") + "/" + r.Station

	railPath, err := p.fetchToCache("rail", railURL)
	if err != nil {
		return giocaltype.DatasetResourcePath{}, err
	}

	stationPath, err := p.fetchToCache("station", stationURL)
	if err != nil {
		return giocaltype.DatasetResourcePath{}, err
	}

	return giocaltype.DatasetResourcePath{Rail: railPath, Station: stationPath}, nil
}

// URL を指定して /tmp にキャッシュしてローカルパスを返す
func (p CachedRemoteProvider) fetchToCache(key, url string) (string, error) {
	if err := os.MkdirAll(p.CacheDir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(p.CacheDir, key+".json")

	// 既にあるなら使う（必要ならETag対応などに拡張）
	if _, err := os.Stat(dst); err == nil {
		return dst, nil
	}

	req, _ := http.NewRequest("GET", url, nil)
	resp, err := p.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("fetch failed: %s (%d)", url, resp.StatusCode)
	}

	f, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}
	return dst, nil
}

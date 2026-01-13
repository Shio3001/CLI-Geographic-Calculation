package testutil

import (
	"path/filepath"
	"runtime"
)

// ProjectRootPath はこのファイル位置を基準に
// プロジェクトルートからの相対パスを解決する
func ProjectRootPath(rel string) string {
	_, filename, _, _ := runtime.Caller(0)
	base := filepath.Dir(filename)

	// testutil → sqlreq → giocal → internal → project root
	root := filepath.Join(base, "..", "..", "..", "..")
	return filepath.Join(root, rel)
}

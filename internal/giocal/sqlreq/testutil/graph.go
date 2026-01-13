package testutil

import (
	"CLI-Geographic-Calculation/internal/giocal/graphstructure"
	"encoding/json"
	"os"
)

func OutputGraph(graph *graphstructure.Graph, name string) {
	// --- JSON化 ---
	b, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		panic(err)
	}
	// ファイルに出力
	err = os.WriteFile(name, b, 0644)
	if err != nil {
		panic(err)
	}
}

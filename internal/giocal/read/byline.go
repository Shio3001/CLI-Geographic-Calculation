package giocal_read

import (
	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal"
	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal/giocaltype"
)

// 条件に当てはまる行番号をint配列で返す

// 指定したカラムの値がtargetValueと等しい行を探す
// internal/giocal/read.goを活用

// それぞれ、対象路線を複数指定して読み込む関数
func ReadGiotypeRailroadSectionForLines(filePath string, targetLines []string) (*giocaltype.GiotypeRailroadSectionFeatureCollection, error) {
	fc, err := giocal.ReadGiotypeRailroadSection(filePath)
	if err != nil {
		return nil, err
	}

	// 対象路線でフィルタリング
	var filteredFeatures []giocaltype.GiotypeRailroadSection
	for _, feature := range fc.Features {
		for _, line := range targetLines {
			if feature.Properties.N02003 == line {
				filteredFeatures = append(filteredFeatures, feature)
				break
			}
		}
	}
	fc.Features = filteredFeatures

	return fc, nil
}

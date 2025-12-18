// giodata内の geojson を読みこみ、変換し、型情報を付与して返すプログラム
//型情報はinternal/giocal/giocaltype/.goパッケージに定義する

package giocal

import (
	"encoding/json"
	"os"

	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal/giocaltype"
)

// RailroadSection
func ReadGiotypeRailroadSection(filePath string) (*giocaltype.GiotypeRailroadSectionFeatureCollection, error) {
	data, err := ReadGeoJSONFile(filePath)
	if err != nil {
		return nil, err
	}

	var fc giocaltype.GiotypeRailroadSectionFeatureCollection
	err = json.Unmarshal(data, &fc)
	if err != nil {
		return nil, err
	}

	return &fc, nil
}

// Station
func ReadGiotypeStation(filePath string) (*giocaltype.GiotypeStationFeatureCollection, error) {
	data, err := ReadGeoJSONFile(filePath)
	if err != nil {
		return nil, err
	}

	var fc giocaltype.GiotypeStationFeatureCollection
	err = json.Unmarshal(data, &fc)
	if err != nil {
		return nil, err
	}

	return &fc, nil
}

// Passengers
func ReadGiotypePassengers(filePath string) (*giocaltype.GiotypePassengers, error) {
	data, err := ReadGeoJSONFile(filePath)
	if err != nil {
		return nil, err
	}

	var fc giocaltype.GiotypePassengers
	err = json.Unmarshal(data, &fc)
	if err != nil {
		return nil, err
	}

	return &fc, nil
}	


//readGeoJSONFile は指定されたファイルパスからGeoJSONデータを読み込むヘルパー関数
func ReadGeoJSONFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// それぞれ、対象路線を複数指定して読み込む関数
func ReadGiotypeRailroadSectionForLines(filePath string, targetLines []string) (*giocaltype.GiotypeRailroadSectionFeatureCollection, error) {
	fc, err := ReadGiotypeRailroadSection(filePath)
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

func ReadGiotypeStationForLines(filePath string, targetLines []string) (*giocaltype.GiotypeStationFeatureCollection, error) {
	fc, err := ReadGiotypeStation(filePath)
	if err != nil {
		return nil, err
	}

	// 対象路線でフィルタリング
	var filteredFeatures []giocaltype.GiotypeStation
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
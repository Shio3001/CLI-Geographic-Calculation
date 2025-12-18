// giodata内の geojson を読みこみ、変換し、型情報を付与して返すプログラム
//型情報はinternal/giocal/giocaltype/.goパッケージに定義する

package giocal

import (
	"encoding/json"
	"os"

	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal/giocaltype"
)

// RailroadSection
func readGiotypeRailroadSection(filePath string) (*giocaltype.GiotypeRailroadSectionFeatureCollection, error) {
	data, err := readGeoJSONFile(filePath)
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
func readGiotypeStation(filePath string) (*giocaltype.GiotypeStationFeatureCollection, error) {
	data, err := readGeoJSONFile(filePath)
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
func readGiotypePassengers(filePath string) (*giocaltype.GiotypePassengers, error) {
	data, err := readGeoJSONFile(filePath)
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
func readGeoJSONFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
// giodata内の geojson を読みこみ、変換し、型情報を付与して返すプログラム
//型情報はinternal/giocal/giocaltype/.goパッケージに定義する

package giocal

import (
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"
	"encoding/json"
	"os"
)


func LoadDatasetResource(path giocaltype.DatasetResourcePath) (*giocaltype.DatasetResource, error) {
	railFC, err := LoadGiotypeRailroadSection(path.Rail)
	if err != nil {
		return nil, err
	}
	stationFC, err := LoadGiotypeStation(path.Station)
	if err != nil {
		return nil,  err
	}
	return &giocaltype.DatasetResource{
		Rail: railFC,
		Station: stationFC,
	}, nil
}

// RailroadSection
func LoadGiotypeRailroadSection(filePath string) (*giocaltype.GiotypeRailroadSectionFeatureCollection, error) {
	data, err := LoadGeoJSONFile(filePath)
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
func LoadGiotypeStation(filePath string) (*giocaltype.GiotypeStationFeatureCollection, error) {
	data, err := LoadGeoJSONFile(filePath)
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
func LoadGiotypePassengers(filePath string) (*giocaltype.GiotypePassengersFeatureCollection, error) {
	data, err := LoadGeoJSONFile(filePath)
	if err != nil {
		return nil, err
	}

	var fc giocaltype.GiotypePassengersFeatureCollection
	err = json.Unmarshal(data, &fc)
	if err != nil {
		return nil, err
	}

	return &fc, nil
}


//loadGeoJSONFile は指定されたファイルパスからGeoJSONデータを読み込むヘルパー関数
func LoadGeoJSONFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func LoadGiotypeStationForLines(filePath string, targetLines []string) (*giocaltype.GiotypeStationFeatureCollection, error) {
	fc, err := LoadGiotypeStation(filePath)
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

// 開業年度

func LoadGiotypeRailHistory(filePath string) (*giocaltype.GiotypeN05RailroadSectionFeatureCollection, error) {
	data, err := LoadGeoJSONFile(filePath)
	if err != nil {
		return nil, err
	}

	var fc giocaltype.GiotypeN05RailroadSectionFeatureCollection
	err = json.Unmarshal(data, &fc)
	if err != nil {
		return nil, err
	}

	return &fc, nil
}


// 


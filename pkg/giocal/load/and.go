package giocal_load

import (
	"CLI-Geographic-Calculation/pkg/giocal"
	"CLI-Geographic-Calculation/pkg/giocal/giocaltype"
)

// それぞれ、鉄道会社と、対象路線を複数指定して読み込む関数(AND条件)
// 対象路線が空の場合は全路線を読み込む
func LoadGiotypeRailroadSectionForCompanyAndLines(filePath string, company string, targetLines []string) (*giocaltype.GiotypeRailroadSectionFeatureCollection, error) {
	fc, err := giocal.LoadGiotypeRailroadSection(filePath)
	if err != nil {
		return nil, err
	}

	// 鉄道会社と対象路線でフィルタリング
	var filteredFeatures []giocaltype.GiotypeRailroadSection
	for _, feature := range fc.Features {
		if feature.Properties.N02004 != company {
			continue
		}
		if len(targetLines) == 0 {
			filteredFeatures = append(filteredFeatures, feature)
			continue
		}
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

func LoadGiotypeStationForCompanyAndLines(filePath string, company string, targetLines []string) (*giocaltype.GiotypeStationFeatureCollection, error) {
	fc, err := giocal.LoadGiotypeStation(filePath)
	if err != nil {
		return nil, err
	}

	// 鉄道会社と対象路線でフィルタリング
	var filteredFeatures []giocaltype.GiotypeStation
	for _, feature := range fc.Features {
		if feature.Properties.N02004 != company {
			continue
		}
		if len(targetLines) == 0 {
			filteredFeatures = append(filteredFeatures, feature)
			continue
		}
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

// 乗降客数データの読み込み関数
func LoadGiotypePassengersForCompanyAndLines(filePath string, company string, targetLines []string) (*giocaltype.GiotypePassengersFeatureCollection, error) {
	fc, err := giocal.LoadGiotypePassengers(filePath)
	if err != nil {
		return nil, err
	}

	// 鉄道会社と対象路線でフィルタリング
	//GiotypePassengersFeatureCollectionで返す必要があるが、

	var filteredFeatures []giocaltype.GiotypePassengersFeature
	for _, feature := range fc.Features {
		if feature.Properties.S12002 != company {
			continue
		}
		if len(targetLines) == 0 {
			filteredFeatures = append(filteredFeatures, feature)
			continue
		}
		for _, line := range targetLines {
			if feature.Properties.S12003 == line {
				filteredFeatures = append(filteredFeatures, feature)
				break
			}
		}
	}
	fc.Features = filteredFeatures

	return fc, nil
}

func LoadGiotypeRailHistoryForCompanyAndLines(filePath string, company string, targetLines []string) (*giocaltype.GiotypeN05RailroadSectionFeatureCollection, error) {
	fc, err := giocal.LoadGiotypeRailHistory(filePath)
	if err != nil {
		return nil, err
	}

	// 鉄道会社と対象路線でフィルタリング
	var filteredFeatures []giocaltype.GiotypeN05RailroadSectionFeature
	for _, feature := range fc.Features {
		if feature.Properties.N05003 != company {
			continue
		}
		if len(targetLines) == 0 {
			filteredFeatures = append(filteredFeatures, feature)
			continue
		}
		for _, line := range targetLines {
			if feature.Properties.N05002 == line {
				filteredFeatures = append(filteredFeatures, feature)
				break
			}
		}
	}
	fc.Features = filteredFeatures

	return fc, nil
}

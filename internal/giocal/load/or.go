// それぞれ、鉄道会社と、対象路線を複数指定して読み込む関数(OR条件)
// 対象路線が空の場合は全路線を読み込む

package giocal_load

import (
	"CLI-Geographic-Calculation/internal/giocal"
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"
)
func LoadGiotypeRailroadSectionForCompanyOrLines(filePath string, company string, targetLines []string) (*giocaltype.GiotypeRailroadSectionFeatureCollection, error) {
	fc, err := giocal.LoadGiotypeRailroadSection(filePath)
	if err != nil {
		return nil, err
	}

	// 鉄道会社または対象路線でフィルタリング
	var filteredFeatures []giocaltype.GiotypeRailroadSection
	for _, feature := range fc.Features {
		if feature.Properties.N02004 == company {
			filteredFeatures = append(filteredFeatures, feature)
			continue
		}
		if len(targetLines) == 0 {
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

func LoadGiotypeStationForCompanyOrLines(filePath string, company string, targetLines []string) (*giocaltype.GiotypeStationFeatureCollection, error)	 {
	fc, err := giocal.LoadGiotypeStation(filePath)
	if err != nil {
		return nil, err
	}

	// 鉄道会社または対象路線でフィルタリング
	var filteredFeatures []giocaltype.GiotypeStation
	for _, feature := range fc.Features {
		if feature.Properties.N02004 == company {
			filteredFeatures = append(filteredFeatures, feature)
			continue
		}
		if len(targetLines) == 0 {
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

func LoadGiotypePassengersForCompanyOrLines(filePath string, company string, targetLines []string) (*giocaltype.GiotypePassengersFeatureCollection, error)	 {
	fc, err := giocal.LoadGiotypePassengers(filePath)
	if err != nil {
		return nil, err
	}

	// 鉄道会社または対象路線でフィルタリング
	var filteredFeatures []giocaltype.GiotypePassengersFeature
	for _, feature := range fc.Features {
		if feature.Properties.S12002 == company {
			filteredFeatures = append(filteredFeatures, feature)
			continue
		}
		if len(targetLines) == 0 {
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
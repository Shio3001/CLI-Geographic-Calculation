// giodata内の geojson を読みこみ、変換し、型情報を付与して返すプログラム
//型情報はinternal/giocal/giocaltype/.goパッケージに定義する

package giocaltype


// RailroadSection
func readGiotypeRailroadSection(filePath string) (*GiotypeRailroadSectionFeatureCollection, error) {
	data, err := readGeoJSONFile(filePath)
	if err != nil {
		return nil, err
	}

	var fc GiotypeRailroadSectionFeatureCollection
	err = json.Unmarshal(data, &fc)
	if err != nil {
		return nil, err
	}

	return &fc, nil
}

// Station
func readGiotypeStation(filePath string) (*GiotypeStationFeatureCollection, error) {
	data, err := readGeoJSONFile(filePath)
	if err != nil {
		return nil, err
	}

	var fc GiotypeStationFeatureCollection
	err = json.Unmarshal(data, &fc)
	if err != nil {
		return nil, err
	}

	return &fc, nil
}

// Passengers
func readGiotypePassengers(filePath string) (*GiotypePassengersFeatureCollection, error) {
	data, err := readGeoJSONFile(filePath)
	if err != nil {
		return nil, err
	}

	var fc GiotypePassengersFeatureCollection
	err = json.Unmarshal(data, &fc)
	if err != nil {
		return nil, err
	}

	return &fc, nil
}	


package linefilter

import "CLI-Geographic-Calculation/pkg/giocal/giocaltype"

func FilterRailroadSectionByProperties(railroadSections *[]giocaltype.GiotypeRailroadSection, property string, value []string) []int {
	// RailroadLinePropertyMap から対応するプロパティコードを取得

	// propertyCodeがvalueに含まれるかどうかをチェックその該当する行番号を返す
	matchingIndices := []int{}
	for i, section := range *railroadSections {

		// 型安全に
		// propValue, exists := section.Properties[propertyCode]できないのでswqitch文で対応
		var propValue string
		switch property {
		// caseはSQL文で使うカラム名に対応させる
		case "company":
			propValue = section.Properties.N02004
		case "line":
			propValue = section.Properties.N02003
		default:
			continue // 未知のプロパティコードの場合はスキップ
		}

		// valueにpropValueが含まれるかチェック
		for _, v := range value {
			if propValue == v {
				matchingIndices = append(matchingIndices, i)
				break
			}
		}
	}
	return matchingIndices
}

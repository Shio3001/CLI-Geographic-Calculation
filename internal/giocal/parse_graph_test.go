package giocal

// テストコード
import (
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"
	"CLI-Geographic-Calculation/internal/giocal/graphstructure"
)

//ConvertGiotypeStationToGraphのテスト用関数
func ConvertGiotypeStationToGraphForTest(stFC *giocaltype.GiotypeStationFeatureCollection, rrFC *giocaltype.GiotypeRailroadSectionFeatureCollection, passengersFC *giocaltype.GiotypePassengersFeatureCollection) *graphstructure.Graph
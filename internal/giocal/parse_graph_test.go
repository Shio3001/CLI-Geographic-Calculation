package giocal

// テストコード
import (
	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal/giocaltype"
	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal/graphstructure"
)

//ConvertGiotypeStationToGraphのテスト用関数
func ConvertGiotypeStationToGraphForTest(stFC *giocaltype.GiotypeStationFeatureCollection, rrFC *giocaltype.GiotypeRailroadSectionFeatureCollection, passengersFC *giocaltype.GiotypePassengersFeatureCollection) *graphstructure.Graph
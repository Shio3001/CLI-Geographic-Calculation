// N02-23_Station.geojson の型情報
package giocaltype
// 駅	「鉄道」の下位クラス。
// 属性情報
// （Station.shp）	属性名
// （かっこ内はshp属性名）	説明	属性の型
// 鉄道区分
// （N02_001）	鉄道路線の種類による区別。	コードリスト「鉄道区分コード」
// 事業者種別
// （N02_002）	鉄道路線の事業者による区別。	コードリスト「事業者種別コード」
// 路線名
// （N02_003）	鉄道路線の名称	文字列型
// 運営会社
// （N02_004）	鉄道路線を運営する会社。	文字列型
// 駅名
// （N02_005）	駅の名称	文字列型
// 駅コード
// （N02_005c）	駅の緯度を降順に並び替えて付加した一意の番号	文字列型
// グループコード
// （N02_005g）	グループコード300m以内の距離にある駅で且つ同じ名称の駅を一つのグループとし、グループの重心に最も近い駅コード	文字列型


// このような形式
// {
//   "type": "FeatureCollection",
//   "name": "N02-23_Station",
//   "features": [
//     {
//       "type": "Feature",
//       "properties": {
//         "N02_001": "11",
//         "N02_002": "2",
//         "N02_003": "指宿枕崎線",
//         "N02_004": "九州旅客鉄道",
//         "N02_005": "二月田",
//         "N02_005c": "010112",
//         "N02_005g": "010112"
//       },
//       "geometry": {
//         "type": "LineString",
//         "coordinates": [
//           [130.63035, 31.25405],
//           [130.62985, 31.25459]
//         ]
//       }
//     },

func GiotypeStation struct {
	Type       string `json:"type"`
	Properties struct {
		N02001 string `json:"N02_001"` // 鉄道区分
		N02002 string `json:"N02_002"` // 事業者種別
		N02003 string `json:"N02_003"` // 路線名
		N02004 string `json:"N02_004"` // 運営会社
		N02005 string `json:"N02_005"` // 駅名
		N02005c string `json:"N02_005c"` // 駅コード
		N02005g string `json:"N02_005g"` // グループコード
	} `json:"properties"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][]float64   `json:"coordinates"`
	} `json:"geometry"`
}

type GiotypeStationFeatureCollection struct {
	Type     string              `json:"type"`			
	Features []GiotypeStation   `json:"features"`
}
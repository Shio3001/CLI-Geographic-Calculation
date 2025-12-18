// N02-23_RailloadSection.geojson の型情報
package giocaltype
// 属性名
// （かっこ内はshp属性名）	説明	属性の型
// 鉄道区分
// （N02_001）	鉄道路線の種類による区別。	コードリスト「鉄道区分コード」
// 事業者種別
// （N02_002）	鉄道路線の事業者による区別。	コードリスト「事業者種別コード」
// 路線名
// （N02_003）	鉄道路線の名称	文字列型
// 運営会社
// （N02_004）	鉄道路線を運営する会社。	文字列型
// 地物情報	地物名	説明
// 鉄道区間	「鉄道」の下位クラス。
// 属性情報	関連役割名	説明	関連先
// 駅	鉄道路線が関連している駅。	地物「駅」
// 地物情報	地物名	説明

// このような形式
// {
// "type": "FeatureCollection",
// "name": "N02-23_RailroadSection",
// "features": [
// { "type": "Feature", "properties": { "N02_001": "23", "N02_002": "5", "N02_003": "沖縄都市モノレール線", "N02_004": "沖縄都市モノレール" }, "geometry": { "type": "LineString", "coordinates": [ [ 127.67948, 26.21454 ], [ 127.6797, 26.21474 ], [ 127.67975, 26.2148 ], [ 127.68217, 26.21728 ], [ 127.68357, 26.21862 ], [ 127.68394, 26.21891 ], [ 127.68419, 26.21905 ] ] } },
// { "type": "Feature", "properties": { "N02_001": "12", "N02_002": "5", "N02_003": "いわて銀河鉄道線", "N02_004": "アイジーアールいわて銀河鉄道" }, "geometry": { "type": "LineString", "coordinates": [ [ 141.29139, 40.3374 ], [ 141.29176, 40.33723 ], [ 141.29243, 40.33692 ], [ 141.29323, 40.33654 ], [ 141.29379, 40.33624 ], [ 141.29411, 40.33608 ], [ 141.2949, 40.33563 ], [ 141.29624, 40.33477 ], [ 141.29813, 40.33354 ], [ 141.29862, 40.33317 ] ] } },
// 	]
// }
type GiotypeRailroadSection struct {
	Type       string `json:"type"`
	Properties struct {
		N02001 string `json:"N02_001"` // 鉄道区分
		N02002 string `json:"N02_002"` // 事業者種別
		N02003 string `json:"N02_003"` // 路線名
		N02004 string `json:"N02_004"` // 運営会社
	} `json:"properties"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][]float64   `json:"coordinates"`
	} `json:"geometry"`
}

type GiotypeRailroadSectionFeatureCollection struct {
	Type     string                      `json:"type"`
	Features []GiotypeRailroadSection   `json:"features"`
}
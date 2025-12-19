// N05-24_RailroadSection2.geojson
package giocaltype

/**

（RailroadSection2.shp）
（Station2.shp）	属性名
（かっこ内はshp属性名）	説明	属性の型
事業者種別
（N05_001）	鉄道路線の事業者を区分するコード	コードリスト「事業者種別コード」
路線名
（N05_002）	鉄道路線の名称	文字列型（CharacterString）
運営会社
（N05_003）	鉄道路線を運営する会社名。	文字列型（CharacterString）
供用開始年
（N05_004）	当該区間の供用が開始された年（西暦年）。
不明な場合は999とする。	時間型（TM_Instant）
設置期間（設置開始）
（N05_005b）
設置期間（設置終了）
（N05_005e）	鉄道路線、駅が設置された年（西暦年）。なお、1950年（昭和25年）以前に設置された場合は1950とする。
鉄道路線、駅が変更・廃止された年の一年前の年（西暦年）。
現存する場合は9999、不明な場合は999とする	時間型（TM_Instant）
関係ID
（N05_006）	路線および駅の属性が変更された場合の同一地物である事を表すグループID（その他の情報欄に詳細を説明）	文字列型（CharacterString）
変遷ID
（N05_007）	同一年次に複数回属性が変更された場合の属性項目を表す識別子	コードリスト「変遷IDコード」
変遷備考
（N05_008）	変遷IDで示した属性の内容を記述する。	文字列型（CharacterString）
備考
（N05_009）	駅位置や路線位置が不明確な場合の図形データの精度に関するコメント。	文字列型（CharacterSt
*/

type GiotypeN05StationFeatureCollection struct {
	Type     string                   `json:"type"`
	Name     string                   `json:"name,omitempty"`
	Crs      *GeoJSONCrs              `json:"crs,omitempty"`
	Features []GiotypeN05StationFeature `json:"features"`
}

type GiotypeN05StationFeature struct {
	Type       string         `json:"type"`
	Properties GiotypeN05Prop `json:"properties"`
	Geometry   GiotypeN05StationGeom `json:"geometry"`
}

// Station: Point が典型
type GiotypeN05StationGeom struct {
	Type        string    `json:"type"`        // "Point"
	Coordinates []float64 `json:"coordinates"` // Point: [lon, lat]
}

// --------------------------------

type GiotypeN05RailroadSectionFeatureCollection struct {
	Type     string                           `json:"type"`
	Name     string                           `json:"name,omitempty"`
	Crs      *GeoJSONCrs                      `json:"crs,omitempty"`
	Features []GiotypeN05RailroadSectionFeature `json:"features"`
}

type GiotypeN05RailroadSectionFeature struct {
	Type       string         `json:"type"`
	Properties GiotypeN05Prop `json:"properties"`
	Geometry   GiotypeN05RailroadSectionGeom `json:"geometry"`
}

type GiotypeN05RailroadSectionGeom struct {
	Type        string          `json:"type"`        // "MultiLineString"
	Coordinates [][]float64   `json:"coordinates"` // MultiLineString: [ [ [lon,lat], ... ] , ... ]
}

type GiotypeN05Prop struct {
  N05001  string  `json:"N05_001"`
  N05002  string  `json:"N05_002"`
  N05003  string  `json:"N05_003"`
  N05004  int     `json:"N05_004,string"`   // ←ここ
  N05005b int     `json:"N05_005b,string"`  // ←実データが文字列なら
  N05005e int     `json:"N05_005e,string"`  // ←実データが文字列なら
  N05006  string  `json:"N05_006"`
  N05007  string  `json:"N05_007"`
  N05008  *string `json:"N05_008"`
  N05009  *string `json:"N05_009"`
}

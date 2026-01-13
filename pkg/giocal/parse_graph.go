package giocal

import (
	"fmt"
	"math"
	"sort"

	"CLI-Geographic-Calculation/pkg/giocal/giocaltype"
	"CLI-Geographic-Calculation/pkg/giocal/graphstructure"
)

func ConvertGiotypeRailwayToGraphByRequired(
	stationFC *giocaltype.GiotypeStationFeatureCollection,
	railroadSectionFC *giocaltype.GiotypeRailroadSectionFeatureCollection,
	stationRequired []int,
	railroadSectionRequired []int,
) *graphstructure.Graph {

	g := &graphstructure.Graph{
		Nodes: map[string]*graphstructure.Node{},
		Edges: []*graphstructure.Edge{},
	}

	coordNodeIDByKey := map[string]string{} // "lon,lat" -> nodeID
	allCoordNodes := make([]coordNodeRef, 0, 4096)

	// 路線(会社+路線名) -> その路線に属する座標ノード一覧
	coordNodesByLineKey := map[string][]coordNodeRef{}

	//rangeRailroadSectionRequired  railroadSectionRequiredの要素数が0の場合、すべての路線区間を対象とする
	rangeRailroadSectionRequired := railroadSectionRequired
	if len(railroadSectionRequired) == 0 {
		rangeRailroadSectionRequired = make([]int, len(railroadSectionFC.Features))
		for i := range railroadSectionFC.Features {
			rangeRailroadSectionRequired[i] = i
		}
	}
	// stationも同様
	rangeStationRequired := stationRequired
	if len(stationRequired) == 0 {
		rangeStationRequired = make([]int, len(stationFC.Features))
		for i := range stationFC.Features {
			rangeStationRequired[i] = i
		}
	}

	for i, index := range rangeRailroadSectionRequired {
		sec := railroadSectionFC.Features[index]
		lineKey := makeLineKey(sec.Properties.N02004, sec.Properties.N02003)

		coords := sec.Geometry.Coordinates // [][]float64 期待: [][lon,lat]
		for j := 0; j < len(coords); j++ {
			lon, lat, ok := pickLonLat2(coords[j])
			if !ok {
				continue
			}

			key := coordKey(lon, lat)
			nodeID, exists := coordNodeIDByKey[key]
			if !exists {
				nodeID = fmt.Sprintf("coord:%s", key)
				coordNodeIDByKey[key] = nodeID

				g.Nodes[nodeID] = &graphstructure.Node{
					ID:   nodeID,
					Kind: "coord",
					Name: "", // 座標ノードは名前なし
					Lon:  lon,
					Lat:  lat,
				}

				ref := coordNodeRef{ID: nodeID, Lon: lon, Lat: lat}
				allCoordNodes = append(allCoordNodes, ref)
				coordNodesByLineKey[lineKey] = append(coordNodesByLineKey[lineKey], ref)
			} else {
				// 既存ノード参照を lineKey に積む（重複は後で軽く除去）
				ref := coordNodeRef{ID: nodeID, Lon: lon, Lat: lat}
				coordNodesByLineKey[lineKey] = append(coordNodesByLineKey[lineKey], ref)
			}

			// 連続点同士をエッジで結ぶ（路線区間エッジ）
			if j > 0 {
				prevLon, prevLat, okPrev := pickLonLat2(coords[j-1])
				if okPrev {
					prevKey := coordKey(prevLon, prevLat)
					fromID := coordNodeIDByKey[prevKey]
					toID := coordNodeIDByKey[key]
					if fromID != "" && toID != "" && fromID != toID {
						// TODO: graphstructure.Edge のフィールド名に合わせて修正

						g.Edges = append(g.Edges, &graphstructure.Edge{
							From:     fromID,
							To:       toID,
							Kind:     "rail",
							WeightKm: distKm(prevLon, prevLat, lon, lat),
							// Optional: 会社/路線などを props に入れるならここ
							// Props: map[string]string{"company": sec.Properties.N02004, "line": sec.Properties.N02003},
							Meta: map[string]string{
								"company": sec.Properties.N02004,
								"line":    sec.Properties.N02003,
								"sec_i":   fmt.Sprintf("%d", i),
							},
						})
					}
				}
			}
		}
	}

	// lineKey ごとの座標参照に重複があるので簡易に uniq 化しておく（距離探索が軽くなる）
	for k, refs := range coordNodesByLineKey {
		coordNodesByLineKey[k] = uniqCoordRefs(refs)
	}

	for _, index := range rangeStationRequired {
		st := stationFC.Features[index]
		stationID := makeStationID(st.Properties.N02005c, st.Properties.N02005g, st.Properties.N02004, st.Properties.N02003, st.Properties.N02005)

		chosenLon, chosenLat, chosenOK := chooseStationRepresentativeLonLat(
			st.Geometry.Coordinates,
			coordNodesByLineKey[makeLineKey(st.Properties.N02004, st.Properties.N02003)],
			allCoordNodes,
		)
		if !chosenOK {
			// 駅に座標が無いケース（通常ないはず）
			continue
		}

		// TODO: graphstructure.Node のフィールド名に合わせて修正
		g.Nodes[stationID] = &graphstructure.Node{
			ID:   stationID,
			Kind: "station",
			Name: st.Properties.N02005,
			Lon:  chosenLon,
			Lat:  chosenLat,
			Meta: map[string]string{
				"station_code": st.Properties.N02005c,
				"group_code":   st.Properties.N02005g,
				"company":      st.Properties.N02004,
				"line":         st.Properties.N02003,
			},
		}

		// 駅ノード → 近い座標ノード を接続（駅が路線上のどこに紐づくか）
		// ※ “路線情報にある座標を選択” は代表座標選びで実現してるが、
		//    ここで座標ノードにもリンクしておくと、後で可視化/探索が楽。
		closestCoordID := findClosestCoordNodeID(chosenLon, chosenLat,
			coordNodesByLineKey[makeLineKey(st.Properties.N02004, st.Properties.N02003)],
			allCoordNodes,
		)
		if closestCoordID != "" {
			g.Edges = append(g.Edges, &graphstructure.Edge{
				From:     stationID,
				To:       closestCoordID,
				Kind:     "station_at",
				WeightKm: distKm(chosenLon, chosenLat, g.Nodes[closestCoordID].Lon, g.Nodes[closestCoordID].Lat),
				Meta: map[string]string{
					"company": st.Properties.N02004,
					"line":    st.Properties.N02003,
				},
			})
		}
	}

	return g
}

// ConvertGiotypeStationToGraph
// 駅や座標をノード、路線区間をエッジとして Graph に変換する。
// 駅が複数座標を持つ場合:
//  1. 同一路線(会社+路線名)の路線座標群に最も近い駅座標を代表にする
//  2. 同一路線が見つからない場合、全路線座標群に最も近い駅座標を代表にする（なければ先頭）
//
// 距離は地球平面仮定：
//
//	dxKm = (lon2-lon1)/LongitudePerKm
//	dyKm = (lat2-lat1)/LatitudePerKm
//	distKm = sqrt(dxKm^2 + dyKm^2)
func ConvertGiotypeStationToGraphPH(
	stationFC *giocaltype.GiotypeStationFeatureCollection,
	railroadSectionFC *giocaltype.GiotypeRailroadSectionFeatureCollection,
	passengersFC *giocaltype.GiotypePassengersFeatureCollection,
	historyFC *giocaltype.GiotypeN05RailroadSectionFeatureCollection,
) *graphstructure.Graph {

	g := &graphstructure.Graph{
		Nodes: map[string]*graphstructure.Node{},
		Edges: []*graphstructure.Edge{},
	}

	coordNodeIDByKey := map[string]string{} // "lon,lat" -> nodeID
	allCoordNodes := make([]coordNodeRef, 0, 4096)

	// 路線(会社+路線名) -> その路線に属する座標ノード一覧
	coordNodesByLineKey := map[string][]coordNodeRef{}

	for i, sec := range railroadSectionFC.Features {
		lineKey := makeLineKey(sec.Properties.N02004, sec.Properties.N02003)

		coords := sec.Geometry.Coordinates // [][]float64 期待: [][lon,lat]
		for j := 0; j < len(coords); j++ {
			lon, lat, ok := pickLonLat2(coords[j])
			if !ok {
				continue
			}

			key := coordKey(lon, lat)
			nodeID, exists := coordNodeIDByKey[key]
			if !exists {
				nodeID = fmt.Sprintf("coord:%s", key)
				coordNodeIDByKey[key] = nodeID

				g.Nodes[nodeID] = &graphstructure.Node{
					ID:   nodeID,
					Kind: "coord",
					Name: "", // 座標ノードは名前なし
					Lon:  lon,
					Lat:  lat,
				}

				ref := coordNodeRef{ID: nodeID, Lon: lon, Lat: lat}
				allCoordNodes = append(allCoordNodes, ref)
				coordNodesByLineKey[lineKey] = append(coordNodesByLineKey[lineKey], ref)
			} else {
				// 既存ノード参照を lineKey に積む（重複は後で軽く除去）
				ref := coordNodeRef{ID: nodeID, Lon: lon, Lat: lat}
				coordNodesByLineKey[lineKey] = append(coordNodesByLineKey[lineKey], ref)
			}

			// 連続点同士をエッジで結ぶ（路線区間エッジ）
			if j > 0 {
				prevLon, prevLat, okPrev := pickLonLat2(coords[j-1])
				if okPrev {
					prevKey := coordKey(prevLon, prevLat)
					fromID := coordNodeIDByKey[prevKey]
					toID := coordNodeIDByKey[key]
					if fromID != "" && toID != "" && fromID != toID {
						// TODO: graphstructure.Edge のフィールド名に合わせて修正

						g.Edges = append(g.Edges, &graphstructure.Edge{
							From:     fromID,
							To:       toID,
							Kind:     "rail",
							WeightKm: distKm(prevLon, prevLat, lon, lat),
							// Optional: 会社/路線などを props に入れるならここ
							// Props: map[string]string{"company": sec.Properties.N02004, "line": sec.Properties.N02003},
							Meta: map[string]string{
								"company":   sec.Properties.N02004,
								"line":      sec.Properties.N02003,
								"sec_i":     fmt.Sprintf("%d", i),
								"open_year": fmt.Sprintf("%d", getOpeningYear(historyFC, sec.Properties.N02004, sec.Properties.N02003)),
							},
						})
					}
				}
			}
		}
	}

	// lineKey ごとの座標参照に重複があるので簡易に uniq 化しておく（距離探索が軽くなる）
	for k, refs := range coordNodesByLineKey {
		coordNodesByLineKey[k] = uniqCoordRefs(refs)
	}

	for _, st := range stationFC.Features {
		stationID := makeStationID(st.Properties.N02005c, st.Properties.N02005g, st.Properties.N02004, st.Properties.N02003, st.Properties.N02005)

		chosenLon, chosenLat, chosenOK := chooseStationRepresentativeLonLat(
			st.Geometry.Coordinates,
			coordNodesByLineKey[makeLineKey(st.Properties.N02004, st.Properties.N02003)],
			allCoordNodes,
		)
		if !chosenOK {
			// 駅に座標が無いケース（通常ないはず）
			continue
		}

		Passengers := getPassengersDataByCode(passengersFC, st.Properties.N02005c)
		if len(Passengers) == 0 {
			Passengers = getPassengersDataByName(passengersFC, st.Properties.N02005)
		}

		// TODO: graphstructure.Node のフィールド名に合わせて修正
		g.Nodes[stationID] = &graphstructure.Node{
			ID:         stationID,
			Kind:       "station",
			Name:       st.Properties.N02005,
			Lon:        chosenLon,
			Lat:        chosenLat,
			Passengers: Passengers,
			Meta: map[string]string{
				"station_code": st.Properties.N02005c,
				"group_code":   st.Properties.N02005g,
				"company":      st.Properties.N02004,
				"line":         st.Properties.N02003,
			},
		}

		// 駅ノード → 近い座標ノード を接続（駅が路線上のどこに紐づくか）
		// ※ “路線情報にある座標を選択” は代表座標選びで実現してるが、
		//    ここで座標ノードにもリンクしておくと、後で可視化/探索が楽。
		closestCoordID := findClosestCoordNodeID(chosenLon, chosenLat,
			coordNodesByLineKey[makeLineKey(st.Properties.N02004, st.Properties.N02003)],
			allCoordNodes,
		)
		if closestCoordID != "" {
			g.Edges = append(g.Edges, &graphstructure.Edge{
				From:     stationID,
				To:       closestCoordID,
				Kind:     "station_at",
				WeightKm: distKm(chosenLon, chosenLat, g.Nodes[closestCoordID].Lon, g.Nodes[closestCoordID].Lat),
				Meta: map[string]string{
					"company": st.Properties.N02004,
					"line":    st.Properties.N02003,
				},
			})
		}
	}

	return g
}

type coordNodeRef struct {
	ID  string
	Lon float64
	Lat float64
}

func makeLineKey(company, line string) string {
	return company + "///" + line
}

// 駅ID: 駅コードが一番安定。無ければ他要素も混ぜる。
func makeStationID(code, group, company, line, name string) string {
	if code != "" {
		return "station:" + code
	}
	// フォールバック
	return fmt.Sprintf("station:%s:%s:%s:%s", sanitize(company), sanitize(line), sanitize(group), sanitize(name))
}

func sanitize(s string) string {
	// ここは軽くでOK（ID用途）
	return s
}

// GeoJSON 座標は基本 [lon,lat] を想定
func pickLonLat2(pair []float64) (lon, lat float64, ok bool) {
	if len(pair) < 2 {
		return 0, 0, false
	}
	return pair[0], pair[1], true
}

// 距離(km): 平面仮定（東京近辺用定数）
func distKm(lon1, lat1, lon2, lat2 float64) float64 {
	dxKm := (lon2 - lon1) / LongitudePerKm
	dyKm := (lat2 - lat1) / LatitudePerKm
	return math.Sqrt(dxKm*dxKm + dyKm*dyKm)
}

// 座標の統合キー（丸め精度は用途に合わせて調整）
func coordKey(lon, lat float64) string {
	// だいたい 0.1m～1m オーダーの丸め（必要なら変えてOK）
	return fmt.Sprintf("%.6f,%.6f", lon, lat)
}

func uniqCoordRefs(refs []coordNodeRef) []coordNodeRef {
	if len(refs) <= 1 {
		return refs
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].ID < refs[j].ID })
	out := refs[:0]
	prev := ""
	for _, r := range refs {
		if r.ID == prev {
			continue
		}
		out = append(out, r)
		prev = r.ID
	}
	return out
}

// 駅の代表座標を選ぶ:
//  1. 同一路線の座標ノード群があるなら、それに最も近い駅座標
//  2. なければ全体座標ノード群に最も近い駅座標
//  3. そもそも座標ノードが無ければ stations の先頭
func chooseStationRepresentativeLonLat(
	stationCoords [][]float64,
	lineCoordNodes []coordNodeRef,
	allCoordNodes []coordNodeRef,
) (lon, lat float64, ok bool) {

	if len(stationCoords) == 0 {
		return 0, 0, false
	}
	if len(stationCoords) == 1 {
		lon, lat, ok = pickLonLat2(stationCoords[0])
		return lon, lat, ok
	}

	// 路線座標があるなら優先
	if len(lineCoordNodes) > 0 {
		return chooseClosestToCoordSet(stationCoords, lineCoordNodes)
	}

	// 全体座標があるなら次点
	if len(allCoordNodes) > 0 {
		return chooseClosestToCoordSet(stationCoords, allCoordNodes)
	}

	// 何も無いなら先頭
	lon, lat, ok = pickLonLat2(stationCoords[0])
	return lon, lat, ok
}

func chooseClosestToCoordSet(
	stationCoords [][]float64,
	target []coordNodeRef,
) (lon, lat float64, ok bool) {

	bestD := math.Inf(1)
	bestLon, bestLat := 0.0, 0.0
	found := false

	for _, sc := range stationCoords {
		sl, sa, ok2 := pickLonLat2(sc)
		if !ok2 {
			continue
		}
		d := minDistToCoordRefs(sl, sa, target)
		if d < bestD {
			bestD = d
			bestLon, bestLat = sl, sa
			found = true
		}
	}
	return bestLon, bestLat, found
}

func minDistToCoordRefs(lon, lat float64, refs []coordNodeRef) float64 {
	best := math.Inf(1)
	for _, r := range refs {
		d := distKm(lon, lat, r.Lon, r.Lat)
		if d < best {
			best = d
		}
	}
	return best
}

// 駅代表座標に一番近い座標ノードIDを返す（同一路線が無ければ全体から）
func findClosestCoordNodeID(
	lon, lat float64,
	line []coordNodeRef,
	all []coordNodeRef,
) string {
	refs := line
	if len(refs) == 0 {
		refs = all
	}
	if len(refs) == 0 {
		return ""
	}
	bestD := math.Inf(1)
	bestID := ""
	for _, r := range refs {
		d := distKm(lon, lat, r.Lon, r.Lat)
		if d < bestD {
			bestD = d
			bestID = r.ID
		}
	}
	return bestID
}

/**
S12001  string  `json:"S12_001"`  // 駅名
S12001c string  `json:"S12_001c"` // 駅コード
S12001g string  `json:"S12_001g"` // グループコード
S12002  string  `json:"S12_002"`  // 運営会社
S12003  string  `json:"S12_003"`  // 路線名
S12004  float64 `json:"S12_004"`  // 鉄道区分
S12005  float64 `json:"S12_005"`  // 事業者種別

S12006 float64 `json:"S12_006"` // 重複コード2011
S12007 float64 `json:"S12_007"` // データ有無コード2011
S12008 *string `json:"S12_008"` // 備考2011
S12009 float64 `json:"S12_009"` // 乗降客数2011（整数相当だがJSONがfloatなので float64）

S12010 float64 `json:"S12_010"` // 重複コード2012
S12011 float64 `json:"S12_011"` // データ有無コード2012
S12012 *string `json:"S12_012"` // 備考2012
S12013 float64 `json:"S12_013"` // 乗降客数2012

S12014 float64 `json:"S12_014"` // 重複コード2013
S12015 float64 `json:"S12_015"` // データ有無コード2013
S12016 *string `json:"S12_016"` // 備考2013
S12017 float64 `json:"S12_017"` // 乗降客数2013

S12018 float64 `json:"S12_018"` // 重複コード2014
S12019 float64 `json:"S12_019"` // データ有無コード2014
S12020 *string `json:"S12_020"` // 備考2014
S12021 float64 `json:"S12_021"` // 乗降客数2014

S12022 float64 `json:"S12_022"` // 重複コード2015
S12023 float64 `json:"S12_023"` // データ有無コード2015
S12024 *string `json:"S12_024"` // 備考2015
S12025 float64 `json:"S12_025"` // 乗降客数2015

S12026 float64 `json:"S12_026"` // 重複コード2016
S12027 float64 `json:"S12_027"` // データ有無コード2016
S12028 *string `json:"S12_028"` // 備考2016
S12029 float64 `json:"S12_029"` // 乗降客数2016

S12030 float64 `json:"S12_030"` // 重複コード2017
S12031 float64 `json:"S12_031"` // データ有無コード2017
S12032 *string `json:"S12_032"` // 備考2017
S12033 float64 `json:"S12_033"` // 乗降客数2017

S12034 float64 `json:"S12_034"` // 重複コード2018
S12035 float64 `json:"S12_035"` // データ有無コード2018
S12036 *string `json:"S12_036"` // 備考2018
S12037 float64 `json:"S12_037"` // 乗降客数2018

S12038 float64 `json:"S12_038"` // 重複コード2019
S12039 float64 `json:"S12_039"` // データ有無コード2019
S12040 *string `json:"S12_040"` // 備考2019
S12041 float64 `json:"S12_041"` // 乗降客数2019

S12042 float64 `json:"S12_042"` // 重複コード2020
S12043 float64 `json:"S12_043"` // データ有無コード2020
S12044 *string `json:"S12_044"` // 備考2020
S12045 float64 `json:"S12_045"` // 乗降客数2020

S12046 float64 `json:"S12_046"` // 重複コード2021
S12047 float64 `json:"S12_047"` // データ有無コード2021
S12048 *string `json:"S12_048"` // 備考2021
S12049 float64 `json:"S12_049"` // 乗降客数2021

S12050 float64 `json:"S12_050"` // 重複コード2022
S12051 float64 `json:"S12_051"` // データ有無コード2022
S12052 *string `json:"S12_052"` // 備考2022
S12053 float64 `json:"S12_053"` // 乗降客数2022

S12054 float64 `json:"S12_054"` // 重複コード2023
S12055 float64 `json:"S12_055"` // データ有無コード2023
S12056 *string `json:"S12_056"` // 備考2023
S12057 float64 `json:"S12_057"` // 乗降客数2023
*/

//	passengersFC *giocaltype.GiotypePassengersFeatureCollectionから連想配列で乗降客数データを取得
//
// 引数で渡す S12001cで駅コードを参照し、GiotypePassengersPropの年度ごとの乗降客数データを数値:数値の連想配列で返す
// フラグも確認する。S12001cが一致する駅コードのデータのみを取得する
func getPassengersDataByCode(
	passengersFC *giocaltype.GiotypePassengersFeatureCollection, stationCode string,
) map[int]float64 {
	passengersData := make(map[int]float64)
	for _, feature := range passengersFC.Features {
		//駅コードが一致しない場合はスキップ
		if feature.Properties.S12001c != stationCode {
			continue
		}
		true_flag := 1
		//フラグを確認し、データ有無コードがtrue_flagの場合のみ乗降客数を取得

		//すべてがint(feature.Properties.S12011) != true_flagで、なおかつ該当数値がすべて0の場合、continueでスキップする
		all_zero := true
		for _, value := range map[int]float64{
			2011: feature.Properties.S12009,
			2012: feature.Properties.S12013,
			2013: feature.Properties.S12017,
			2014: feature.Properties.S12021,
			2015: feature.Properties.S12025,
			2016: feature.Properties.S12029,
			2017: feature.Properties.S12033,
			2018: feature.Properties.S12037,
			2019: feature.Properties.S12041,
			2020: feature.Properties.S12045,
			2021: feature.Properties.S12049,
			2022: feature.Properties.S12053,
			2023: feature.Properties.S12057,
		} {
			if value != 0 {
				all_zero = false
				break
			}
		}
		if all_zero {
			continue
		}

		if int(feature.Properties.S12007) == true_flag {
			passengersData[2011] = feature.Properties.S12009
		}
		if int(feature.Properties.S12011) == true_flag {
			passengersData[2012] = feature.Properties.S12013
		}
		if int(feature.Properties.S12015) == true_flag {
			passengersData[2013] = feature.Properties.S12017
		}
		if int(feature.Properties.S12019) == true_flag {
			passengersData[2014] = feature.Properties.S12021
		}
		if int(feature.Properties.S12023) == true_flag {
			passengersData[2015] = feature.Properties.S12025
		}
		if int(feature.Properties.S12027) == true_flag {
			passengersData[2016] = feature.Properties.S12029
		}
		if int(feature.Properties.S12031) == true_flag {
			passengersData[2017] = feature.Properties.S12033
		}
		if int(feature.Properties.S12035) == true_flag {
			passengersData[2018] = feature.Properties.S12037
		}
		if int(feature.Properties.S12039) == true_flag {
			passengersData[2019] = feature.Properties.S12041
		}
		if int(feature.Properties.S12043) == true_flag {
			passengersData[2020] = feature.Properties.S12045
		}
		if int(feature.Properties.S12047) == true_flag {
			passengersData[2021] = feature.Properties.S12049
		}
		if int(feature.Properties.S12051) == true_flag {
			passengersData[2022] = feature.Properties.S12053
		}
		if int(feature.Properties.S12055) == true_flag {
			passengersData[2023] = feature.Properties.S12057
		}
	}
	return passengersData
}

func getPassengersDataByName(
	passengersFC *giocaltype.GiotypePassengersFeatureCollection, stationName string,
) map[int]float64 {
	passengersData := make(map[int]float64)
	for _, feature := range passengersFC.Features {
		//駅名が一致しない場合はスキップ
		if feature.Properties.S12001 != stationName {
			continue
		}
		true_flag := 1
		//フラグを確認し、データ有無コードがtrue_flagの場合のみ乗降客数を取得
		if int(feature.Properties.S12007) == true_flag {
			passengersData[2011] = feature.Properties.S12009
		}
		if int(feature.Properties.S12011) == true_flag {
			passengersData[2012] = feature.Properties.S12013
		}
		if int(feature.Properties.S12015) == true_flag {
			passengersData[2013] = feature.Properties.S12017
		}
		if int(feature.Properties.S12019) == true_flag {
			passengersData[2014] = feature.Properties.S12021
		}
		if int(feature.Properties.S12023) == true_flag {
			passengersData[2015] = feature.Properties.S12025
		}
		if int(feature.Properties.S12027) == true_flag {
			passengersData[2016] = feature.Properties.S12029
		}
		if int(feature.Properties.S12031) == true_flag {
			passengersData[2017] = feature.Properties.S12033
		}
		if int(feature.Properties.S12035) == true_flag {
			passengersData[2018] = feature.Properties.S12037
		}
		if int(feature.Properties.S12039) == true_flag {
			passengersData[2019] = feature.Properties.S12041
		}
		if int(feature.Properties.S12043) == true_flag {
			passengersData[2020] = feature.Properties.S12045
		}
		if int(feature.Properties.S12047) == true_flag {
			passengersData[2021] = feature.Properties.S12049
		}
		if int(feature.Properties.S12051) == true_flag {
			passengersData[2022] = feature.Properties.S12053
		}
		if int(feature.Properties.S12055) == true_flag {
			passengersData[2023] = feature.Properties.S12057
		}
	}
	return passengersData
}

// そのエッジの開業年度を取得する関数
func getOpeningYear(
	historyFC *giocaltype.GiotypeN05RailroadSectionFeatureCollection, company string, line string,
) int {
	for _, feature := range historyFC.Features {
		if feature.Properties.N05003 != company {
			continue
		}
		if feature.Properties.N05002 != line {
			continue
		}
		return int(feature.Properties.N05004)
	}
	return 0 //見つからなかった場合は0を返す
}

package giocal

import (
	"fmt"
	"math"
	"sort"

	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal/giocaltype"
	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal/graphstructure"
)

// ConvertGiotypeStationToGraph
// 駅や座標をノード、路線区間をエッジとして Graph に変換する。
// 駅が複数座標を持つ場合:
//  1) 同一路線(会社+路線名)の路線座標群に最も近い駅座標を代表にする
//  2) 同一路線が見つからない場合、全路線座標群に最も近い駅座標を代表にする（なければ先頭）
//
// 距離は地球平面仮定：
//   dxKm = (lon2-lon1)/LongitudePerKm
//   dyKm = (lat2-lat1)/LatitudePerKm
//   distKm = sqrt(dxKm^2 + dyKm^2)
func ConvertGiotypeStationToGraph(
	stationFC *giocaltype.GiotypeStationFeatureCollection,
	railroadSectionFC *giocaltype.GiotypeRailroadSectionFeatureCollection,
) *graphstructure.Graph {

	// TODO: graphstructure に合わせて初期化を調整
	g := &graphstructure.Graph{
		Nodes: map[string]*graphstructure.Node{},
		Edges: []*graphstructure.Edge{},
	}

	// =========================
	// 1) 路線区間の座標を「座標ノード」として作る
	//    同一座標は統合（キーは丸め）
	// =========================
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

				// TODO: graphstructure.Node のフィールド名に合わせて修正
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

	// =========================
	// 2) 駅ノードを作る（代表座標を選ぶ）
	// =========================
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

// =========================
// helper
// =========================

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
//  1) 同一路線の座標ノード群があるなら、それに最も近い駅座標
//  2) なければ全体座標ノード群に最も近い駅座標
//  3) そもそも座標ノードが無ければ stations の先頭
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

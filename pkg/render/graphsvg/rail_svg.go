package graphsvg

import (
	"bytes"
	"fmt"
	"html"
	"math"
	"sort"
	"strings"

	"CLI-Geographic-Calculation/pkg/giocal/graphstructure"
)

type Options struct {
	Width        int
	Height       int
	Padding      int
	DrawStations bool
	DrawLabels   bool
}

func (o Options) withDefaults() Options {
	if o.Width <= 0 {
		o.Width = 1200
	}
	if o.Height <= 0 {
		o.Height = 800
	}
	if o.Padding < 0 {
		o.Padding = 0
	}
	return o
}

// RenderRailGraphSVG: rail(graph.Edges Kind="rail") を座標として描画する簡易 SVG
func RenderRailGraphSVG(graph any, opt Options) (string, error) {
	opt = opt.withDefaults()

	g, ok := graph.(*graphstructure.Graph)
	if !ok {
		return "", fmt.Errorf("unexpected graph type: %T (expected *graphstructure.Graph)", graph)
	}

	// 座標ノード収集（rail のエッジに登場する coord ノードを優先）
	coordIDs := collectCoordIDsFromRailEdges(g)

	// fallback: 何も取れなければ Nodes の coord 全部
	if len(coordIDs) == 0 {
		for id, n := range g.Nodes {
			if n.Kind == "coord" {
				coordIDs = append(coordIDs, id)
			}
		}
	}

	if len(coordIDs) == 0 {
		return "", fmt.Errorf("no coord nodes found")
	}

	// bbox
	minLon, minLat := math.Inf(1), math.Inf(1)
	maxLon, maxLat := math.Inf(-1), math.Inf(-1)
	for _, id := range coordIDs {
		n, ok := g.Nodes[id]
		if !ok {
			continue
		}
		minLon = math.Min(minLon, n.Lon)
		maxLon = math.Max(maxLon, n.Lon)
		minLat = math.Min(minLat, n.Lat)
		maxLat = math.Max(maxLat, n.Lat)
	}

	// bbox が潰れている場合の保険
	if !(minLon < maxLon) {
		maxLon = minLon + 1e-6
	}
	if !(minLat < maxLat) {
		maxLat = minLat + 1e-6
	}

	pad := float64(opt.Padding)
	w := float64(opt.Width)
	h := float64(opt.Height)

	// lon/lat -> x/y（lat は上が大きいので反転）
	project := func(lon, lat float64) (x, y float64) {
		x = pad + (lon-minLon)/(maxLon-minLon)*(w-2*pad)
		y = pad + (maxLat-lat)/(maxLat-minLat)*(h-2*pad)
		return
	}

	// rail edge をセグメントとして描画
	railEdges := make([]graphstructure.Edge, 0, len(g.Edges))
	for _, e := range g.Edges {
		if e.Kind == "rail" {
			railEdges = append(railEdges, *e)
		}
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`+"\n",
		opt.Width, opt.Height, opt.Width, opt.Height))
	buf.WriteString(`<rect x="0" y="0" width="100%" height="100%" fill="white"/>` + "\n")

	// 線（rail）
	buf.WriteString(`<g fill="none" stroke="#111" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" opacity="0.9">` + "\n")
	for _, e := range railEdges {
		a, okA := g.Nodes[e.From]
		b, okB := g.Nodes[e.To]
		if !okA || !okB {
			continue
		}
		if a.Kind != "coord" || b.Kind != "coord" {
			continue
		}
		x1, y1 := project(a.Lon, a.Lat)
		x2, y2 := project(b.Lon, b.Lat)
		buf.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f"/>`+"\n", x1, y1, x2, y2))
	}
	buf.WriteString(`</g>` + "\n")

	// station
	if opt.DrawStations {
		stations := collectStations(g)

		// 駅を上に描く
		buf.WriteString(`<g>` + "\n")
		for _, s := range stations {
			x, y := project(s.Lon, s.Lat)
			buf.WriteString(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="3" fill="#d00" stroke="#fff" stroke-width="1"/>`+"\n", x, y))
			if opt.DrawLabels && s.Name != "" {
				// 文字が被るので右上に少しずらす
				label := html.EscapeString(s.Name)
				buf.WriteString(fmt.Sprintf(`<text x="%.2f" y="%.2f" font-size="12" fill="#111" stroke="white" stroke-width="3" paint-order="stroke">%s</text>`+"\n",
					x+6, y-6, label))
				buf.WriteString(fmt.Sprintf(`<text x="%.2f" y="%.2f" font-size="12" fill="#111">%s</text>`+"\n",
					x+6, y-6, label))
			}
		}
		buf.WriteString(`</g>` + "\n")
	}

	buf.WriteString(`</svg>` + "\n")
	return buf.String(), nil
}

func collectCoordIDsFromRailEdges(g *graphstructure.Graph) []string {
	set := map[string]struct{}{}
	for _, e := range g.Edges {
		if e.Kind != "rail" {
			continue
		}
		if strings.HasPrefix(e.From, "coord:") {
			set[e.From] = struct{}{}
		}
		if strings.HasPrefix(e.To, "coord:") {
			set[e.To] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

type stationInfo struct {
	Name string
	Lon  float64
	Lat  float64
}

func collectStations(g *graphstructure.Graph) []stationInfo {
	out := []stationInfo{}
	for _, n := range g.Nodes {
		if n.Kind != "station" {
			continue
		}
		out = append(out, stationInfo{
			Name: n.Name,
			Lon:  n.Lon,
			Lat:  n.Lat,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

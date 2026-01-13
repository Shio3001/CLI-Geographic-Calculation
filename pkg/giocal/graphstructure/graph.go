package graphstructure

// Graph はノード集合とエッジ集合を保持する最小構造。
// Nodes は ID -> Node の辞書。
// Edges は有向エッジの配列（必要なら無向扱いは上位で両方向を追加）。
type Graph struct {
	Nodes map[string]*Node
	Edges []*Edge
}

type Node struct {
	ID   string
	Kind string // "station" | "coord" | ...
	Name string // station名など（coordは空でもOK）

	Lon float64 // 経度
	Lat float64 // 緯度

	// 任意のメタデータ（駅コード/会社/路線など）
	Meta map[string]string

	//乗降客数(keyを年度の数値 , valueを人数の文字列とする)
	Passengers map[int]float64
}

type Edge struct {
	From string
	To   string

	Kind string // "rail" | "station_at" | ...
	// 平面近似距離(km)やコストなど
	WeightKm float64

	// 任意のメタデータ（会社/路線/区間indexなど）
	Meta map[string]string
}

// NewGraph は空の Graph を返す。
func NewGraph() *Graph {
	return &Graph{
		Nodes: map[string]*Node{},
		Edges: []*Edge{},
	}
}

// AddNode はノードを追加/上書きする。
func (g *Graph) AddNode(n *Node) {
	if g.Nodes == nil {
		g.Nodes = map[string]*Node{}
	}
	g.Nodes[n.ID] = n
}

// AddEdge はエッジを追加する。
func (g *Graph) AddEdge(e *Edge) {
	g.Edges = append(g.Edges, e)
}

// Go package that extracts the graph data of nodes and edges from directed graph edited in PPTX.
// PPTX で書かれたグラフデータを読み込み、ノード・エッジの情報を抽出し保持するパッケージ
package graph_pptx

import (
	"fmt"
	//"strings"
)

const (
	Unknown_shape = "unknown"
)

// ノードシェイプの型
type Shape string

// ノードシェイプの型の値を文字列にするメソッド
func (s Shape) String() string {
	return string(s)
}

// 新しいノードシェイプの値を文字列から作る関数
func NewShape(str string) (r Shape, ok bool) {

	// TODO: str をチェックする

	r = Shape(str)
	ok = true

	return
}

// ノードの型（ラベル＋ノードシェイプ）
type Node struct {
	Label string
	Shape Shape
}

// ノードのオブジェクトを文字列にするメソッド
func (n Node) String() string {
	return fmt.Sprintf("Node{Label:\"%s\", Shape:%s}", n.Label, n.Shape.String())
}

// ノードのマップの型（文字列→ノード）
type NodeMap map[string]Node

// ノードID nid が nodes に含まれるかをチェックするメソッド
func (nm NodeMap) HasNid(nid string) bool {
	_,ok := nm[nid]
	return ok
}

// すべてのノードをダンプするメソッド
func (nm NodeMap) Dump() {
	for nid,node := range(nm) {
		fmt.Println(nid, "= ", node)		
	}	
}

// エッジの型（起点となるノードと終点となるノード）
type Edge struct {
	Src string
	Dst string
}

// エッジのオブジェクトを文字列にするメソッド
func (n Edge) String() string {
	return fmt.Sprintf("Edge{Src:%s, Dst:%s}", n.Src, n.Dst)
}

// エッジのマップの型（文字列→エッジ）
type EdgeMap map[string]Edge

// すべてのエッジをダンプするメソッド
func (em EdgeMap) Dump() {
	for eid,edge := range(em) {
		fmt.Println(eid, "= ", edge)
	}	
}

// ノードのマップを格納するグローバル変数
var nodes NodeMap

// エッジのマップを格納するグローバル変数
var edges EdgeMap

// すべてのノードに対して処理を行う関数
// 関数 f の処理結果が false となったらそこで中断する。
func DoWithNodes(f func(string,Node)bool) {
	for nid,node := range nodes {
		if !f(nid, node) {
			break
		}
	}
}

// すべてのエッジに対して処理を行う関数
// 関数 f の処理結果が false となったらそこで中断する。
func DoWithEdges(f func(string,Edge)bool) {
	for eid,edge := range edges {
		if !f(eid, edge) {
			break
		}
	}
}

// 指定したIDのノードを取得する関数
func GetNodeOf(nid string) Node {
	return nodes[nid]
}

// 指定したIDのエッジを取得する関数
func GetEdgeOf(eid string) Edge {
	return edges[eid]
}

// 初期化関数。必須。
func Init() {
	nodes = NodeMap{}
	edges = EdgeMap{}
	init_pptx()
}

// ファイルから、ノードとエッジを抽出する関数。
// グローバル変数 nodes と edges に抽出結果を保存する。
func Parse(filename string) error {
	return parse_pptx(filename)
}

// すべてのノードとエッジをダンプするデバッグ用関数
func Dump() {
	fmt.Println("# nodes")
	nodes.Dump()
	fmt.Println("# edges")
	edges.Dump()
}

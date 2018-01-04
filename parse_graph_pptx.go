package graph_pptx

// Requires: gopkg.in/xmlpath.v2

/*
pptx のドキュメントから、XPath を使って、ノードとエッジを取得する。

○ノード
・ノード群の位置
  /sld/cSld/spTree/sp

  ・ノードＩＤ：
  /sld/cSld/spTree/sp/nvSpPr/cNvPr/@id

・ノードシェイプ：
  /sld/cSld/spTree/sp/spPr/prstGeom/@prst

・ノードラベル：
  /sld/cSld/spTree/sp/txBody/p/r/t を concatenate する

○エッジ

・エッジ群の位置
  graphml/graph/edge

・エッジＩＤ：
  graphml/graph/edge/@id

・始点となるノードＩＤ：
  graphml/graph/edge/@source

・終点となるノードＩＤ：
  graphml/graph/edge/@target

 */


import (
	//"fmt"
	//"io"
	//"os"
	"archive/zip"
	//"strconv"
	"strings"
	"gopkg.in/xmlpath.v2"
)

const (

	// pptx ファイル中の処理対象となるパス
	target_path = "ppt/slides/slide1.xml"

	// pptx 用の XPath。
	// ほぼ固定になると思われるので設定ファイルにはしない。
	pptx_xpath_node = `/sld/cSld/spTree/sp`
	pptx_xpath_node_id = `nvSpPr/cNvPr/@id`
	pptx_xpath_node_shape = `spPr/prstGeom/@prst`
	pptx_xpath_node_label = "txBody/p/r/t"

	pptx_xpath_edge = `/sld/cSld/spTree/cxnSp`
	pptx_xpath_edge_id = `nvCxnSpPr/cNvPr/@id`
	pptx_xpath_edge_src = `nvCxnSpPr/cNvCxnSpPr/stCxn/@id`
	pptx_xpath_edge_dst = `nvCxnSpPr/cNvCxnSpPr/endCxn/@id`
)

// XPath オブジェクトを格納する変数
var px_node, px_node_id, px_node_shape, px_node_label *xmlpath.Path
var px_edge, px_edge_id, px_edge_src, px_edge_dst *xmlpath.Path

// 初期化関数
func init_pptx() {
	px_node = xmlpath.MustCompile(
		pptx_xpath_node)
	px_node_id = xmlpath.MustCompile(
		pptx_xpath_node_id)

	px_node_shape = xmlpath.MustCompile(
		pptx_xpath_node_shape)
		
	px_node_label = xmlpath.MustCompile(
		pptx_xpath_node_label)
	
	px_edge = xmlpath.MustCompile(
		pptx_xpath_edge)
	px_edge_id = xmlpath.MustCompile(
		pptx_xpath_edge_id)
	px_edge_src = xmlpath.MustCompile(
		pptx_xpath_edge_src)
	px_edge_dst = xmlpath.MustCompile(
		pptx_xpath_edge_dst)
}

// pptx のノードから必要な属性情報（label,shape）を取り出す関数
func px_node_props(n *xmlpath.Node) (nlabel string, nshape Shape, ok bool) {
	//fmt.Println("node_props!")
	ok = false
	nlabel = ""
	iter := px_node_label.Iter(n)
	for iter.Next() {
		nlabel += strings.TrimSpace(iter.Node().String())
	}
	if tmp, okk := px_node_shape.String(n); okk {
		tmp = strings.TrimSpace(tmp)
		nshape, ok = NewShape(tmp)
		//fmt.Println("shape =", nshape)
	}
	return
}

// ノード登録処理
// ノードの属性情報（id,label,shape） を抽出し、
// グローバル変数 nodes に保存する。
func process_px_node(p *xmlpath.Node) {
	nid := ""
	if tmp, ok := px_node_id.String(p); ok {
		nid = "n" + tmp
	}

	// id 以外の属性情報（ラベルとタイプ）を抽出する
	if nlabel, nshape, ok := px_node_props(p); ok {
		nodes[nid] = Node{Label:nlabel, Shape:nshape}
	}

}

// すべてのノード登録処理
func pickup_px_nodes(p *xmlpath.Node) error {
	iter := px_node.Iter(p)
	for iter.Next() {
		n := iter.Node()
		process_px_node(n)
	}
	return nil
}

// すべてのエッジ登録処理
// エッジの属性情報（id,src,dst ）を抽出し、
// グローバル変数 edges に保存する。
func pickup_px_edges(n *xmlpath.Node) error {
	iter := px_edge.Iter(n)
	for iter.Next() {
		e := iter.Node()
		eid := ""
		esrc := ""
		edst := ""
		if tmp, ok := px_edge_id.String(e); ok {
			eid = "e" + tmp
		}
		if tmp, ok := px_edge_src.String(e); ok {
			esrc = "n" + tmp
		}
		if tmp, ok := px_edge_dst.String(e); ok {
			edst = "n" + tmp
		}

		// esrc と edst のノードが存在するなら、エッジ情報を登録
		edges[eid] = Edge{Src:esrc, Dst:edst}
	}
	return nil
}

// PPTXファイルの中から、最初のスライドのXMLに含まれる
// ノードとエッジを XPath を使って抽出し、
// グローバル変数 nodes と edges に抽出結果を保存する。
func parse_pptx(filename string) error {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}

	var file *zip.File
	for _, file = range reader.File { // [+++]

		//fmt.Println(file.Name)
		if file.Name == target_path {
			break
		}
		
	} // [+++]

	if file == nil {
		return nil
	}

	// target_file をオープンする
	fileReader, err := file.Open()
	if err != nil {
		return err
	}

	// XMLドキュメントとしてパースする
	root, err := xmlpath.Parse(fileReader)
	if err != nil {
		return err
	}

	err = pickup_px_nodes(root)
	if err != nil {
		return err
	}

	return pickup_px_edges(root)
}

package parser

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Parser 
type Parser struct {
	doc *goquery.Document
}

func NewParserFromReader(reader io.Reader) (*Parser, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	return &Parser{doc: doc}, nil
}

func NewParser(text string) (*Parser, error) {
	return NewParserFromReader(strings.NewReader(text))
}

// Each 遍历指定的元素
// 	root: 根节点，遍历从此开始
//	sel: 遍历的条件
func (p *Parser) Each(parent, sel string) []*html.Node {
	nodes := make([]*html.Node, 0)

	p.doc.Find(parent).Each(func(i int, selection *goquery.Selection) {
		nodes = append(nodes, selection.Find(sel).Nodes...)
	})

	return nodes
}

// For 遍历指定的元素，并获取元素的文本内容和 *html.Node
func (p *Parser) For(parent, sel string, fn func(text string, node *html.Node)) {
	p.doc.Find(parent).Find(sel).Each(func(i int, selection *goquery.Selection) {
		fn(selection.Text(), selection.Nodes[0])
	})
}

// Text 返回指定元素的文本内容
func (p *Parser) Text(sel string) string {
	return p.doc.Find(sel).Text()
}

// Html 返回指定元素的 html 内容
func (p *Parser) Html(sel string) (string, error) {
	return p.doc.Find(sel).Html()
}
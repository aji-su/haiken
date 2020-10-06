package ikku

import "github.com/shogo82148/go-mecab"

type Parser struct {
	mecab mecab.MeCab
}

func NewParser() (*Parser, error) {
	mecab, err := mecab.New(make(map[string]string))
	if err != nil {
		return nil, err
	}
	return &Parser{
		mecab: mecab,
	}, nil
}

func (p *Parser) Destroy() {
	p.mecab.Destroy()
}

func (p *Parser) Parse(s string) ([]*Node, error) {
	node, err := p.mecab.ParseToNode(s)
	if err != nil {
		return nil, err
	}
	var nodes []*Node
	for ; !node.IsZero(); node = node.Next() {
		n := NewNode(node)
		if err != nil {
			return nil, err
		}
		if n.IsAnalyzable() {
			nodes = append(nodes, n)
		}
	}
	return nodes, nil
}

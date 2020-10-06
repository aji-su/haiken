package ikku

import "strings"

type Rule []int

var defaultRule = Rule{5, 7, 5}

type Song struct {
	nodes   []*Node
	exactly bool
	rule    Rule
	phrases [][]*Node
}

func NewSong(n []*Node, e bool, r Rule) (*Song, error) {
	if len(r) < 1 {
		r = defaultRule
	}
	return &Song{
		nodes:   n,
		exactly: e,
		rule:    r,
	}, nil
}

func (s *Song) String() string {
	var phrases []string
	for _, p := range s.phrases {
		var phrase []string
		for _, n := range p {
			phrase = append(phrase, n.node.Surface())
		}
		phrases = append(phrases, strings.Join(phrase, ""))
	}
	return strings.Join(phrases, ",")
}

func (s *Song) Phrases() [][]*Node {
	if s.phrases != nil {
		return s.phrases
	}
	scanner := NewScanner(s.nodes, s.exactly, s.rule)
	s.phrases = scanner.Scan()
	return s.phrases
}

func (s *Song) IsValid() bool {
	if len(s.Phrases()) < 1 {
		return false
	}
	// TODO: odd parentheses チェックは未実装
	return true
}

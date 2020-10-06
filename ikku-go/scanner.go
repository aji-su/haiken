package ikku

import "github.com/aji-su/haiken/ikku-go/util"

type Scanner struct {
	nodes   []*Node
	exactly bool
	rule    Rule
	phrases [][]*Node
	count   int
}

func NewScanner(n []*Node, e bool, r Rule) *Scanner {
	p := make([][]*Node, len(r))
	return &Scanner{
		nodes:   n,
		exactly: e,
		rule:    r,
		phrases: p,
	}
}

func (s *Scanner) Scan() [][]*Node {
	if s.nodes[0].FirstOfIkku() {
		for _, node := range s.nodes {
			if s.consume(node) {
				if s.isSatisfied() {
					if !s.exactly {
						return s.phrases
					}
				}
			} else {
				return nil
			}
		}
		if s.isSatisfied() {
			return s.phrases
		}
	}
	return nil
}

func (s *Scanner) consume(node *Node) bool {
	if node.PronunciationLength() > s.maxConsumableLength() {
		return false
	}
	if !node.ElementOfIkku() {
		return false
	}
	if s.firstOfPhrase() && !node.FirstOfPhrase() {
		return false
	}
	if node.PronunciationLength() == s.maxConsumableLength() && !node.LastOfPhrase() {
		return false
	}
	if s.phrases[s.phraseIndex()] == nil {
		s.phrases[s.phraseIndex()] = []*Node{}
	}
	s.phrases[s.phraseIndex()] = append(s.phrases[s.phraseIndex()], node)
	s.count += node.PronunciationLength()
	return true
}

func (s *Scanner) maxConsumableLength() int {
	return util.Sum(s.rule[0:(s.phraseIndex()+1)]) - s.count
}

func (s *Scanner) phraseIndex() int {
	for i := range s.rule {
		if s.count < util.Sum(s.rule[0:(i+1)]) {
			return i
		}
	}
	return len(s.rule) - 1
}

func (s *Scanner) firstOfPhrase() bool {
	var ns []int
	for i, length := range s.rule {
		if i == 0 {
			ns = append(ns, length)
		} else {
			ns = append(ns, ns[len(ns)-1]+length)
		}
	}
	return util.IContains(ns, s.count)
}

func (s *Scanner) isSatisfied() bool {
	return s.count == util.Sum(s.rule) && s.hasValidLastNode()
}

func (s *Scanner) hasValidLastNode() bool {
	last := s.phrases[len(s.phrases)-1]
	return last[len(last)-1].LastOfIkku()
}

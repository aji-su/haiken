package ikku

type Reviewer struct {
	parser *Parser
	rule   Rule
}

func NewReviewer(p *Parser, r Rule) *Reviewer {
	return &Reviewer{
		parser: p,
		rule:   r,
	}
}

func (r *Reviewer) Search(text string) ([]*Node, []*Song, error) {
	nodes, err := r.parser.Parse(text)
	if err != nil {
		return nil, nil, err
	}
	var songs []*Song
	for i := range nodes {
		song, err := NewSong(nodes[i:], false, r.rule)
		if err != nil {
			return nil, nil, err
		}
		if song.IsValid() {
			songs = append(songs, song)
		}
	}
	return nodes, songs, nil
}

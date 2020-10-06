package ikku

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/aji-su/haiken/ikku-go/util"
	"github.com/shogo82148/go-mecab"
	"gopkg.in/guregu/null.v3"
)

type Node struct {
	node                mecab.Node
	feature             []string
	pronunciationLength null.Int
}

func NewNode(n mecab.Node) *Node {
	return &Node{
		node: n,
	}
}

func (n *Node) String() string {
	return n.node.String()
}

// GetFeature returns the feature
// 0:品詞
// 1:品詞細分類1
// 2:品詞細分類2
// 3:品詞細分類3
// 4:活用型
// 5:活用形
// 6:原形
// 7:読み
// 8:発音
func (n *Node) GetFeature(i int) string {
	if n.feature == nil {
		n.feature = strings.Split(n.node.Feature(), ",")
	}
	if i < len(n.feature) {
		return n.feature[i]
	}
	return ""
}

func (n *Node) IsAnalyzable() bool {
	return n.node.Stat() != mecab.BOSNode && n.node.Stat() != mecab.EOSNode
}

func (n *Node) Surface() string {
	return n.node.Surface()
}

func (n *Node) FirstOfIkku() bool {
	switch {
	case !n.FirstOfPhrase():
		return false
	case n.Type() == "記号" && !util.SContains([]string{"括弧開", "括弧閉"}, n.SubType1()):
		return false
	default:
		return true
	}
}

func (n *Node) FirstOfPhrase() bool {
	switch {
	case util.SContains([]string{"助詞", "助動詞"}, n.Type()):
		return false
	case util.SContains([]string{"非自立", "接尾"}, n.SubType1()):
		return false
	case n.SubType1() == "自立" && util.SContains([]string{"する", "できる"}, n.RootForm()):
		return false
	default:
		return true
	}
}

func (n *Node) LastOfIkku() bool {
	if util.SContains([]string{"名詞接続", "格助詞", "係助詞", "連体化", "接続助詞", "並立助詞", "副詞化", "数接続", "連体詞"}, n.Type()) {
		return false
	}
	if n.Conjugation2() == "連用タ接続" {
		return false
	}
	if n.Conjugation1() == "サ変・スル" && n.Conjugation2() == "連用形" {
		return false
	}
	if n.Type() == "動詞" && util.SContains([]string{"仮定形", "未然形"}, n.Conjugation2()) {
		return false
	}
	if n.Type() == "名詞" && n.SubType1() == "非自立" && n.Pronunciation() == "ン" {
		return false
	}
	return true
}

func (n *Node) LastOfPhrase() bool {
	return n.Type() != "接頭詞"
}

func (n *Node) Type() string {
	return n.GetFeature(0)
}

func (n *Node) SubType1() string {
	return n.GetFeature(1)
}

func (n *Node) SubType2() string {
	return n.GetFeature(2)
}

func (n *Node) SubType3() string {
	return n.GetFeature(3)
}

func (n *Node) Conjugation1() string {
	return n.GetFeature(4)
}

func (n *Node) Conjugation2() string {
	return n.GetFeature(5)
}

func (n *Node) RootForm() string {
	return n.GetFeature(6)
}

func (n *Node) Pronunciation() string {
	return n.GetFeature(8)
}

func (n *Node) PronunciationLength() int {
	if !n.pronunciationLength.Valid {
		if n.Pronunciation() != "" {
			n.pronunciationLength = null.IntFrom(int64(utf8.RuneCountInString(n.pronunciationMora())))
		} else {
			n.pronunciationLength = null.IntFrom(0)
		}
	}
	return int(n.pronunciationLength.ValueOrZero())
}

var (
	charFrom = []rune{'ぁ', 'ゔ'}
	charTo   = []rune{'ァ', 'ヴ'}
	replacer = regexp.MustCompile("[^アイウエオカ-モヤユヨラ-ロワヲンヴー]")
)

func (n *Node) pronunciationMora() string {
	if n.Pronunciation() != "" {
		s := tr(charFrom, charTo, n.Pronunciation())
		return replacer.ReplaceAllString(s, "")
	}
	return ""
}

func tr(charFrom []rune, charTo []rune, in string) string {
	sub := charTo[0] - charFrom[0]
	out := make([]rune, 0, len(in))
	for _, c := range in {
		if charFrom[0] <= c && c <= charTo[0] {
			out = append(out, c+sub)
		} else {
			out = append(out, c)
		}
	}
	return string(out)
}

func (n *Node) ElementOfIkku() bool {
	return n.node.Stat() == mecab.NormalNode
}

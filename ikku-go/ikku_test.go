package ikku_test

import (
	"testing"

	ikku "github.com/aji-su/haiken/ikku-go"
	"github.com/google/go-cmp/cmp"
)

func TestSearch(t *testing.T) {
	p, err := ikku.NewParser()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Destroy()

	tests := []struct {
		name string
		text string
		want []string
		rule ikku.Rule
	}{
		{
			name: "with invalid song",
			text: "test",
			want: nil,
		},
		{
			name: "with valid song",
			text: "古池や蛙飛びこむ水の音",
			want: []string{"古池や,蛙飛びこむ,水の音"},
		},
		{
			name: "with text including song",
			text: "ああ古池や蛙飛びこむ水の音ああ",
			want: []string{"古池や,蛙飛びこむ,水の音"},
		},
		{
			name: "with text including song ending with 連用タ接続",
			text: "リビングでコーヒー飲んでだめになってる",
			want: nil,
		},
		{
			name: "with song ending with 仮定形",
			text: "その人に金をあげたい人がいれば",
			want: nil,
		},
		{
			name: "with song ending with 未然形 (い)",
			text: "学会に多分ネイティブほとんどいない",
			want: nil,
		},
		{
			name: "with song ending with ん as 非自立名詞",
			text: "古池や蛙飛びこむかかったんだ",
			want: nil,
		},
		{
			name: "with empty text",
			text: "",
			want: nil,
		},
		{
			name: "with rule option and invalid song on neologd",
			text: "すもももももももものうち",
			want: nil,
			rule: ikku.Rule{4, 3, 5},
		},
		{
			name: "with rule option and valid song",
			text: "人生楽ありゃ苦もあるさ",
			want: []string{"人生,楽ありゃ,苦もあるさ"},
			rule: ikku.Rule{4, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ikku.NewReviewer(p, tt.rule)
			nodes, songs, err := r.Search(tt.text)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(nodes, songs)
			var got []string
			for _, song := range songs {
				got = append(got, song.String())
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("songs mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

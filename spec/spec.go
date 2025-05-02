package spec

type HiraganaRomajiPair struct {
	Hiragana string
	Romaji   string
}

type QS struct {
	Hiragana string   `query:"hiragana" json:"hiragana"`
	Answer   string   `query:"answer" json:"answer"`
	Choices  []string `query:"choices" json:"choices"`
}

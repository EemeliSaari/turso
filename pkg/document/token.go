package document

// TokenType ...
type TokenType string

const (
	// TextToken ...
	TextToken TokenType = "text"
	// HTMLToken ...
	HTMLToken TokenType = "html"
	// SpecialToken ...
	SpecialToken TokenType = "special"
	// PuncToken ...
	PuncToken TokenType = "punc"
)

// Token ...
type Token struct {
	Start   int       `json:"start"`
	End     int       `json:"end"`
	Content string    `json:"content"`
	Idx     int       `json:"index"`
	Type    TokenType `json:"type"`
}

// CompressedToken ...
type CompressedToken struct {
	Token
	Content string `json:"-"`
}

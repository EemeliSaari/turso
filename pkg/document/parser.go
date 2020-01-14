package document

import (
	"errors"
	"fmt"
	"strings"
)

// TokenizerMode ...
type mode int

const (
	// NullMode ...
	NullMode mode = iota
	// HTMLMode ...
	HTMLMode
	// TextMode ...
	TextMode
	// SkipMode ...
	SkipMode
	// PuncMode ...
	PuncMode
)

// Tokenizer ...
type Tokenizer struct {
	offset   int
	tokenIdx int
	mode     mode
	buffer   *[]byte
}

// NewTokenizer ...
func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		offset:   0,
		tokenIdx: 0,
		mode:     NullMode,
	}
}

// TokenReader ...
func (t *Tokenizer) TokenReader(data *[]byte) <-chan *Token {
	t.buffer = data
	token := make(chan *Token)
	go func() {
		var read bool
		var b byte

		stop := false
		i := 0
		for i < len(*t.buffer) {
			b = (*t.buffer)[i]
			read = true

			switch b {
			case 60, 38: // "<", "&"
				if t.mode == TextMode && !stop {
					stop = true
				} else {
					end, err := t.handleHTML(i, b)
					if err == nil {
						i += end
						t.offset += end
						t.mode = HTMLMode
					} else {
						t.mode = PuncMode
					}
					stop = false
				}
				if t.mode != SkipMode {
					read = false
				}
			case 32: // " "
				if t.mode == TextMode || t.mode == PuncMode {
					read = false
				} else {
					t.mode = SkipMode
				}
			// .,?!:;()"[]{}#'+<>/\
			case 46, 44, 63, 33, 58, 59, 40, 41, 34, 91, 93, 123, 125, 35, 39, 43, 62, 47, 92:
				if t.mode != HTMLMode && !stop && t.mode != NullMode {
					stop = true
				} else if stop {
					stop = false
					t.mode = PuncMode
				} else if t.mode == NullMode {
					t.mode = PuncMode
				}
				if t.mode != SkipMode {
					read = false
				}
			case 10: // "\n"
				t.mode = SkipMode
			default:
				t.mode = TextMode
			}

			if t.mode != SkipMode {
				if read {
					t.offset++
				} else {
					token <- t.read(i)
				}
			} else {
				t.mode = NullMode
			}

			if stop {
				continue
			} else {
				i++
			}
		}
		close(token)
	}()
	return token
}

func (t *Tokenizer) read(end int) *Token {
	var endIdx int
	var tp TokenType

	if t.mode == HTMLMode || t.mode == PuncMode {
		endIdx = end + 1
	} else {
		endIdx = end
	}
	switch t.mode {
	case HTMLMode:
		tp = HTMLToken
	case TextMode:
		tp = TextToken
	case PuncMode:
		tp = PuncToken
	}

	token := Token{
		Start:   end - t.offset,
		End:     endIdx,
		Content: string((*t.buffer)[end-t.offset : endIdx]),
		Idx:     t.tokenIdx,
		Type:    tp,
	}
	t.tokenIdx++
	t, _ = t.reset()
	return &token
}

func (t *Tokenizer) reset() (*Tokenizer, error) {
	t.mode = NullMode
	t.offset = 0
	return t, nil
}

func (t *Tokenizer) handleHTML(idx int, value byte) (int, error) {
	var err error
	var end int

	switch value {
	case 60:
		end, err = t.handleTag(idx, value)
		if err == nil {
			err = t.validateTag(idx, idx+end)
		}
	case 38:
		end, err = t.handleSymbol(idx, value)
		if err == nil {
			err = t.validateSymbol(idx, idx+end)
		}
	default:
		err = fmt.Errorf("Received unknown character: %s", string(value))
	}
	return end, err
}

func (t *Tokenizer) handleTag(idx int, value byte) (int, error) {
	var end int
	var err error

	for i, b := range (*t.buffer)[idx:] {
		if b == 62 { // ">"
			end = i
			break
		} else if b == 60 && i > 0 { // "<"
			err = errors.New("invalid character \"<\" encountered within HTML tag")
		} else {
			continue
		}
	}
	return end, err
}

func (t *Tokenizer) validateSymbol(start int, end int) error {
	var err error
	str := strings.Split(string((*t.buffer)[start+1:end]), " ")[0]
	if ok := inspectEntity(str); ok {
		err = fmt.Errorf("received unknown symbol: %s", str)
	}
	return err
}

func (t *Tokenizer) validateTag(start int, end int) error {
	var err error
	shiftl := 1
	if (*t.buffer)[start+1] == 47 { // "/"
		shiftl++
	}
	str := strings.Split(string((*t.buffer)[start+shiftl:end]), " ")[0]
	if !HTMLTagSet[str] {
		err = fmt.Errorf("received unknown tag: %s", str)
	}
	return err
}

func (t *Tokenizer) handleSymbol(idx int, value byte) (int, error) {
	var end int
	var err error

	for i, b := range (*t.buffer)[idx:] {
		if b == 59 {
			end = i
			break
		} else if b == 32 {
			err = errors.New("invalid symbol encountered")
		} else {
			continue
		}
	}
	return end, err
}

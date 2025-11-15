package tokenizer

import (
	"encoding/json"
	"os"
)

// Config construct configuration for creating Tokenizer.
type Config struct {
	Version       string                 `json:"version"`
	Truncation    map[string]interface{} `json:"truncation"`
	Padding       map[string]interface{} `json:"padding"`
	AddedTokens   []TokenConfig          `json:"added_tokens"`
	Normalizer    map[string]interface{} `json:"normalizer"`
	PreTokenizer  map[string]interface{} `json:"pre_tokenizer"`
	PostProcessor map[string]interface{} `json:"post_processor"`
	Decoder       map[string]interface{} `json:"decoder"`
	Model         map[string]interface{} `json:"model"`
}

type TokenConfig struct {
	Id         int64  `json:"id"`
	Content    string `json:"content"`
	SingleWord bool   `json:"single_word"`
	Lstrip     bool   `json:"lstrip"`
	Rstrip     bool   `json:"rstrip"`
	Normalized bool   `json:"normalized"`
	Special    bool   `json:"special"`
}

type NormalizerConfig struct {
	Type        string                   `json:"type"`
	Normalizers []map[string]interface{} `json:"normalizers"`
}
type PreTokenizerConfig struct{}
type PostProcessorConfig struct {
	Type          string                   `json:"type"`
	Single        []map[string]interface{} `json:"single"`
	Pair          []map[string]interface{} `json:"pair"`
	SpecialTokens map[string]interface{}   `json:"speical_tokens"`
}

type DecoderConfig struct {
	Type     string                   `json:"type"`
	Decoders []map[string]interface{} `json:"decoders"`
}

type ModelConfig struct {
	Type                    string         `json:"type"`
	Dropout                 interface{}    `json:"dropout"`
	UnkToken                string         `json:"unk_token"`
	ContinuingSubwordPrefix interface{}    `json:"continuing_subword_prefix"`
	EndOfWordSuffix         interface{}    `json:"end_of_word_suffix"`
	FuseUnk                 bool           `json:"fuse_unk"`
	ByteFallback            bool           `json:"byte_fallback"`
	Vocab                   map[string]int `json:"vocab"`
	Merges                  []string       `json:"merges"`
	MaxInputCharsPerWord    interface{}    `json:"max_input_chars_per_word"`
}

// ConfigFromFile loads config from file.
func ConfigFromFile(file string) (*Config, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	normalized, err := normalizeConfigJSON(raw)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(normalized, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func normalizeConfigJSON(raw []byte) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	if model, ok := data["model"].(map[string]interface{}); ok {
		if merges, ok := model["merges"].([]interface{}); ok && len(merges) > 0 {
			flattened := make([]string, 0, len(merges))
			for _, item := range merges {
				switch v := item.(type) {
				case []interface{}:
					if len(v) != 2 {
						continue
					}
					left, lok := v[0].(string)
					right, rok := v[1].(string)
					if lok && rok {
						flattened = append(flattened, left+" "+right)
					}
				case string:
					flattened = append(flattened, v)
				}
			}
			model["merges"] = flattened
		}
	}

	return json.Marshal(data)
}

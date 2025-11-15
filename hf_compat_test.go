package tokenizer_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

func loadFixtureTokenizer(t *testing.T, name string) *tokenizer.Tokenizer {
	t.Helper()
	path := filepath.Join("testdata", name)
	tk, err := pretrained.FromFile(path)
	if err != nil {
		t.Fatalf("load fixture %s: %v", name, err)
	}
	return tk
}

func TestHuggingFaceBPEPriority(t *testing.T) {
	tk := loadFixtureTokenizer(t, "bpe_priority_tokenizer.jsonc")
	enc, err := tk.EncodeSingle("Ġ1000", false)
	if err != nil {
		t.Fatalf("encode fixture input: %v", err)
	}
	want := []string{"Ġ1000"}
	if !reflect.DeepEqual(enc.Tokens, want) {
		t.Fatalf("unexpected tokens: want %v, got %v", want, enc.Tokens)
	}
}

func TestHuggingFaceByteFallback(t *testing.T) {
	tk := loadFixtureTokenizer(t, "bpe_byte_fallback_tokenizer.jsonc")
	enc, err := tk.EncodeSingle("😀", false)
	if err != nil {
		t.Fatalf("encode fixture input: %v", err)
	}
	want := []string{"<0xF0>", "<0x9F>", "<0x98>", "<0x80>"}
	if !reflect.DeepEqual(enc.Tokens, want) {
		t.Fatalf("unexpected tokens: want %v, got %v", want, enc.Tokens)
	}
}

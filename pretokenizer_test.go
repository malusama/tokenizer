package tokenizer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sugarme/tokenizer/normalizer"
)

func TestBytesToCharConverter(t *testing.T) {

	sequence := "Löwe 老虎 Léopard"
	converter := NewBytesToCharOffsetConverter(sequence)

	want := map[int]int{
		0:  0,
		1:  1,
		2:  1,
		3:  2,
		4:  3,
		5:  4,
		6:  5,
		7:  5,
		8:  5,
		9:  6,
		10: 6,
		11: 6,
		12: 7,
		13: 8,
		14: 9,
		15: 9,
		16: 10,
		17: 11,
		18: 12,
		19: 13,
		20: 14,
	}

	got := converter.b2c

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v\n", want, got)
	}
}

func TestPreTokenizedStringNormalizeKeepsExistingTokens(t *testing.T) {
	first := normalizer.NewNormalizedFrom("hi")
	second := normalizer.NewNormalizedFrom("there")
	pt := &PreTokenizedString{
		original: "hi there",
		splits: []Split{
			{normalized: first, tokens: []Token{NewToken(1, "hi", []int{0, 2})}},
			{normalized: second, tokens: nil},
		},
	}

	upper := func(ns *normalizer.NormalizedString) *normalizer.NormalizedString {
		return normalizer.NewNormalizedFrom(strings.ToUpper(ns.GetNormalized()))
	}

	pt.Normalize(upper)

	if len(pt.splits) != 2 {
		t.Fatalf("normalize dropped splits with tokens: %+v", pt.splits)
	}

	if len(pt.splits[0].tokens) != 1 {
		t.Fatalf("existing tokens were removed during normalize")
	}

	if got := pt.splits[1].normalized.GetNormalized(); got != "THERE" {
		t.Fatalf("second split was not normalized, want THERE got %s", got)
	}
}

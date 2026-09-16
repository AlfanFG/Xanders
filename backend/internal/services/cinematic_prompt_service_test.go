package services

import "testing"

func TestUIInputFromPromptInputPreservesTemplateAndAdditionalInstructions(t *testing.T) {
	input := UIInputFromPromptInput(map[string]interface{}{
		"template":     "automotive",
		"lighting":     "golden-hour",
		"camera":       "drone",
		"audio":        "ambient",
		"personModel":  "No Model (Product Only)",
		"customPrompt": "A red car races through rain before stopping under neon signs.",
	})

	if input.Product != "automotive" {
		t.Fatalf("Product = %q, want selected template", input.Product)
	}
	if input.AdditionalInstructions != "A red car races through rain before stopping under neon signs." {
		t.Fatalf("AdditionalInstructions = %q, want custom direction", input.AdditionalInstructions)
	}
	if input.Motion != "drone" {
		t.Fatalf("Motion = %q, want selected camera", input.Motion)
	}
}

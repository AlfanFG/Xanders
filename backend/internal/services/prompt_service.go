package services

import (
	"fmt"
	"strings"
)

type PromptService struct{}

func NewPromptService() *PromptService {
	return &PromptService{}
}

// Template descriptors for dynamic prompt building
var templateDescriptors = map[string]string{
	"fashion":     "high-fashion product showcase with a sleek garment as the hero",
	"accessories": "luxury accessories product reveal — watch, bag, or jewelry as the focal point",
	"pc-gaming":   "cutting-edge gaming rig and peripherals with RGB lighting",
	"electronics": "premium consumer electronics product on a polished surface",
	"automotive":  "exotic sports car gleaming under studio lights",
	"real-estate": "stunning architectural exterior of a luxury property",
	"aviation":    "private jet on a tarmac at golden hour",
}

var lightingDescriptors = map[string]string{
	"golden-hour": "warm golden hour sunlight, long shadows, soft bokeh",
	"cyberpunk":   "neon-lit cyberpunk environment, vibrant pinks and purples, volumetric fog",
	"studio":      "clean soft-box studio lighting, white seamless background, professional product photography",
	"dramatic":    "dramatic high-contrast chiaroscuro, deep blacks, single hard key light",
}

var cameraDescriptors = map[string]string{
	"slow-track":  "slow cinematic horizontal tracking shot, smooth dolly movement",
	"drone":       "sweeping high-angle drone flyover, revealing aerial perspective",
	"static":      "perfectly composed static tripod shot, locked-off, razor-sharp focus",
	"orbit":       "elegant 360-degree orbit around the subject, smooth circular pan",
}

// AssemblePrompt transforms wizard selections into a rich, detailed cinematic prompt.
func (s *PromptService) AssemblePrompt(input map[string]interface{}) string {
	// Check for a full custom prompt override first
	if customPrompt, ok := input["customPrompt"].(string); ok && strings.TrimSpace(customPrompt) != "" {
		return customPrompt
	}

	template, _ := input["template"].(string)
	lighting, _ := input["lighting"].(string)
	camera, _ := input["camera"].(string)
	audio, _ := input["audio"].(string)

	// Resolve descriptors with fallbacks
	templateDesc := templateDescriptors[template]
	if templateDesc == "" {
		templateDesc = fmt.Sprintf("sleek %s product", template)
	}

	lightingDesc := lightingDescriptors[lighting]
	if lightingDesc == "" {
		lightingDesc = "cinematic natural lighting"
	}

	cameraDesc := cameraDescriptors[camera]
	if cameraDesc == "" {
		cameraDesc = "cinematic tracking shot"
	}

	// Build audio/sound note
	audioNote := ""
	switch audio {
	case "epic":
		audioNote = " Epic orchestral score with soaring strings."
	case "ambient":
		audioNote = " Ambient electronic underscore, calm and immersive."
	case "corporate":
		audioNote = " Upbeat corporate background music, energetic and professional."
	case "none":
		audioNote = " No audio, video only."
	}

	prompt := fmt.Sprintf(
		"Cinematic 4K commercial video: %s. %s. %s. Photorealistic, 35mm anamorphic lens, shallow depth of field, hyperdetailed, award-winning commercial photography quality.%s",
		templateDesc,
		lightingDesc,
		cameraDesc,
		audioNote,
	)

	return prompt
}

package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/genai"
)

const CinematicPromptModel = "gemini-3-flash-preview"

type UIInput struct {
	Product                string `json:"product"`
	Environment            string `json:"environment"`
	Style                  string `json:"style"`
	Motion                 string `json:"motion"`
	PersonModel            string `json:"person_model"`
	AdditionalInstructions string `json:"additional_instructions"`
}

func UIInputFromPromptInput(input map[string]interface{}) UIInput {
	template, _ := input["template"].(string)
	lighting, _ := input["lighting"].(string)
	camera, _ := input["camera"].(string)
	audio, _ := input["audio"].(string)
	customPrompt, _ := input["customPrompt"].(string)
	personModel, _ := input["personModel"].(string)

	return UIInput{
		Product:                strings.TrimSpace(template),
		Environment:            strings.TrimSpace(lighting),
		Style:                  strings.TrimSpace(audio),
		Motion:                 strings.TrimSpace(camera),
		PersonModel:            strings.TrimSpace(personModel),
		AdditionalInstructions: strings.TrimSpace(customPrompt),
	}
}

type CinematicPromptService struct{}

func NewCinematicPromptService() *CinematicPromptService {
	return &CinematicPromptService{}
}

func (s *CinematicPromptService) GenerateCinematicPrompt(ctx context.Context, input UIInput) (string, error) {
	// Explicitly use the Gemini Developer API (API key) to avoid conflict with
	// Vertex AI env vars (GOOGLE_CLOUD_PROJECT / GOOGLE_CLOUD_LOCATION).
	apiKey := os.Getenv("GOOGLE_API_KEY")
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
		APIKey:  apiKey,
	})
	if err != nil {
		return "", fmt.Errorf("create Gemini client: %w", err)
	}

	// Build person model context — injected into both prompt and instruction
	personCtx := ""
	if input.PersonModel != "" && input.PersonModel != "No Model (Product Only)" {
		personCtx = fmt.Sprintf("\nPerson/Model in video: %s", input.PersonModel)
	}

	additionalInstructions := "None provided."
	if input.AdditionalInstructions != "" {
		additionalInstructions = input.AdditionalInstructions
	}

	rawPrompt := fmt.Sprintf(
		"Product category: %s\nLighting & mood: %s\nCamera motion: %s\nBackground audio style: %s\nAdditional creator direction (required): %s%s",
		input.Product,
		input.Environment,
		input.Motion,
		input.Style,
		additionalInstructions,
		personCtx,
	)

	personInstruction := ""
	if personCtx != "" {
		personInstruction = fmt.Sprintf(
			"The video must prominently feature a %s naturally interacting with or using the product — show genuine emotion, lifestyle context, and human energy.",
			input.PersonModel,
		)
	} else {
		personInstruction = "No person appears — focus entirely on the product with creative camera choreography."
	}

	systemInstruction := strings.Join([]string{
		"You are a world-class commercial director and Veo 3 prompt engineer at a top advertising agency.",
		"Your job is to transform product category inputs into ONE rich, ultra-descriptive video generation prompt that produces premium-quality, broadcast-ready commercial footage.",
		"When additional creator direction is provided, treat every compatible detail in it as a required visual constraint. Preserve its subject, action, setting, and requested changes in the final prompt; do not replace it with a generic product showcase.",

		"CRITICAL RULES — follow every one or the output is rejected:",
		"1. NEVER describe a product simply rotating or spinning on a pedestal — that is low quality and forbidden.",
		"2. ALWAYS include dynamic, purposeful camera movement: slow push-ins, elegant orbits, sweeping crane shots, rack focus pulls, macro close-ups transitioning to wide establishing shots, or hand-held tracking.",
		"3. Describe the environment in depth: surface materials, reflections, particle effects (steam, mist, floating dust, bokeh light orbs), depth of field, and background atmosphere.",
		"4. Use cinematography vocabulary: anamorphic lens flares, shallow depth of field, golden-hour rim light, practical light sources (candles, neon tubes, LED panels), color grading (teal-and-orange LUT, desaturated matte, vibrant pop art).",
		"5. Output exactly ONE flowing paragraph — no lists, no labels, no markdown, no line breaks, no explanations.",
		"6. The total video duration is 8 seconds — make every second count with a clear visual arc: establish → reveal → detail → emotional close.",
		personInstruction,
	}, " ")

	resp, err := client.Models.GenerateContent(
		ctx,
		CinematicPromptModel,
		genai.Text(rawPrompt),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: systemInstruction}},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("generate cinematic prompt: %w", err)
	}

	optimizedPrompt := strings.TrimSpace(resp.Text())
	if optimizedPrompt == "" {
		return "", errors.New("Gemini returned an empty optimized prompt")
	}

	log.Printf("[CinematicPrompt] Generated prompt: %s", optimizedPrompt)
	return optimizedPrompt, nil
}

func CallVeoVideoService(prompt string) (string, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", errors.New("prompt is required")
	}

	return "veo-task-" + uuid.NewString(), nil
}

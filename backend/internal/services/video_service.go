package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"xanders-gen-video/internal/config"
	"xanders-gen-video/internal/models"

	"cloud.google.com/go/storage"
	"google.golang.org/genai"
)

type VideoService struct {
	cfg *config.Config
}

func NewVideoService(cfg *config.Config) *VideoService {
	return &VideoService{cfg: cfg}
}

// GenerateVideo triggers Veo video generation and returns the GCS URI of the completed video.
func (s *VideoService) GenerateVideo(ctx context.Context, job *models.Job) (string, error) {
	log.Printf("Starting Vertex AI Veo 3 generation for job: %s", job.ID)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend:  genai.BackendVertexAI,
		Project:  s.cfg.GoogleCloudProject,
		Location: s.cfg.GoogleCloudLocation,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gen AI client: %v", err)
	}

	source := &genai.GenerateVideosSource{
		Prompt: job.AssembledPrompt,
	}

	// If there's an image uploaded locally, we upload it to Gen AI or pass its bytes
	if job.ImageURL != nil && *job.ImageURL != "" {
		// Our ImageURL is something like /uploads/filename.ext
		// For the SDK, if we have local bytes, we can read them and pass them as *genai.Image
		localPath := "." + *job.ImageURL // e.g., ./uploads/filename.ext
		imageBytes, err := os.ReadFile(localPath)
		if err == nil {
			source.Image = &genai.Image{
				ImageBytes: imageBytes,
				// MIME type could be inferred, assume jpeg or png based on ext
				MIMEType: "image/jpeg",
			}
			if strings.HasSuffix(localPath, ".png") {
				source.Image.MIMEType = "image/png"
			}
			log.Printf("Attached reference image: %s", localPath)
		} else {
			log.Printf("Warning: failed to read image from disk %s: %v", localPath, err)
		}
	}

	duration := int32(8)
	generateAudio := true
	videoConfig := &genai.GenerateVideosConfig{
		OutputGCSURI:     s.cfg.GoogleCloudGCSBucket,
		AspectRatio:      "16:9",
		DurationSeconds:  &duration,
		Resolution:       "1080p",
		EnhancePrompt:    true,
		NumberOfVideos:   1,
		PersonGeneration: "allow_adult",
		GenerateAudio:    &generateAudio,
	}

	op, err := client.Models.GenerateVideosFromSource(ctx, "veo-3.0-generate-001", source, videoConfig)
	if err != nil {
		return "", fmt.Errorf("GenerateVideosFromSource failed: %v", err)
	}

	log.Printf("Video generation operation created: %s", op.Name)

	// Poll until completion (Veo generation takes minutes)
	for !op.Done {
		log.Printf("Polling operation %s...", op.Name)
		time.Sleep(15 * time.Second)
		op, err = client.Operations.GetVideosOperation(ctx, op, nil)
		if err != nil {
			log.Printf("Warning: failed to poll operation: %v", err)
			continue
		}
	}

	if op.Error != nil {
		return "", fmt.Errorf("operation failed: %v", op.Error)
	}

	// Log the response to see exactly what Google returned
	responseJSON, _ := json.MarshalIndent(op.Response, "", "  ")
	log.Printf("Operation completed! Full Response:\n%s", string(responseJSON))

	// Extract the output URI if available
	if op.Response != nil && len(op.Response.GeneratedVideos) > 0 {
		vid := op.Response.GeneratedVideos[0]
		if vid.Video != nil && vid.Video.URI != "" {
			uri := vid.Video.URI
			if strings.HasPrefix(uri, "gs://") {
				trimmed := strings.TrimPrefix(uri, "gs://")
				parts := strings.SplitN(trimmed, "/", 2)
				if len(parts) == 2 {
					bucket := parts[0]
					object := parts[1]

					client, err := storage.NewClient(ctx)
					if err == nil {
						opts := &storage.SignedURLOptions{
							Scheme:  storage.SigningSchemeV4,
							Method:  "GET",
							Expires: time.Now().Add(7 * 24 * time.Hour),
						}
						signedURL, err := client.Bucket(bucket).SignedURL(object, opts)
						client.Close()
						
						if err == nil {
							return signedURL, nil
						}
						log.Printf("Warning: failed to generate signed URL via client: %v. Falling back to public URL.", err)
					} else {
						log.Printf("Warning: failed to initialize storage client: %v. Falling back to public URL.", err)
					}
				}
				uri = strings.Replace(uri, "gs://", "https://storage.googleapis.com/", 1)
			}
			return uri, nil
		}
	}

	// If we can't parse it out dynamically, we return the folder we gave it
	log.Println("Could not parse URI from operation response. Video might be missing or in a different field.")
	return s.cfg.GoogleCloudGCSBucket, nil
}

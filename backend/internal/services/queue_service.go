package services

import (
	"context"
	"log"

	"xanders-gen-video/internal/database"
	"xanders-gen-video/internal/models"
)

type QueueService struct {
	videoService *VideoService
}

func NewQueueService(videoService *VideoService) *QueueService {
	return &QueueService{videoService: videoService}
}

// EnqueueVideoJob is a placeholder for Google Cloud Tasks.
// For now, we process it async using a Goroutine.
func (s *QueueService) EnqueueVideoJob(ctx context.Context, jobID string) error {
	log.Printf("Enqueued job %s to local worker", jobID)

	go func() {
		bgCtx := context.Background()
		log.Printf("[Worker] Starting to process job %s", jobID)

		// 1. Fetch the job from the DB
		var job models.Job
		if err := database.DB.WithContext(bgCtx).Where("id = ?", jobID).First(&job).Error; err != nil {
			log.Printf("Worker Error: could not find job %s: %v", jobID, err)
			return
		}

		// 2. Set to processing
		if err := database.DB.WithContext(bgCtx).Model(&job).Update("status", "processing").Error; err != nil {
			log.Printf("Error updating job status: %v", err)
			return
		}

		// 3. Call Veo 3 Generation via Gen AI SDK
		// We pass the job reference in case the VideoService needs the ImageURL
		gcsURI, err := s.videoService.GenerateVideo(bgCtx, &job)
		if err != nil {
			log.Printf("[Worker] Job %s failed: %v", jobID, err)
			// Update to failed
			errorMessage := err.Error()
			database.DB.WithContext(bgCtx).Model(&job).Updates(map[string]interface{}{
				"status":        "failed",
				"error_message": &errorMessage,
			})
			return
		}

		// 4. Set to completed with the real video URL
		err = database.DB.WithContext(bgCtx).Model(&job).Updates(map[string]interface{}{
			"status":        "completed",
			"gcs_video_url": &gcsURI,
		}).Error
		
		if err != nil {
			log.Printf("Error completing job: %v", err)
		} else {
			log.Printf("[Worker] Job %s completed successfully!", jobID)
		}
	}()

	return nil
}

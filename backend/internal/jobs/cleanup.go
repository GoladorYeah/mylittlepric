package jobs

import (
	"context"
	"log"
	"time"

	"mylittleprice/internal/services"
)

// CleanupJob handles periodic cleanup of expired anonymous search history
type CleanupJob struct {
	searchHistoryService *services.SearchHistoryService
	interval             time.Duration
	ctx                  context.Context
	cancel               context.CancelFunc
}

// NewCleanupJob creates a new cleanup job instance
func NewCleanupJob(shs *services.SearchHistoryService) *CleanupJob {
	ctx, cancel := context.WithCancel(context.Background())
	return &CleanupJob{
		searchHistoryService: shs,
		interval:             24 * time.Hour,
		ctx:                  ctx,
		cancel:               cancel,
	}
}

// Start begins the cleanup job ticker
func (j *CleanupJob) Start() {
	ticker := time.NewTicker(j.interval)
	go func() {
		// Run cleanup immediately on start
		log.Println("🧹 Running initial cleanup...")
		j.runCleanup()

		// Then run on schedule
		for {
			select {
			case <-ticker.C:
				j.runCleanup()
			case <-j.ctx.Done():
				ticker.Stop()
				log.Println("🛑 Cleanup job ticker stopped")
				return
			}
		}
	}()
	log.Println("🧹 Cleanup job started (runs every 24h)")
}

// runCleanup executes the cleanup operation
func (j *CleanupJob) runCleanup() {
	startTime := time.Now()
	count, err := j.searchHistoryService.CleanupExpiredAnonymousHistory(j.ctx)
	duration := time.Since(startTime)

	if err != nil {
		log.Printf("❌ Cleanup job failed after %v: %v", duration, err)
	} else {
		if count > 0 {
			log.Printf("✅ Cleanup job completed in %v: %d records deleted", duration, count)
		} else {
			log.Printf("✅ Cleanup job completed in %v: no expired records found", duration)
		}
	}
}

// Stop gracefully stops the cleanup job
func (j *CleanupJob) Stop() {
	j.cancel()
	log.Println("🛑 Cleanup job stopped")
}

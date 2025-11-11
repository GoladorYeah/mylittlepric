package jobs

import (
	"context"
	"log/slog"
	"time"

	"mylittleprice/internal/services"
	"mylittleprice/internal/utils"
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
		utils.LogInfo(j.ctx, "running initial cleanup")
		j.runCleanup()

		// Then run on schedule
		for {
			select {
			case <-ticker.C:
				j.runCleanup()
			case <-j.ctx.Done():
				ticker.Stop()
				utils.LogInfo(j.ctx, "cleanup job ticker stopped")
				return
			}
		}
	}()
	utils.LogInfo(j.ctx, "cleanup job started", slog.Duration("interval", j.interval))
}

// runCleanup executes the cleanup operation
func (j *CleanupJob) runCleanup() {
	startTime := time.Now()
	count, err := j.searchHistoryService.CleanupExpiredAnonymousHistory(j.ctx)
	duration := time.Since(startTime)

	if err != nil {
		utils.LogError(j.ctx, "cleanup job failed", err,
			slog.Duration("duration", duration),
		)
	} else {
		if count > 0 {
			utils.LogInfo(j.ctx, "cleanup job completed",
				slog.Duration("duration", duration),
				slog.Int64("records_deleted", count),
			)
		} else {
			utils.LogInfo(j.ctx, "cleanup job completed - no expired records found",
				slog.Duration("duration", duration),
			)
		}
	}
}

// Stop gracefully stops the cleanup job
func (j *CleanupJob) Stop() {
	j.cancel()
	utils.LogInfo(j.ctx, "cleanup job stopped")
}

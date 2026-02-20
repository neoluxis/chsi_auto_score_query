package service

import (
	"time"

	"chsi-auto-score-query/internal/logger"
	"chsi-auto-score-query/internal/repo"
	"chsi-auto-score-query/pkg/config"

	"gorm.io/gorm"
)

type Scheduler struct {
	db           *gorm.DB
	userRepo     *repo.UserRepo
	queryService *QueryService
	interval     time.Duration
	stopChan     chan struct{}
	isRunning    bool
}

func NewScheduler(db *gorm.DB, cfg *config.Config) *Scheduler {
	return &Scheduler{
		db:           db,
		userRepo:     repo.NewUserRepo(db),
		queryService: NewQueryService(cfg),
		interval:     time.Duration(cfg.QueryInterval) * time.Second,
		stopChan:     make(chan struct{}),
		isRunning:    false,
	}
}

// Start begins the background query scheduler
func (s *Scheduler) Start() {
	if s.isRunning {
		logger.Warn("Scheduler is already running")
		return
	}

	s.isRunning = true
	logger.Info("Background scheduler started with interval: %v", s.interval)

	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		// Run first query immediately
		s.queryPendingUsers()

		for {
			select {
			case <-ticker.C:
				s.queryPendingUsers()
			case <-s.stopChan:
				logger.Info("Background scheduler stopped")
				s.isRunning = false
				return
			}
		}
	}()
}

// Stop stops the background query scheduler
func (s *Scheduler) Stop() {
	if !s.isRunning {
		logger.Warn("Scheduler is not running")
		return
	}
	close(s.stopChan)
}

// queryPendingUsers queries scores for all pending users
func (s *Scheduler) queryPendingUsers() {
	logger.Info("=== Starting background score query batch ===")

	// Get all users with pending scores
	users, err := s.userRepo.FindPending()
	if err != nil {
		logger.Error("Failed to fetch pending users: %v", err)
		return
	}

	if len(users) == 0 {
		logger.Debug("No pending users to query")
		return
	}

	logger.Info("Found %d pending user(s) to process", len(users))

	successCount := 0
	failureCount := 0

	// Query score for each user
	for i, user := range users {
		logger.Info("[%d/%d] Processing user: %s (%s)", i+1, len(users), user.Name, user.Email)

		// Query and email result
		if err := s.queryService.QueryAndEmail(&user); err != nil {
			failureCount++
			logger.Error("     ❌ Query result: Failed - %v", err)
			// Update notice field with error message
			user.Notice = "查询失败：" + err.Error()
		} else {
			successCount++
			// Mark user as queried by setting LastQueryAt
			user.LastQueryAt = time.Now()
			if user.Score != "" {
				logger.Info("     ✅ Query result: Score found and email sent")
			} else {
				logger.Info("     ⏳ Query result: Score not yet available or pending")
			}
		}

		// Update user record
		if err := s.userRepo.Update(&user); err != nil {
			logger.Error("     ⚠️  Failed to update user record: %v", err)
		}
	}

	logger.Info("=== Background score query batch completed [Success: %d, Failed: %d, Total: %d] ===",
		successCount, failureCount, len(users))
}

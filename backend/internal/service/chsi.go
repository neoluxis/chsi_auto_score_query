package service

import (
"chsi-auto-score-query/internal/logger"
"chsi-auto-score-query/internal/model"
"chsi-auto-score-query/pkg/config"
)

type QueryService struct {
	chsiClient *ChsiClient
	emailSvc   *EmailService
	cfg        *config.Config
}

func NewQueryService(cfg *config.Config) *QueryService {
	return &QueryService{
		chsiClient: NewChsiClient(cfg),
		emailSvc:   NewEmailService(cfg),
		cfg:        cfg,
	}
}

// QueryAndEmail performs login, query, and email operations
func (s *QueryService) QueryAndEmail(user *model.User) error {
	logger.Info("Starting score query for user: %s", user.Email)

	// Step 1: Login
	if err := s.chsiClient.Login(); err != nil {
		logger.Error("Login failed for user %s: %v", user.Email, err)
		return s.emailSvc.SendError(user.Email, "登录学信网失败，请稍后重试")
	}

	// Step 2: Query score
	htmlContent, err := s.chsiClient.QueryScore(user)
	if err != nil {
		logger.Error("Query failed for user %s: %v", user.Email, err)
		return s.emailSvc.SendError(user.Email, "查询成绩失败，请确保信息正确")
	}

	// Step 3: Parse score
	score, err := s.chsiClient.ParseScore(htmlContent)
	if err != nil {
		logger.Error("Parse failed for user %s: %v", user.Email, err)
		return nil // Not an error if score doesn't exist yet
	}

	if score == "" {
		logger.Info("Score not available yet for user %s", user.Email)
		return nil
	}

	// Step 4: Send email
	if err := s.emailSvc.SendScore(user.Email, user.Name, score); err != nil {
		logger.Error("Email send failed for user %s: %v", user.Email, err)
		return err
	}

	logger.Info("Successfully completed score query and email for user: %s", user.Email)
	return nil

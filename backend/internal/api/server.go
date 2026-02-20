package api

import (
"fmt"
"net/http"

"chsi-auto-score-query/internal/logger"
"chsi-auto-score-query/internal/repo"
"chsi-auto-score-query/internal/service"
"chsi-auto-score-query/pkg/config"
"gorm.io/gorm"
)

type Server struct {
	cfg       *config.Config
	db        *gorm.DB
	userRepo  *repo.UserRepo
	scheduler *service.Scheduler
	mux       *http.ServeMux
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{
		cfg:       cfg,
		db:        db,
		userRepo:  repo.NewUserRepo(db),
		scheduler: service.NewScheduler(db, cfg),
		mux:       http.NewServeMux(),
	}
}

func (s *Server) Start() error {
	s.registerRoutes()

	// Start background scheduler
	s.scheduler.Start()

	addr := fmt.Sprintf(":%s", s.cfg.Port)
	logger.Info("Server listening on %s", addr)

	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) Stop() {
	s.scheduler.Stop()
	logger.Info("Server stopped")
}

func (s *Server) registerRoutes() {
	// API routes
	s.mux.HandleFunc("GET /", s.handleIndex)
	s.mux.HandleFunc("POST /api/submit", s.handleSubmit)
	s.mux.HandleFunc("GET /api/score/{email}", s.handleQueryScore)
	s.mux.HandleFunc("GET /api/health", s.handleHealth)
}

package api

import (
"crypto/md5"
"encoding/json"
"fmt"
"io"
"net/http"
"time"

"chsi-auto-score-query/internal/logger"
"chsi-auto-score-query/internal/model"
)

type SubmitRequest struct {
	Name       string `json:"name"`
	IDCard     string `json:"id_card"`
	ExamID     string `json:"exam_id"`
	Email      string `json:"email"`
	SchoolCode string `json:"school_code"`
}

type ScoreResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "CHSI Auto Score Query System"})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req SubmitRequest
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Error("Failed to parse request body: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// 验证必填字段
	if req.Name == "" || req.IDCard == "" || req.ExamID == "" || req.Email == "" || req.SchoolCode == "" {
		respondError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// 生成InfoHash以防止重复
	hash := md5.Sum([]byte(req.Name + ":" + req.IDCard + ":" + req.ExamID))
	infoHash := fmt.Sprintf("%x", hash)

	// 创建用户记录
	user := &model.User{
		Name:       req.Name,
		IDCard:     req.IDCard,
		ExamID:     req.ExamID,
		Email:      req.Email,
		SchoolCode: req.SchoolCode,
		InfoHash:   infoHash,
		CreatedAt:  time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("Failed to create user: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to save user info")
		return
	}

	logger.Info("User submitted: %s (%s)", req.Name, req.Email)
	respondSuccess(w, map[string]interface{}{
"user_id": user.ID,
"message": "Your info has been submitted. We'll send you the score when available.",
})
}

func (s *Server) handleQueryScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	email := r.PathValue("email")
	if email == "" {
		respondError(w, http.StatusBadRequest, "Email is required")
		return
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		logger.Error("Failed to query user: %v", err)
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if user == nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	respondSuccess(w, map[string]interface{}{
"name":       user.Name,
"email":      user.Email,
"score":      user.Score,
"notice":     user.Notice,
"query_time": user.LastQueryAt,
})
}

func respondSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := ScoreResponse{Code: 0, Message: "success", Data: data}
	json.NewEncoder(w).Encode(resp)
}

func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := ScoreResponse{Code: code, Message: message}
	json.NewEncoder(w).Encode(resp)
}

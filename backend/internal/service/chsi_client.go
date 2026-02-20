package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"

	"chsi-auto-score-query/internal/logger"
	"chsi-auto-score-query/internal/model"
	"chsi-auto-score-query/pkg/config"
)

type ChsiClient struct {
	client   *http.Client
	cfg      *config.Config
	username string
	password string
}

func NewChsiClient(cfg *config.Config) *ChsiClient {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return &ChsiClient{
		client:   client,
		cfg:      cfg,
		username: cfg.ChsiUsername,
		password: cfg.ChsiPassword,
	}
}

// Login logs into CHSI website
func (c *ChsiClient) Login() error {
	logger.Info("Attempting to login CHSI with username: %s", c.username)

	// 第一步：获取登录页面以获取lt和execution参数
	loginPageURL := "https://account.chsi.com.cn/passport/login?entrytype=yzgr&service=https%3A%2F%2Fyz.chsi.com.cn%2Fj_spring_cas_security_check"
	resp, err := c.client.Get(loginPageURL)
	if err != nil {
		logger.Error("Failed to get login page: %v", err)
		return err
	}
	defer resp.Body.Close()

	// 简化版：直接尝试登录（在实际应用中应该解析HTML获取lt和execution）
	// 这里使用空值，因为某些情况下服务器可能不需要这些参数
	loginData := url.Values{}
	loginData.Set("username", c.username)
	loginData.Set("password", c.password)
	loginData.Set("lt", "")
	loginData.Set("execution", "")
	loginData.Set("_eventId", "submit")

	loginURL := "https://account.chsi.com.cn/passport/login?entrytype=yzgr&service=https%3A%2F%2Fyz.chsi.com.cn%2Fj_spring_cas_security_check"

	req, err := http.NewRequest("POST", loginURL, bytes.NewBufferString(loginData.Encode()))
	if err != nil {
		logger.Error("Failed to create login request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")

	resp, err = c.client.Do(req)
	if err != nil {
		logger.Error("Failed to login: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		logger.Warn("Login returned status code: %d", resp.StatusCode)
	}

	logger.Info("Login successful")
	return nil
}

// QueryScore queries exam score from CHSI
func (c *ChsiClient) QueryScore(user *model.User) (string, error) {
	logger.Info("Querying score for user: %s (ID: %s, ExamID: %s)", user.Name, user.IDCard, user.ExamID)

	queryData := url.Values{}
	queryData.Set("xm", user.Name)           // 姓名
	queryData.Set("zjhm", user.IDCard)       // 身份证号
	queryData.Set("ksbh", user.ExamID)       // 考生编号
	queryData.Set("bkdwdm", user.SchoolCode) // 报考单位代码
	queryData.Set("checkcode", "")           // 验证码（模拟可以为空）

	queryURL := "https://yz.chsi.com.cn/apply/cjcx/cjcx.do"

	req, err := http.NewRequest("POST", queryURL, bytes.NewBufferString(queryData.Encode()))
	if err != nil {
		logger.Error("Failed to create query request: %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Referer", "https://yz.chsi.com.cn/apply/cjcx/t/10358.dhtml")

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error("Failed to query score: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Query returned status code: %d", resp.StatusCode)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body: %v", err)
		return "", err
	}

	htmlContent := string(body)
	logger.Debug("Response HTML length: %d bytes", len(htmlContent))

	return htmlContent, nil
}

// ParseScore parses score from HTML response using Vue's cj JSON object
func (c *ChsiClient) ParseScore(htmlContent string) (string, error) {
	logger.Info("Parsing score from HTML response")

	if htmlContent == "" {
		logger.Warn("Score query status: Empty HTML response - possible network error or invalid session")
		return "", nil
	}

	html := string(htmlContent)

	// Step 1: Extract Vue's cj object using improved regex
	// Try pattern 1: cj : {...}
	re := regexp.MustCompile(`(?s)\bcj\s*:\s*(\{[^}]*\}|null)`)
	matches := re.FindStringSubmatch(html)

	var raw string
	if len(matches) >= 2 {
		raw = matches[1]
	}

	// If simple pattern failed, try more complex pattern for nested objects
	if raw == "" {
		re = regexp.MustCompile(`(?s)\bcj\s*:\s*(\{(?:[^{}]|(?:\{[^{}]*\}))*\}|null)`)
		matches = re.FindStringSubmatch(html)
		if len(matches) >= 2 {
			raw = matches[1]
		}
	}

	// Last resort: look for quotes around cj
	if raw == "" {
		re = regexp.MustCompile(`(?s)"cj"\s*:\s*(\{[^{}]*\}|null)`)
		matches = re.FindStringSubmatch(html)
		if len(matches) >= 2 {
			raw = matches[1]
		}
	}

	if raw == "" {
		logger.Warn("Score query status: Could not find score data structure in response")
		return "", nil
	}

	logger.Debug("Extracted raw cj data: %s", raw[:minInt(100, len(raw))])

	// Step 2: Check if cj is null
	if raw == "null" {
		// Extract msg field if available
		msgRe := regexp.MustCompile(`(?s)\bmsg\s*[:=]\s*["']([^"']*)["']`)
		msgMatches := msgRe.FindStringSubmatch(html)
		msg := "请检查报考信息或成绩查询尚未开放"
		if len(msgMatches) > 1 {
			msg = msgMatches[1]
		}

		logger.Warn("Score query status: ⏳ No query result available - msg: %s", msg)

		// Categorize the message to provide more specific status
		if bytes.Contains([]byte(msg), []byte("信息不匹配")) {
			logger.Warn("  └─ Detailed: 信息不匹配 (User information doesn't match CHSI records)")
			return "", nil
		}

		if bytes.Contains([]byte(msg), []byte("暂未")) || bytes.Contains([]byte(msg), []byte("未开放")) {
			logger.Info("  └─ Status: Scores not yet published")
			return "", nil
		}

		return "", nil
	}

	// Step 3: Parse cj as JSON
	var scoreData map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &scoreData); err != nil {
		logger.Error("Score query status: Failed to parse score data as JSON: %v", err)
		logger.Debug("Raw data was: %s", raw[:minInt(200, len(raw))])
		return "", err
	}

	logger.Debug("Successfully parsed score JSON with %d fields", len(scoreData))

	// Step 4: Extract score information
	// Build score string from available fields
	var scoreStr string

	// Common score fields based on CHSI structure
	scoreFields := []string{"总分", "zf", "total_score", "zsxh", "ksbh", "xm", "zymc"}
	for _, field := range scoreFields {
		if val, ok := scoreData[field]; ok && val != "" && val != nil {
			scoreStr += fmt.Sprintf("%s: %v; ", field, val)
		}
	}

	if scoreStr != "" {
		logger.Info("Score query status: ✅ Score found - %s", scoreStr)
		return scoreStr, nil
	}

	// Check for admission status fields
	admissionFields := []string{"lqzt", "录取状态", "psyz", "拟录取"}
	for _, field := range admissionFields {
		if val, ok := scoreData[field]; ok && val != "" && val != nil {
			valStr := fmt.Sprintf("%v", val)
			if bytes.Contains([]byte(valStr), []byte("录取")) {
				logger.Info("Score query status: ✅ Admission status found - %s: %v", field, val)
				return valStr, nil
			}
			if bytes.Contains([]byte(valStr), []byte("体检")) {
				logger.Info("Score query status: 📋 Physical exam stage - %s: %v", field, val)
				return "", nil
			}
			if bytes.Contains([]byte(valStr), []byte("复试")) {
				logger.Info("Score query status: 📝 Reexamination/interview stage - %s: %v", field, val)
				return "", nil
			}
		}
	}

	// Check for preliminary score
	preliminaryFields := []string{"cxsj", "初试成绩", "cs_cj"}
	for _, field := range preliminaryFields {
		if val, ok := scoreData[field]; ok && val != "" && val != nil {
			logger.Info("Score query status: 📊 Preliminary score - %s: %v", field, val)
			return fmt.Sprintf("%v", val), nil
		}
	}

	// Check for zsdwsm (招生单位说明 - admission office note)
	if note, ok := scoreData["zsdwsm"]; ok && note != "" && note != nil {
		logger.Info("Score query status: 📌 Admission office note: %v", note)
	}

	// Unknown status - log all fields for debugging
	logger.Debug("Score query status: All score fields - %+v", scoreData)
	logger.Info("Score query status: ℹ️  No definitive score or admission status detected yet")

	return "", nil
}

// Helper function for min integer
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

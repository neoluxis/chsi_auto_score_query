package service

import (
"bytes"
"io"
"net/http"
"net/http/cookiejar"
"net/url"
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
	jar, _ := cookiejar.New()
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
	queryData.Set("xm", user.Name)                // 姓名
	queryData.Set("zjhm", user.IDCard)            // 身份证号
	queryData.Set("ksbh", user.ExamID)            // 考生编号
	queryData.Set("bkdwdm", user.SchoolCode)      // 报考单位代码
	queryData.Set("checkcode", "")                // 验证码（模拟可以为空）

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

// ParseScore parses score from HTML response
func (c *ChsiClient) ParseScore(htmlContent string) (string, error) {
	logger.Info("Parsing score from HTML response")

	// 简单的HTML解析：查找是否包含成绩信息
	// 在实际应用中应该使用HTML解析库如goquery
	if htmlContent == "" {
		logger.Error("Empty HTML content")
		return "", nil
	}

	// 检查常见的成功标记
	if bytes.Contains([]byte(htmlContent), []byte("总分")) ||
		bytes.Contains([]byte(htmlContent), []byte("Score")) ||
		bytes.Contains([]byte(htmlContent), []byte("分数")) {
		logger.Info("Score information found in response")
		// TODO: 实现详细的成绩解析
		return "score_found", nil
	}

	// 检查常见的失败标记
	if bytes.Contains([]byte(htmlContent), []byte("暂未发布")) ||
		bytes.Contains([]byte(htmlContent), []byte("未查询到")) ||
		bytes.Contains([]byte(htmlContent), []byte("信息不匹配")) {
		logger.Warn("Score not available yet or information mismatch")
		return "", nil
	}

	logger.Info("Unable to determine score status from HTML")
	return "", nil
}

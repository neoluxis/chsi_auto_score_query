package service

import (
"fmt"
"net/smtp"

"chsi-auto-score-query/internal/logger"
"chsi-auto-score-query/pkg/config"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

// SendScore sends exam score to user email
func (s *EmailService) SendScore(toEmail string, name string, score string) error {
	logger.Info("Preparing to send score email to: %s", toEmail)

	// Build email content
	subject := "考研成绩已发布"
	body := fmt.Sprintf(`<html><body>
<h2>尊敬的 %s：</h2>
<p>您的考研成绩已发布，请登录学信网查看详情。</p>
<p><strong>成绩信息：</strong> %s</p>
<p>祝贺您！</p>
<p>此邮件由自动查询系统发送，请勿回复。</p>
</body></html>`, name, score)

	return s.sendSMTPEmail(toEmail, subject, body)
}

// SendError sends error notification to user
func (s *EmailService) SendError(toEmail string, errMsg string) error {
	logger.Info("Preparing to send error email to: %s", toEmail)

	subject := "成绩查询失败通知"
	body := fmt.Sprintf(`<html><body>
<h2>成绩查询失败</h2>
<p>您的考研成绩查询失败，原因如下：</p>
<p><strong>错误信息：</strong> %s</p>
<p>请检查您的个人信息是否正确，或稍后重试。</p>
<p>此邮件由自动查询系统发送，请勿回复。</p>
</body></html>`, errMsg)

	return s.sendSMTPEmail(toEmail, subject, body)
}

// sendSMTPEmail sends email via SMTP
func (s *EmailService) sendSMTPEmail(toEmail string, subject string, body string) error {
	if s.cfg.SMTPUser == "" || s.cfg.SMTPPass == "" {
		logger.Warn("SMTP configuration incomplete, skipping email send to %s", toEmail)
		return nil
	}

	// Build email message
	message := fmt.Sprintf("To: %s\r\n"+
"Subject: %s\r\n"+
"Content-Type: text/html; charset=UTF-8\r\n"+
"\r\n"+
"%s", toEmail, subject, body)

	// Setup SMTP server
	host := s.cfg.SMTPServer
	port := fmt.Sprintf(":%d", s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPServer)

	// Send email
	logger.Debug("Sending email to %s via %s:%d", toEmail, host, s.cfg.SMTPPort)
	err := smtp.SendMail(
host+port,
auth,
s.cfg.SMTPUser,
[]string{toEmail},
[]byte(message),
)

	if err != nil {
		logger.Error("Failed to send email to %s: %v", toEmail, err)
		return err
	}

	logger.Info("Email sent successfully to: %s", toEmail)
	return nil
}

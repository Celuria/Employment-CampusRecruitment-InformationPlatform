package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/employment-center/campus-recruitment/config"
)

// Sender 邮件发送器
type Sender struct {
	cfg      config.EmailConfig
	tlsCfg   *tls.Config
}

// NewSender 创建邮件发送器
func NewSender(cfg config.EmailConfig) *Sender {
	return &Sender{
		cfg: cfg,
		tlsCfg: &tls.Config{
			ServerName:         cfg.SMTPHost,
			InsecureSkipVerify: false,
		},
	}
}

// Enabled 返回邮件功能是否已启用
func (s *Sender) Enabled() bool {
	return s.cfg.Enabled
}

// Send 发送邮件。subject 为主题，body 为纯文本正文，to 为收件人列表
func (s *Sender) Send(to []string, subject, body string) error {
	if !s.cfg.Enabled {
		return nil
	}

	from := s.cfg.Username
	fromName := s.cfg.FromName
	if fromName == "" {
		fromName = "就业中心"
	}

	// 拼接符合 RFC 822 的邮件内容
	msg := buildMessage(from, fromName, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.SMTPHost)

	// 先尝试 STARTTLS
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.SMTPHost)
	if err != nil {
		return fmt.Errorf("smtp new client: %w", err)
	}
	defer client.Quit()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(s.tlsCfg); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}

	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return fmt.Errorf("smtp rcpt to %s: %w", addr, err)
		}
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	_, err = fmt.Fprint(wc, msg)
	if err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	return wc.Close()
}

func buildMessage(from, fromName string, to []string, subject, body string) string {
	sb := &strings.Builder{}
	sb.WriteString("From: ")
	if fromName != "" {
		sb.WriteString(fmt.Sprintf("%s <%s>", fromName, from))
	} else {
		sb.WriteString(from)
	}
	sb.WriteString("\r\n")
	sb.WriteString("To: ")
	sb.WriteString(strings.Join(to, ", "))
	sb.WriteString("\r\n")
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(body)
	return sb.String()
}

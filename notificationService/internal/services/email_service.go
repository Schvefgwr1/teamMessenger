package services

import (
	"bytes"
	"common/config"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"path/filepath"

	"common/models"
)

type EmailService struct {
	config    *config.EmailConfig
	templates map[models.NotificationType]*template.Template
	dialer    *gomail.Dialer
}

func NewEmailService(cfg *config.EmailConfig) (*EmailService, error) {
	// Валидация конфигурации
	if cfg.SMTPHost == "" {
		return nil, fmt.Errorf("SMTP host is required")
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("SMTP username is required - set SMTP_USERNAME environment variable")
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("SMTP password is required - set SMTP_PASSWORD environment variable")
	}

	dialer := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.Username, cfg.Password)

	service := &EmailService{
		config:    cfg,
		templates: make(map[models.NotificationType]*template.Template),
		dialer:    dialer,
	}

	// Загружаем шаблоны
	if err := service.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}

	return service, nil
}

func (s *EmailService) loadTemplates() error {
	templateFiles := map[models.NotificationType]string{
		models.NotificationNewTask: "new_task.html",
		models.NotificationNewChat: "new_chat.html",
		models.NotificationLogin:   "login.html",
	}

	for notificationType, filename := range templateFiles {
		path := filepath.Join(s.config.TemplatePath, filename)
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", filename, err)
		}
		s.templates[notificationType] = tmpl
	}

	return nil
}

func (s *EmailService) SendNotification(notification interface{}) error {
	var (
		email        string
		subject      string
		templateData interface{}
		tmplType     models.NotificationType
	)

	switch n := notification.(type) {
	case *models.NewTaskNotification:
		email = n.Email
		subject = fmt.Sprintf("Новая задача: %s", n.TaskTitle)
		templateData = n
		tmplType = models.NotificationNewTask

	case *models.NewChatNotification:
		email = n.Email
		if n.IsGroup {
			subject = fmt.Sprintf("Вы добавлены в группу: %s", n.ChatName)
		} else {
			subject = fmt.Sprintf("Новый чат с %s", n.CreatorName)
		}
		templateData = n
		tmplType = models.NotificationNewChat

	case *models.LoginNotification:
		email = n.Email
		subject = "Вход в аккаунт TeamMessenger"
		templateData = n
		tmplType = models.NotificationLogin

	default:
		return fmt.Errorf("unknown notification type: %T", notification)
	}

	// Генерируем HTML содержимое
	htmlBody, err := s.generateHTML(tmplType, templateData)
	if err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Создаем и отправляем email
	return s.sendEmail(email, subject, htmlBody)
}

func (s *EmailService) generateHTML(notificationType models.NotificationType, data interface{}) (string, error) {
	tmpl, exists := s.templates[notificationType]
	if !exists {
		return "", fmt.Errorf("template not found for notification type: %s", notificationType)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (s *EmailService) sendEmail(to, subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	return nil
}

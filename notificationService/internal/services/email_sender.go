package services

import "gopkg.in/gomail.v2"

// EmailSender интерфейс для отправки email (создан для мокирования в тестах)
type EmailSender interface {
	DialAndSend(m ...*gomail.Message) error
}

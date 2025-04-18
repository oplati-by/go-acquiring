package oacquiring

import "fmt"

type (
	// ServerError - ошибка сервера Оплати
	ServerError struct {
		StatusCode   string // Код ошибки
		InternalCode string // Внутренний код ошибки
		Message      string // Сообщение
		UserMessage  string // Сообщение для пользователя
	}
)

func (s *ServerError) Error() string {
	return fmt.Sprintf("OPLATI error %s: %s", s.InternalCode, s.Message)
}

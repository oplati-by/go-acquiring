package oacquiring

import "net/http"

type (
	// Client - клиент для использования API. Для инициализации используйте NewClient
	Client struct {
		baseUrl string

		cashboxRegNumber string
		cashboxPassword  string

		httpClient http.Client
	}
)

// NewClient возвращает новый Client.
//   - baseUrl - Базовый URL сервера Оплати, например https://oplati-cashboxapi.lwo-dev.by/ms-pay
//   - cashboxRegNumber - Регистрационный номер кассы, например OPL000011111
//   - cashboxPassword - Пароль для интернет-кассы
//   - opts - Дополнительные настройки: WithCustomHTTPClient
func NewClient(baseUrl, cashboxRegNumber, cashboxPassword string, opts ...ClientOpt) Client {
	c := Client{
		baseUrl:          baseUrl,
		cashboxRegNumber: cashboxRegNumber,
		cashboxPassword:  cashboxPassword,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

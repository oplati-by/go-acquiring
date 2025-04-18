package oacquiring

import "net/http"

type (
	// ClientOpt - дополнительные параметры Client
	ClientOpt func(*Client)
)

// WithCustomHTTPClient - позволяет переопределить http.Client, используемый для отправки запросов к серверу Оплати
func WithCustomHTTPClient(client http.Client) ClientOpt {
	return func(c *Client) {
		c.httpClient = client
	}
}

package oacquiring

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type (
	// PaymentNotificationHandler - интерфейс для обработки данных платежа, полученного от системы Оплати.
	PaymentNotificationHandler interface {
		// HandlePayment должен обработать платеж, полученный от системы Оплати. Это может быть обновление данных в БД,
		// отправка события и т.п. Если обработка завершилась с ошибкой, метод должен вернуть эту ошибку. В
		// этом случае серверу Оплати будет отправлен ответ с кодом "500 Internal Server Error", и через некоторое время
		// сервер Оплати сделает повторный запрос.
		HandlePayment(PaymentInfo) error
	}

	// HTTPNotificationHandler - обработчик HTTP уведомления от сервера Оплати, реализует интерфейс
	// http.Handler. Осуществляет:
	//  1. Проверку подписи Server-Sign. В случае, если запрос подписан неверно, клиент получит ответ "401 Unauthorized"
	//  2. Преобразования тела запроса в PaymentInfo. В случае, если получен некорректный json, клиент получит
	//  ответ "400 Bad Request"
	//  3. Выполнение логики PaymentNotificationHandler.HandlePayment с корректным PaymentInfo
	//  4. Отправка ответа клиенту в зависимости от успеха выполнения шага 3
	// Для инициализации используйте NewHTTPNotificationHandler.
	HTTPNotificationHandler struct {
		publicKey *rsa.PublicKey
		handler   PaymentNotificationHandler
	}
)

// NewHTTPNotificationHandler возвращает новый HTTPNotificationHandler для получения HTTP уведомлений от сервера Оплати.
//   - publicKey - Публичный ключ, используемый для проверки подписи Server-Sign, полученный в личном кабинете Оплати.Бизнес.
//   - paymentHandler - обработчик для выполнения каких-либо действий с полученным платежом.
func NewHTTPNotificationHandler(publicKey string, paymentHandler PaymentNotificationHandler) (HTTPNotificationHandler, error) {
	rawKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return HTTPNotificationHandler{}, fmt.Errorf("public key base64 decoding failed: %w", err)
	}

	key, err := x509.ParsePKIXPublicKey(rawKey)
	if err != nil {
		return HTTPNotificationHandler{}, fmt.Errorf("public key parsing failed: %w", err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return HTTPNotificationHandler{}, errors.New("provided public key is not RSA")
	}

	if paymentHandler == nil {
		return HTTPNotificationHandler{}, errors.New("nil handler is not allowed")
	}

	return HTTPNotificationHandler{
		publicKey: rsaKey,
		handler:   paymentHandler,
	}, nil
}

func (nh *HTTPNotificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decodedSignature, err := base64.StdEncoding.DecodeString(r.Header.Get("Server-Sign"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sum := sha256.Sum256(body)

	err = rsa.VerifyPKCS1v15(nh.publicKey, crypto.SHA256, sum[:], decodedSignature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var rawPaymentInfo paymentInfoResponse
	err = json.Unmarshal(body, &rawPaymentInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paymentInfo, err := makePaymentInfoFromRaw(rawPaymentInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = nh.handler.HandlePayment(paymentInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

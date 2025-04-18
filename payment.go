package oacquiring

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// PaymentItemTypeProduct - Товар
	PaymentItemTypeProduct PaymentItemType = 1
	// PaymentItemTypeService - Услуга
	PaymentItemTypeService PaymentItemType = 2
)

type (
	// PaymentItemType - тип позиции в чеке. Варианты PaymentItemTypeService, PaymentItemTypeProduct
	PaymentItemType int

	// Payment - данные, необходимые для создания платежа
	Payment struct {
		Shift             string        // Смена. Например дата в формате ДДММГГГГ
		OrderNumber       string        // Уникальный номер заказа
		Items             []PaymentItem // Список позиций в чеке
		ReceiptFooterText string        // Дополнительная информация в конце чека
		SuccessUrl        string        // URL для перехода после успешной оплаты
		FailureUrl        string        // URL для перехода после неуспешной оплаты
		NotificationUrl   string        // URL для отправки уведомления об оплате со стороны Оплати
	}

	// PaymentItem - Позиция в чеке
	PaymentItem struct {
		Type PaymentItemType // Тип (товар/услуга)
		Name string          // Наименование
		Cost int64           // Стоимость в копейках. Например, 545 ~ 5.45 BYN
	}

	// SuccessfulPayment - Результат успешного создания платежа
	SuccessfulPayment struct {
		PaymentId   int64  // Идентификатор платежа в Оплати
		RedirectUrl string // URL страницы для оплаты
	}
)

// CreatePayment - создание платежа на стороне Оплати. Возвращает уникальный номер платежа в системе Оплати и URL для
// выполнения оплаты. Используется запрос POST /pos/webPayments/v2.
//
// В случае, если сервер вернул ответ отличный от 200 OK, возвращаемый error можно попробовать привести к *ServerError
// для получения дополнительных данных об ошибке.
func (a *Client) CreatePayment(ctx context.Context, payment Payment) (SuccessfulPayment, error) {
	request := a.makePaymentRequest(payment)

	body, err := json.Marshal(&request)
	if err != nil {
		return SuccessfulPayment{}, fmt.Errorf("request encoding failed: %w", err)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseUrl+"/pos/webPayments/v2", bytes.NewReader(body))
	if err != nil {
		return SuccessfulPayment{}, fmt.Errorf("request initialization failed: %w", err)
	}

	r.Header.Set("RegNum", a.cashboxRegNumber)
	r.Header.Set("Password", a.cashboxPassword)
	r.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(r)
	if err != nil {
		return SuccessfulPayment{}, fmt.Errorf("request execution failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return SuccessfulPayment{}, fmt.Errorf("decoding error failed: %w (status code %d)", err, resp.StatusCode)
		}
		return SuccessfulPayment{}, &ServerError{
			StatusCode:   errResp.Code,
			InternalCode: errResp.InternalCode,
			Message:      errResp.DevMessage,
			UserMessage:  errResp.UserMessage.LangRu,
		}
	}

	var successfulPayment newPaymentResponse
	err = json.NewDecoder(resp.Body).Decode(&successfulPayment)
	if err != nil {
		return SuccessfulPayment{}, fmt.Errorf("decoding response failed: %w", err)
	}

	return SuccessfulPayment{
		PaymentId:   successfulPayment.PaymentId,
		RedirectUrl: successfulPayment.RedirectUrl,
	}, nil
}

package oacquiring

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type (
	// PaymentReversal - данные, необходимые для полного или частичного возврата платежа
	PaymentReversal struct {
		Shift             string        // Смена. Например дата в формате ДДММГГГГ
		OrderNumber       string        // Уникальный номер заказа
		Items             []PaymentItem // Список позиций в чеке
		ReceiptFooterText string        // Дополнительная информация в конце чека
	}
)

// ReversePayment - Выполнение отмены (полной или частичной) любой операции, выполненной в течение периода
// предопределенного системой Оплати. Используется запрос POST /pos/payments/{paymentId}/reversals.
//
// В случае, если сервер вернул ответ отличный от 200 OK, возвращаемый error можно попробовать привести к *ServerError
// для получения дополнительных данных об ошибке.
func (a *Client) ReversePayment(ctx context.Context, paymentId int64, payment PaymentReversal) (PaymentInfo, error) {
	request := a.makeReversePaymentRequest(payment)

	body, err := json.Marshal(&request)
	if err != nil {
		return PaymentInfo{}, fmt.Errorf("request encoding failed: %w", err)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseUrl+"/pos/payments/"+strconv.FormatInt(paymentId, 10)+"/reversals", bytes.NewReader(body))
	if err != nil {
		return PaymentInfo{}, fmt.Errorf("request initialization failed: %w", err)
	}

	r.Header.Set("RegNum", a.cashboxRegNumber)
	r.Header.Set("Password", a.cashboxPassword)
	r.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(r)
	if err != nil {
		return PaymentInfo{}, fmt.Errorf("request execution failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return PaymentInfo{}, fmt.Errorf("decoding error failed: %w (status code %d)", err, resp.StatusCode)
		}
		return PaymentInfo{}, &ServerError{
			StatusCode:   errResp.Code,
			InternalCode: errResp.InternalCode,
			Message:      errResp.DevMessage,
			UserMessage:  errResp.UserMessage.LangRu,
		}
	}

	var rawPaymentInfo paymentInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&rawPaymentInfo)
	if err != nil {
		return PaymentInfo{}, fmt.Errorf("decoding response failed: %w", err)
	}

	paymentInfo, err := makePaymentInfoFromRaw(rawPaymentInfo)
	if err != nil {
		return PaymentInfo{}, fmt.Errorf("handling response failed: %w", err)
	}

	return paymentInfo, nil
}

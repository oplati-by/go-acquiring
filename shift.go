package oacquiring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetPaymentsOnShift - Получение списка платежей для сверки итогов по смене. Используется запрос GET /pos/paymentReports.
//
// В случае, если сервер вернул ответ отличный от 200 OK, возвращаемый error можно попробовать привести к *ServerError
// для получения дополнительных данных об ошибке.
func (a *Client) GetPaymentsOnShift(ctx context.Context, shift string) ([]PaymentInfo, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseUrl+"/pos/paymentReports?shift="+shift, nil)
	if err != nil {
		return nil, fmt.Errorf("request initialization failed: %w", err)
	}

	r.Header.Set("RegNum", a.cashboxRegNumber)
	r.Header.Set("Password", a.cashboxPassword)

	resp, err := a.httpClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, fmt.Errorf("decoding error failed: %w (status code %d)", err, resp.StatusCode)
		}
		return nil, &ServerError{
			StatusCode:   errResp.Code,
			InternalCode: errResp.InternalCode,
			Message:      errResp.DevMessage,
			UserMessage:  errResp.UserMessage.LangRu,
		}
	}

	var rawPayments []paymentInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&rawPayments)
	if err != nil {
		return nil, fmt.Errorf("decoding response failed: %w", err)
	}

	payments := make([]PaymentInfo, len(rawPayments))
	for i, rawPayment := range rawPayments {
		payments[i], err = makePaymentInfoFromRaw(rawPayment)
		if err != nil {
			return nil, fmt.Errorf("handling response failed: %w", err)
		}
	}

	return payments, nil
}

package oacquiring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type (
	// PaymentType - Тип кассовой операции. Варианты: PaymentTypeSell, PaymentTypeBuy, PaymentItemTypeSellReverse,
	// PaymentItemTypeBuyReverse
	PaymentType int

	// PaymentStatus - Статус платежа. Варианты: PaymentStatusInProgress, PaymentStatusDone, PaymentStatusDeclined,
	// PaymentStatusNotEnoughMoney, PaymentStatusTimeout, PaymentStatusTechCancel
	PaymentStatus int

	// PaymentInfo - Информация о платеже в системе Оплати
	PaymentInfo struct {
		Id            int64         // Идентификатор платежа в Оплати
		Type          PaymentType   // Тип платежа
		Sum           int64         // Сумма в копейках. Например, 545 ~ 5.45 BYN
		Status        PaymentStatus // Статус платежа
		CreatedDate   time.Time     // Дата создания платежа
		PaidDate      time.Time     // Дата выполнения оплаты. Для новых платежей совпадает с CreatedDate
		OrderNumber   string        // Уникальный номер заказа
		PursePublicId string        // Публичный идентификатор кошелька. Может быть указан для платежей со статусом PaymentStatusDone
	}
)

// GetPaymentInfo - Получение информации о платеже в системе Оплати. Используется запрос GET /pos/payments/{paymentId}.
//
// В случае, если сервер вернул ответ отличный от 200 OK, возвращаемый error можно попробовать привести к *ServerError
// для получения дополнительных данных об ошибке.
func (a *Client) GetPaymentInfo(ctx context.Context, paymentId int64) (PaymentInfo, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseUrl+"/pos/payments/"+strconv.FormatInt(paymentId, 10), nil)
	if err != nil {
		return PaymentInfo{}, fmt.Errorf("request initialization failed: %w", err)
	}

	r.Header.Set("RegNum", a.cashboxRegNumber)
	r.Header.Set("Password", a.cashboxPassword)

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

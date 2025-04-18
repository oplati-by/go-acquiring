package oacquiring

import (
	"fmt"
	"math"
	"time"
)

func makePaymentItems(items []PaymentItem) ([]paymentRequestDetailsItem, float64) {
	var sum int64
	paymentItems := make([]paymentRequestDetailsItem, len(items))

	for i, item := range items {
		paymentItems[i] = paymentRequestDetailsItem{
			Type: int(item.Type),
			Name: item.Name,
			Cost: float64(item.Cost) / 100,
		}
		sum += item.Cost
	}

	return paymentItems, float64(sum) / 100
}

func (a *Client) makePaymentRequest(payment Payment) newPaymentRequest {
	items, sum := makePaymentItems(payment.Items)

	return newPaymentRequest{
		Shift:       payment.Shift,
		Sum:         sum,
		OrderNumber: payment.OrderNumber,
		RegNum:      a.cashboxRegNumber,
		Details: paymentRequestDetails{
			RegNum:      a.cashboxRegNumber,
			Items:       items,
			AmountTotal: sum,
			FooterInfo:  payment.ReceiptFooterText,
		},
		SuccessUrl:      payment.SuccessUrl,
		FailureUrl:      payment.FailureUrl,
		NotificationUrl: payment.NotificationUrl,
	}
}

func (a *Client) makeReversePaymentRequest(payment PaymentReversal) reversePaymentRequest {
	items, sum := makePaymentItems(payment.Items)

	return reversePaymentRequest{
		Shift:       payment.Shift,
		Sum:         sum,
		OrderNumber: payment.OrderNumber,
		RegNum:      a.cashboxRegNumber,
		Details: paymentRequestDetails{
			RegNum:      a.cashboxRegNumber,
			Items:       items,
			AmountTotal: sum,
			FooterInfo:  payment.ReceiptFooterText,
		},
	}
}

func makePaymentInfoFromRaw(rawPaymentInfo paymentInfoResponse) (PaymentInfo, error) {
	paymentInfo := PaymentInfo{
		Id:            rawPaymentInfo.PaymentId,
		Type:          PaymentType(rawPaymentInfo.PaymentType),
		Sum:           int64(math.Round(rawPaymentInfo.Sum * 100)),
		Status:        PaymentStatus(rawPaymentInfo.Status),
		OrderNumber:   rawPaymentInfo.OrderNumber,
		PursePublicId: rawPaymentInfo.PursePublicId,
	}

	createdDate, err := time.Parse(time.RFC3339, rawPaymentInfo.CreatedDate)
	if err != nil {
		return PaymentInfo{}, fmt.Errorf("bad payment createdDate: %w", err)
	}
	paymentInfo.CreatedDate = createdDate

	paidDate, err := time.Parse(time.RFC3339, rawPaymentInfo.PaidDate)
	if err != nil {
		return PaymentInfo{}, fmt.Errorf("bad payment paidDate: %w", err)
	}
	paymentInfo.PaidDate = paidDate

	return paymentInfo, nil
}

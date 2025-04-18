package oacquiring

const (
	// PaymentStatusInProgress - Платеж ожидает подтверждения, которое должно быть выполнено клиентом
	// на мобильном устройстве.
	//
	// Код: IN_PROGRESS
	PaymentStatusInProgress = 0

	// PaymentStatusDone - Платеж совершен, можно выдать товар клиенту.
	//
	// Код: OK
	PaymentStatusDone = 1

	// PaymentStatusDeclined - Отказ от платежа. Клиент не подтвердил платеж.
	//
	// Код: DECLINE
	PaymentStatusDeclined = 2

	// PaymentStatusNotEnoughMoney - Недостаточно средств на кошельке клиента.
	//
	// Код: NOT_ENOUGH
	PaymentStatusNotEnoughMoney = 3

	// PaymentStatusTimeout - Клиент не подтвердил платеж в течение предопределенного системой Оплати отрезка времени.
	// Равносильно отказу от оплаты.
	//
	// Код: TIMEOUT
	PaymentStatusTimeout = 4

	// PaymentStatusTechCancel - Операция была отменена либо кассой, либо системой, когда не смогла получить информацию
	// о статусе платежа в течение предопределенного системой Оплати отрезка времени.
	//
	// Код: TECHNICAL_CANCELLING
	PaymentStatusTechCancel = 5
)

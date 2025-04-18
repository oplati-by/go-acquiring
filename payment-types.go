package oacquiring

const (
	// PaymentTypeSell - Продажа (приходная кассовая операция)
	PaymentTypeSell = 1

	// PaymentTypeBuy - Покупка (расходная кассовая операция)
	PaymentTypeBuy = 2

	// PaymentItemTypeSellReverse - Возврат продажи (расходная кассовая операция)
	PaymentItemTypeSellReverse = 3

	// PaymentItemTypeBuyReverse - Возврат покупки (приходная кассовая операция)
	PaymentItemTypeBuyReverse = 4
)

// Package oacquiring содержит инструменты для интеграции с системой платежей Оплати.
// Для взаимодействия с API используйте Client. Для получения уведомлений от Оплати - HTTPNotificationHandler.
// Примеры использования:
//
// # Инициализация клиента
//
//	import oacquiring "github.com/oplati-by/go-acquiring"
//	// ...
//	oplatiClient := oacquiring.NewClient("https://oplati-cashboxapi.lwo-dev.by/ms-pay", "OPL000011111", "1111")
//
// # Создание платежа:
//
//	paymentData := oacquiring.Payment{
//		     Shift:       "14092001",
//		     OrderNumber: "AA-1111",
//		     Items: []oacquiring.PaymentItem{
//		         {
//		             Type: oacquiring.PaymentItemTypeService,
//		             Name: "Консультация продавца",
//		             Cost: 499,
//		         },
//		         {
//		             Type: oacquiring.PaymentItemTypeProduct,
//		             Name: "Товар",
//		             Cost: 5999,
//		         },
//		     },
//		     ReceiptFooterText: "Спасибо за покупку!",
//		     SuccessUrl:        "https://my.shop.by/me/orders/AA-1111",
//		     FailureUrl:        "https://my.shop.by/payment-failed",
//		     NotificationUrl:   "https://my.shop.by/api/webhook/orders/AA-1111",
//	}
//
//	result, err := oplatiClient.CreatePayment(context.Background(), paymentData)
//	// ...
//
// # Проверка статуса платежа
//
//	paymentInfo, err := oplatiClient.GetPaymentInfo(context.Background(), 123456)
//	// ...
//
// # Отмена платежа (частичная или полная)
//
//	paymentData := oacquiring.PaymentReversal{
//	    Shift:       "14092001",
//	    OrderNumber: "AA-1111",
//	    Items: []oacquiring.PaymentItem{
//	        {
//	            Type: oacquiring.PaymentItemTypeProduct,
//	            Name: "Товар",
//	            Cost: 5999,
//	        },
//	    },
//	    ReceiptFooterText: "Будем рады видеть вас снова!",
//	}
//
//	paymentInfo, err := oplatiClient.ReversePayment(context.Background(), 123456, paymentData)
//	// ...
//
// # Получение списка продаж за смену
//
//	payments, err := oplatiClient.GetPaymentsOnShift(context.Background(), "15042025")
//	// ...
//
// # Обработка ошибок
//
//	// ...
//	if err != nil {
//	    if oplatiErr := (*oacquiring.ServerError)(nil); errors.As(err, &oplatiErr) {
//	        // Do something with oplatiErr
//	    }
//	}
//
// # Получение уведомлений от сервера Оплати
//
// Реализация интерфейса PaymentNotificationHandler:
//
//	type Handler struct {
//	    // DB connection, etc.
//	}
//
//	func (p *Handler) HandlePayment(payment oacquiring.PaymentInfo) error {
//	    // Do something with payment: update record in DB, send event, etc...
//	    return nil
//	}
//
// Регистрации обработчика HTTP уведомлений:
//
//	 key := `
//	 MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0fAzD+LmFQAKd0MLwnMQ
//	 75h+dGpDyGfneE2TIXMd4FxhlCjXEgQSp+kyqUswYDcy1t12pinuZHlZx0JnDVWB
//	 sfU7COK0bT/LEAOzoGhThqowP3qvxXTq2xWleZvxVYXwVXjIF4FFzieh0SoE8XaV
//	 GkqFLpDjDk5CYWHvoQ1FCeOmd5cVsXIQBEYJda45HRXdo9GcLwRDjpJDZZku6RIH
//	 sA6HpPa0Neo5THIpACa2noIcRF4IJkZDoU3bKE5qKNzSgpEQYp7M6Vgheh7VhLgy
//	 1Bv7+ABxuTn3CysTsT8C4IVRqsC3OmZ4wBsl/YkwZLnI0AMX911xagEweWXp9jz7
//	 EQIDAQAB
//	 `
//		handler, err := oacquiring.NewHTTPNotificationHandler(key, &Handler{})
//		if err != nil {
//	     // Handle err!
//		}
//
//	 // Register HTTP handler, for example:
//	 // http.Handle("/oplati/notification", &handler)
//	 // http.ListenAndServe(":8080", nil)
package oacquiring

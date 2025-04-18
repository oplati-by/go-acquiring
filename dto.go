package oacquiring

type (
	errorResponse struct {
		Code         string                   `json:"code"`
		InternalCode string                   `json:"internalCode"`
		DevMessage   string                   `json:"devMessage"`
		UserMessage  errorResponseUserMessage `json:"userMessage"`
	}

	errorResponseUserMessage struct {
		LangRu string `json:"lang_ru"`
		LangEn string `json:"lang_en"`
	}
)

type (
	newPaymentRequest struct {
		Shift           string                `json:"shift,omitempty"`
		Sum             float64               `json:"sum"`
		OrderNumber     string                `json:"orderNumber"`
		RegNum          string                `json:"regNum"`
		Details         paymentRequestDetails `json:"details"`
		SuccessUrl      string                `json:"successUrl"`
		FailureUrl      string                `json:"failureUrl"`
		NotificationUrl string                `json:"notificationUrl"`
	}

	paymentRequestDetailsItem struct {
		Type int     `json:"type"`
		Name string  `json:"name"`
		Cost float64 `json:"cost"`
	}

	paymentRequestDetails struct {
		RegNum      string                      `json:"regNum"`
		Items       []paymentRequestDetailsItem `json:"items"`
		AmountTotal float64                     `json:"amountTotal"`
		FooterInfo  string                      `json:"footerInfo"`
	}

	newPaymentResponse struct {
		PaymentId   int64  `json:"paymentId"`
		RedirectUrl string `json:"redirectUrl"`
	}
)

type (
	paymentInfoResponse struct {
		PaymentId     int64   `json:"paymentId"`
		PaymentType   int     `json:"paymentType"`
		Sum           float64 `json:"sum"`
		Status        int     `json:"status"`
		CreatedDate   string  `json:"createdDate"`
		PaidDate      string  `json:"paidDate"`
		OrderNumber   string  `json:"orderNumber"`
		PursePublicId string  `json:"pursePublicId"`
	}
)

type (
	reversePaymentRequest struct {
		Shift       string                `json:"shift,omitempty"`
		Sum         float64               `json:"sum"`
		OrderNumber string                `json:"orderNumber"`
		RegNum      string                `json:"regNum"`
		Details     paymentRequestDetails `json:"details"`
	}
)

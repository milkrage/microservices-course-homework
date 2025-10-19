package model

type Order struct {
	OrderID       string
	UserID        string
	PartIDs       []string
	TotalPrice    float64
	TransactionID *string
	PaymentMethod *OrderPaymentMethod
	Status        string
}

type OrderPaymentMethod struct {
	name   string
	number int32
}

func NewOrderPaymentMethod(name string) *OrderPaymentMethod {
	const (
		UnknownPaymentMethod       = "UNKNOWN"
		CardPaymentMethod          = "CARD"
		SBPPaymentMethod           = "SPB"
		CreditCardPaymentMethod    = "CREDIT_CARD"
		InvestorMoneyPaymentMethod = "INVESTOR_MONEY"
	)

	enum := map[string]int32{
		UnknownPaymentMethod:       0,
		CardPaymentMethod:          1,
		SBPPaymentMethod:           2,
		CreditCardPaymentMethod:    3,
		InvestorMoneyPaymentMethod: 4,
	}

	number, ok := enum[name]
	if !ok {
		number = enum[UnknownPaymentMethod]
		name = UnknownPaymentMethod
	}

	return &OrderPaymentMethod{
		number: number,
		name:   name,
	}
}

func (opm *OrderPaymentMethod) Number() int32 {
	return opm.number
}

func (opm *OrderPaymentMethod) Name() string {
	return opm.name
}

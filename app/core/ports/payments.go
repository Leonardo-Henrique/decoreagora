package ports

type PaymentHandler interface {
	Init()
	CreateCustomer(email string) (Customer, error)
	StartSubscription(customerID, priceID, successURL, cancelURL string) (CheckoutSession, error)
	CreateBillingPortalSession(customerID, returnURL string) (string, error)
}

type Customer struct {
	ID    string
	Email string
}

type CheckoutSession struct {
	ID  string
	URL string
}

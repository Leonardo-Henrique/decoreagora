package ports

type PaymentHandler interface {
	CreateCustomer(email string, name string) (Customer, error)
	StartSubscription(customerID, priceID, successURL, cancelURL string) (CheckoutSession, error)
	CreateBillingPortalSession(customerID, returnURL string) (string, error)
	CreatePaymentSession(customerID, priceID, successURL, cancelURL string) (CheckoutSession, error)
}

type Customer struct {
	ID    string
	Email string
}

type CheckoutSession struct {
	ID  string
	URL string
}

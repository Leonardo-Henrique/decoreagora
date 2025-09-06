package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"github.com/Leonardo-Henrique/decoreagora/app/core/config/logger"
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"go.uber.org/zap"
)

type PaymentUsecase struct {
	db      ports.Database
	payment ports.PaymentHandler
}

func NewPaymentUsecase(db ports.Database, payment ports.PaymentHandler) *PaymentUsecase {
	return &PaymentUsecase{
		db:      db,
		payment: payment,
	}
}

func (p *PaymentUsecase) CreateCustomer(req models.CreateSessionRequest) (ports.Customer, error) {
	// try to find customer
	sub, err := p.db.GetSubscriptionByEmail(req.Email)
	if err != nil {
		log.Println(err)
		return ports.Customer{}, errors.New("error when querying subscription")
	}

	log.Println("sub ->", sub)

	// if user doesnt exist, create
	if sub.StripeCostumerID == "" {

		user, err := p.db.GetUserByEmail(req.Email)
		if err != nil {
			log.Println(err)
			return ports.Customer{}, errors.New("error when querying user")
		}

		customer, err := p.payment.CreateCustomer(req.Email, user.Name)
		if err != nil {
			log.Println(err)
			return ports.Customer{}, errors.New("error when creating customer")
		}

		if err := p.db.UpdateUserCustomerID(user.ID, customer.ID); err != nil {
			log.Println(err)
			return ports.Customer{}, errors.New("error when updating customer id")
		}

		return customer, nil
	}

	customer := ports.Customer{
		ID:    sub.StripeCostumerID,
		Email: sub.Email,
	}

	return customer, nil
}

func (p *PaymentUsecase) CreatePaymentSession(req models.CreatePaymentSessionRequest) (ports.CheckoutSession, error) {
	return p.payment.CreatePaymentSession(req.StripeCustomerID, req.PriceID, req.SuccessUrl, req.CancelUrl)
}

func (p *PaymentUsecase) HandleWebhookEvent(event stripe.Event, body []byte) error {

	logger.Logging.Info("Event received", zap.Any("type", event.Type))

	switch event.Type {
	case "checkout.session.completed":
		return p.handleCheckoutSessionCompleted(event)
	case "payment_intent.payment_failed":
		return errors.New("payment failed")
	}

	return nil
}

func (p *PaymentUsecase) handleCheckoutSessionCompleted(event stripe.Event) error {
	logger.Logging.Info("Starting handleCheckoutSessionCompleted")

	var session stripe.CheckoutSession

	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		logger.Logging.Error("Error when trying to parse checkout.session.completed body", err)
		return err
	}

	priceID := session.Metadata["price_id"]
	if priceID == "" {
		logger.Logging.Error("There is no priceID in session:", fmt.Errorf("%v", session))
		return fmt.Errorf("There is no priceID in session")
	}

	credits := p.defineUserCreditsToDeposit(priceID)
	if credits == 0 {
		logger.Logging.Error("There is no credit deposit rule for this priceID:", fmt.Errorf("%v", priceID))
		return fmt.Errorf("there is no credit deposit rule for this priceID")
	}

	customerID := session.Customer.ID
	if customerID == "" {
		logger.Logging.Error("There is no customerID in session:", fmt.Errorf("%v", session))
		return fmt.Errorf("there is no credit deposit rule for this priceID")
	}

	sub, err := p.db.GetSubscriptionByCustomerID(customerID)
	if err != nil {
		logger.Logging.Error("Error when GetSubscriptionByCustomerID", fmt.Errorf("customerID: %s, err: %v", customerID, err))
		return fmt.Errorf("error when GetSubscriptionByCustomerID")
	}

	if sub.ID == 0 {
		logger.Logging.Error("No sub found for customerID:", fmt.Errorf("%s", customerID))
		return fmt.Errorf("no sub found")
	}

	if err := p.db.UpdateUserTier(customerID, "discrete_credits"); err != nil {
		logger.Logging.Error("Error when updating user tier:", err)
		return fmt.Errorf("error when updating user tier")
	}

	paymentEntry := models.PaymentHistory{
		PublicID:        uuid.NewString(),
		CustomerID:      customerID,
		ProcessedAt:     time.Now(),
		StripePriceID:   priceID,
		AmountPaid:      int(session.AmountTotal),
		CreditsReceived: credits,
	}
	if err := p.db.CreatePaymentHistoryEntry(paymentEntry); err != nil {
		logger.Logging.Error(fmt.Sprintf("Error when updating creating user payment history entry. Details: %v | Error: %v", paymentEntry, err), err)
	}

	return p.db.IncrementUserCreditsByCustomerID(sub.UserID, credits)
}

func (p *PaymentUsecase) defineUserCreditsToDeposit(priceID string) int {
	validCredits := map[string]int{
		config.C.PKG_30_LAUNCH:  30,
		config.C.PKG_100_LAUNCH: 100,
		config.C.PKG_200_LAUNCH: 200,
	}
	credits, ok := validCredits[priceID]
	if !ok {
		return 0
	}
	return credits
}

func (p *PaymentUsecase) SelectPriceIDByPlan(plan string) string {
	validPlans := map[string]string{
		"pkg_30_launch":  config.C.PKG_30_LAUNCH,
		"pkg_100_launch": config.C.PKG_100_LAUNCH,
		"pkg_200_launch": config.C.PKG_200_LAUNCH,
	}

	value, ok := validPlans[plan]
	if !ok {
		return ""
	}

	return value

}

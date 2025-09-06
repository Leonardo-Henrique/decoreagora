package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"github.com/Leonardo-Henrique/decoreagora/app/core/config/logger"
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/usecases"
	"github.com/Leonardo-Henrique/decoreagora/app/core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82/webhook"
)

type PaymentController struct {
	payment usecases.PaymentUsecase
}

func NewPaymentController(payment usecases.PaymentUsecase) *PaymentController {
	return &PaymentController{
		payment: payment,
	}
}

func (p *PaymentController) CreateSession(c *fiber.Ctx) error {

	var request models.CreateSessionRequest

	if err := c.BodyParser(&request); err != nil {
		log.Println("Error parsing request for session", err)
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, errors.New("invalid request body"))
	}

	sub, err := p.payment.CreateCustomer(request)
	if err != nil {
		log.Println("Error parsing request for session", err)
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, errors.New("error when creating customer"))

	}

	priceId := p.payment.SelectPriceIDByPlan(request.Plan)
	if priceId == "" {
		log.Println("Error when evaluating priceId", nil)
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, errors.New("no plan found"))
	}

	paymentSessionRequest := models.CreatePaymentSessionRequest{
		StripeCustomerID: sub.ID,
		PriceID:          priceId,
		SuccessUrl:       "https://www.google.com",
		CancelUrl:        "https://www.youtube.com",
	}

	checkoutSession, err := p.payment.CreatePaymentSession(paymentSessionRequest)
	if err != nil {
		log.Println("Error when evaluating creating checkout session", err)
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, errors.New("error when creating checkout session"))
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"session_id":  checkoutSession.ID,
		"session_url": checkoutSession.URL,
	})
}

func (p *PaymentController) Webhook(c *fiber.Ctx) error {
	logger.Logging.Info("Starting Webhook controller")
	const MaxBodyBytes = int64(65536)

	body := c.Body()
	if len(body) > int(MaxBodyBytes) {
		logger.Logging.Error("Stripe Webhook Body was too large", nil)
		return utils.ErrorResponse(c, fiber.StatusRequestEntityTooLarge, nil)
	}

	sigHeader := c.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(body, sigHeader, config.C.STRIPE_WEBHOOK_SECRET)
	if err != nil {
		logger.Logging.Error("Error verifying Stripe Webhook Signature", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	if err := p.payment.HandleWebhookEvent(event, body); err != nil {
		logger.Logging.Error("Error when handling webhook", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return nil

}

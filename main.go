package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Leonardo-Henrique/decoreagora/app/adapters/ai"
	"github.com/Leonardo-Henrique/decoreagora/app/adapters/controllers"
	"github.com/Leonardo-Henrique/decoreagora/app/adapters/middlewares"
	"github.com/Leonardo-Henrique/decoreagora/app/adapters/repositories"
	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"github.com/Leonardo-Henrique/decoreagora/app/core/config/logger"
	"github.com/Leonardo-Henrique/decoreagora/app/core/infra"
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/usecases"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	logger.InitLogger()

	app := fiber.New()

	dsn := repositories.MakeDSNString(models.DSN{
		Host:         config.C.DB_HOST,
		User:         config.C.DB_USER,
		Pass:         config.C.DB_PASS,
		Port:         config.C.DB_PORT,
		DatabaseName: config.C.DB_DATABASE,
	})

	db, err := repositories.ConnectToDatabase(dsn)
	if err != nil {
		log.Fatalf("FATAL ERROR WHENN CONNECTING TO DB %v", err)
	}
	defer db.Close()

	jwt := repositories.NewJWT()

	s3 := repositories.NewS3Repository(*infra.NewAWSConfig())

	genai := ai.NewGemini(context.Background())

	stripe := repositories.NewStripeClient()

	ses := repositories.NewSESRepository()

	/*
		REPOSITORIES
	*/
	mysql := repositories.NewMySQLRepository(db)

	middleware := middlewares.NewMiddleware(mysql, jwt)

	/*
		USECASES
	*/
	userUC := usecases.NewUserUsecase(mysql)
	loginUC := usecases.NewLoginUsecase(mysql, jwt)
	imgUC := usecases.NewImagesUsecase(mysql, s3, genai)
	creditsUC := usecases.NewCreditsUsecase(mysql)
	subsUC := usecases.NewSubscriptionUsecase(mysql)
	payUC := usecases.NewPaymentUsecase(mysql, stripe)
	emailUC := usecases.NewEmailUseCase(ses)

	/*
		CONTROLLERS
	*/
	userCtrl := controllers.NewUserController(*loginUC, *userUC, *subsUC, *emailUC)
	loginCtrl := controllers.NewLoginController(*loginUC, *emailUC)
	imgCtrl := controllers.NewImageController(*creditsUC, *imgUC)
	payCtrl := controllers.NewPaymentController(*payUC)

	/*
		ROUTES
	*/
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return origin == "https://app.decoreagora.com.br" || origin == "http://localhost:5174" || origin == "http://localhost:4173"
		},
		AllowCredentials: true,
		AllowHeaders:     "Content-Type, Authorization, X-Requested-With, Accept, Origin, X-CSRF-Token",
	}))

	app.Get("/users/me", middleware.AuthMiddleware(userCtrl.GetMe))
	app.Post("/users", userCtrl.NewUser)
	app.Post("/login", loginCtrl.Login)
	app.Post("/authenticate", loginCtrl.AuthenticateCode)

	app.Get("/images/credits", middleware.AuthMiddleware(imgCtrl.GetUserCredits))
	app.Post("/image/create", middleware.AuthMiddleware(middleware.CreditsMiddleware(imgCtrl.CreateNewImage)))
	app.Get("/images", middleware.AuthMiddleware(imgCtrl.GetUserImages))

	app.Post("/payments/session", middleware.AuthMiddleware(payCtrl.CreateSession))
	app.Post("/payments/webhook", payCtrl.Webhook)
	app.Get("/health/warmup", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	log.Println("App Started")

	log.Fatal(app.Listen(":" + config.C.AppPort))

}

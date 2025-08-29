package main

import (
	"log"

	"github.com/Leonardo-Henrique/decoreagora/app/adapters/controllers"
	"github.com/Leonardo-Henrique/decoreagora/app/adapters/repositories"
	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/usecases"
	"github.com/gofiber/fiber/v2"
)

func main() {

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

	/*
		REPOSITORIES
	*/
	mysql := repositories.NewMySQLRepository(db)

	/*
		USECASES
	*/
	userUC := usecases.NewUserUsecase(mysql)

	/*
		CONTROLLERS
	*/
	userCtrl := controllers.NewUserController(*userUC)

	/*
		ROUTES
	*/
	app.Post("/users", userCtrl.NewUser)

	log.Println("App Started")

	log.Fatal(app.Listen(":" + config.C.AppPort))

}

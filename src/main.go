package main

import (
	basicauth "bot/src/controllers/basicauth"
	config "bot/src/controllers/config"
	munro "bot/src/controllers/munro"
	routes "bot/src/controllers/routes"
	storage "bot/src/controllers/storage"

	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	/*
		Init
	*/
	// Conf
	conf := config.Read()

	// Kitsu auth
	// TODO: check for token expiration
	JWTToken := basicauth.AuthForJWTToken(conf.Kitsu.Hostname+"api/auth/login", conf.Kitsu.Email, conf.Kitsu.Password)
	os.Setenv("JWTToken", JWTToken)

	// Connect to DB
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migration
	db.AutoMigrate(&storage.TaskRecord{})

	/*
		Telegram Bot
	*/
	bot, err := tgbotapi.NewBotAPI(conf.Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	// Create update
	updates := munro.InitWebhookOrUpdate(bot, conf)
	updates.Clear()

	// Debug
	bot.Debug = conf.Bot.Debug
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Bot updates
	go munro.ListenBotUpdates(bot, updates, conf)

	// Parse statuses
	go func() {
		for x := range time.Tick(time.Duration(conf.Messaging.PollDuration) * time.Minute) {
			munro.ParseTaskStatuses(bot, conf, x, db)
		}
	}()

	/*
		Routing
	*/
	app := fiber.New()

	// CORS
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     conf.CORS.AllowOrigins,
		AllowMethods:     conf.CORS.AllowMethods,
		AllowHeaders:     conf.CORS.AllowHeaders,
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	// API routes
	// give response when at /api/v1
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Hello",
		})
	})

	routes.APIRoutes(app, bot, db)

	app.Listen(":3001")
}

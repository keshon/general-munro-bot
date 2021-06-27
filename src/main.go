package main

import (
	"bot/src/controllers/basicauth"
	"bot/src/controllers/config"
	"bot/src/controllers/munro"
	"bot/src/controllers/routes"
	"bot/src/controllers/storage"
	"bot/src/utils/remove"
	"fmt"

	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	db.AutoMigrate(&storage.Task{})
	db.AutoMigrate(&storage.Attachment{})

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

	// Parse notifications
	go func() {
		for x := range time.Tick(time.Duration(conf.Notification.PollDuration) * time.Minute) {
			fmt.Println("START checking Task statuses at " + time.Now().String())
			munro.ParseTaskStatuses(bot, conf, x, db)
			fmt.Println("DONE checking Task statuses at " + time.Now().String())
		}
	}()

	// Parse attachments
	go func() {
		for x := range time.Tick(time.Duration(conf.Backup.PollDuration) * time.Minute) {
			if conf.Backup.FastDelete != true {
				remove.RemoveContents(conf.Backup.LocalStorage + "trash/")
				fmt.Println("DONE clearing Trash at " + time.Now().String())
			}

			fmt.Println("START checking Attachments at " + time.Now().String())
			munro.ParseAttachments(bot, conf, x, db)
			fmt.Println("DONE checking Attachments at " + time.Now().String())

		}
	}()

	/*
		Routing
	*/
	// Config app
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Error type
			type Error struct {
				StatusCode int    `json:"statusCode"`
				Error      string `json:"error"`
			}
			// Default 500 statuscode
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
			}
			// Set Content-Type: text/plain; charset=utf-8
			ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			// Return statuscode with error message
			fmt.Println(err)
			return ctx.Status(code).JSON(Error{code, err.Error()})
		},
	})

	// CORS
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     conf.CORS.AllowOrigins,
		AllowMethods:     conf.CORS.AllowMethods,
		AllowHeaders:     conf.CORS.AllowHeaders,
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           100,
	}))

	// Middlewares
	// Recover middleware
	app.Use(recover.New())

	// API routes
	// give response when at /api/v1
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Hello",
		})
	})

	routes.APIRoutes(app, bot, db, conf)

	app.Listen(conf.Kitsu.ListenHostname)
}

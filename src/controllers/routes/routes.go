package routes

import (
	"bot/src/controllers/config"
	"bot/src/controllers/kitsu"
	"bot/src/controllers/munro"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

func APIRoutes(app *fiber.App, bot *tgbotapi.BotAPI, db *gorm.DB) {
	// API
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// give response when at /api/v1
	v1.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the api endpoint 😉",
		})
	})

	// Routes
	// Kitsu
	KitsuRoute(v1, bot, db)
}

func KitsuRoute(route fiber.Router, bot *tgbotapi.BotAPI, db *gorm.DB) {
	route.Post("/kitsu", func(c *fiber.Ctx) error {
		// Conf
		conf := config.Read()

		type KitsuRequest struct {
			EntityType   string   `json:"entitytype,omitempty"`
			OriginURL    string   `json:"originurl,omitempty"`
			OriginServer string   `json:"originserver,omitempty"`
			Selection    []string `json:"selection,omitempty"`
			ProductionID string   `json:"productionid"`
			UserID       string   `json:"userid"`
			UserEmail    string   `json:"useremail"`
		}

		req := new(KitsuRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}

		for _, elem := range req.Selection {

			currentTask := kitsu.GetTask(elem)
			munro.ParseTaskStatus(bot, conf, db, currentTask)
		}

		return c.JSON("OK")
	})
}

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

	// give response when at /api
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

			/*
				var assigneePhone string
				if len(currentTask.Assignees) > 0 {
					for _, elem := range currentTask.Assignees {
						currentAssingnee := kitsu.GetPerson(elem)
						if currentAssingnee.Phone != "" {
							assigneePhone = assigneePhone + currentAssingnee.Phone + ", "
						}
					}
				}

				// Compose message
				var messageTemplate string = ""
				initiate := kitsu.GetPerson(req.UserID)

				if initiate.Phone != "" {
					messageTemplate = i18n.Tr(conf.Bot.Language, "from") + initiate.Phone + "\n"
				}

				messageTemplate = messageTemplate + assigneePhone + i18n.Tr(conf.Bot.Language, "status") + " <b>" + strings.ToUpper(currentTaskStatus.ShortName) + "</b> " + i18n.Tr(conf.Bot.Language, "for") + " " + currentEntity.Name

				// Send message by Role matching
				var messageSent = false
				for _, elem := range conf.Credentials.ChatIDByRoles {
					role := strings.ToLower(strings.Split(elem, ":")[0]) // extract role name and make it lowercase
					currentTaskStatusName := strings.ToLower(currentTaskStatus.ShortName)
					if role == currentTaskStatusName {
						// get all chat ids
						chatIDs := strings.Split(elem, ":")[1]
						if len(strings.Split(chatIDs, "|")) > 0 {
							chatID, _ := strconv.ParseInt(strings.Split(chatIDs, "|")[0], 10, 64)
							msg := tgbotapi.NewMessage(chatID, messageTemplate)
							msg.ParseMode = "html"
							bot.Send(msg)
							messageSent = true

							// confirmation
							chatID, _ = strconv.ParseInt(strings.Split(chatIDs, "|")[1], 10, 64)
							msg = tgbotapi.NewMessage(chatID, "\xF0\x9F\x91\x8D")
							msg.ParseMode = "html"
							bot.Send(msg)
							messageSent = true
						} else {
							chatID, _ := strconv.ParseInt(chatIDs, 10, 64)
							msg := tgbotapi.NewMessage(chatID, messageTemplate)
							msg.ParseMode = "html"
							bot.Send(msg)
							messageSent = true
						}

					}
				}

				// Send message to Admin if no role matching was done successfuly
				if messageSent == false {
					chatID, _ := strconv.ParseInt(conf.Credentials.AdminChatID, 10, 64)
					messageTemplate = i18n.Tr(conf.Bot.Language, "unknown-status") + "\n" + messageTemplate
					msg := tgbotapi.NewMessage(chatID, messageTemplate)
					msg.ParseMode = "html"
					bot.Send(msg)
				}
			*/
		}

		return c.JSON("OK")
	})
}

package main

import (
	basicauth "bot/src/controllers/basicauth"
	config "bot/src/controllers/config"
	kitsu "bot/src/controllers/kitsu"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hokaccha/go-prettyjson"
)

func listenBotUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	// Conf
	conf := config.Read()

	for update := range updates {

		// Debug
		if conf.Bot.Debug == true {

			fmt.Println("--start--")

			dt := time.Now()
			fmt.Println(dt.String())

			resp, err := prettyjson.Marshal(update)
			fmt.Println(string(resp))
			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Println("--end--\n")
		}

		// Commands
		if update.Message != nil {

			var chatID int64

			if _, err := strconv.Atoi(conf.Credentials.AdminChatID); err == nil {
				chatID, _ = strconv.ParseInt(conf.Credentials.AdminChatID, 10, 64)
			} else {
				chatID = int64(update.Message.From.ID)
			}

			// Listen for admin only command in groups or in private
			if update.Message.Chat.Type == "group" {
				if int64(update.Message.From.ID) == chatID {
					// Listen for bot's UserName
					botUserName := "@" + bot.Self.UserName
					if strings.HasPrefix(update.Message.Text, botUserName) {
						// Listen for specific word (command)
						command := strings.Split(update.Message.Text, botUserName+" ")
						if len(command) > 1 {
							if command[1] == "lookup" {
								resp, err := json.MarshalIndent(update, "", "  ")
								if err != nil {
									log.Fatal(err)
									return
								}
								msg := tgbotapi.NewMessage(chatID, "<pre>"+string(resp)+"</pre>")
								msg.ParseMode = "html"
								bot.Send(msg)
							}
						}
					}
				}
			} else {
				command := update.Message.Text
				if command == "lookup" {
					resp, err := json.MarshalIndent(update, "", "  ")
					if err != nil {
						log.Fatal(err)
						return
					}
					msg := tgbotapi.NewMessage(chatID, "<pre>"+string(resp)+"</pre>")
					msg.ParseMode = "html"
					bot.Send(msg)
				}
			}
		}
	}
}

func main() {
	// Conf
	conf := config.Read()

	// Kitsu auth
	JWTToken := basicauth.AuthForJWTToken(conf.Kitsu.Hostname+"api/auth/login", conf.Kitsu.Email, conf.Kitsu.Password)
	os.Setenv("JWTToken", JWTToken)

	// Bot
	// Create instance
	bot, err := tgbotapi.NewBotAPI(conf.Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	// Webhook or GetUpdate
	var updates tgbotapi.UpdatesChannel
	if conf.Bot.Webhook == true {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(conf.Bot.Hostname + bot.Token))
		if err != nil {
			log.Fatal(err)
		}
		info, err := bot.GetWebhookInfo()
		if err != nil {
			log.Fatal(err)
		}
		if info.LastErrorDate != 0 {
			log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		}
		updates = bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServe(conf.Bot.ListenHostname, nil)
	} else {
		// Delete Webhook
		_, err = bot.RemoveWebhook()
		if err != nil {
			log.Fatal(err)
		}
		// Create Update
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 100

		updates, _ = bot.GetUpdatesChan(u)
	}
	updates.Clear()

	// Debug
	bot.Debug = conf.Bot.Debug
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Updates
	go listenBotUpdates(bot, updates)

	// Fiber
	// Create instance
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

	// Routes
	app.Post("/", func(c *fiber.Ctx) error {

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
			TaskStatusID := currentTask.TaskStatusID
			currentTaskStatus := kitsu.GetTaskStatus(TaskStatusID)

			currentEntity := kitsu.GetEntity(currentTask.EntityID)
			var assigneePhone = ""
			fmt.Println(currentTask.Assignees)
			if len(currentTask.Assignees) > 0 {
				for _, elem := range currentTask.Assignees {
					currentAssingnee := kitsu.GetPerson(elem)
					if currentAssingnee.Phone != "" {
						assigneePhone = assigneePhone + currentAssingnee.Phone + ", "
					}
					fmt.Println(currentAssingnee.Phone)
				}
			}

			var messageTemplate string = ""
			initiate := kitsu.GetPerson(req.UserID)

			if initiate.Phone != "" {
				messageTemplate = "Уведомление от " + initiate.Phone + "\n"
			}
			fmt.Println(initiate.Phone)
			messageTemplate = messageTemplate + assigneePhone + "статус по задаче " + currentEntity.Name + ": <b>" + strings.ToUpper(currentTaskStatus.ShortName) + "</b>"
			fmt.Println(string(messageTemplate))

			// Role matching
			var messageSent = false
			for _, elem := range conf.Credentials.ChatIDByRoles {
				role := strings.ToLower(strings.Split(elem, "=")[0]) // extract role name and make it lowercase
				currentTaskStatusName := strings.ToLower(currentTaskStatus.ShortName)

				fmt.Println(role + " == " + currentTaskStatusName)

				if role == currentTaskStatusName {
					chatID, _ := strconv.ParseInt(strings.Split(elem, "=")[1], 10, 64)
					msg := tgbotapi.NewMessage(chatID, messageTemplate)
					msg.ParseMode = "html"
					bot.Send(msg)
					messageSent = true
				}
			}

			// Send message to Admin if no role matching was done successfuly
			if messageSent == false {
				chatID, _ := strconv.ParseInt(conf.Credentials.AdminChatID, 10, 64)
				messageTemplate = "<b>Внимание!</b> Настройки для данного статуса не найдены.\n" + messageTemplate
				msg := tgbotapi.NewMessage(chatID, messageTemplate)
				msg.ParseMode = "html"
				bot.Send(msg)
			}
		}

		return c.JSON("OK")
	})

	app.Listen(":3001")
}

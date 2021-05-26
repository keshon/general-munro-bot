package munro

import (
	"bot/src/controllers/config"
	"bot/src/controllers/kitsu"
	"bot/src/controllers/storage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hokaccha/go-prettyjson"

	"gorm.io/gorm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kataras/i18n"
)

func InitWebhookOrUpdate(bot *tgbotapi.BotAPI, conf config.Config) tgbotapi.UpdatesChannel {

	var updates tgbotapi.UpdatesChannel

	if conf.Bot.Webhook == true {
		_, err := bot.SetWebhook(tgbotapi.NewWebhook(conf.Bot.Hostname + bot.Token))
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

		return updates
	}

	// Delete Webhook
	_, err := bot.RemoveWebhook()
	if err != nil {
		log.Fatal(err)
	}

	// Create Update
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 100

	updates, _ = bot.GetUpdatesChan(u)

	return updates
}

func ListenBotUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, conf config.Config) {
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

			// Check if admin chat is defined
			if _, err := strconv.Atoi(conf.Credentials.AdminChatID); err == nil {
				chatID, _ = strconv.ParseInt(conf.Credentials.AdminChatID, 10, 64)
			} else {
				chatID = int64(update.Message.From.ID)
			}

			// Listen command in group of private chat
			var command string
			if update.Message.Chat.Type == "group" {
				if int64(update.Message.From.ID) == chatID {
					// Listen for bot's UserName
					botUserName := "@" + bot.Self.UserName
					if strings.HasPrefix(update.Message.Text, botUserName) {
						splitted := strings.Split(update.Message.Text, botUserName+" ")
						if len(splitted) > 1 {
							command = splitted[1]
						}
					}
				}
			} else {
				command = update.Message.Text
			}

			switch command {
			case "/lookup":
				resp, err := json.MarshalIndent(update, "", "  ")
				if err != nil {
					log.Fatal(err)
					return
				}
				msg := tgbotapi.NewMessage(chatID, "<pre>"+string(resp)+"</pre>")
				msg.ParseMode = "html"
				bot.Send(msg)
			case "/about":
				msg := tgbotapi.NewMessage(chatID, i18n.Tr(conf.Bot.Language, "about"))
				msg.ParseMode = "html"
				bot.Send(msg)
			}
		}
	}
}

func ParseTaskStatuses(bot *tgbotapi.BotAPI, conf config.Config, t time.Time, db *gorm.DB) {
	// Get all Tasks
	allTasks := kitsu.GetTasks()
	if len(allTasks.Each) > 0 {
		for _, task := range allTasks.Each {

			// Get human readbale status
			currentTaskStatus := kitsu.GetTaskStatus(task.TaskStatusID)

			// Get entity name (Top Task)
			currentEntity := kitsu.GetEntity(task.EntityID)
			entityName := currentEntity.Name

			// Get assignee for the Task and his phone data (we store Telegram nicknames there)
			currentDetailedTask := kitsu.GetTask(task.ID)

			var assigneePhone = ""
			if len(currentDetailedTask.Assignees) > 0 {
				for _, elem := range currentDetailedTask.Assignees {
					assingnee := kitsu.GetPerson(elem)
					if assingnee.Phone != "" {
						assigneePhone = assigneePhone + assingnee.Phone + ", "
					}
				}
			}

			// Get comment
			var commentMessage = ""
			currentComments := kitsu.GetComment(currentDetailedTask.ID)
			if len(currentComments.Each) > 0 {
				// find the most recent comment in array
				sort.Slice(currentComments.Each, func(i, j int) bool {
					return currentComments.Each[i].UpdatedAt > currentComments.Each[j].UpdatedAt
				})

				if currentComments.Each[0].Text != "" {
					commentAuthor := kitsu.GetPerson(currentComments.Each[0].PersonID)
					commentMessage = " " + i18n.Tr(conf.Bot.Language, "with-comment") + "\n<pre>" + commentAuthor.FullName + ":\n" + currentComments.Each[0].Text + "</pre>"
				}
			}

			// Compose message
			messageTemplate := assigneePhone + i18n.Tr(conf.Bot.Language, "status") + " <b>" + strings.ToUpper(currentTaskStatus.ShortName) + "</b> " + i18n.Tr(conf.Bot.Language, "for") + " " + entityName

			if commentMessage != "" {
				messageTemplate += commentMessage
			}

			// Decision making
			result := storage.FindRecord(db, task.ID)
			if len(result.TaskID) <= 0 {
				// create
				storage.CreateRecord(db, task.ID, currentTaskStatus.ShortName)
				// say
				sendMessage(bot, conf, messageTemplate, currentTaskStatus.ShortName)
			} else {
				// TODO or date mismatch
				if result.TaskStatus != currentTaskStatus.ShortName {
					// update
					storage.UpdateRecord(db, task.ID, currentTaskStatus.ShortName)
					// say
					sendMessage(bot, conf, messageTemplate, currentTaskStatus.ShortName)
				}
			}
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, conf config.Config, message, taskStatus string) {

	var messageSent = false
	for _, elem := range conf.Credentials.ChatIDByRoles {
		role := strings.ToLower(elem) // extract role name and make it lowercase
		currentTaskStatusName := strings.ToLower(taskStatus)

		// find chat id
		var chatID int64
		if role == currentTaskStatusName {
			//chatID, _ = strconv.ParseInt(elem, 10, 64)
			chatIDs := strings.Split(elem, ":")[1]
			chatID, _ = strconv.ParseInt(chatIDs, 10, 64)
			// send
			msg := tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "html"
			bot.Send(msg)
			messageSent = true
		}
	}
	// Send message to Admin if no role matching was done successfuly
	if messageSent == false {
		chatID, _ := strconv.ParseInt(conf.Credentials.AdminChatID, 10, 64)
		message = i18n.Tr(conf.Bot.Language, "unknown-status") + "\n" + message
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = "html"
		bot.Send(msg)
	}

}

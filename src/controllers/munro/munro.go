package munro

import (
	"bot/src/controllers/config"
	"bot/src/controllers/kitsu"
	"bot/src/controllers/storage"
	"bot/src/controllers/wasabi"
	"bot/src/utils/truncate"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
			if _, err := strconv.Atoi(conf.Notification.AdminChatID); err == nil {
				chatID, _ = strconv.ParseInt(conf.Notification.AdminChatID, 10, 64)
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
	array := kitsu.GetTasks()
	if len(array.Each) > 0 {
		for _, elem := range array.Each {
			ParseTaskStatus(bot, conf, db, elem)
		}
	}
}

func ParseTaskStatus(bot *tgbotapi.BotAPI, conf config.Config, db *gorm.DB, task kitsu.Task) {

	// Get human readbale status
	currentTaskStatus := kitsu.GetTaskStatus(task.TaskStatusID)

	// Get entity name (Top Task)
	currentEntity := kitsu.GetEntity(task.EntityID)
	entityName := currentEntity.Name

	// Get assingee for the Task and his phone data (we store Telegram nicknames there)
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
	var commentID = ""
	var commentUpdatedAt = ""

	//var debug = ""

	currentComments := kitsu.GetComment(currentDetailedTask.ID)
	if len(currentComments.Each) > 0 {
		// find the most recent comment in array
		sort.Slice(currentComments.Each, func(i, j int) bool {
			layout := "2006-01-02T15:04:05"
			a, err := time.Parse(layout, currentComments.Each[i].UpdatedAt)
			if err != nil {
				fmt.Println(err)
			}
			b, err := time.Parse(layout, currentComments.Each[j].UpdatedAt)
			if err != nil {
				fmt.Println(err)
			}
			return a.Unix() > b.Unix()
		})

		//debug = "Newest comment is \n - updated: " + currentComments.Each[0].UpdatedAt + "\n - text: " + currentComments.Each[0].Text

		commentID = currentComments.Each[0].ID
		commentUpdatedAt = currentComments.Each[0].UpdatedAt
		//commentAuthor = kitsu.GetPerson(currentComments.Each[0].PersonID)

		truncatedComment := truncate.TruncateString(currentComments.Each[0].Text, 128)
		if truncatedComment != currentComments.Each[0].Text {
			truncatedComment += "..."
		}

		if currentComments.Each[0].Text != "" {
			commentAuthor := kitsu.GetPerson(currentComments.Each[0].PersonID)
			commentMessage = "\n<pre>" + commentAuthor.FullName + ":\n" + truncatedComment + "</pre>"
		}
	}

	// Decision making
	var messageTemplate = ""
	//messageTemplate += "<pre>" + debug + "</pre>\n"

	result := storage.FindTask(db, task.ID)
	if len(result.TaskID) > 0 {
		// check if status is different or last comment date don't match
		if result.TaskStatus != currentTaskStatus.ShortName || result.TaskUpdatedAt != task.UpdatedAt {
			// update
			storage.UpdateTask(db, task.ID, task.UpdatedAt, currentTaskStatus.ShortName, commentID, commentUpdatedAt)

			// say
			if conf.Notification.SilentUpdate != true {

				// Same status or not
				if result.TaskStatus != currentTaskStatus.ShortName {
					messageTemplate += assigneePhone + i18n.Tr(conf.Bot.Language, "updated-status") + " <b>" + strings.ToUpper(currentTaskStatus.ShortName) + "</b> (" + i18n.Tr(conf.Bot.Language, "prev-status") + " " + strings.ToLower(result.TaskStatus) + ") " + i18n.Tr(conf.Bot.Language, "for-task") + " " + entityName
				} else {
					messageTemplate += assigneePhone + i18n.Tr(conf.Bot.Language, "status") + " <b>" + strings.ToUpper(currentTaskStatus.ShortName) + "</b> " + i18n.Tr(conf.Bot.Language, "for-task") + " " + entityName
				}

				if commentMessage != "" {
					layout := "2006-01-02T15:04:05"
					db, err := time.Parse(layout, result.CommentUpdatedAt)
					if err != nil {
						fmt.Println(err)
					}
					comment, err := time.Parse(layout, commentUpdatedAt)
					if err != nil {
						fmt.Println(err)
					}

					if comment.Unix() > db.Unix() {
						messageTemplate += commentMessage
					}
				}

				sendMessage(bot, conf, messageTemplate, currentTaskStatus.ShortName)
			}
		}

	} else {
		// create
		storage.CreateTask(db, task.ID, task.UpdatedAt, currentTaskStatus.ShortName, commentID, commentUpdatedAt)
		// say
		if conf.Notification.SilentUpdate != true {

			// Compose message
			messageTemplate += assigneePhone + i18n.Tr(conf.Bot.Language, "new-status") + " <b>" + strings.ToUpper(currentTaskStatus.ShortName) + "</b> " + i18n.Tr(conf.Bot.Language, "for-task") + " <i>" + entityName + "</i>"

			if commentMessage != "" {
				messageTemplate += commentMessage
			}
			sendMessage(bot, conf, messageTemplate, currentTaskStatus.ShortName)
		}
	}
}

// Bot sending message
func sendMessage(bot *tgbotapi.BotAPI, conf config.Config, message, taskStatus string) {

	var messageSent = false
	for _, elem := range conf.Notification.ChatIDByRoles {
		role := strings.ToLower(strings.Split(elem, ":")[0]) // extract role name and make it lowercase
		currentTaskStatusName := strings.ToLower(taskStatus)

		// find chat id
		var chatID int64
		if role == currentTaskStatusName {
			//chatID, _ = strconv.ParseInt(elem, 10, 64)
			chatIDs := strings.Split(elem, ":")[1]
			chatID, _ = strconv.ParseInt(chatIDs, 10, 64)

			if conf.Notification.SilentUpdate != true {
				// send status (not working :/)
				status := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
				bot.Send(status)

				// Calling Sleep method
				time.Sleep(5 * time.Second)
			}

			// send message
			msg := tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "html"
			bot.Send(msg)

			messageSent = true
		}
	}

	// Send message to Admin if no role matching was done successfuly and supress is disabled
	if messageSent == false && conf.Notification.SuppressUndefinedRoles != true {
		chatID, _ := strconv.ParseInt(conf.Notification.AdminChatID, 10, 64)

		bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
		message = i18n.Tr(conf.Bot.Language, "unknown-status") + "\n" + message
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = "html"
		bot.Send(msg)
	}

}

func ParseAttachments(bot *tgbotapi.BotAPI, conf config.Config, t time.Time, db *gorm.DB) {
	// Get all Attachments
	array := kitsu.GetAttachments()
	if len(array.Each) > 0 {
		var count int
		for _, elem := range array.Each {

			resp := ParseAttachment(bot, conf, db, elem)
			if resp == true {
				count++
			}
		}

		if count > 0 {
			chatID, _ := strconv.ParseInt(conf.Notification.AdminChatID, 10, 64)
			bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
			strCount := strconv.Itoa(count)
			message := i18n.Tr(conf.Bot.Language, "backup-finished") + strCount
			msg := tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "html"
			bot.Send(msg)
		}
	}
}

func ParseAttachment(bot *tgbotapi.BotAPI, conf config.Config, db *gorm.DB, attachment kitsu.Attachment) bool {
	fmt.Println("Proccesing " + attachment.Name)

	if attachment.ID == "" {
		return false
	}

	for _, elem := range conf.Backup.Ignore {
		if attachment.Extension == elem {
			fmt.Println("Skipping ignored extension: " + elem)
			return false
		}
	}

	localPath := conf.Backup.LocalStorage + attachment.ID

	s3Path := ""
	if attachment.Comment.ObjectID != "" {

		task := kitsu.GetTask(attachment.Comment.ObjectID)

		// Get entity name (Top Task)
		entity := kitsu.GetEntity(task.EntityID)
		entityName := entity.Name

		// Get task type (Sub Task)
		taskType := kitsu.GetTaskType(task.TaskTypeID)
		taskTypeName := ""
		if taskType.Name != "" {
			taskTypeName = taskType.Name + "/"
		}
		project := kitsu.GetProject(task.ProjectID)
		//projectStatus := kitsu.GetProjectStatus(project.ProjectStatusID)

		s3Path = conf.Backup.S3.RootFolderName + "/" + project.Name + "/" + entityName + "/" + taskTypeName + attachment.Name
	} else {
		s3Path = conf.Backup.S3.RootFolderName + "/" + "LOST.FILES" + "/" + attachment.ID + "/" + attachment.Name
	}

	result := storage.FindAttachment(db, attachment.ID)
	if len(result.AttachmentID) > 0 {
		// check if status is different or last comment date don't match
		if result.AttachmentStatus != "done" || result.AttachmentUpdatedAt != attachment.UpdatedAt {
			// update
			storage.UpdateAttachment(db, attachment.ID, attachment.UpdatedAt, "new")
			kitsu.DownloadAttachment(localPath, attachment.ID, attachment.Name, conf)

			// Read file from local dir
			content, err := ioutil.ReadFile(localPath + "/" + attachment.Name)
			if err != nil {
				panic(err)
			}

			// Upload file to S3 storage
			wasabi.UploadFile(s3Path, string(content), conf)
			storage.UpdateAttachment(db, attachment.ID, attachment.UpdatedAt, "done")
		} else {
			fmt.Println("Skipping existing attachment: " + attachment.Name)
			return false
		}

	} else {
		// create
		// Download file from Kitsu
		storage.CreateAttachment(db, attachment.ID, attachment.UpdatedAt, "new")
		kitsu.DownloadAttachment(localPath, attachment.ID, attachment.Name, conf)

		// Read file from local dir
		content, err := ioutil.ReadFile(localPath + "/" + attachment.Name)
		if err != nil {
			panic(err)
		}

		// Upload file to S3 storage
		wasabi.UploadFile(s3Path, string(content), conf)
		storage.UpdateAttachment(db, attachment.ID, attachment.UpdatedAt, "done")
	}

	// Cleaning
	if conf.Backup.FastDelete != true {
		trashPath := conf.Backup.LocalStorage + "trash/" + attachment.ID
		doneAttachment := storage.FindAttachment(db, attachment.ID)
		if doneAttachment.AttachmentStatus == "done" {
			if _, err := os.Stat(localPath); !os.IsNotExist(err) {
				err := os.Rename(localPath, trashPath) // rename directory
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		os.RemoveAll(localPath)
		fmt.Println("DONE deleting at" + time.Now().String())
	}

	fmt.Println("")

	return true
}

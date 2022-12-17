package main

import (
	"fmt"
	"github.com/glebarez/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"strings"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("LIST", "LIST"),
	),
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("ddg_email.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")
	} else {
		log.SetOutput(file)
	}
}
func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN")) //os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	_ = db.AutoMigrate(Token{})
	_ = db.AutoMigrate(WaitOTP{})

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {

		// Check if we've gotten a message update.
		if update.Message != nil {
			if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Not Supported \n Pls DM")
				if _, err = bot.Send(msg); err != nil {
					log.Error(err)
					continue
				}
				continue
			}
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "欢迎食用DuckDuckGOEmailBOT\n,可以借助本BOT创建匿名邮箱，请先到DuckDuckGoEmail https://duckduckgo.com/email 先行注册 \n Thanks for Using This Bot, U might to get a username from https://duckduckgo.com/email firstly")
				msg.ReplyMarkup = numericKeyboard
				if _, err = bot.Send(msg); err != nil {
					log.Error(err)
				}
			case "del":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "输入存在问题 请/del 用户名 输入\n Error Input, Please follow /del Username style, try again")
				args := strings.Split(update.Message.Text, " ")
				if len(args) > 1 {
					username := args[1]
					tokens := []Token{}
					db.Where(&Token{UserName: username}).Find(&tokens)
					for _, value := range tokens {
						db.Delete(&value)
					}
					msg.Text = "已删除\nDeleted"
				}
				if _, err = bot.Send(msg); err != nil {
					log.Error(err)
				}
			case "add":
				args := strings.Split(update.Message.Text, " ")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "输入存在问题 请/add 用户名 输入\n Error Input, Please follow /add Username style, try again")
				if len(args) > 1 {
					username := args[1]
					token := Token{}
					db.Where(&Token{UserName: username}).Last(&token)
					if token.UserName != "" && token.TID != update.Message.Chat.ID {
						msg.Text = "该用户名已经被他人绑定 \n This username has already been bound by someone else"
					} else if token.UserName != "" && token.TID == update.Message.Chat.ID {
						msg.Text = "该用户名已经被你绑定 \nThis username has already been bound by you"
					} else {
						db.Create(&WaitOTP{TID: update.Message.Chat.ID, UserName: username})
						RequestOTP(username)
						msg.Text = "请输入邮箱中的一次性密码\nPlease enter the one-time password in the email"
					}

				}
				if _, err = bot.Send(msg); err != nil {
					log.Error(err)
				}
			}
			if !update.Message.IsCommand() {
				opts := strings.Split(update.Message.Text, " ")
				if len(opts) == 4 {
					waittop := WaitOTP{}
					db.Where(&WaitOTP{TID: update.Message.Chat.ID}).Last(&waittop)
					if waittop.TID != 0 && waittop.UserName != "" && !waittop.Expire {
						waittop.Expire = true
						db.Save(&waittop)
						token, err := GetToken(strings.Join(opts, "-"), waittop.UserName)
						if err != nil {
							log.Error(err)
						}
						db.Create(&Token{Token: token, TID: update.Message.Chat.ID, UserName: waittop.UserName})
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("添加 %s 成功 请点击List查看\n Add %s ! Please Use List Button to find it", waittop.UserName, waittop.UserName))
						msg.ReplyMarkup = numericKeyboard
						if _, err = bot.Send(msg); err != nil {
							log.Error(err)
						}
					}
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "暂未支持该输入\nDon't Support Yet")
					if _, err = bot.Send(msg); err != nil {
						log.Error(err)
					}
				}
			}

		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Error(err)
				continue
			}
			if update.CallbackData() == "LIST" {
				msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "请点击下方用户名，生成邮箱\nPlease Press The Following ID For Correspond UserName:")
				Tokens := []Token{}
				db.Where(&Token{TID: update.CallbackQuery.Message.Chat.ID}).Find(&Tokens)
				if len(Tokens) == 0 {
					msg.Text = "尚未添加用户名，请使用/add 用户名\n Don't Have any Username, Please Use /add Username to add one"
				} else {
					row := []tgbotapi.InlineKeyboardButton{}
					for index, value := range Tokens {
						msg.Text = msg.Text + fmt.Sprintf("\n(%d): `%s`", index, value.UserName)
						row = append(row, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", index), value.UserName))
					}
					row = append([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("⬅️", "Home")}, row...)
					key := tgbotapi.NewInlineKeyboardMarkup(row)
					msg.ReplyMarkup = &key
					msg.ParseMode = tgbotapi.ModeMarkdown
				}
				if _, err := bot.Send(msg); err != nil {
					log.Error(err)
					continue
				}
			} else if update.CallbackData() == "Home" {
				msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "欢迎食用DuckDuckGOEmailBOT\n,可以借助本BOT创建匿名邮箱，请先到DuckDuckGoEmail先行注册")
				msg.ReplyMarkup = &numericKeyboard
				if _, err := bot.Send(msg); err != nil {
					log.Error(err)
					continue
				}
			} else {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "空Token \nEmpty Token\n 请尝试重新绑定 Please Re Add UserName Again")
				token := Token{}
				db.Where(&Token{UserName: update.CallbackData()}).Find(&token)
				if token.Token != "" {
					email, err := Generate(token)
					if err != nil {
						log.Error(err)
					} else {
					}
					msg.Text = fmt.Sprintf("生成的邮箱是 `%s`", email)
					msg.ParseMode = tgbotapi.ModeMarkdown
				}
				if _, err := bot.Send(msg); err != nil {
					log.Error(err)
					continue
				}

			}

		}
	}
}

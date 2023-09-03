package telegram

import (
	"DobroBot/model"
	"DobroBot/store"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var mainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		//tgbotapi.NewKeyboardButton("о себе"),
		tgbotapi.NewKeyboardButton("благотворительность"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("о фонде"),
	),
)

// var userKeyboard = tgbotapi.NewReplyKeyboard(
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("сделать что-то"),
// 		tgbotapi.NewKeyboardButton("сделать еще что-то"),
// 		tgbotapi.NewKeyboardButton("назад"),
// 	),
// )

// var backButton = tgbotapi.NewOneTimeReplyKeyboard(
// 	tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("назад"),
// 	),
// )

type Telegram struct {
	s        store.Store
	wannaPay map[int64]struct{}
	ch       chan (model.Discont)
	//registration map[int]model.User
}

func NewTelegram(s store.Store, ch chan (model.Discont)) *Telegram {
	return &Telegram{
		s:  s,
		ch: ch,
	}
}

const textAbout = `Очень крутой фонд
всем советую!`

func (t *Telegram) Run(token string) {
	t.wannaPay = make(map[int64]struct{}, 10)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	go t.checkForDisconts(bot)
	//t.registration = make(map[int]model.User, 10)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		update.Message.IsCommand()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// То, что закомментированно это регистрация пользователя, я решил что это не очень то и нужно
		// в model тоже нужно поля недостающие добавить
		//
		// if k, ok := t.registration[int(update.Message.From.ID)]; ok {
		// 	if update.Message.Text == "назад" {
		// 		msg.ReplyMarkup = mainMenuKeyboard
		// 		delete(t.registration, int(update.Message.From.ID))
		// 		continue
		// 	}

		// 	if k.Firstname == "" {
		// 		k.Firstname = update.Message.Text
		// 		t.registration[int(update.Message.From.ID)] = k

		// 		msg.Text = "Пожалуйста введите фамилию:"
		// 		msg.ReplyMarkup = backButton

		// 	} else {
		// 		k.Lastname = update.Message.Text
		// 		delete(t.registration, int(update.Message.From.ID))

		// 		msg.Text = fmt.Sprintf("Успешная регистрация, %v %v", k.Lastname, k.Firstname)
		// 		msg.ReplyMarkup = userKeyboard

		// 		t.s.Add(k)
		// 	}
		// } else {
		if _, ok := t.wannaPay[update.Message.From.ID]; ok {
			delete(t.wannaPay, update.Message.From.ID)

			count, err := strconv.Atoi(update.Message.Text)
			if err != nil {
				log.Printf("cant parse %v: %v", update.Message.Text, err)
				continue
			}

			// какой то сервис оплаты надо чтоли
			t.s.UpdateDontes(int(update.Message.From.ID), count)
			msg.Text = "Спасибо за пожертвование;)"
		} else {
			user, err := t.s.Get(int(update.Message.From.ID))
			if err != nil {
				user = model.User{
					Id:       int(update.Message.From.ID),
					Username: update.Message.From.UserName,
				}
				t.s.Add(user)
			}
			switch update.Message.Text {
			case "/start":
				msg.Text = fmt.Sprintf("Привет, %v!", user.Username)
				msg.ReplyMarkup = mainMenuKeyboard
				// case "о себе":
				// 	t.registration[int(update.Message.From.ID)] = model.User{Id: int(update.Message.From.ID), Username: update.Message.From.UserName}
				// 	msg.Text = "Пожалуйста введите имя:"
				// 	msg.ReplyMarkup = backButton
			case "благотворительность":
				msg.Text = "Сколько вы готовы пожертвовать?"
				t.wannaPay[update.Message.From.ID] = struct{}{}
			case "о фонде":
				msg.Text = textAbout
			default:
				msg.Text = "Пожалуйста, используйте команды для взаимодействия с ботом."
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

	//}
}

func (t *Telegram) checkForDisconts(bot *tgbotapi.BotAPI) {
	for val := range t.ch {
		result, _ := t.s.GetAllWithDonates(val.ForDonate)

		for _, v := range result {
			msg := tgbotapi.NewMessage(int64(v), val.Text)
			msg.ParseMode = "HTML"

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
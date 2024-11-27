package Telegram

import (
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Telegram struct {
}

func (Telegram) Get(c *gin.Context) {
	// 初始化机器人，使用你自己的 Telegram Bot Token
	bot, err := tgbotapi.NewBotAPI("YOUR_BOT_API_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true // 开启调试模式，查看详细的日志

	// 获取更新（消息）
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	// 创建一个自定义菜单
	menu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Option 1"),
			tgbotapi.NewKeyboardButton("Option 2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Help"),
		),
	)

	// 监听更新
	for update := range updates {
		if update.Message == nil { // 如果没有消息，则跳过
			continue
		}

		// 如果收到 /start 命令，显示菜单
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "欢迎使用菜单机器人！请选择一个选项。")
			msg.ReplyMarkup = menu // 显示菜单
			bot.Send(msg)
		}

		// 用户选择菜单项时的处理
		if update.Message.Text == "Option 1" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你选择了 Option 1！")
			bot.Send(msg)
		} else if update.Message.Text == "Option 2" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你选择了 Option 2！")
			bot.Send(msg)
		} else if update.Message.Text == "Help" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "这是帮助信息：请选择一个菜单项来继续。")
			bot.Send(msg)
		}
	}
}

func (t Telegram) GetPull(context *gin.Context) {

}

//// SendTelegramMessage 发送Telegram机器人消息
//func SendTelegramMessage(text string) {
//	//os.Setenv("HTTPS_PROXY", proxyAddress)
//	token := "6966979790:AAFLTrtKbxdtDPXguvsRgscLGz2ftYaCFWY"
//	chatId := "-1002072043275"
//	_, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", token, chatId, text))
//
//	return
//}

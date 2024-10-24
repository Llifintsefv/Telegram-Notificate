package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

var (
	mu sync.Mutex
	ctxMap = make(map[int64]context.CancelFunc)
)



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err, "Ошибка создания бота")
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	defer bot.StopLongPolling()                                       

	for update := range updates {
		handleUpdate(bot, update) 
	}
}

func handleUpdate(bot *telego.Bot, update telego.Update) {
	if update.Message != nil {
		handleMessage(bot, update.Message)
	}
}

func handleMessage(bot *telego.Bot, message *telego.Message) {
	chatID := tu.ID(message.Chat.ID)
	go sendNotificate(bot, chatID)
	mu.Lock()
	defer mu.Unlock()

	switch message.Text {
	case "/startM":
		if cancel, ok := ctxMap[message.Chat.ID]; ok  {
			cancel()
		}
		ctx,cancel := context.WithCancel(context.Background())
		ctxMap[message.Chat.ID] = cancel

		GenerateVisual(ctx,bot,chatID)	
	case "/stop":
		if cancel, ok := ctxMap[message.Chat.ID]; ok  {
			cancel()
			delete(ctxMap, message.Chat.ID)
		}

	default:
		_, _ = bot.CopyMessage(
			tu.CopyMessage(
				chatID,
				chatID,
				message.MessageID,
			),
		)
	}
}

func sendNotificate(bot *telego.Bot, chatID telego.ChatID) {
	for {
			currentTime := time.Now()
			currentHour := currentTime.Hour()
			if currentHour == 13 || currentHour == 16 || currentHour == 23 {
				_, _ = bot.SendMessage(&telego.SendMessageParams{
					ChatID: chatID,
					Text:   "САМОЕ ВРЕМЯ ДЛЯ ПОЗИТИВНОЙ ПСИХОЛОГИИ ДРУЖИЩЕ",
				})
			}
			time.Sleep(60 * time.Minute)
		}
	}


	func GenerateVisual(ctx context.Context, bot *telego.Bot, chatID telego.ChatID) {
	for {
		select {
		case <-ctx.Done():
			// Контекст отменен, завершаем работу функции
			return
		default:
			// Ожидаем новые сообщения
			updates, _ := bot.GetUpdates(nil) 
			for _, update := range updates {
				if update.Message != nil && tu.ID(update.Message.Chat.ID) == chatID {
					// Обрабатываем сообщение, отправленное в чат, 
					// для которого запущена функция GenerateVisual

					switch update.Message.Text {
					case "/stop":
						if cancel, ok := ctxMap[update.Message.Chat.ID]; ok {
							cancel()
							delete(ctxMap, update.Message.Chat.ID)
						}
						return // Завершаем функцию GenerateVisual после /stop
					default:
						_, _ = bot.SendMessage(&telego.SendMessageParams{
					ChatID: chatID,
					Text:   "САМОЕ ВРЕМЯ ДЛЯ ПОЗИТИВНОЙ ПСИХОЛОГИИ ДРУЖИЩЕ",
				})
					}
				}
				
              bot.GetUpdates(&telego.GetUpdatesParams{Offset: update.UpdateID + 1})	
			}

			time.Sleep(time.Second) // Небольшая пауза, чтобы не нагружать API Telegram
		}
	}
}

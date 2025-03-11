package main

import (
	"asset_tracker/common"
	"asset_tracker/types"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// # TODO list
// 2. goroutine for each asset check in AssetWatcher
// 4. Implement panic?

func ValidateEnvs() []string {
	return []string{"TG_TOKEN", "DB_PATH", "TIINGO_API_TOKEN"}
}

func AssetWatcher(bot *tgbotapi.BotAPI) {
	var currentPrice float64
	var savedAssets []types.SavedAsset
	var err error
	for {
		savedAssets = common.DBRead()
		fmt.Println(savedAssets)
		for _, elem := range savedAssets {
			if elem.AssetType == "stock" {
				//#TODO some error handling maybe?
				currentPrice, err = common.FetchStockPrice(elem.AssetName)
			} else {
				currentPrice, err = common.FetchCryptoPrice(elem.AssetName)
			}
			if err != nil {
				log.Printf("Filed to fetch price for " + elem.AssetName)
				continue
			}
			if (elem.InitPrice < elem.TargetPrice && elem.TargetPrice < currentPrice) ||
				(elem.InitPrice > elem.TargetPrice && elem.TargetPrice > currentPrice) {

				common.DeleteDBRecord(elem.UserID, elem.AssetName)
				msg := tgbotapi.NewMessage(elem.UserID, "Reached target price ("+fmt.Sprint(elem.TargetPrice)+") for "+elem.AssetName)
				bot.Send(msg)
			}
		}
		time.Sleep(1 * time.Hour)
	}
}

func main() {
	var err error
	for _, env := range ValidateEnvs() {
		_, status := os.LookupEnv(env)
		if status == false {
			log.Printf("%s env is missing.", env)
			os.Exit(1)
		}
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	var userChats []types.SavedChat
	var savedAssets []types.SavedAsset
	var ChatPath string
	var ChatStage int8
	var replyText string
	var validationBool bool
	var usedKeyboard tgbotapi.ReplyKeyboardMarkup

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Create chan for telegram updates
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)
	go AssetWatcher(bot)

	for update := range updates {
		if update.Message != nil {
			UserName := update.Message.From.UserName

			ChatID := update.Message.Chat.ID

			Text := update.Message.Text

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			usedKeyboard = common.CompileDefaultKeyboard()

			if Text == "info" {
				common.EndChat(&userChats, ChatID)
			}
			ChatPath, ChatStage = common.FetchUser(&userChats, ChatID)
			fmt.Println(ChatPath, ChatStage) // debug

			// Straight up ignore any input longer than 30 symbols since atm there is so case when it might be needed
			if ChatPath != "" && len(Text) < 30 {

				common.EndChat(&userChats, ChatID)
				switch {
				case ChatPath == "/crypto" || ChatPath == "/stock":
					replyText = common.ConvertPrice(ChatPath[1:], Text)
				case ChatPath == "/delete":
					if common.ValidateUniqAsset(ChatID, Text) == false {
						common.DeleteDBRecord(ChatID, Text)
						replyText = "Asset successfully deleted"
					} else {
						replyText = "You are not tracking " + Text
					}
				// Watch conversation path
				case ChatPath == "/watch":
					switch ChatStage {
					case 0:
						if Text == "crypto" || Text == "stock" {
							userChats = append(userChats, types.SavedChat{ChatID, ChatPath, 1})
							savedAssets = append(savedAssets, types.SavedAsset{UserID: ChatID, AssetType: Text})
							replyText = "Provide asset name"
						} else {
							replyText = "Unknown asset type, please write crypto/stock or use buttons"
						}
					case 1:
						replyText, validationBool = common.UpdateAssetNPrice(&savedAssets, ChatID, Text)
						if validationBool == true {
							userChats = append(userChats, types.SavedChat{ChatID, ChatPath, 2})
						}
					case 2:
						// Avoid really big numbers
						if len(Text) > 20 {
							replyText = "The target price provided is too long, try again"
						} else {
							replyText = common.AssetTrargetPrice(&savedAssets, ChatID, Text)
						}
					}
				}
			} else {
				switch Text {
				case "/start":
					replyText = "Welcome :)"
				case "/crypto":
					userChats = append(userChats, types.SavedChat{ChatID, Text, 0})
					replyText = "Provide coin name"
				case "/stock":
					userChats = append(userChats, types.SavedChat{ChatID, Text, 0})
					replyText = "Provide stock name"
				case "/watch":
					userChats = append(userChats, types.SavedChat{ChatID, Text, 0})
					replyText = "Provide asset type (stock/crypto)"
					usedKeyboard = common.CompileTypeKeyboard()
				case "/delete":
					userChats = append(userChats, types.SavedChat{ChatID, Text, 0})
					replyText = "Provide asset name"
				case "/wipedb":
					if ChatID == 88770025 {
						common.WipeDB()
						replyText = "DB wiped"
					}
				case "/assets":
					replyText = common.CompileUserAssets(common.GetAssets(ChatID))
				case "info":
					replyText = "Info :)"
				default:
					replyText = "Command not recognized or the input was too long"

				}
			}

			fmt.Println(userChats)   // debug
			fmt.Println(savedAssets) // debug

			msg := tgbotapi.NewMessage(ChatID, replyText)
			msg.ReplyMarkup = usedKeyboard
			bot.Send(msg)
		}
	}
}

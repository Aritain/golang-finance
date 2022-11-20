package common

import (
    "github.com/piquette/finance-go/quote"
    "github.com/Syfaro/telegram-bot-api"
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "asset_tracker/types"
    "regexp"
    "strconv"
)

// Sort of const
func getCoinBaseUrl() []string {
	return []string{"https://api.coinbase.com/v2/prices/", "-USD/spot"}
}

/* Old way of handling stock prices
func getStockBaseUrl() []string {
	return []string{"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=", "&apikey=demo"}
}*/

func isTicker(ticker string) bool {
    var isTicker = regexp.MustCompile(`^[A-Za-z0-9]+$`)
    return isTicker.MatchString(ticker)
}

func isNumber(number string) bool {
    var isNumber = regexp.MustCompile(`^[0-9,.]+$`)
    return isNumber.MatchString(number)
}
    
func FetchCryptoPrice(coinName string) (float64, error) {
    var coinPrice float64
	var coinBase []string = getCoinBaseUrl()
	var cryptoUrl string = coinBase[0] + coinName + coinBase[1]

    res, err := http.Get(cryptoUrl)
    if err != nil {
        return 0.0, errors.New("Failed to fetch price, please check crypto name provided")
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
	fetchedJson := types.Crypto{}
	jsonError := json.Unmarshal(body, &fetchedJson)

	if jsonError != nil || fetchedJson.Data.Amount == ""  {
		return 0.0, errors.New("Failed to fetch price, please check crypto name provided")
	}

    coinPrice, _ = strconv.ParseFloat(fetchedJson.Data.Amount, 64)

    return coinPrice, nil
}

func FetchStockPrice(stockName string) (float64, error) {

    stockPrice, err := quote.Get(stockName)

    if err != nil {
        return 0.0, errors.New("Failed to fetch asset price, please check asset name provided")
    }
    return stockPrice.RegularMarketPrice, nil
}

/* Old way of handling stock prices
func FetchStockPrice(shareName string) string {
	var shareBase []string = getStockBaseUrl()
	var shareUrl string = shareBase[0] + shareName + shareBase[1]

    res, err := http.Get(shareUrl)
    if err != nil {
        return "Error while fetching the data"
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
	fetchedJson := types.Stocks{}
	jsonError := json.Unmarshal(body, &fetchedJson)

	if jsonError != nil || fetchedJson.GlobalQuote.Price == "" {
		return "Failed to fetch stock price, please check stock name provided"
	}

    return fetchedJson.GlobalQuote.Price
} */

func ValidateUniqAsset(userID int64, assetName string) bool {
    var userAssets []types.SavedAsset = GetAssets(userID)
    for _, elem := range userAssets {
        if elem.AssetName == assetName {
            return false
        }
    }
    return true
}


func CompileTypeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var keyboard = tgbotapi.NewReplyKeyboard(
    	tgbotapi.NewKeyboardButtonRow(
       		tgbotapi.NewKeyboardButton("stock"),
    	),
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("crypto"),
        ),
    )
    keyboard.OneTimeKeyboard = true
    return keyboard
}

func CompileDefaultKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var keyboard = tgbotapi.NewReplyKeyboard(
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("info"),
        ),
    )
    keyboard.OneTimeKeyboard = false
    return keyboard
}

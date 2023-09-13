package types

type Crypto struct {
    Data struct {
		Amount string
		Base string
		Currency string
	}
}

/* Old way of handling stocks
type Stocks struct {
	GlobalQuote struct {
		Symbol string `json:"01. symbol"`
		Open string `json:"02. open"`
		High string `json:"03. high"`
		Low string `json:"04. low"`
		Price string `json:"05. price"`
		Volume string `json: "06. volume"`
		LastTrade string `json:"07. latest trading day"`
		PrevClose string `json:"08. previous close"`
		Change string `json:"09. change"`
		ChangePercent string `json:"10. change percent"`
	} `json:"Global Quote"`
} */

// TiingoResponse struct
type TiingoResponse []struct {
	Price float64 `json:"last"`
}

type SavedChat struct {
	UserID int64
	ChatPath string
	ChatStage int8
}

type SavedAsset struct {
	UserID int64
	AssetType string
	AssetName string
	InitPrice float64
	TargetPrice float64
}

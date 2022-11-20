package common

import (
    "strconv"
    "asset_tracker/types"
	"fmt"
	"strings"
)

func FetchUser(userChats *[]types.SavedChat, userID int64) (string, int8) {
	for _, elem := range *userChats {
		if elem.UserID == userID {
			return elem.ChatPath, elem.ChatStage
		}
	}
	return "", 0
}

func EndChat(userChats *[]types.SavedChat, userID int64) {
	for index, elem := range *userChats {
		if elem.UserID == userID {
			// Magically remove element from the list
			*userChats = append((*userChats)[:index], (*userChats)[index+1:]...)
		}
	}	
}

func DeleteWatchCache(savedAssets *[]types.SavedAsset, userID int64) {
	for index, elem := range *savedAssets {
		if elem.UserID == userID {
			// Magically remove element from the list
			*savedAssets = append((*savedAssets)[:index], (*savedAssets)[index+1:]...)
		}
	}	
}

func UpdateAssetNPrice(savedAssets *[]types.SavedAsset, userID int64, assetName string) (string, bool) {
	var assetPrice float64
	var err error
	if isTicker(assetName) == false {
		DeleteWatchCache(savedAssets, userID)
		return "Error in asset name provided, please try again", false
	}
	for index, elem := range  *savedAssets {
		if elem.UserID == userID {
			if ValidateUniqAsset(userID, assetName) == false {
				DeleteWatchCache(savedAssets, userID)
				return "You are already tracking this asset", false
			}
			if elem.AssetType == "crypto" {
				assetPrice, err = FetchCryptoPrice(assetName)
			} else if elem.AssetType == "stock" {
				assetPrice, err = FetchStockPrice(assetName)
			}
			if err != nil {
				DeleteWatchCache(savedAssets, userID)
				return err.Error(), false
			}
			(*savedAssets)[index].AssetName = assetName
			(*savedAssets)[index].InitPrice = assetPrice

		}
	}
	return "You chose "+assetName+" (current price is "+fmt.Sprintf("%.2f", assetPrice)+"). Please provide target price for the asset", true
}

func AssetTrargetPrice(savedAssets *[]types.SavedAsset, userID int64, targetPrice string) string {
	var err error
	if isNumber(targetPrice) == false {
		DeleteWatchCache(savedAssets, userID)
		return "Error in target price provided, please try again"
	}
	for index, elem := range  *savedAssets {
		if elem.UserID == userID {
			(*savedAssets)[index].TargetPrice, err = strconv.ParseFloat(targetPrice, 64)
			writeAsset := types.SavedAsset{
				(*savedAssets)[index].UserID,
				(*savedAssets)[index].AssetType,
				(*savedAssets)[index].AssetName,
				(*savedAssets)[index].InitPrice,
				(*savedAssets)[index].TargetPrice,
			}
			DBWriteAsset(writeAsset)
			DeleteWatchCache(savedAssets, userID)
			if err != nil {
				return "Failed add target price, check input and try again"
			}
		}
	}
	return "Asset succesfully added!"
}


func CompileUserAssets(userAssets []types.SavedAsset) string {
    var compiledAssets string
    var fetchedPrice float64
    if len(userAssets) == 0 {
        compiledAssets = "Currently you are not tracking any assets"
    } else {
       compiledAssets =  "Your current tracked assets:\nName | Init Price | Target Price | Current Price\n\n"
    }
	for _, elem := range userAssets {
        fetchedPrice, _ = FetchStockPrice(elem.AssetName)
    	compiledAssets = compiledAssets + strings.ToUpper(elem.AssetName) +
    	" | " + fmt.Sprint(elem.InitPrice) +
    	" | " + fmt.Sprint(elem.TargetPrice) + 
        " | " + fmt.Sprintf("%.2f", fetchedPrice) + "\n"
    }
    return compiledAssets
}


func ConvertPrice(assetType string, assetName string) string {
    var err error
    var assetPrice float64
    if assetType == "crypto" {
        assetPrice, err = FetchCryptoPrice(assetName)
    } else if assetType == "stock" {
        assetPrice, err = FetchStockPrice(assetName)
    }
    if err != nil {
        return err.Error()
    }
    return fmt.Sprintf("%.2f", assetPrice)
}

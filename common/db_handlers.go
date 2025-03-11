package common

import (
	"asset_tracker/types"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func getDBName() string {
	return os.Getenv("DB_PATH")
}

// TODO: add error handling
func DBRead() []types.SavedAsset {
	var savedAssets []types.SavedAsset
	db, _ := sql.Open("sqlite3", getDBName())
	defer db.Close()
	rows, _ := db.Query("SELECT * FROM watch")
	defer rows.Close()

	var uid int
	for rows.Next() {
		var tA types.SavedAsset
		rows.Scan(&uid, &tA.UserID, &tA.AssetType, &tA.AssetName, &tA.InitPrice, &tA.TargetPrice)
		savedAssets = append(savedAssets, tA)
	}

	return savedAssets
}

func WipeDB() {
	db, _ := sql.Open("sqlite3", getDBName())
	defer db.Close()
	dbDriver, _ := db.Prepare("DELETE FROM watch")
	dbDriver.Exec()
}

func DeleteDBRecord(userID int64, AssetName string) {
	db, _ := sql.Open("sqlite3", getDBName())
	defer db.Close()
	dbDriver, _ := db.Prepare("DELETE FROM watch WHERE userid = ? AND assetname = ?")
	dbDriver.Exec(userID, AssetName)
}

func DBWriteAsset(assetCache types.SavedAsset) {
	db, _ := sql.Open("sqlite3", getDBName())
	defer db.Close()
	dbDriver, err := db.Prepare("INSERT INTO watch(userid, assettype, assetname, initprice, targetprice) values(?,?,?,?,?)")
	//TODO make me better (erroring)
	if err != nil {
		fmt.Println("prepare failed")
	}
	dbDriver.Exec(
		assetCache.UserID,
		assetCache.AssetType,
		assetCache.AssetName,
		assetCache.InitPrice,
		assetCache.TargetPrice,
	)
	//checkErr(errr)
}

func GetAssets(userID int64) []types.SavedAsset {
	var savedAssets []types.SavedAsset
	var userIDstr string = fmt.Sprint(userID)

	db, _ := sql.Open("sqlite3", getDBName())
	defer db.Close()
	rows, _ := db.Query("SELECT * FROM watch WHERE userid = " + userIDstr)
	defer rows.Close()

	var uid int
	for rows.Next() {
		var tA types.SavedAsset
		rows.Scan(&uid, &tA.UserID, &tA.AssetType, &tA.AssetName, &tA.InitPrice, &tA.TargetPrice)
		savedAssets = append(savedAssets, tA)
	}

	return savedAssets
}

// #TODO redo me
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

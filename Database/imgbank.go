package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func AddImg(message *discordgo.MessageCreate, category string, imgUrl string) error {

	// 1. Get the Gif Category
	imgCat, err := GetImgCategory(message.GuildID, category)
	if err != nil {
		return errors.New("unable to get gif category")
	}

	// 2. Get the Image Storage ID
	storageId, err := GetImgStorage(message.GuildID, imgUrl)
	if err != nil {
		return errors.New("unable to get gif category")
	}

	// 3. Insert the Link
	err = InsertImgGuildLink(storageId, imgCat.ID, message.GuildID, message.Author.ID)

	return err

}

// Returns the Img Category object from the database, using the Cache if available
func GetImgCategory(guildId string, category string) (cache.ImgCategory, error) {
	var imgCat cache.ImgCategory = cache.ImgCategory{
		ID:       0,
		Category: "",
	}

	//	=> Is it in the Cache?
	for _, cat := range cache.ImgCategories {
		if cat.Category == category {
			imgCat = cat
			break
		}
	}

	//	=> Is it in the Database?
	var ID sql.NullInt32
	var DbCategory sql.NullString
	err := db.QueryRow("SELECT ID, Category FROM ImgCategories WHERE Category = ?", strings.ToLower(category)).Scan(&ID, &DbCategory)
	noRows := false
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error(guildId, err)
			return imgCat, err
		} else {
			noRows = true
		}
	}

	if !noRows {
		if !ID.Valid {
			err = errors.New("db id was not valid")
			logger.Error(guildId, err)
			return imgCat, err
		}

		if !DbCategory.Valid {
			err = errors.New("db category was not valid")
			logger.Error(guildId, err)
			return imgCat, err
		}

		imgCat.ID = int(ID.Int32)
		imgCat.Category = DbCategory.String
		return imgCat, nil
	}

	// => Not in Database, create
	stmt, err := db.Prepare("INSERT INTO ImgCategories (Category) VALUES (?)")
	if err != nil {
		logger.Error(guildId, err)
		return imgCat, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(strings.ToLower(category))
	if err != nil {
		logger.Error(guildId, err)
		return imgCat, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error(guildId, err)
		return imgCat, err
	}

	imgCat.ID = int(lastInsertID)
	imgCat.Category = strings.ToLower(category)
	return imgCat, nil

}

func GetImgStorage(guildId string, imgUrl string) (int, error) {

	// 1. Is it in the Database?
	var ID sql.NullInt32
	err := db.QueryRow("SELECT ID FROM ImgStorage WHERE URL = ?", imgUrl).Scan(&ID)
	noRows := false
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error(guildId, err)
			return 0, err
		} else {
			noRows = true
		}
	}

	if !noRows {
		if !ID.Valid {
			err = errors.New("db id was not valid")
			logger.Error(guildId, err)
			return 0, err
		}
		return int(ID.Int32), nil
	}

	// 2. Not in Database, create
	stmt, err := db.Prepare("INSERT INTO ImgStorage (URL) VALUES (?)")
	if err != nil {
		logger.Error(guildId, err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(imgUrl)
	if err != nil {
		logger.Error(guildId, err)
		return 0, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error(guildId, err)
		return 0, err
	}

	return int(lastInsertID), nil

}

func InsertImgGuildLink(storageId int, categoryId int, guildId string, userId string) error {

	query := `INSERT INTO ImgGuildLink (StorageID, CategoryID, GuildID, AddedByUserID)
          SELECT * FROM (SELECT ? AS StorageID, ? AS CategoryID, ? AS GuildID, ? AS AddedByUserID) AS tmp
          WHERE NOT EXISTS (
              SELECT ID FROM ImgGuildLink WHERE StorageID = ? AND CategoryID = ? AND GuildID = ?
          ) LIMIT 1`

	_, err := db.Exec(query, storageId, categoryId, guildId, userId, storageId, categoryId, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	return nil

}

func GetRandomImage(guildId string, categoryId int) (string, error) {

	query := `SELECT S.URL FROM ImgGuildLink L
			INNER JOIN ImgStorage S ON S.ID = L.StorageID
			WHERE L.CategoryID = ? AND L.GuildID = ?
			ORDER BY RAND() LIMIT 1;`

	var dbImg sql.NullString

	err := db.QueryRow(query, categoryId, guildId).Scan(&dbImg)
	if err != nil {
		logger.Error(guildId, err)
		return "", err
	}

	if !dbImg.Valid {
		err = errors.New("invalid url returned from db")
		logger.Error(guildId, err)
		return "", err
	}

	return dbImg.String, nil

}

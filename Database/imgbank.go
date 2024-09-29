package database

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
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
	err = InsertImgGuildLink(storageId.ID, imgCat.ID, message.GuildID, message.Author.ID)

	return err
}

// Returns the Img Category object from the database, using the Cache if available
func GetImgCategory(guildId string, category string) (imgwork.ImgCategory, error) {
	var imgCat imgwork.ImgCategory = imgwork.ImgCategory{}

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
	err := Db.QueryRow("SELECT ID, Category FROM ImgCategories WHERE Category = ?", strings.ToLower(category)).Scan(&ID, &DbCategory)
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
	stmt, err := Db.Prepare("INSERT INTO ImgCategories (Category) VALUES (?)")
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

func GetImgStorage(guildId string, imgUrl string) (imgwork.ImgStorage, error) {
	var imgStorage imgwork.ImgStorage = imgwork.ImgStorage{
		LastChecked: helpers.GetNullDateTime(),
	}

	// 1. Is it in the Database?
	var ID sql.NullInt32
	var URL sql.NullString
	var LastCheck sql.NullTime

	err := Db.QueryRow("SELECT ID, URL, LastCheckDateTime FROM ImgStorage WHERE URL = ?", imgUrl).Scan(&ID, &URL, &LastCheck)
	noRows := false
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error(guildId, err)
			return imgStorage, err
		} else {
			noRows = true
		}
	}

	if !noRows {
		if !ID.Valid {
			err = errors.New("db id was not valid")
			logger.Error(guildId, err)
			return imgStorage, err
		}

		if !URL.Valid {
			err = errors.New("db url was not valid")
			logger.Error(guildId, err)
			return imgStorage, err
		}

		imgStorage.ID = int(ID.Int32)
		imgStorage.URL = URL.String
		if LastCheck.Valid {
			imgStorage.LastChecked = LastCheck.Time
		} else {
			imgStorage.LastChecked = helpers.GetNullDateTime()
		}

		return imgStorage, nil
	}

	// 2. Not in Database, create
	stmt, err := Db.Prepare("INSERT INTO ImgStorage (URL) VALUES (?)")
	if err != nil {
		logger.Error(guildId, err)
		return imgStorage, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(imgUrl)
	if err != nil {
		logger.Error(guildId, err)
		return imgStorage, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error(guildId, err)
		return imgStorage, err
	} else if lastInsertID == 0 {
		logger.Error(guildId, errors.New("last insert id was 0"))
		return imgStorage, err
	}

	imgStorage.ID = int(lastInsertID)
	imgStorage.URL = imgUrl
	return imgStorage, nil
}

func GetImgGuildLink(guildId string, category imgwork.ImgCategory, storage imgwork.ImgStorage) (imgwork.ImgGuildLink, error) {
	var imgGuildLink imgwork.ImgGuildLink = imgwork.ImgGuildLink{}

	// 1. Is it in the Database?
	var ID sql.NullInt32
	var StorageID, CategoryID sql.NullInt32
	var GuildID, AddedByUserID sql.NullString
	var AddedDateTime sql.NullTime

	query := `SELECT ID, StorageID, CategoryID, GuildID, AddedByUserID, AddedDateTime 
			FROM ImgGuildLink WHERE StorageID = ? AND CategoryID = ? AND GuildID = ?
			`
	err := Db.QueryRow(query, storage.ID, category.ID, guildId).Scan(&ID, &StorageID, &CategoryID, &GuildID, &AddedByUserID, &AddedDateTime)
	noRows := false
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error(guildId, err)
			return imgGuildLink, err
		} else {
			noRows = true
		}
	}

	if !noRows {
		if !ID.Valid {
			err = errors.New("db id was not valid")
			logger.Error(guildId, err)
			return imgGuildLink, err
		}

		if !StorageID.Valid {
			err = errors.New("db storageId was not valid")
			logger.Error(guildId, err)
			return imgGuildLink, err
		}

		if !CategoryID.Valid {
			err = errors.New("db categoryId was not valid")
			logger.Error(guildId, err)
			return imgGuildLink, err
		}

		if !GuildID.Valid {
			err = errors.New("db guildId was not valid")
			logger.Error(guildId, err)
			return imgGuildLink, err
		}

		if !AddedByUserID.Valid {
			err = errors.New("db addedByUserId was not valid")
			logger.Error(guildId, err)
			return imgGuildLink, err
		}

		if !AddedDateTime.Valid {
			err = errors.New("db addedDateTime was not valid")
			logger.Error(guildId, err)
			return imgGuildLink, err
		}

		imgGuildLink.ID = int(ID.Int32)
		imgGuildLink.Storage = storage
		imgGuildLink.Category = category
		imgGuildLink.GuildID = GuildID.String
		imgGuildLink.AddedByUserID = AddedByUserID.String

		if AddedDateTime.Valid {
			imgGuildLink.AddedDateTime = AddedDateTime.Time
		} else {
			imgGuildLink.AddedDateTime = helpers.GetNullDateTime()
		}

		return imgGuildLink, nil
	}

	return imgGuildLink, errors.New("no rows found")
}

func InsertImgGuildLink(storageId int, categoryId int, guildId string, userId string) error {
	query := `INSERT INTO ImgGuildLink (StorageID, CategoryID, GuildID, AddedByUserID)
          SELECT * FROM (SELECT ? AS StorageID, ? AS CategoryID, ? AS GuildID, ? AS AddedByUserID) AS tmp
          WHERE NOT EXISTS (
              SELECT ID FROM ImgGuildLink WHERE StorageID = ? AND CategoryID = ? AND GuildID = ?
          ) LIMIT 1`

	_, err := Db.Exec(query, storageId, categoryId, guildId, userId, storageId, categoryId, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	return nil
}

func DeleteGuildLink(guildLink imgwork.ImgGuildLink) error {
	// 1. Delete from the Link table
	query := `DELETE FROM ImgGuildLink WHERE ID = ?`
	_, err := Db.Exec(query, guildLink.ID)
	if err != nil {
		logger.Error(guildLink.GuildID, err)
		return err
	}

	// 2. Tidy up orphaned Storage items
	return TidyImgStorage(guildLink.GuildID)
}

func GetRandomImage(guildId string, categoryId int) (string, error) {
	query := `SELECT S.URL FROM ImgGuildLink L
			INNER JOIN ImgStorage S ON S.ID = L.StorageID
			WHERE L.CategoryID = ? AND L.GuildID = ?
			ORDER BY RAND() LIMIT 1;`

	var dbImg sql.NullString

	err := Db.QueryRow(query, categoryId, guildId).Scan(&dbImg)
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

func TidyImgStorage(guildId string) error {
	queryStorage := `SELECT * FROM ImgStorage
	WHERE ID NOT IN (
		SELECT StorageID FROM ImgGuildLink
	)`
	_, err := Db.Exec(queryStorage)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	return err
}

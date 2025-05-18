package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type UserWowStats struct {
	MaxWow        int
	MaxWowUpdated time.Time
	Effects       []UserWowStatEffect
}

type UserWowStatEffect struct {
	Type          string
	Name          string
	Emoji         string
	Count         int
	LastTriggered time.Time
}

func GetUserWowStats(guildId string, userId string) (*UserWowStats, error) {
	query := `SELECT S.MaxWow, S.MaxWowUpdated, E.EffectType, E.EffectName, E.EffectEmoji, E.Count, E.LastTriggered 
			FROM WowStats S LEFT JOIN WowEffects E ON E.UserID = S.UserID AND E.GuildID = S.GuildID
			WHERE S.GuildID = ? AND S.UserID = ? AND E.EffectName NOT LIKE 'Pokémon%' AND E.EffectName NOT LIKE 'A Wild%' ORDER BY EffectType, E.Count DESC`
	rows, err := Db.Query(query, guildId, userId)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer func(g string) {
		err := rows.Close()
		if err != nil {
			logger.Error(g, err)
		}
	}(guildId)

	var r UserWowStats

	statsObtained := false
	for rows.Next() {

		var MaxWow, Count sql.NullInt32
		var EffectType, EffectName, EffectEmoji sql.NullString
		var MaxWowUpdated, LastTriggered sql.NullTime

		err := rows.Scan(&MaxWow, &MaxWowUpdated, &EffectType, &EffectName, &EffectEmoji, &Count, &LastTriggered)
		if err != nil {
			logger.Error(guildId, err)
			return nil, err
		}

		if !statsObtained {
			if !MaxWow.Valid || !MaxWowUpdated.Valid {
				return nil, errors.New("stats invalid")
			}
			r.MaxWow = int(MaxWow.Int32)
			r.MaxWowUpdated = MaxWowUpdated.Time

			statsObtained = true
		}

		if !EffectType.Valid || !EffectName.Valid || !Count.Valid || !LastTriggered.Valid {
			continue
		}

		r.Effects = append(r.Effects, UserWowStatEffect{
			Type:          EffectType.String,
			Name:          EffectName.String,
			Emoji:         EffectEmoji.String,
			Count:         int(Count.Int32),
			LastTriggered: LastTriggered.Time,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &r, nil
}

type WowLeaderboard struct {
	UserID        string
	MaxWow        int
	MaxWowUpdated time.Time
}

func GetWowLeaderboard(guildId string) ([]WowLeaderboard, error) {
	var wows []WowLeaderboard

	var rows *sql.Rows
	var err error
	query := `SELECT UserID, MaxWow, MaxWowUpdated FROM WowStats WHERE GuildID = ? ORDER BY MaxWow DESC LIMIT 6`
	rows, err = Db.Query(query, guildId)

	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer func(g string) {
		err := rows.Close()
		if err != nil {
			logger.Error(g, err)
		}
	}(guildId)

	// Iterate over the rows
	for rows.Next() {

		var MaxWow sql.NullInt32
		var UserID sql.NullString
		var MaxWowUpdated sql.NullTime

		err := rows.Scan(&UserID, &MaxWow, &MaxWowUpdated)
		if err != nil {
			return nil, err
		}

		if !UserID.Valid {
			logger.Error(guildId, errors.New("invalid userid"))
			continue
		}

		if !MaxWow.Valid {
			logger.Error(guildId, errors.New("invalid maxwow"))
			continue
		}

		if !MaxWowUpdated.Valid {
			logger.Error(guildId, errors.New("invalid maxwowupdated"))
			continue
		}

		wows = append(wows, WowLeaderboard{
			UserID:        UserID.String,
			MaxWow:        int(MaxWow.Int32),
			MaxWowUpdated: MaxWowUpdated.Time,
		})
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return wows, nil

}

func GetUserWowBalance(guildId string, userId string) (int, error) {
	var balance sql.NullInt64
	query := "SELECT Currency FROM WowCurrency WHERE GuildID = ? AND UserID = ? LIMIT 1"
	err := Db.QueryRow(query, guildId, userId).Scan(&balance)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return 0, nil
		}
		logger.Error(guildId, err)
		return 0, err
	}

	ret := 0
	if balance.Valid {
		ret = int(balance.Int64)
	}

	return ret, nil
}

func UpdateWowBalance(guildId string, userId string, amount int, add bool) error {
	char := "-"
	if add {
		char = "+"
	}
	// 1. Try an Update
	updateQuery := fmt.Sprintf(`UPDATE WowCurrency
		SET Currency = Currency%s%d, LastUpdated = NOW()
		WHERE GuildID = ? AND UserID = ?`, char, amount)

	result, err := Db.Exec(updateQuery, guildId, userId)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	// 2. Check if we affected a row
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	if rowsAffected == 0 {

		insAmount := amount
		if !add {
			insAmount = -amount
		}

		// 3. If not, insert the new row
		insertQuery := `
			INSERT INTO WowCurrency
			(GuildID, UserID, Currency, LastUpdated)
			VALUES
			(?, ?, ?, NOW())
		`
		_, err = Db.Exec(insertQuery, guildId, userId, insAmount)
		if err != nil {
			logger.Error(guildId, err)
			return err
		}
	}

	return nil
}

func CountWowPurchase(guildId string, userId string, shopItemId int) {
	res, err := Db.Exec(`UPDATE WowPurchases SET Count = Count + 1, LastPurchased = NOW() WHERE GuildID = ? AND UserID = ? AND ShopItemID =?`,
		guildId, userId, shopItemId)
	if err != nil {
		logger.ErrorText("WOW", "Error counting Wow Purchase for [%s|%s|%d]: %s", guildId, userId, shopItemId, err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.ErrorText("WOW", "Error getting affected rows for [%s|%s|%d]: %s", guildId, userId, shopItemId, err)
		return
	}

	if rowsAffected == 0 {
		_, err = Db.Exec(`INSERT INTO WowPurchases (GuildID, UserID, ShopItemID) VALUES (?, ?, ?)`,
			guildId, userId, shopItemId)
		if err != nil {
			logger.ErrorText("WOW", "Error inserting currency for [%s|%s|%d]: %s", guildId, userId, shopItemId, err)
			return
		}
	}
}

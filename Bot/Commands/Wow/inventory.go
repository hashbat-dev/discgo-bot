package wow

import (
	"database/sql"
	"fmt"
	"slices"
	"sync"
	"time"

	database "github.com/hashbat-dev/discgo-bot/Database"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	dataInventoryLock   sync.RWMutex
	dataUserInventories = make(map[string][]InventoryItem)
)

type InventoryItem struct {
	ShopItem   WowShopItem
	Expires    time.Time
	DatabaseID int
}

type InventoryDatabase struct {
	ID         int
	GuildID    string
	UserID     string
	ShopItemID int
	ExpiryTime time.Time
}

func getDataWowInventories() {
	// First check if any items in the inventory have duplicate IDs
	var ids []int
	for _, item := range ShopItems {
		if slices.Contains(ids, item.ID) {
			panic("cannot have WowShopItems with duplicate IDs")
		}
		ids = append(ids, item.ID)
	}

	// Now get inventories
	query := `SELECT
					ID, GuildID, UserID, ShopItemID, ExpiryTime
				FROM
					WowInventory`

	var dbItems []InventoryDatabase

	rows, err := database.Db.Query(query)
	if err != nil {
		logger.Error("WOW", err)
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error("WOW", err)
		}
	}()

	for rows.Next() {
		var rowId, rowShopItemId sql.NullInt32
		var rowGuildId, rowUserId sql.NullString
		var rowExpiryTime sql.NullTime

		err := rows.Scan(&rowId, &rowGuildId, &rowUserId, &rowShopItemId, &rowExpiryTime)
		if err != nil {
			continue
		}

		if !rowId.Valid || !rowGuildId.Valid || !rowUserId.Valid || !rowShopItemId.Valid || !rowExpiryTime.Valid {
			continue
		}

		dbItems = append(dbItems, InventoryDatabase{
			ID:         int(rowId.Int32),
			GuildID:    rowGuildId.String,
			UserID:     rowUserId.String,
			ShopItemID: int(rowShopItemId.Int32),
			ExpiryTime: rowExpiryTime.Time,
		})
	}

	if err = rows.Err(); err != nil {
		logger.ErrorText("WOW", "Error in row loop: %s", err.Error())
	}

	var newCache = make(map[string][]InventoryItem)
	for _, db := range dbItems {
		var item WowShopItem
		found := false
		for _, i := range ShopItems {
			if i.ID == db.ShopItemID {
				item = i
				found = true
				break
			}
		}

		if !found {
			_ = deleteWowItemFromDb("WOW", db.ID)
			logger.ErrorText("WOW", "Couldn't find ShopItemID: %d", db.ShopItemID)
			continue
		}

		if !item.OneTimeUse && time.Now().After(db.ExpiryTime) {
			_ = deleteWowItemFromDb("WOW", db.ID)
			continue
		}

		newCache[db.GuildID+"|"+db.UserID] = append(newCache[db.GuildID+"|"+db.UserID], InventoryItem{
			ShopItem:   item,
			Expires:    db.ExpiryTime,
			DatabaseID: db.ID,
		})
	}

	dataInventoryLock.Lock()
	dataUserInventories = newCache
	dataInventoryLock.Unlock()
	logger.Info("WOW", "Imported %d Wow Inventory items into the Cache", len(dbItems))
}

func addToWowInventory(guildId string, userId string, item WowShopItem) error {
	// Create inventory item
	invItem := InventoryItem{
		ShopItem: item,
	}
	if item.OneTimeUse {
		invItem.Expires = time.Now().AddDate(5, 0, 0)
	} else {
		invItem.Expires = time.Now().Add(item.Duration)
	}

	// Add to DB
	result, err := database.Db.Exec(`INSERT INTO WowInventory (GuildID, UserID, ShopItemID, ExpiryTime) VALUES (?, ?, ?, ?)`,
		guildId, userId, item.ID, invItem.Expires)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	invItem.DatabaseID = int(insertedID)

	// Add to Cache
	key := fmt.Sprintf("%s|%s", guildId, userId)
	dataInventoryLock.Lock()
	dataUserInventories[key] = append(dataUserInventories[key], invItem)
	dataInventoryLock.Unlock()

	logger.Info(guildId, "Added WowItem. GuildID: %s, UserID: %s, ID: %d, DBID: %d", guildId, userId, item.ID, insertedID)
	return nil
}

func wowInventoryItemCount(guildId string, userId string, itemId int) int {
	var items []InventoryItem
	key := guildId + "|" + userId

	dataInventoryLock.RLock()
	items = dataUserInventories[key]
	dataInventoryLock.RUnlock()

	count := 0
	for _, i := range items {
		if i.ShopItem.ID == itemId {
			count++
		}
	}

	return count
}

func deleteFromWowInventory(guildId string, userId string, dbId int) {
	var userInv []InventoryItem
	key := fmt.Sprintf("%s|%s", guildId, userId)
	found := true

	// Get the User's Cached Inventory
	dataInventoryLock.RLock()
	if inv, exists := dataUserInventories[key]; exists {
		userInv = inv
	} else {
		found = false
	}
	dataInventoryLock.RUnlock()

	if !found {
		logger.ErrorText(guildId, "Unable to find Inventory for GuildID: %s,  UserID: %s", guildId, userId)
		return
	}

	// Find the Item
	var newInv []InventoryItem
	found = false
	for i, inv := range userInv {
		if inv.DatabaseID == dbId {
			newInv = append(userInv[:i], userInv[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		logger.ErrorText(guildId, "Unable to find DBID in Inventory for GuildID: %s,  UserID: %s", guildId, userId)
		return
	}

	// Delete from Database
	err := deleteWowItemFromDb(guildId, dbId)
	if err != nil {
		return
	}

	// Update Cache
	dataInventoryLock.Lock()
	dataUserInventories[key] = newInv
	dataInventoryLock.Unlock()

	logger.Debug(guildId, "Removed WowItem. GuildID: %s, UserID: %s, DBID: %d", guildId, userId, dbId)
}

func deleteWowItemFromDb(guildId string, dbId int) error {
	_, err := database.Db.Exec(`DELETE FROM WowInventory WHERE ID = ?`, dbId)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}
	return nil
}

package wow

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	database "github.com/hashbat-dev/discgo-bot/Database"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	dbWowStats    map[string]DbWowStats   = make(map[string]DbWowStats)
	dbEffectStats map[string]DbWowEffects = make(map[string]DbWowEffects)

	recordLocks sync.Map
)

type DbWowStats struct {
	ID            int
	GuildID       string
	UserID        string
	MaxWow        int
	MaxWowUpdated time.Time
}

type DbWowEffects struct {
	ID            int
	GuildID       string
	UserID        string
	EffectType    string
	EffectName    string
	EffectEmoji   string
	Count         int
	LastTriggered time.Time
}

func postToDatabase(i *Generation) {
	cacheKey := fmt.Sprintf("%s|%s", i.Message.GuildID, i.Message.Author.ID)
	lock := getLock(cacheKey)
	// === dbo.WowStats =====================================================
	lock.Lock()
	recordStats(cacheKey, i)
	lock.Unlock()

	// === dbo.WowEffects =====================================================
	// 1. Roll Effects
	for _, roll := range i.DiceRolls {
		if len(roll.Effects) > 0 {
			for _, effect := range roll.Effects {
				lock.Lock()
				recordEffect(cacheKey, "roll", i, effect)
				lock.Unlock()
			}
		}
	}

	// 2. Static Effects
	for _, effect := range i.StaticEffects {
		lock.Lock()
		recordEffect(cacheKey, "static", i, effect)
		lock.Unlock()
	}
}

func getLock(key string) *sync.Mutex {
	actual, _ := recordLocks.LoadOrStore(key, &sync.Mutex{})
	return actual.(*sync.Mutex)
}

func recordStats(cacheKey string, i *Generation) {
	if _, ok := dbWowStats[cacheKey]; !ok {
		// -> Not in cache, try get the DB entry and if it doesn't exist insert a new one.
		new, err := dbTryGetWowStats(i.Message.GuildID, i.Message.Author.ID)
		if err != nil {
			return
		}
		if new == nil {
			new = &DbWowStats{
				GuildID:       i.Message.GuildID,
				UserID:        i.Message.Author.ID,
				MaxWow:        i.OCount,
				MaxWowUpdated: time.Now(),
			}
		}
		dbWowStats[cacheKey] = *new
	}

	if val, ok := dbWowStats[cacheKey]; ok {
		// -> Get the WowStats out of the cache and update them
		if i.OCount >= val.MaxWow {
			val.MaxWow = i.OCount
			val.MaxWowUpdated = time.Now()
			err := dbUpdateWowStats(cacheKey, val)
			if err != nil {
				return
			}
		}
	} else {
		logger.Error(i.Message.GuildID, errors.New("unable to get WowStats from cache"))
		return
	}
}

func recordEffect(cacheKey string, effectType string, i *Generation, effect Effect) {
	effectCacheKey := fmt.Sprintf("%s|%s|%s", cacheKey, effectType, effect.Name)
	if _, ok := dbEffectStats[effectCacheKey]; !ok {
		// -> Not in cache, try get the DB entry and if it doesn't exist insert a new one.
		new, err := dbTryGetWowEffects(i.Message.GuildID, i.Message.Author.ID, effectType, effect.Name)
		if err != nil {
			return
		}
		if new == nil {
			new = &DbWowEffects{
				GuildID:       i.Message.GuildID,
				UserID:        i.Message.Author.ID,
				EffectType:    effectType,
				EffectName:    effect.Name,
				EffectEmoji:   effect.Emoji,
				Count:         0,
				LastTriggered: time.Now(),
			}
		}
		dbEffectStats[effectCacheKey] = *new
	}

	if val, ok := dbEffectStats[effectCacheKey]; ok {
		// -> Get the WowEffect out of the cache and update it
		val.Count++
		val.LastTriggered = time.Now()
		err := dbUpdateWowEffects(effectCacheKey, val)
		if err != nil {
			return
		}
	} else {
		logger.Error(i.Message.GuildID, fmt.Errorf("unable to get WowEffect from cache: %s", effectCacheKey))
		return
	}
}

func dbTryGetWowStats(guildId string, userId string) (*DbWowStats, error) {
	var ID, MaxWow sql.NullInt32
	var GuildID, UserID sql.NullString
	var MaxWowUpdated sql.NullTime
	var r DbWowStats
	query := `SELECT ID, GuildID, UserID, MaxWow, MaxWowUpdated FROM WowStats WHERE GuildID = ? AND UserID = ?`
	err := database.Db.QueryRow(query, guildId, userId).Scan(&ID, &GuildID, &UserID, &MaxWow, &MaxWowUpdated)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			logger.Error(guildId, err)
			return nil, err
		} else {
			return nil, nil
		}
	}
	if !ID.Valid {
		return nil, err
	} else {
		r.ID = int(ID.Int32)
	}
	if GuildID.Valid {
		r.GuildID = GuildID.String
	}
	if UserID.Valid {
		r.UserID = UserID.String
	}
	if MaxWow.Valid {
		r.MaxWow = int(MaxWow.Int32)
	}
	if MaxWowUpdated.Valid {
		r.MaxWowUpdated = MaxWowUpdated.Time
	}

	return &r, nil
}

func dbUpdateWowStats(cacheKey string, s DbWowStats) error {
	id, err := database.Upsert(s.GuildID, "WowStats", "ID", s.ID, "GuildID", s.GuildID, "UserID", s.UserID, "MaxWow", s.MaxWow, "MaxWowUpdated", s.MaxWowUpdated)
	if err != nil {
		logger.Error(s.GuildID, err)
		return err
	}

	if s.ID != int(id) {
		s.ID = int(id)
	}
	dbWowStats[cacheKey] = s

	return nil
}

func dbTryGetWowEffects(guildId string, userId string, effectType string, effectName string) (*DbWowEffects, error) {
	var ID, Count sql.NullInt32
	var GuildID, UserID, EffectType, EffectName, EffectEmoji sql.NullString
	var LastTriggered sql.NullTime
	var r DbWowEffects
	query := `SELECT ID, GuildID, UserID, EffectType, EffectName, EffectEmoji, Count, LastTriggered FROM WowEffects WHERE GuildID = ? AND UserID = ? AND EffectType = ? AND EffectName = ?`
	err := database.Db.QueryRow(query, guildId, userId, effectType, effectName).Scan(&ID, &GuildID, &UserID, &EffectType, &EffectName, &EffectEmoji, &Count, &LastTriggered)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			logger.Error(guildId, err)
			return nil, err
		} else {
			return nil, nil
		}
	}
	if !ID.Valid {
		return nil, err
	} else {
		r.ID = int(ID.Int32)
	}
	if GuildID.Valid {
		r.GuildID = GuildID.String
	}
	if UserID.Valid {
		r.UserID = UserID.String
	}
	if EffectType.Valid {
		r.EffectType = EffectType.String
	}
	if EffectName.Valid {
		r.EffectName = EffectName.String
	}
	if EffectEmoji.Valid {
		r.EffectEmoji = EffectEmoji.String
	}
	if Count.Valid {
		r.Count = int(Count.Int32)
	}
	if LastTriggered.Valid {
		r.LastTriggered = LastTriggered.Time
	}

	return &r, nil
}

func dbUpdateWowEffects(cacheKey string, s DbWowEffects) error {
	id, err := database.Upsert(s.GuildID, "WowEffects", "ID", s.ID, "GuildID", s.GuildID, "UserID", s.UserID, "EffectType", s.EffectType,
		"EffectName", s.EffectName, "EffectEmoji", s.EffectEmoji, "Count", s.Count, "LastTriggered", s.LastTriggered)
	if err != nil {
		logger.Error(s.GuildID, err)
		return err
	}

	if s.ID != int(id) {
		s.ID = int(id)
	}
	dbEffectStats[cacheKey] = s

	return nil
}

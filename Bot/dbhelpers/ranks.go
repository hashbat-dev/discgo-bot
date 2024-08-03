package dbhelper

import (
	"database/sql"
	"strings"
)

// Returns: (Rank Text, Rank Image, Bool: Is this a new Rank?, Error)
func GetRankFromCount(count int, category string) (string, string, bool, error) {

	var rankStart, rankEnd, rankImage sql.NullString
	var newRank sql.NullInt32

	err := GetDB().QueryRow(`
		SELECT (SELECT RankStart FROM RankStart WHERE RankCount <= ? ORDER BY RankCount DESC LIMIT 1) AS RankStart,
		(SELECT RankEnd FROM RankEnd WHERE RankCategory = @Category AND RankCount <= ? ORDER BY RankCount DESC LIMIT 1) AS RankEnd,
		(SELECT RankImage FROM RankImage WHERE RankCount >= ? LIMIT 1) AS RankImage,
		((SELECT COUNT(*) FROM RankStart WHERE RankCount = ?) + 
		(SELECT COUNT(*) FROM RankEnd WHERE RankCategory = ? AND RankCount = ?)) AS NewRank
	`, count, count, count, count, category, count).Scan(&rankStart, &rankEnd, &rankImage, &newRank)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			err = nil
		}
		return "", "", false, err
	}

	if !rankStart.Valid || !rankEnd.Valid {
		return "", "", false, nil
	}

	rank := rankStart.String + " " + rankEnd.String
	image := ""
	if rankImage.Valid {
		image = rankImage.String
	}

	isNewRank := false
	if newRank.Valid {
		isNewRank = newRank.Int32 > 0
	}

	return rank, image, isNewRank, nil

}

// Returns: (Rank Text, Rank Image, Error)
func GetRankFromUser(user string, category string) (int, string, string, error) {

	var rankStart, rankEnd, rankImage sql.NullString
	var rankCount sql.NullInt32

	err := GetDB().QueryRow(`SELECT TrackWords.Count FROM TrackWords WHERE UserID = ? AND Phrase = ? LIMIT 1`, user, category).Scan(&rankCount)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			err = nil
		}
		return 0, "", "", err
	}

	count := 0
	if rankCount.Valid {
		count = int(rankCount.Int32)
	}

	err = GetDB().QueryRow(`
	SELECT (SELECT RankStart FROM RankStart WHERE RankCount <= ? ORDER BY RankCount DESC LIMIT 1) AS RankStart,
	(SELECT RankEnd FROM RankEnd WHERE RankCategory = ? AND RankCount <= ? ORDER BY RankCount DESC LIMIT 1) AS RankEnd,
	(SELECT RankImage FROM RankImage WHERE RankCount >= ? LIMIT 1) AS RankImage
	`, count, category, count, count).Scan(&rankStart, &rankEnd, &rankImage)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			err = nil
		}
		return 0, "", "", err
	}

	if !rankStart.Valid || !rankEnd.Valid {
		return count, "", "", nil
	}

	rank := rankStart.String + " " + rankEnd.String
	image := ""
	if rankImage.Valid {
		image = rankImage.String
	}

	return count, rank, image, nil

}

func GetDistinctRanks() ([]string, error) {

	var retArray []string
	query := "SELECT DISTINCT RankCategory FROM BotDB.RankEnd ORDER BY RankCategory;"
	rows, err := GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {

		var dbRet sql.NullString

		err := rows.Scan(&dbRet)
		if err != nil {
			return nil, err
		} else if dbRet.Valid {
			retArray = append(retArray, dbRet.String)
		}

	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return retArray, nil

}

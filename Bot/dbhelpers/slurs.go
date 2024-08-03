package dbhelper

import (
	"database/sql"
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/datasets"
)

func GetARandomSlur() (datasets.SlurEntry, error) {

	var id sql.NullInt32
	var slur sql.NullString
	var slurTarget sql.NullString
	var slurDescription sql.NullString

	var retSet datasets.SlurEntry

	err := db.QueryRow("SELECT ID, Slur, SlurTarget, SlurDescription FROM SlurBank ORDER BY RAND() LIMIT 1; ").Scan(&id, &slur, &slurTarget, &slurDescription)
	if err != nil {
		audit.Error(err)
	} else {
		if id.Valid {
			retSet.ID = int(id.Int32)
		}
		if slur.Valid {
			retSet.Slur = slur.String
		}
		if slurTarget.Valid {
			retSet.SlurTarget = slurTarget.String
		}
		if slurDescription.Valid {
			retSet.SlurDescription = slurDescription.String
		}
	}

	return retSet, err
}

func GetARandomNationality() (datasets.NationalityEntry, error) {

	var id sql.NullInt32
	var nationality sql.NullString

	var retSet datasets.NationalityEntry

	err := db.QueryRow("SELECT ID, Nationality FROM Nationalities ORDER BY RAND() LIMIT 1; ").Scan(&id, &nationality)
	if err != nil {
		audit.Error(err)
	} else {
		if id.Valid {
			retSet.ID = int(id.Int32)
		}
		if nationality.Valid {
			retSet.Nationality = nationality.String
		}
	}

	return retSet, err
}

func GetARandomJobTitle() (datasets.JobTitleEntry, error) {

	var id sql.NullInt32
	var jobTitle sql.NullString

	var retSet datasets.JobTitleEntry

	err := db.QueryRow("SELECT ID, JobTitle FROM JobTitles ORDER BY RAND() LIMIT 1; ").Scan(&id, &jobTitle)
	if err != nil {
		audit.Error(err)
	} else {
		if id.Valid {
			retSet.ID = int(id.Int32)
		}
		if jobTitle.Valid {
			retSet.JobTitle = jobTitle.String
		}
	}

	return retSet, err
}

func GetSlurDefinition(slur string) (string, string, error) {
	var target, definition sql.NullString
	err := GetDB().QueryRow("SELECT SlurTarget, SlurDescription FROM SlurBank WHERE Slur = ? LIMIT 1;", strings.TrimSpace(slur)).Scan(&target, &definition)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		err = nil
	} else if err != nil {
		return "", "", err
	}

	retTarget := ""
	retDefinition := ""

	if target.Valid {
		retTarget = target.String
	}

	if definition.Valid {
		retDefinition = definition.String
	}
	return retTarget, retDefinition, err
}

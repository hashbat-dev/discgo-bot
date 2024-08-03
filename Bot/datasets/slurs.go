package datasets

type SlurEntry struct {
	ID              int
	Slur            string
	SlurTarget      string
	SlurDescription string
}

type NationalityEntry struct {
	ID          int
	Nationality string
}

type JobTitleEntry struct {
	ID       int
	JobTitle string
}

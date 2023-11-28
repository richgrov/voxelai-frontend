package jobs

import (
	"database/sql"
	"time"
)

type SqlJobsService struct {
	db *sql.DB
}

func NewSqlServce(dbUrl string) (*SqlJobsService, error) {
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS jobs (id text NOT NULL PRIMARY KEY, prompt text NOT NULL, status text, result text)")
	if err != nil {
		return nil, err
	}

	return &SqlJobsService{
		db: db,
	}, nil
}

func (service *SqlJobsService) GetJob(id string) (*Job, error) {
	var prompt string
	var status sql.NullString
	var result sql.NullString
	err := service.db.QueryRow("SELECT prompt, status, result FROM jobs WHERE id=? LIMIT 1", id).Scan(&prompt, &status, &result)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var statusEnum Status
	switch status.String {
	case "COMPLETED":
		statusEnum = StatusSucceess
	case "FAILED":
		statusEnum = StatusFailed
	default:
		statusEnum = StatusInProgress
	}

	return &Job{
		Prompt: prompt,
		Result: result.String,
		Status: statusEnum,
	}, nil
}

func (service *SqlJobsService) StartJob(id string, prompt string) error {
	_, err := service.db.Exec("INSERT INTO jobs (id, prompt) VALUES (?, ?)", id, prompt)
	return err
}

func (service *SqlJobsService) UpdateJobStatus(id string, status Status, result string) error {
	_, err := service.db.Exec("UPDATE jobs SET status=?, result=? WHERE id=?", status.String(), result, id)
	return err
}

func (service *SqlJobsService) WaitForCompletion(id string, timeout time.Duration) (*Job, error) {
	endTime := time.Now().Add(timeout)
	for time.Now().Compare(endTime) < 0 {
		var prompt string
		var status sql.NullString
		var result sql.NullString
		err := service.db.QueryRow("SELECT prompt, status, result FROM jobs WHERE id=? LIMIT 1", id).Scan(&prompt, &status, &result)

		if err != nil {
			return nil, err
		}

		statusEnum := statusFromStr(status.String)
		if statusEnum != StatusInProgress {
			return &Job{
				Prompt: prompt,
				Result: result.String,
				Status: statusEnum,
			}, nil
		}

		time.Sleep(2 * time.Second)
	}

	return nil, ErrTimeout
}

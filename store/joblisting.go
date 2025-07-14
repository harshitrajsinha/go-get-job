package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/harshitrajsinha/go-get-job/graph/model"
)

type JobStore struct {
	db *sql.DB
}

func NewJobStore(db *sql.DB) *JobStore {
	return &JobStore{db: db}
}

func capitalize(word *string) string {
	firstChar := strings.ToUpper(string((*word)[0]))
	capitalizedWord := strings.Join([]string{firstChar, (*word)[1:]}, "")
	return (capitalizedWord)
}

func (j JobStore) GetJobs(limit int32, offset int32) (*model.JobData, error) {

	lists := make([]*model.JobListing, 0)
	total_records := model.TotalRecords{
		TRec: 0,
	}

	jobData := model.JobData{
		Rows:         lists,
		TotalRecords: &total_records,
	}

	if limit <= 0 {
		limit = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	rows, err := j.db.QueryContext(ctx, "SELECT id, title, company, url, description, job_id, experience, job_type, city, country,  count(*) over() as total_records FROM joblisting ORDER BY posted_on LIMIT $1 OFFSET $2", limit, offset)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("error getting data, no rows found")
			return nil, errors.New("no records found")
		}
		log.Println("error getting data", err)
		return nil, errors.New("error getting data")
	}
	defer rows.Close()

	// Get each row data into a slice
	for rows.Next() {

		queryData := model.JobListing{
			Location: &model.Location{},
		}

		err := rows.Scan(
			&queryData.ID, &queryData.Title, &queryData.Company, &queryData.URL, &queryData.Description, &queryData.JobID, &queryData.Experience, &queryData.JobType, &queryData.Location.City, &queryData.Location.Country, &total_records.TRec)
		if err != nil {
			log.Println("error getting data", err)
			return nil, errors.New("error getting data")
		}
		queryData.Company = capitalize(&queryData.Company)
		lists = append(lists, &queryData)
	}
	if len(lists) == 0 {
		log.Println("error getting data, no rows")
		return nil, errors.New("no records found")
	}

	jobData.Rows = lists
	jobData.TotalRecords = &total_records
	return &jobData, nil
}

func (j JobStore) GetJobByTitle(title string, limit int32, offset int32) (*model.JobData, error) {

	lists := make([]*model.JobListing, 0)
	total_records := model.TotalRecords{
		TRec: 0,
	}

	jobData := model.JobData{
		Rows:         lists,
		TotalRecords: &total_records,
	}

	if limit <= 0 {
		limit = 10
	}

	if title == "" {
		return nil, errors.New("job title is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	rows, err := j.db.QueryContext(ctx, "SELECT id, title, company, url, description, job_id, experience, job_type, city, country, count(*) over() as total_records FROM joblisting WHERE title ~* $1 ORDER BY posted_on LIMIT $2 OFFSET $3", title, limit, offset)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("error getting data, no match")
			return nil, fmt.Errorf("no records found for Job Title: %s", title)
		}
		log.Println("error getting data", err)
		return nil, errors.New("error getting data")
	}
	defer rows.Close()

	// Get each row data into a slice
	for rows.Next() {

		queryData := model.JobListing{
			Location: &model.Location{},
		}

		err := rows.Scan(
			&queryData.ID, &queryData.Title, &queryData.Company, &queryData.URL, &queryData.Description, &queryData.JobID, &queryData.Experience, &queryData.JobType, &queryData.Location.City, &queryData.Location.Country, &total_records.TRec)
		if err != nil {
			log.Println("error getting data", err)
			return nil, errors.New("error getting data")
		}
		queryData.Company = capitalize(&queryData.Company)
		lists = append(lists, &queryData)
	}
	if len(lists) == 0 {
		log.Println("error getting data, no match")
		return nil, fmt.Errorf("no records found for Job Title: %s", title)
	}

	jobData.Rows = lists
	jobData.TotalRecords = &total_records
	return &jobData, nil
}

func (j JobStore) GetJobByCompany(company string, limit int32, offset int32) (*model.JobData, error) {

	lists := make([]*model.JobListing, 0)
	total_records := model.TotalRecords{
		TRec: 0,
	}

	jobData := model.JobData{
		Rows:         lists,
		TotalRecords: &total_records,
	}

	if limit <= 0 {
		limit = 10
	}

	if company == "" {
		return nil, errors.New("company is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	rows, err := j.db.QueryContext(ctx, "SELECT id, title, company, url, description, job_id, experience, job_type, city, country, count(*) over() as total_records FROM joblisting WHERE company ~* $1 GROUP BY id, title, company, url, description, job_id, experience, job_type, city, country ORDER BY posted_on LIMIT $2 OFFSET $3", company, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("error getting data, no match")
			return nil, fmt.Errorf("no records found for company: %s", company)
		}
		log.Println("error getting data", err)
		return nil, errors.New("error getting data")
	}
	defer rows.Close()

	// slice to store all rows

	// Get each row data into a slice
	for rows.Next() {

		queryData := model.JobListing{
			Location: &model.Location{},
		}

		err := rows.Scan(
			&queryData.ID, &queryData.Title, &queryData.Company, &queryData.URL, &queryData.Description, &queryData.JobID, &queryData.Experience, &queryData.JobType, &queryData.Location.City, &queryData.Location.Country, &total_records.TRec)
		if err != nil {
			log.Println("error getting data", err)
			return nil, errors.New("error getting data")
		}
		queryData.Company = capitalize(&queryData.Company)
		lists = append(lists, &queryData)
	}

	if len(lists) == 0 {
		log.Println("error getting data, no match")
		return nil, fmt.Errorf("no records found for company: %s", company)
	}

	jobData.Rows = lists
	jobData.TotalRecords = &total_records
	return &jobData, nil
}

func (j JobStore) GetJobByID(job_id int32) (*model.JobListing, error) {

	queryData := model.JobListing{
		Location: &model.Location{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	err := j.db.QueryRowContext(ctx, "SELECT id, title, company, url, description, job_id, experience, job_type, city, country FROM joblisting WHERE job_id=$1", job_id).Scan(
		&queryData.ID, &queryData.Title, &queryData.Company, &queryData.URL, &queryData.Description, &queryData.JobID, &queryData.Experience, &queryData.JobType, &queryData.Location.City, &queryData.Location.Country)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("error getting data, no match")
			return nil, fmt.Errorf("no records found for Job ID: %d", job_id)
		}
		log.Println("error getting data", err)
		return nil, errors.New("error getting data")
	}
	queryData.Company = capitalize(&queryData.Company)

	return &queryData, nil
}

package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/harshitrajsinha/go-get-job/graph/model"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Jobs is the resolver for the jobs field.
func (r *queryResolver) Jobs(ctx context.Context, limit int32, offset int32) (*model.JobData, error) {

	// cache for large data set
	if limit == 0 && offset == 0 {
		key := fmt.Sprintf("jobs:limit:%d:offset:%d", limit, offset)
		cached, err := r.rdb.Get(ctx, key).Result()
		if err == nil {
			log.Printf("Cache hit for jobs:limit:%d:offset:%d", limit, offset)
			// Cache hit - deserialize JSON
			var data *model.JobData
			err := json.Unmarshal([]byte(cached), &data)
			if err != nil {
				return nil, gqlerror.Errorf("%s", err)
			}
			return data, nil
		}

		// Cache miss - fetch from DB
		log.Printf("Cache miss for jobs:limit:%d:offset:%d", limit, offset)
		data, err := r.res.GetJobs(ctx, limit, offset)
		if err != nil {
			return nil, gqlerror.Errorf("%s", err)
		}

		// Store in Redis for 10 minutes
		jsonData, _ := json.Marshal(data)
		log.Printf("Cache store for jobs:limit:%d:offset:%d", limit, offset)
		r.rdb.Set(ctx, key, jsonData, 10*time.Minute)

		return data, nil
	}

	data, err := r.res.GetJobs(ctx, limit, offset)
	if err != nil {
		return nil, gqlerror.Errorf("%s", err)
	}
	return data, nil

}

// JobByTitle is the resolver for the jobByTitle field.
func (r *queryResolver) JobByTitle(ctx context.Context, title string, limit int32, offset int32) (*model.JobData, error) {

	// cache for large data set
	if limit == 0 && offset == 0 {
		key := fmt.Sprintf("jobs:title:%s:limit:%d:offset:%d", title, limit, offset)
		cached, err := r.rdb.Get(ctx, key).Result()
		if err == nil {
			log.Printf("Cache hit for jobs:title:%s:limit:%d:offset:%d", title, limit, offset)
			// Cache hit - deserialize JSON
			var data *model.JobData
			err := json.Unmarshal([]byte(cached), &data)
			if err != nil {
				return nil, gqlerror.Errorf("%s", err)
			}
			return data, nil
		}

		// Cache miss - fetch from DB
		log.Printf("Cache miss for jobs:title:%s:limit:%d:offset:%d", title, limit, offset)
		data, err := r.res.GetJobByTitle(ctx, title, limit, offset)
		if err != nil {
			return nil, gqlerror.Errorf("%s", err)
		}

		// Store in Redis for 10 minutes
		jsonData, _ := json.Marshal(data)
		log.Printf("Cache store for jobs:title:%s:limit:%d:offset:%d", title, limit, offset)
		r.rdb.Set(ctx, key, jsonData, 10*time.Minute)

		return data, nil
	}

	data, err := r.res.GetJobByTitle(ctx, title, limit, offset)
	if err != nil {
		return nil, gqlerror.Errorf("%s", err)
	}
	return data, nil
}

// JobByID is the resolver for the jobByID field.
func (r *queryResolver) JobByID(ctx context.Context, jobID int32) (*model.JobListing, error) {
	if len(strconv.Itoa(int(jobID))) != 6 {
		return nil, gqlerror.Errorf("job id must be exactly 6 digits")
	}

	key := fmt.Sprintf("job:%d", jobID)
	cached, err := r.rdb.Get(ctx, key).Result()
	if err == nil {
		log.Printf("Cache hit for job_id: %d", jobID)
		// Cache hit - deserialize JSON
		var data *model.JobListing
		err := json.Unmarshal([]byte(cached), &data)
		if err != nil {
			return nil, gqlerror.Errorf("%s", err)
		}
		return data, nil
	}

	// Cache miss - fetch from DB
	log.Printf("Cache miss for job_id: %d", jobID)
	data, err := r.res.GetJobByID(ctx, jobID)
	if err != nil {
		return nil, gqlerror.Errorf("%s", err)
	}

	// Store in Redis for 10 minutes
	jsonData, _ := json.Marshal(data)
	log.Printf("Cache store for job_id: %d", jobID)
	r.rdb.Set(ctx, key, jsonData, 10*time.Minute)

	return data, nil
}

// JobByCompany is the resolver for the jobByCompany field.
func (r *queryResolver) JobByCompany(ctx context.Context, company string, limit int32, offset int32) (*model.JobData, error) {

	// cache for large data set
	if limit == 0 && offset == 0 {
		key := fmt.Sprintf("jobs:company:%s:limit:%d:offset:%d", company, limit, offset)
		cached, err := r.rdb.Get(ctx, key).Result()
		if err == nil {
			log.Printf("Cache hit for jobs:title:%s:limit:%d:offset:%d", company, limit, offset)
			// Cache hit - deserialize JSON
			var data *model.JobData
			err := json.Unmarshal([]byte(cached), &data)
			if err != nil {
				return nil, gqlerror.Errorf("%s", err)
			}
			return data, nil
		}

		// Cache miss - fetch from DB
		log.Printf("Cache miss for jobs:title:%s:limit:%d:offset:%d", company, limit, offset)
		data, err := r.res.GetJobByTitle(ctx, company, limit, offset)
		if err != nil {
			return nil, gqlerror.Errorf("%s", err)
		}

		// Store in Redis for 10 minutes
		jsonData, _ := json.Marshal(data)
		log.Printf("Cache store for jobs:title:%s:limit:%d:offset:%d", company, limit, offset)
		r.rdb.Set(ctx, key, jsonData, 10*time.Minute)

		return data, nil
	}

	data, err := r.res.GetJobByCompany(ctx, company, limit, offset)
	if err != nil {
		return nil, gqlerror.Errorf("%s", err)
	}
	return data, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

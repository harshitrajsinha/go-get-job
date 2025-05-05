package graph

import (
	"github.com/harshitrajsinha/go-get-job/store"
	"github.com/redis/go-redis/v9"
)

type Resolver struct {
	res *store.JobStore
	rdb *redis.Client
}

func NewGQLQueryResolver(res *store.JobStore, rdb *redis.Client) *Resolver {
	return &Resolver{res: res, rdb: rdb}
}

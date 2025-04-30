package graph

import "github.com/harshitrajsinha/go-get-job/store"

type Resolver struct {
	res *store.JobStore
}

func NewGQLQueryResolver(res *store.JobStore) *Resolver {
	return &Resolver{res: res}
}

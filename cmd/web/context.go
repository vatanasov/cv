package main

import (
	"context"
	"net/http"
)

type contextKey string

const candidateIdContextKey = contextKey("candidateId")

func (app *application) contextSetCandidateId(r *http.Request, candidateId int64) *http.Request {
	ctx := context.WithValue(r.Context(), candidateIdContextKey, candidateId)
	return r.WithContext(ctx)
}

func (app *application) contextGetCandidateId(r *http.Request) int64 {
	cid, ok := r.Context().Value(candidateIdContextKey).(int64)
	if !ok {
		panic("missing candidateId in request context")
	}
	return cid
}

package apis

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/whisdom/server/logs"
)

type logsResponse struct {
	Logs       []logs.Record `json:"logs"`
	NextCursor string        `json:"next_cursor"`
}

// handleListLogs implements GET /api/v1/logs.
func (a *App) handleListLogs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	opts := logs.QueryOptions{
		Keyword: q.Get("keyword"),
		Sort:    q.Get("sort"),
		Order:   q.Get("order"),
	}

	if v := q.Get("status"); v != "" {
		opts.Statuses = strings.Split(v, ",")
	}
	if v := q.Get("start"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			opts.Start = &t
		}
	}
	if v := q.Get("end"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			opts.End = &t
		}
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts.Limit = n
		}
	}
	opts.Cursor = q.Get("cursor")
	if v := q.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts.Page = n
		}
	}
	if v := q.Get("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts.PageSize = n
		}
	}

	result, err := a.Logs.Query(opts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read logs")
		return
	}

	records := result.Logs
	if records == nil {
		records = []logs.Record{}
	}
	writeJSON(w, http.StatusOK, logsResponse{Logs: records, NextCursor: result.NextCursor})
}

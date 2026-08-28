package logs

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// QueryOptions filters and paginates a log read. Page/PageSize take
// precedence over Limit/Cursor when both are supplied, matching the two
// pagination styles the UI.md log table needs (offset pages for a grid,
// cursor for infinite scroll).
type QueryOptions struct {
	Start    *time.Time
	End      *time.Time
	Statuses []string
	Keyword  string
	Sort     string // "timestamp" or "status"
	Order    string // "asc" or "desc"

	Limit  int
	Cursor string

	Page     int
	PageSize int
}

// QueryResult is the response payload for GET /api/v1/logs.
type QueryResult struct {
	Logs       []Record
	NextCursor string
	Total      int
}

// Query loads every retained record, filters and sorts it in memory, and
// returns one page. Retention is capped at RetentionDays, so this stays
// small enough to read fully on every request without an index.
func (s *Service) Query(opts QueryOptions) (QueryResult, error) {
	records, err := s.readAll()
	if err != nil {
		return QueryResult{}, err
	}

	records = filterRecords(records, opts)
	sortRecords(records, opts)

	total := len(records)
	offset, limit := resolvePagination(opts, total)

	end := offset + limit
	if end > total {
		end = total
	}
	if offset > total {
		offset = total
	}
	page := records[offset:end]

	nextCursor := ""
	if end < total {
		nextCursor = strconv.Itoa(end)
	}

	return QueryResult{Logs: page, NextCursor: nextCursor, Total: total}, nil
}

func (s *Service) readAll() ([]Record, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	var records []Record
	for _, entry := range entries {
		if entry.IsDir() || !dailyFileName.MatchString(entry.Name()) {
			continue
		}
		file, err := os.Open(filepath.Join(s.dir, entry.Name()))
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var record Record
			if err := json.Unmarshal(line, &record); err != nil {
				continue
			}
			records = append(records, record)
		}
		file.Close()
	}
	return records, nil
}

func filterRecords(records []Record, opts QueryOptions) []Record {
	statusSet := make(map[string]bool, len(opts.Statuses))
	for _, s := range opts.Statuses {
		statusSet[s] = true
	}
	keyword := strings.ToLower(opts.Keyword)

	filtered := records[:0]
	for _, r := range records {
		if len(statusSet) > 0 && !statusSet[r.Status] {
			continue
		}
		if opts.Start != nil || opts.End != nil {
			ts, err := time.Parse(timestampLayout, r.Timestamp)
			if err == nil {
				if opts.Start != nil && ts.Before(*opts.Start) {
					continue
				}
				if opts.End != nil && ts.After(*opts.End) {
					continue
				}
			}
		}
		if keyword != "" {
			haystack := strings.ToLower(r.Content + " " + r.UserID + " " + r.IP)
			if !strings.Contains(haystack, keyword) {
				continue
			}
		}
		filtered = append(filtered, r)
	}
	return filtered
}

func sortRecords(records []Record, opts QueryOptions) {
	desc := opts.Order != "asc"
	byStatus := opts.Sort == "status"

	key := func(r Record) string {
		if byStatus {
			return r.Status
		}
		return r.Timestamp
	}

	sort.SliceStable(records, func(i, j int) bool {
		if desc {
			return key(records[i]) > key(records[j])
		}
		return key(records[i]) < key(records[j])
	})
}

func resolvePagination(opts QueryOptions, total int) (offset, limit int) {
	if opts.Page > 0 || opts.PageSize > 0 {
		pageSize := opts.PageSize
		if pageSize <= 0 {
			pageSize = 100
		}
		page := opts.Page
		if page <= 0 {
			page = 1
		}
		return (page - 1) * pageSize, pageSize
	}

	limit = opts.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	offset = 0
	if opts.Cursor != "" {
		if v, err := strconv.Atoi(opts.Cursor); err == nil && v > 0 {
			offset = v
		}
	}
	return offset, limit
}

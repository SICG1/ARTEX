package agent

import "github.com/Autumn-27/artex/db"

// eachActivityPage walks a worker's activity rows in id order, one page at a time, until
// the rows run out. Callers keep only what they need, so a long run is never held whole.
func eachActivityPage(fetch func(sinceID int64, limit int) ([]db.Activity, int64, error), visit func([]db.Activity)) error {
	for since := int64(0); ; {
		acts, next, err := fetch(since, 1000)
		if err != nil {
			return err
		}
		visit(acts)
		if len(acts) < 1000 || next <= since {
			return nil
		}
		since = next
	}
}

// lastWorkerOutput returns the final result row of a worker's activity, or the last
// text row when there is none. It reads every page, so a long run is not cut off after
// its first 1000 rows.
func lastWorkerOutput(fetch func(sinceID int64, limit int) ([]db.Activity, int64, error)) (*db.Activity, error) {
	var chosen, fallback *db.Activity
	err := eachActivityPage(fetch, func(acts []db.Activity) {
		for i := range acts {
			switch acts[i].Kind {
			case "result":
				chosen = &acts[i]
				fallback = &acts[i]
			case "text":
				fallback = &acts[i]
			}
		}
	})
	if err != nil {
		return nil, err
	}
	if chosen != nil {
		return chosen, nil
	}
	return fallback, nil
}

package agent

import (
	"errors"
	"testing"

	"github.com/Autumn-27/artex/db"
)

// The final result sits on the second page, past the first 1000 rows: the walk must reach it.
func TestLastWorkerOutputReadsPastFirstPage(t *testing.T) {
	var calls []int64
	fetch := func(since int64, limit int) ([]db.Activity, int64, error) {
		calls = append(calls, since)
		if since == 0 {
			page := make([]db.Activity, 1000)
			for i := range page {
				page[i] = db.Activity{ID: int64(i + 1), Kind: "text"}
			}
			return page, 1000, nil
		}
		return []db.Activity{{ID: 1001, Kind: "text"}, {ID: 1002, Kind: "result"}, {ID: 1003, Kind: "text"}}, 1003, nil
	}
	got, err := lastWorkerOutput(fetch)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != 1002 {
		t.Fatalf("pick = %+v, want result row 1002", got)
	}
	if len(calls) != 2 || calls[1] != 1000 {
		t.Fatalf("fetch calls = %v, want [0 1000]", calls)
	}
}

func TestLastWorkerOutputPropagatesReadError(t *testing.T) {
	_, err := lastWorkerOutput(func(int64, int) ([]db.Activity, int64, error) {
		return nil, 0, errors.New("read failed")
	})
	if err == nil {
		t.Fatal("expected the read error to be returned")
	}
}

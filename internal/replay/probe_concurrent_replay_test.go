package replay

import (
	"sync"
	"testing"
)

func TestProbe_ConcurrentReplayUpdatesStateSafely(t *testing.T) {
	r := New()
	req := Request{Method: "GET", Path: "/concurrent"}
	if err := r.Record(req, Response{Status: 200}); err != nil {
		t.Fatal(err)
	}

	const workers = 100
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, err := r.Replay(req)
			if err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
	if got := r.Stats().Hits; got != workers {
		t.Fatalf("Hits = %d, want %d", got, workers)
	}
}

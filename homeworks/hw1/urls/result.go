package urls

import "sync"

type Result struct {
	mu      sync.RWMutex
	results map[int]int
	done    bool
}

func newResult() *Result {
	return &Result{
		results: make(map[int]int),
	}
}

func (r *Result) GetResponsesCount(code int) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	if _, ok := r.results[code]; ok {
		count = r.results[code]
	}

	return count
}

func (r *Result) GetResult() map[int]int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make(map[int]int)
	for i := range r.results {
		res[i] = r.results[i]
	}

	return res
}

func (r *Result) Done() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.done
}

func (r *Result) addResponse(code int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.results[code]++
}

func (r *Result) setDone() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.done = true
}

package x_test

import (
	"context"
	"errors"
	"iter"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"go.spris.dev/tassert"
	"go.spris.dev/x"
)

func TestGoCollectsResultsAndErrorsByInputIndex(t *testing.T) {
	wantErr := errors.New("failed")
	var calls atomic.Int32

	results, errs := x.Iter([]int{1, 2, 3}).Go(func(v int) (int, error) {
		calls.Add(1)
		if v == 2 {
			return 0, wantErr
		}
		return v * 10, nil
	})

	if calls.Load() != 3 {
		t.Fatalf("f called %d times, want 3", calls.Load())
	}
	tassert.EmptyDiff(t, []int{10, 0, 30}, results)
	tassert.EmptyDiff(t, []error{nil, wantErr, nil}, errs, cmpopts.EquateErrors())
}

func TestGoReturnsEmptySlicesForEmptyInput(t *testing.T) {
	var called atomic.Bool
	results, errs := x.Iter[int](nil).Go(func(v int) (int, error) {
		called.Store(true)
		return v, nil
	})

	if called.Load() {
		t.Fatal("f called for empty input")
	}
	if len(results) != 0 || len(errs) != 0 {
		t.Fatalf("got results %v and errors %v, want empty slices", results, errs)
	}
}

func TestGoPreservesInputOrderWhenCallsFinishOutOfOrder(t *testing.T) {
	finished := [4]chan struct{}{}
	for i := range finished {
		finished[i] = make(chan struct{})
	}
	var timedOut atomic.Bool

	results, errs := x.Iter([]int{1, 2, 3}).Go(func(v int) (int, error) {
		if v < 3 {
			select {
			case <-finished[v+1]:
			case <-time.After(2 * time.Second):
				timedOut.Store(true)
			}
		}
		close(finished[v])
		return v, nil
	})

	if timedOut.Load() {
		t.Fatal("calls did not run concurrently")
	}
	tassert.EmptyDiff(t, []int{1, 2, 3}, results)
	tassert.EmptyDiff(t, []error{nil, nil, nil}, errs)
}

func TestGoCtxWithLimitRejectsZero(t *testing.T) {
	_, err := x.Iter([]int{1}).GoCtxWithLimit(context.Background(), 0, func(_ context.Context, v int) (int, error) {
		return v, nil
	})
	if err == nil {
		t.Fatal("expected an error for a zero limit")
	}
}

func TestGoCtxWithLimitChecksContextAfterCollect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	seq := iter.Seq[int](func(yield func(int) bool) {
		cancel()
		yield(1)
	})

	results, err := x.IterSeq(seq).GoCtxWithLimit(ctx, 1, func(_ context.Context, v int) (int, error) {
		return v, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want context.Canceled", err)
	}
	if results != nil {
		t.Fatalf("got results %v, want nil", results)
	}
}

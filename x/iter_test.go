package x_test

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"go.spris.dev/tassert"
	"go.spris.dev/x"
)

func TestMap(t *testing.T) {
	t.Run("transforms values", func(t *testing.T) {
		in := []int{1, 2, 3, 4}

		t.Run("Map", func(t *testing.T) {
			want := []int{2, 4, 6, 8}
			got := x.Iter(in).
				Map(func(v int) int { return v * 2 }).
				Collect()

			tassert.EmptyDiff(t, want, got)
		})

		t.Run("MapIndex", func(t *testing.T) {
			want := []int{1, 3, 5, 7}
			got := x.Iter(in).
				MapIndex(func(v, i int) int { return v + i }).
				Collect()

			tassert.EmptyDiff(t, want, got)
		})
	})

	t.Run("changes type", func(t *testing.T) {
		in := []int{1, 2, 3}

		t.Run("Map", func(t *testing.T) {
			want := []string{"a", "b", "c"}
			got := x.Iter(in).
				Map(func(v int) string { return string(rune('a' + v - 1)) }).
				Collect()

			tassert.EmptyDiff(t, want, got)
		})

		t.Run("MapIndex", func(t *testing.T) {
			want := []string{"a", "b", "c"}
			got := x.Iter(in).
				MapIndex(func(v int, _ int) string { return string(rune('a' + v - 1)) }).
				Collect()

			tassert.EmptyDiff(t, want, got)
		})
	})

	t.Run("empty input", func(t *testing.T) {
		want := []int{}

		t.Run("Map", func(t *testing.T) {
			got := x.Iter[int](nil).
				Map(func(v int) int { return v + 1 }).
				Collect()

			tassert.EmptyDiff(t, want, got, cmpopts.EquateEmpty())
		})

		t.Run("MapIndex", func(t *testing.T) {
			got := x.Iter[int](nil).
				MapIndex(func(v int, _ int) int { return v + 1 }).
				Collect()

			tassert.EmptyDiff(t, want, got, cmpopts.EquateEmpty())
		})
	})
}

func FuzzMap(f *testing.F) {
	f.Add(1, 2, 3, 4)
	f.Add(0, -1, 1, 2)
	f.Add(0, 0, 0, 0)

	f.Fuzz(func(t *testing.T, a, b, c, d int) {
		in := []int{a, b, c, d}
		want := make([]int, len(in))
		for i, v := range in {
			want[i] = v * 2
		}

		got := x.Iter(in).
			Map(func(v int) int { return v * 2 }).
			Collect()

		tassert.EmptyDiff(t, want, got)
	})
}

func TestFilter(t *testing.T) {
	t.Run("keeps matching values", func(t *testing.T) {
		in := []int{1, 2, 3, 4}

		t.Run("Filter", func(t *testing.T) {
			want := []int{2, 4}
			got := x.Iter(in).
				Filter(func(v int) bool { return v%2 == 0 }).
				Collect()

			tassert.EmptyDiff(t, want, got)
		})

		t.Run("FilterIndex", func(t *testing.T) {
			want := []int{1, 3}
			got := x.Iter(in).
				FilterIndex(func(_ int, i int) bool { return i%2 == 0 }).
				Collect()

			tassert.EmptyDiff(t, want, got)
		})
	})

	t.Run("no matches", func(t *testing.T) {
		in := []int{1, 3, 5}
		want := []int{}

		t.Run("Filter", func(t *testing.T) {
			got := x.Iter(in).
				Filter(func(v int) bool { return v%2 == 0 }).
				Collect()

			tassert.EmptyDiff(t, want, got, cmpopts.EquateEmpty())
		})

		t.Run("FilterIndex", func(t *testing.T) {
			got := x.Iter(in).
				FilterIndex(func(_ int, i int) bool { return i < 0 }).
				Collect()

			tassert.EmptyDiff(t, want, got, cmpopts.EquateEmpty())
		})
	})
}

func FuzzFilter(f *testing.F) {
	f.Add(1, 2, 3, 4)
	f.Add(0, -1, 1, 2)
	f.Add(1, 3, 5, 7)

	f.Fuzz(func(t *testing.T, a, b, c, d int) {
		in := []int{a, b, c, d}
		var want []int
		for _, v := range in {
			if v%2 == 0 {
				want = append(want, v)
			}
		}

		got := x.Iter(in).
			Filter(func(v int) bool { return v%2 == 0 }).
			Collect()

		tassert.EmptyDiff(t, want, got, cmpopts.EquateEmpty())
	})
}

func TestFold(t *testing.T) {
	t.Run("accumulates values", func(t *testing.T) {
		in := []int{1, 2, 3, 4}

		t.Run("Fold", func(t *testing.T) {
			want := 10
			got := x.Iter(in).
				Fold(0, func(sum, v int) int { return sum + v })

			tassert.Equal(t, want, got)
		})

		t.Run("FoldIndex", func(t *testing.T) {
			want := 20
			got := x.Iter(in).
				FoldIndex(0, func(sum, v, i int) int { return sum + v*i })

			tassert.Equal(t, want, got)
		})
	})

	t.Run("uses initial value for empty input", func(t *testing.T) {
		want := 10

		t.Run("Fold", func(t *testing.T) {
			got := x.Iter[int](nil).
				Fold(10, func(sum, v int) int { return sum + v })

			tassert.Equal(t, want, got)
		})

		t.Run("FoldIndex", func(t *testing.T) {
			got := x.Iter[int](nil).
				FoldIndex(10, func(sum, v, _ int) int { return sum + v })

			tassert.Equal(t, want, got)
		})
	})

	t.Run("preserves iteration order", func(t *testing.T) {
		in := []string{"a", "b", "c"}

		t.Run("Fold", func(t *testing.T) {
			want := "abc"
			got := x.Iter(in).
				Fold("", func(result, v string) string { return result + v })

			tassert.Equal(t, want, got)
		})

		t.Run("FoldIndex", func(t *testing.T) {
			want := "abc"
			got := x.Iter(in).
				FoldIndex("", func(result, v string, _ int) string { return result + v })

			tassert.Equal(t, want, got)
		})
	})
}

func TestSome(t *testing.T) {
	t.Run("returns true when a value matches", func(t *testing.T) {
		got := x.Iter([]int{1, 2, 3}).Some(func(v int) bool { return v == 2 })

		tassert.Equal(t, true, got)
	})

	t.Run("returns false when no value matches", func(t *testing.T) {
		got := x.Iter([]int{1, 2, 3}).Some(func(v int) bool { return v == 4 })

		tassert.Equal(t, false, got)
	})

	t.Run("returns false for empty input", func(t *testing.T) {
		got := x.Iter[int](nil).Some(func(v int) bool { return true })

		tassert.Equal(t, false, got)
	})
}

func TestNone(t *testing.T) {
	t.Run("returns true when no value matches", func(t *testing.T) {
		got := x.Iter([]int{1, 2, 3}).None(func(v int) bool { return v == 4 })

		tassert.Equal(t, true, got)
	})

	t.Run("returns false when a value matches", func(t *testing.T) {
		got := x.Iter([]int{1, 2, 3}).None(func(v int) bool { return v == 2 })

		tassert.Equal(t, false, got)
	})

	t.Run("returns true for empty input", func(t *testing.T) {
		got := x.Iter[int](nil).None(func(v int) bool { return true })

		tassert.Equal(t, true, got)
	})
}

func TestEvery(t *testing.T) {
	t.Run("returns true when every value matches", func(t *testing.T) {
		got := x.Iter([]int{2, 4, 6}).Every(func(v int) bool { return v%2 == 0 })

		tassert.Equal(t, true, got)
	})

	t.Run("returns false when a value does not match", func(t *testing.T) {
		got := x.Iter([]int{2, 3, 6}).Every(func(v int) bool { return v%2 == 0 })

		tassert.Equal(t, false, got)
	})

	t.Run("returns true for empty input", func(t *testing.T) {
		got := x.Iter[int](nil).Every(func(v int) bool { return false })

		tassert.Equal(t, true, got)
	})
}

package loadavg

import (
	"math/rand"
	"testing"
	"time"
)

func calcAvg(items []Value, i int, num int) []int {
	s := 0
	n := 0

	avg := make([]int, num)
	for j := range num {
		if i < 0 {
			i = num - 1
		}

		if i >= num {
			i = i % num
		}

		s += items[len(items)-j-1].One
		n++

		avg[i] = s / n

		i--
	}

	return avg
}

func compareAvg(t *testing.T, c *LoadAvg, avg []int, i int, num int) {
	t.Helper()

	for range num {
		if i < 0 {
			i = num - 1
		}

		if i >= num {
			i = i % num
		}

		a, b := c.buf[i].One, avg[i]*c.prec

		if a != b && (max(a, b)-min(a, b)) > c.prec {
			t.Fatalf(
				"value [%d] is not equal: wants %d got %d (diff %d)",
				i,
				avg[i]*c.prec,
				c.buf[i].One,
				(max(a, b) - min(a, b)),
			)
		}

		i--
	}
}

func TestAdd(t *testing.T) {
	const l = 60
	c := NewAggregator(l)
	items := []Value{}

	tm := time.Now()
	for i := range 100_000 {
		n := rand.Intn(1000)

		val := Value{
			One: n,
			T:   tm.Add(time.Duration(i) * time.Second),
		}

		items = append(items, val)
		avg := calcAvg(items, i, min(i+1, l+1))

		c.Add(val)

		compareAvg(t, c, avg, i, min(i+1, l+1))
	}
}

func TestGrow(t *testing.T) {
	maxL := 10
	c := NewAggregator(1)
	items := []Value{}

	tm := time.Now()
	for n := range maxL {
		val := Value{
			One: n * 10,
			T:   tm.Add(time.Duration(n) * time.Second),
		}

		items = append(items, val)
		avg := calcAvg(items, n, n+1)

		c.Add(val)

		compareAvg(t, c, avg, n, n+1)
		c.Grow(n + 2)
	}
}

package ketama

import (
	"fmt"
	"strconv"
	"testing"
)

func TestBasicCompat(t *testing.T) {

	var compatTest = []Bucket{
		{"server1", 8699},
		{"server10", 9462},
		{"server2", 10885},
		{"server3", 9980},
		{"server4", 10237},
		{"server5", 9099},
		{"server6", 10997},
		{"server7", 10365},
		{"server8", 10380},
		{"server9", 9896},
	}

	var buckets []Bucket

	for i := 1; i <= 10; i++ {
		b := &Bucket{Label: fmt.Sprintf("server%d", i), Weight: 1}
		buckets = append(buckets, *b)
	}

	k, _ := New(buckets)

	m := make(map[string]int)

	for i := 0; i < 100000; i++ {
		s := k.Hash("foo" + strconv.Itoa(i))
		m[s]++
	}

	for _, tt := range compatTest {
		if m[tt.Label] != tt.Weight {
			t.Errorf("basic compatibility check failed key=%s expected=%d got=%d", tt.Label, tt.Weight, m[tt.Label])
		}
	}
}

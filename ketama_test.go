package ketama

import (
	"encoding/hex"
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

func TestSegfault(t *testing.T) {

	// perl -Mblib -MAlgorithm::ConsistentHash::Ketama -wE 'my $ketama = Algorithm::ConsistentHash::Ketama->new(); $ketama->add_bucket( "r01", 100 ); $ketama->add_bucket( "r02", 100 ); my $key = $ketama->hash( pack "H*", "37292b669dd8f7c952cf79ca0dc6c5d7" ); say $key'

	buckets := []Bucket{Bucket{Label: "r01", Weight: 100}, Bucket{Label: "r02", Weight: 100}}
	k, _ := New(buckets)

	tests := []struct {
		key string
		b   string
	}{
		{"161c6d14dae73a874ac0aa0017fb8340", "r01"},
		{"37292b669dd8f7c952cf79ca0dc6c5d7", "r01"},
	}

	for _, tt := range tests {
		key, _ := hex.DecodeString(tt.key)
		b := k.Hash(string(key))
		if b != tt.b {
			t.Errorf("k.Hash(%v)=%v, want %v", tt.key, b, tt.b)
		}
	}

}

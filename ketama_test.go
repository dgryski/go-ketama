package ketama

import (
	"fmt"
	"strconv"
	"testing"
)

/*
use Algorithm::ConsistentHash::Ketama;

my $ketama = Algorithm::ConsistentHash::Ketama->new();

for my $i (0..9) {
  $ketama->add_bucket( "server$i", 1);
}

my %m;

for (my $i=0; $i < 100000; $i++) {
  $m{$ketama->hash("foo$i")}++
}

for my $key (sort keys %m) {
  print "$key $m{$key}\n";
}
*/

func TestKetama(t *testing.T) {

	var kvs = []Bucket{
		{"server0", 9478},
		{"server1", 9029},
		{"server2", 10863},
		{"server3", 10011},
		{"server4", 9734},
		{"server5", 9444},
		{"server6", 10756},
		{"server7", 10021},
		{"server8", 10762},
		{"server9", 9902},
	}

	var buckets []Bucket

	for i := 0; i < 10; i++ {
		b := &Bucket{Label: fmt.Sprintf("server%d", i), Weight: 1}
		buckets = append(buckets, *b)
	}

	k, _ := New(buckets)

	m := make(map[string]int)

	for i := 0; i < 100000; i++ {
		s := k.Hash("foo" + strconv.Itoa(i))
		m[s]++
	}

	for _, tt := range kvs {
		if m[tt.Label] != tt.Weight {
			t.Errorf("compatibility check failed for key=%s expected=%d got=%d", tt.Label, tt.Weight, m[tt.Label])
		}
	}
}

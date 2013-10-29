// Package ketama implements consistent hashing compatible with Algorithm::ConsistentHash::Ketama
/*
This implementation draws from the Daisuke Maki's Perl module, which itself is
based on the original libketama code.  That code was licensed under the GPLv2,
and thus so it this.

The major API change from libketama is that Algorithm::ConsistentHash::Ketama allows hashing
arbitrary strings, instead of just memcached server IP addresses.
*/
package ketama

import (
	"crypto/md5"
	"fmt"
	"math"
	"sort"
)

type Bucket struct {
	Label  string
	Weight int
}

type continuumPoint struct {
	bucket Bucket
	point  uint
}

type Continuum []continuumPoint

func (c Continuum) Less(i, j int) bool { return c[i].point < c[j].point }
func (c Continuum) Len() int           { return len(c) }
func (c Continuum) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func md5Digest(in string) []byte {
	h := md5.New()
	h.Write([]byte(in))
	return h.Sum(nil)
}

func hashString(in string) uint {
	digest := md5Digest(in)
	return uint(digest[3])<<24 | uint(digest[2])<<16 | uint(digest[1])<<8 | uint(digest[0])
}

func New(buckets []Bucket) (Continuum, error) {

	numbuckets := len(buckets)

	if numbuckets == 0 {
		// let them error when they try to use it
		return Continuum(nil), nil
	}

	ket := make([]continuumPoint, 0, numbuckets*160)

	totalweight := float64(0)
	for _, b := range buckets {
		totalweight += float64(b.Weight)
	}

	for i, b := range buckets {
		pct := float64(b.Weight) / totalweight

		limit := int(math.Floor(pct * 40.0 * float64(numbuckets)))

		for k := 0; k < limit; k++ {
			/* 40 hashes, 4 numbers per hash = 160 points per bucket */
			ss := fmt.Sprintf("%s-%d", b.Label, k)
			digest := md5Digest(ss)

			for h := 0; h < 4; h++ {
				point := continuumPoint{
					point:  uint(digest[3+h*4])<<24 | uint(digest[2+h*4])<<16 | uint(digest[1+h*4])<<8 | uint(digest[h*4]),
					bucket: buckets[i],
				}
				ket = append(ket, point)
			}
		}
	}

	cont := Continuum(ket)

	sort.Sort(cont)

	return cont, nil
}

func (cont Continuum) Hash(thing string) string {

	if len(cont) == 0 {
		return ""
	}

	h := hashString(thing)
	i := sort.Search(len(cont), func(i int) bool { return cont[i].point >= h })
	if i >= len(cont) {
		i = 0
	}
	return cont[i].bucket.Label
}

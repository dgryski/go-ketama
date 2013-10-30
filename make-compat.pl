#!/usr/bin/perl

use strict;
use warnings;

use Algorithm::ConsistentHash::Ketama;

my $ketama = Algorithm::ConsistentHash::Ketama->new();

print <<EOGO;
package ketama

import (
        "fmt"
        "strconv"
        "testing"
)

var compatTests = [][]Bucket{
EOGO

for my $bucket (1..500) {

    # print STDERR "bucket=$bucket\n";

    $ketama->add_bucket( "server$bucket", 1 );

    my %m;

    for (my $i=0; $i < 100000; $i++) {
        $m{$ketama->hash("foo$i")}++
    }

    print "\t{\n";
    for my $key (sort keys %m) {
        print "\t\t{\"$key\", $m{$key}},\n";
    }
    print "\t},\n";
}
print "}\n";

print <<EOGO;

func TestKetama(t *testing.T) {

        var buckets []Bucket

BUCKET:
        for bucket := 1; bucket <= len(compatTests); bucket++ {

                b := &Bucket{Label: fmt.Sprintf("server%d", bucket), Weight: 1}
                buckets = append(buckets, *b)

                k, _ := New(buckets)

                m := make(map[string]int)

                for i := 0; i < 100000; i++ {
                        s := k.Hash("foo" + strconv.Itoa(i))
                        m[s]++
                }

                for _, tt := range compatTests[bucket-1] {
                        if m[tt.Label] != tt.Weight {
                                t.Errorf("compatibility check failed for buckets=%d key=%s expected=%d got=%d", bucket, tt.Label, tt.Weight, m[tt.Label])
                                continue BUCKET
                        }
                }
        }
}

EOGO

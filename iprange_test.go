// Copyright 2017 Weborama. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cidr_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/Weborama/cidr"
	"github.com/Weborama/uint128"
	"github.com/davecgh/go-spew/spew"
)

func TestRangeNotIPv4(t *testing.T) {
	startIP := net.ParseIP("fc00::1")
	endIP := net.ParseIP("fc00::2")
	result := cidr.IPv4Range2CIDR(startIP, endIP)
	if result != nil {
		t.Fatal("Providing IPv6 addresses should return nil")
	}
}

func TestRangeWrongOrder(t *testing.T) {
	startIP := net.ParseIP("192.168.0.1")
	endIP := net.ParseIP("192.168.0.2")
	result := cidr.IPv4Range2CIDR(endIP, startIP)
	if result != nil {
		t.Fatal("Providing addresses in the wrong order should return nil")
	}

}

func TestIPv4Range2CIDR(t *testing.T) {
	testCases := []struct {
		startIP     net.IP
		endIP       net.IP
		numExpected int
		expected    []string
	}{
		{
			// Simple case
			startIP:     net.ParseIP("192.168.0.0"),
			endIP:       net.ParseIP("192.168.0.1"),
			numExpected: 1,
			expected: []string{
				"192.168.0.0/31",
			},
		},
		{
			// Simple case (for code coverage)
			startIP:     net.ParseIP("0.0.0.0"),
			endIP:       net.ParseIP("0.0.0.1"),
			numExpected: 1,
			expected: []string{
				"0.0.0.0/31",
			},
		},
		{
			startIP:     net.ParseIP("192.168.0.5"),
			endIP:       net.ParseIP("192.168.0.254"),
			numExpected: 13,
		},
		{
			startIP:     net.ParseIP("192.168.0.6"),
			endIP:       net.ParseIP("192.168.0.254"),
			numExpected: 12,
		},
		{
			startIP:     net.ParseIP("192.168.0.7"),
			endIP:       net.ParseIP("192.168.0.254"),
			numExpected: 12,
		},
		{
			startIP:     net.ParseIP("192.168.0.8"),
			endIP:       net.ParseIP("192.168.0.253"),
			numExpected: 10,
		},
		{
			startIP:     net.ParseIP("192.168.0.126"),
			endIP:       net.ParseIP("192.168.0.132"),
			numExpected: 3,
		},
		{
			startIP:     net.ParseIP("192.168.0.148"),
			endIP:       net.ParseIP("192.168.0.157"),
			numExpected: 3,
		},
		{
			// Worst case
			startIP:     net.ParseIP("0.0.0.1"),
			endIP:       net.ParseIP("255.255.255.254"),
			numExpected: 62,
		},
	}

	var cidrs []net.IPNet

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Range_%s-%s-%d", testCase.startIP, testCase.endIP, testCase.numExpected), func(t *testing.T) {
			cidrs = cidr.IPv4Range2CIDR(testCase.startIP, testCase.endIP)
			if testCase.numExpected != len(cidrs) {
				t.Fatalf("CIDRs expected %d, got %d\n%s", testCase.numExpected, len(cidrs), spew.Sdump(cidrs))
			} else if len(testCase.expected) != 0 && len(testCase.expected) == len(cidrs) {
				for i, expected := range testCase.expected {
					if expected != cidrs[i].String() {
						t.Fatalf("CIDR number %d, expected %s, got %s", i, expected, cidrs[i].String())
					}
				}
			}
		})
	}
}

func TestIPv6Range2CIDR(t *testing.T) {
	testCases := []struct {
		startIP     net.IP
		endIP       net.IP
		numExpected int
		expected    []string
	}{
		{
			// Simple case
			startIP:     net.ParseIP("192.168.0.0"),
			endIP:       net.ParseIP("192.168.0.1"),
			numExpected: 1,
			expected: []string{
				"192.168.0.0/31",
			},
		},
		{
			// Simple case (for code coverage)
			startIP:     net.ParseIP("0.0.0.0"),
			endIP:       net.ParseIP("0.0.0.1"),
			numExpected: 1,
			expected: []string{
				"0.0.0.0/31",
			},
		},
		{
			// Worst IPv4 case
			startIP:     net.ParseIP("0.0.0.1"),
			endIP:       net.ParseIP("255.255.255.254"),
			numExpected: 62,
		},
		{
			// Worst IPv6 case
			startIP:     net.ParseIP("::1"),
			endIP:       net.ParseIP("FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFE"),
			numExpected: 254,
		},
	}

	var cidrs []net.IPNet

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Range_%s-%s-%d", testCase.startIP, testCase.endIP, testCase.numExpected), func(t *testing.T) {
			cidrs = cidr.IPv6Range2CIDR(testCase.startIP, testCase.endIP)
			if testCase.numExpected != len(cidrs) {
				t.Fatalf("CIDRs expected %d, got %d\n%s", testCase.numExpected, len(cidrs), spew.Sdump(cidrs))
			} else if len(testCase.expected) != 0 && len(testCase.expected) == len(cidrs) {
				for i, expected := range testCase.expected {
					if expected != cidrs[i].String() {
						t.Fatalf("CIDR number %d, expected %s, got %s", i, expected, cidrs[i].String())
					}
				}
			}
		})
	}
}

func BenchmarkIPv4ToUint32(b *testing.B) {
	ip := net.ParseIP("128.128.128.128")
	var i uint32
	for n := 0; n < b.N; n++ {
		i = cidr.IPv4ToUint32(ip)
	}
	_ = i
}

func BenchmarkIPv6ToUint128(b *testing.B) {
	ip := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	var i uint128.Uint128
	for n := 0; n < b.N; n++ {
		i = cidr.IPv6ToUint128(ip)
	}
	_ = i
}

// func BenchmarkIPv6ToUint128Alt(b *testing.B) {
// 	ip := net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
// 	var i uint128.Uint128
// 	for n := 0; n < b.N; n++ {
// 		i = IPv6ToUint128Alt(ip)
// 	}
// 	_ = i
// }

func BenchmarkUint32ToIPv4(b *testing.B) {
	var i uint32 = 325551551
	var ip net.IP
	for n := 0; n < b.N; n++ {
		ip = cidr.Uint32ToIPv4(i)
	}
	_ = ip
}

func BenchmarkUint128ToIPv6(b *testing.B) {
	i := uint128.Uint128{
		H: 2306139570357600256,
		L: 151930230829876,
	}
	var ip net.IP
	for n := 0; n < b.N; n++ {
		ip = cidr.Uint128ToIPv6(i)
	}
	_ = ip
}

func BenchmarkIPv4Range2CIDR(b *testing.B) {
	benchmarkCases := []struct {
		startIP net.IP
		endIP   net.IP
	}{
		{
			// Simple case
			net.ParseIP("0.0.0.0"),
			net.ParseIP("0.0.0.1"),
		},
		{
			// Worst case
			net.ParseIP("0.0.0.1"),
			net.ParseIP("255.255.255.254"),
		},
	}
	var cidrs []net.IPNet
	for _, benchmarkCase := range benchmarkCases {
		b.Run(fmt.Sprintf("Range_%s-%s", benchmarkCase.startIP, benchmarkCase.endIP), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				cidrs = cidr.IPv4Range2CIDR(benchmarkCase.startIP, benchmarkCase.endIP)
			}
		})
	}
	_ = cidrs
}

func BenchmarkIPv6Range2CIDR(b *testing.B) {
	benchmarkCases := []struct {
		startIP net.IP
		endIP   net.IP
	}{
		{
			// Simple case
			net.ParseIP("0.0.0.0"),
			net.ParseIP("0.0.0.1"),
		},
		{
			// Worst IPv4 case
			net.ParseIP("0.0.0.1"),
			net.ParseIP("255.255.255.254"),
		},
		{
			// Worst IPv6 case
			net.ParseIP("::1"),
			net.ParseIP("FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFE"),
		},
	}
	var cidrs []net.IPNet
	for _, benchmarkCase := range benchmarkCases {
		b.Run(fmt.Sprintf("Range_%s-%s", benchmarkCase.startIP, benchmarkCase.endIP), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				cidrs = cidr.IPv6Range2CIDR(benchmarkCase.startIP, benchmarkCase.endIP)
			}
		})
	}
	_ = cidrs
}
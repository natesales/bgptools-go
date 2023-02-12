package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tcs := []struct {
		in          string
		as          uint32
		ip          string
		prefix      string
		country     string
		registry    string
		allocated   string
		asName      string
		shouldError bool
	}{
		{`AS      | IP               | BGP Prefix          | CC | Registry | Allocated  | AS Name
968     |                  |                     | US | ARIN     | 2022-06-22 | Packetframe`, 968, "", "", "US", "ARIN", "2022-06-22", "Packetframe", false},
		{`AS      | IP               | BGP Prefix          | CC | Registry | Allocated  | AS Name
13335   | 1.1.1.1          | 1.1.1.0/24          | US | ARIN     | 0001-01-01 | Cloudflare, Inc.`, 13335, "1.1.1.1", "1.1.1.0/24", "US", "ARIN", "0001-01-01", "Cloudflare, Inc.", false},
		{"", 0, "", "", "", "", "", "", true},
	}

	for _, tc := range tcs {
		out, err := Parse(tc.in)
		if tc.shouldError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.as, out.AS)
			assert.Equal(t, tc.ip, out.IP)
			assert.Equal(t, tc.prefix, out.Prefix)
			assert.Equal(t, tc.country, out.Country)
			assert.Equal(t, tc.registry, out.Registry)
			assert.Equal(t, tc.allocated, out.Allocated)
			assert.Equal(t, tc.asName, out.ASName)
		}
	}
}

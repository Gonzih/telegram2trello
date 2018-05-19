package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexp(t *testing.T) {
	samples := []string{
		"https://powerdns.org/hello-dns/",
		"http://doyoutrustthiscomputer.org/watch",
		"Url with title check this out: http://doyoutrustthiscomputer.org/watch",
		"http://doyoutrustthiscomputer.org/watch - this awesome url",
	}

	for _, sample := range samples {
		assert.True(t, urlRegexp.MatchString(sample))
	}
}

func TestUrlExtract(t *testing.T) {
	samples := map[string]string{
		"https://powerdns.org/hello-dns/":                                        "https://powerdns.org/hello-dns/",
		"Url with title check this out: http://doyoutrustthiscomputer.org/watch": "http://doyoutrustthiscomputer.org/watch",
		"http://doyoutrustthiscomputer.org/watch - this awesome url":             "http://doyoutrustthiscomputer.org/watch",
	}

	for sample, match := range samples {
		assert.Equal(t, match, extractUrl(sample))
	}
}

func TestCardName(t *testing.T) {
	samples := []struct {
		title  string
		desc   string
		result string
	}{
		{"short title", "", "short title"},
		{"short title", "with a medium descr", "short title - with a medium descr"},
		{"short title", "with a very long description that barely fits the limit", "short title - with a very long description that barely fits the limit"},
		{"short title", "with a very long description that does not fit the limit for sure but maybe who knows for yadada", "short title - with a very long description that does not fit the limit fo..."},
	}

	for _, sample := range samples {
		assert.Equal(t, sample.result, generateCardName(sample.title, sample.desc))
	}
}

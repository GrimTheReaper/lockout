package network

import (
	"os"
	"testing"
)

func TestPullingRecord(t *testing.T) {
	// Lets pull our record
	ok(t, pullRecord())

	// Make sure the file is there.
	file, err := os.Open(countryMMDB)
	ok(t, err)
	ok(t, file.Close())

	// Cleanup after this test.
	equals(t, nil, os.Remove(countryMMDB))
}

func TestWhitelist(t *testing.T) {
	// TODO Add more test.

	ok(t, pullRecord())

	// Test 1: Google's DNS _should_ be in the US
	whitelisted, err := checkIP("8.8.8.8", []string{"US"})
	ok(t, err)
	equals(t, true, whitelisted)

	// Test 2: IDK What this is, just added an 8. Resolves to ES
	whitelisted, err = checkIP("88.8.8.8", []string{"US"})
	ok(t, err)
	equals(t, false, whitelisted)

	// checkIP calls PullRecord.
	// Cleanup after this test.
	equals(t, nil, os.Remove(countryMMDB))
}

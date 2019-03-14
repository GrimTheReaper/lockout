package network

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	// TODO: Slim the Libs.
	geoip2 "github.com/oschwald/geoip2-golang"
)

var (
	fileLock = &sync.Mutex{}
	lastPull time.Time
	pullTick = 24 * time.Hour

	// NOTE: No idea on the rate limit here.
	// TODO Figure out a better way to retrieve it, given every new container will pull it.
	recordURI = "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.tar.gz"
)

const countryMMDB = "GeoLite2-Country.mmdb"

// CheckIP will check the API against the GeoLite2-Country Database.
// TODO Pull down the GeoLite2-Country.mmdb after so long.
func checkIP(ipAddress string, countries []string) (bool, error) {
	if ipAddress == "" {
		return false, nil
	}

	// Look at the file, see if it needs to be updated.
	fileLock.Lock()
	if time.Now().After(lastPull.Add(pullTick)) {
		fileLock.Unlock()
		if err := pullRecord(); err != nil {
			return false, err
		}
		fileLock.Lock()
	}
	defer fileLock.Unlock()

	db, err := geoip2.Open(countryMMDB)
	if err != nil {
		return false, err
	}
	defer db.Close()

	ip := net.ParseIP(ipAddress)
	record, err := db.Country(ip)
	if err != nil {
		return false, err
	}

	// TODO: Verify that we're using ISOCode.
	for _, country := range countries {
		// Lets match to the ISOCode. Seems like a good idea.
		if country == record.Country.IsoCode {
			return true, nil
		}
	}

	return false, nil
}

// PullRecord lets us deal with this file without having to pack it in the container.
func pullRecord() error {
	response, err := http.Get(recordURI)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// I need to close this gzip reader, otherwise I'd put it in the `pullFileFromTar` function.
	gzr, err := gzip.NewReader(response.Body)
	if err != nil {
		return err
	}
	defer gzr.Close()

	// Lets deal with this file in memory. It's only a couple of MBs.
	mmDB, err := pullFileFromTar(gzr)
	if err != nil {
		return err
	}

	// Hopefully this will never happen, but if it does at least we still have the old one.
	if mmDB == nil {
		return fmt.Errorf("Failed to find the file in the tar.gz from GeoLite2")
	}

	// Don't lock until we know that we need to.
	fileLock.Lock()
	defer fileLock.Unlock()

	// Lets remove the old one.
	// TODO: Move the file instead of delete it, incase we can't create/write to the new one.
	err = os.Remove(countryMMDB)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	mmDBFile, err := os.Create(countryMMDB)
	if err != nil {
		return err
	}
	defer mmDBFile.Close()

	// Copy over contents
	if _, err := io.Copy(mmDBFile, mmDB); err != nil {
		return err
	}

	lastPull = time.Now()

	return nil
}

// We're gonna handle this function in memory since it's pretty small.
func pullFileFromTar(gzr *gzip.Reader) (io.Reader, error) {
	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil, nil
			// return any other error
		case err != nil:
			return nil, err
		case header == nil:
			continue
		}

		switch header.Typeflag {
		// if it's a file, lets check it.
		case tar.TypeReg:
			if header.FileInfo().Name() == countryMMDB {
				return tr, nil
			}
		}
	}
}

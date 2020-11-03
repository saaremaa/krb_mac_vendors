package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type macAddress struct {
	value string
	valid bool
}

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "info"
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}

	// set global log level
	logrus.SetLevel(ll)

}

func main() {

	if err := checkArgsCount(os.Args); err != nil {
		logrus.Fatalf("%s \n", err)
	}

	mac, err := isValidMac(os.Args[1])
	if err != nil {
		logrus.Fatalf("%s \n", err)
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	if mac.valid {
		vendors := make(map[string]string)

		filename := fmt.Sprintf("%s/macdb/oui.csv", exPath)

		// Open CSV file
		f, err := os.Open(filename)
		if err != nil {
			logrus.Fatalf("%s \n", err)
		}
		defer f.Close()

		// Read File into memory
		lines, err := csv.NewReader(f).ReadAll()
		if err != nil {
			logrus.Fatalf("%s \n", err)
		}

		// Loop through lines & store into  map
		for _, line := range lines {
			vendors[line[1]] = line[2]
		}
		vendorMAC := extractVendor(strings.ToUpper(mac.value))

		logrus.Infof("%s → %s", mac.value, vendors[vendorMAC])
	}

}

// Checking how many args incomming
func checkArgsCount(args []string) (err error) {
	if len(args) != 2 {
		return fmt.Errorf("strange number of arguments received")
	}
	return
}

// Checking the MAC address for validity
func isValidMac(mac interface{}) (macAddr *macAddress, err error) {
	macAddr = new(macAddress)
	macAddr.valid = false
	macAddr.value = mac.(string)

	// check that the MAC address is in hexadecimal
	matched, _ := regexp.MatchString("^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$", macAddr.value)
	if !matched {
		return macAddr, fmt.Errorf("Wrong MAC-address format: %s", macAddr.value)
	}

	macAddr.valid = true

	return macAddr, nil
}

// remove delimiters and get first 6 bytes of mac
func extractVendor(userMAC string) (vendorMAC string) {
	r := strings.NewReplacer(":", "", "-", "")
	vendorMAC = r.Replace(userMAC)
	return vendorMAC[:6]
}

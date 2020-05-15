package main

import (
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRunMain(t *testing.T) {

	t.Run("mac 00:04:56:ee:aa:07", mainTestFunc("mac 00:04:56:ee:aa:07", false))
	t.Run("mac AA:AA:AA:AA:AA:AA aa", mainTestFunc("mac AA:AA:AA:AA:AA:AA aa", true))
	t.Run("mac GG:04:56:ee:aa:07", mainTestFunc("mac GG:04:56:ee:aa:07", true))
}

func mainTestFunc(cli string, fatal bool) func(*testing.T) {
	return func(t *testing.T) {
		var isfatal bool
		logrus.StandardLogger().ExitFunc = func(int) { isfatal = true }
		os.Args = strings.Split(cli, " ")
		main()
		assert.Equal(t, fatal, isfatal)
	}
}

func Test_checkArgsCount(t *testing.T) {
	t.Run("Normal", checkArgsCountTestFunc("mac FF:FF:FF:FF:FF:FF", ""))
	t.Run("With wrong argv", checkArgsCountTestFunc("mac FF:FF:FF:FF:FF:FF -AAAAAAAAAAAAAAAAA", "strange number of arguments received"))
}

func checkArgsCountTestFunc(cli string, errResult string) func(*testing.T) {
	return func(t *testing.T) {
		argv := strings.Split(cli, " ")
		err := checkArgsCount(argv)
		if err != nil {
			assert.EqualError(t, err, errResult)
		}
	}
}

func Test_ValidMac(t *testing.T) {
	t.Run("00:04:56:ee:aa:07", isValidMacTestFunc("00:04:56:ee:aa:07", true, ""))
	t.Run("00-04-56-ee-aa-07", isValidMacTestFunc("00-04-56-ee-aa-07", true, ""))
	t.Run("0004:506:ee:aa:07", isValidMacTestFunc("0004:506:ee:aa:07", false, "Wrong MAC-address format: 0004:506:ee:aa:07"))
	t.Run("00+04-56:ee:aa:07", isValidMacTestFunc("00+04-56:ee:aa:07", false, "Wrong MAC-address format: 00+04-56:ee:aa:07"))
	t.Run("00:04:56:ee:aa:7", isValidMacTestFunc("00:04:56:ee:aa:7", false, "Wrong MAC-address format: 00:04:56:ee:aa:7"))
	t.Run("00:04:56:ee:aa:00HTY7", isValidMacTestFunc("00:04:56:ee:aa:00HTY7", false, "Wrong MAC-address format: 00:04:56:ee:aa:00HTY7"))
	t.Run("00:04:G6:ee:aa:07", isValidMacTestFunc("00:04:G6:ee:aa:07", false, "Wrong MAC-address format: 00:04:G6:ee:aa:07"))
	t.Run("00:04:56:ee:Ja:07", isValidMacTestFunc("00:04:56:ee:Ja:07", false, "Wrong MAC-address format: 00:04:56:ee:Ja:07"))
	t.Run("000456eeaa07", isValidMacTestFunc("000456eeaa07", false, "Wrong MAC-address format: 000456eeaa07"))
}

func isValidMacTestFunc(mac string, expected bool, errResult string) func(*testing.T) {
	return func(t *testing.T) {
		mac, err := isValidMac(mac)
		if err != nil {
			assert.EqualError(t, err, errResult)
		}
		assert.Equal(t, expected, mac.valid)
	}
}

func Test_extractVendor(t *testing.T) {
	result := extractVendor("00:04:56:ee:aa:07")
	assert.Equal(t, "000456", result, "they should be equal")
}

package main

import (
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
	"syscall"
)

var myEnv = regexp.MustCompile(`(?s)\AICINGADB_(\w+?)_(\w+)=(.*)\z`)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	log.Fatal(runDaemon())
}

func runDaemon() error {
	log.Debug("translating env vars to .ini config")

	cfg := ini.Empty()
	for _, env := range os.Environ() {
		if match := myEnv.FindStringSubmatch(env); match != nil {
			_, errNK := cfg.Section(strings.ToLower(match[1])).NewKey(
				strings.ToLower(match[2]), strings.ToLower(match[3]),
			)
			if errNK != nil {
				return errNK
			}
		}
	}

	if errST := cfg.SaveTo("/etc/icingadb/config.ini"); errST != nil {
		return errST
	}

	log.Debug("starting actual daemon via exec(3)")

	return syscall.Exec("/icingadb", []string{"/icingadb", "-config", "/etc/icingadb/config.ini"}, os.Environ())
}

// Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

package main

import (
	"compress/bzip2"
	"database/sql"
	icingaSql "github.com/Icinga/go-libs/sql"
	icingadb "github.com/Icinga/icingadb/config"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"syscall"
)

const config = "/etc/icingadb/config.ini"

var myEnv = regexp.MustCompile(`(?s)\AICINGADB_(\w+?)_(\w+)=(.*)\z`)
var sqlComment = regexp.MustCompile(`(?m)^--.*`)
var sqlStmtSep = regexp.MustCompile(`(?m);$`)

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
			_, errNK := cfg.Section(strings.ToLower(match[1])).NewKey(strings.ToLower(match[2]), match[3])
			if errNK != nil {
				return errNK
			}
		}
	}

	if errST := cfg.SaveTo(config); errST != nil {
		return errST
	}

	if errID := initDb(); errID != nil {
		return errID
	}

	log.Debug("starting actual daemon via exec(3)")

	return syscall.Exec("/icingadb", []string{"/icingadb", "-config", config}, os.Environ())
}

func initDb() error {
	log.Debug("checking SQL database")

	if errPC := icingadb.ParseConfig(config); errPC != nil {
		return errPC
	}

	mi := icingadb.GetMysqlInfo()

	db, errOp := sql.Open("mysql", mi.User+":"+mi.Password+"@tcp("+mi.Host+":"+mi.Port+")/"+mi.Database)
	if errOp != nil {
		return errOp
	}

	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	hasSchema, errHS := dbHasSchema(db, mi.Database)
	if errHS != nil {
		return errHS
	}

	if !hasSchema {
		log.Debug("importing schema into SQL database")

		file, errOp := os.Open("/mysql.schema.sql.bz2")
		if errOp != nil {
			return errOp
		}

		defer file.Close()

		schema, errRA := ioutil.ReadAll(bzip2.NewReader(file))
		if errRA != nil {
			return errRA
		}

		for _, ddl := range sqlStmtSep.Split(string(sqlComment.ReplaceAll(schema, nil)), -1) {
			if ddl = strings.TrimSpace(ddl); ddl != "" {
				if _, errEx := db.Exec(ddl); errEx != nil {
					return errEx
				}
			}
		}
	}

	return nil
}

func dbHasSchema(db *sql.DB, dbName string) (bool, error) {
	type one struct {
		One uint8
	}

	rows, errQr := db.Query(
		"SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=? AND TABLE_NAME='icingadb_schema'", dbName,
	)
	if errQr != nil {
		return false, errQr
	}

	defer rows.Close()

	res, errFR := icingaSql.FetchRowsAsStructSlice(rows, one{}, -1)
	if errFR != nil {
		return false, errFR
	}

	return len(res.([]one)) > 0, nil
}

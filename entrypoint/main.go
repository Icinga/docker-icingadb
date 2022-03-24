// Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

package main

import (
	"compress/bzip2"
	"database/sql"
	"encoding/json"
	icingaSql "github.com/Icinga/go-libs/sql"
	icingadb "github.com/icinga/icingadb/pkg/config"
	"github.com/icinga/icingadb/pkg/logging"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

const config = "/etc/icingadb/config.ini"

var myEnv = regexp.MustCompile(`(?s)\AICINGADB_(\w+?)_(\w+)=(.*)\z`)
var sqlComment = regexp.MustCompile(`(?m)^--.*`)
var sqlStmtSep = regexp.MustCompile(`(?m);$`)

var log = func() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return logger
}()

func main() {
	log.Fatal("failed", zap.Error(runDaemon()))
}

func runDaemon() error {
	log.Debug("translating env vars to YAML config")

	cfg := map[string]map[string]interface{}{}
	for _, env := range os.Environ() {
		if match := myEnv.FindStringSubmatch(env); match != nil {
			section := strings.ToLower(match[1])

			sectionCfg, ok := cfg[section]
			if !ok {
				sectionCfg = map[string]interface{}{}
				cfg[section] = sectionCfg
			}

			rawValue := match[3]
			var value interface{}

			if parsed, err := strconv.ParseInt(rawValue, 10, 64); err == nil {
				value = parsed
			} else {
				value = rawValue
			}

			sectionCfg[strings.ToLower(match[2])] = value
		}
	}

	yml, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(config, yml, 0600); err != nil {
		return err
	}

	if errID := initDb(); errID != nil {
		return errID
	}

	log.Debug("starting actual daemon via exec(3)")

	return syscall.Exec("/icingadb", []string{"/icingadb", "-c", config}, os.Environ())
}

func initDb() error {
	log.Debug("checking SQL database")

	cfg, errFY := icingadb.FromYAMLFile(config)
	if errFY != nil {
		return errFY
	}

	idb, errDB := cfg.Database.Open(&logging.Logger{SugaredLogger: log.Sugar()})
	if errDB != nil {
		return errDB
	}

	defer idb.Close()

	db := idb.DB.DB

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	hasSchema, errHS := dbHasSchema(db, cfg.Database.Database)
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

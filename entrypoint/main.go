// Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

package main

import (
	"compress/bzip2"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	icingaSql "github.com/Icinga/go-libs/sql"
	icingadb "github.com/icinga/icingadb/pkg/config"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
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
	log.Debug("translating env vars to YaML config")

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

	if sleep := os.Getenv("ICINGADB_SLEEP"); sleep != "" {
		seconds, err := strconv.ParseUint(sleep, 10, 64)
		if err != nil {
			return err
		}

		log.Info("sleeping")
		time.Sleep(time.Duration(seconds) * time.Second)
	}

	if errID := initDb(); errID != nil {
		return errID
	}

	log.Debug("starting actual daemon via exec(3)")

	return syscall.Exec("/icingadb", []string{"/icingadb", "-c", config, "--datadir", "/data"}, os.Environ())
}

type schemaUpgrade struct {
	version uint64
	file    string
}

func initDb() error {
	log.Debug("checking SQL database")

	cfg, errFY := icingadb.FromYAMLFile(config)
	if errFY != nil {
		return errFY
	}

	idb, errDB := cfg.Database.Open(log.Sugar())
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

	if hasSchema {
		version, errGV := getDbSchemaVersion(db)
		if errGV != nil {
			return errGV
		}

		files, errGl := filepath.Glob("/mysql-[0-9]*.sql.bz2")
		if errGl != nil {
			return errGl
		}

		var upgrades []schemaUpgrade
		for _, file := range files {
			fileVersion, err := strconv.ParseUint(
				strings.TrimSuffix(strings.TrimPrefix(file, "/mysql-"), ".sql.bz2"), 10, 64,
			)
			if err != nil {
				return err
			}

			if fileVersion > version {
				upgrades = append(upgrades, schemaUpgrade{fileVersion, file})
			}
		}

		if len(upgrades) > 0 {
			log.Info("upgrading SQL database schema")

			sort.Slice(upgrades, func(i, j int) bool {
				return upgrades[i].version < upgrades[j].version
			})

			for _, upgrade := range upgrades {
				log.Info(fmt.Sprintf(".. to v%d", upgrade.version))

				if err := importDdl(upgrade.file, db); err != nil {
					return err
				}
			}
		}
	} else {
		log.Info("importing schema into SQL database")

		if err := importDdl("/mysql-schema.sql.bz2", db); err != nil {
			return err
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

func getDbSchemaVersion(db *sql.DB) (uint64, error) {
	type icingadbSchema struct {
		Version uint64
	}

	rows, errQr := db.Query(
		"SELECT version FROM icingadb_schema ORDER BY `timestamp` DESC LIMIT 1",
	)
	if errQr != nil {
		return 0, errQr
	}

	defer rows.Close()

	res, errFR := icingaSql.FetchRowsAsStructSlice(rows, icingadbSchema{}, -1)
	if errFR != nil {
		return 0, errFR
	}

	versions := res.([]icingadbSchema)
	if len(versions) < 1 {
		return 0, errors.New("unknown SQL database schema version")
	}

	return versions[0].Version, nil
}

func importDdl(file string, db *sql.DB) error {
	f, errOp := os.Open(file)
	if errOp != nil {
		return errOp
	}

	defer f.Close()

	schema, errRA := ioutil.ReadAll(bzip2.NewReader(f))
	if errRA != nil {
		return errRA
	}

	for _, ddl := range sqlStmtSep.Split(string(sqlComment.ReplaceAll(schema, nil)), -1) {
		if ddl = strings.TrimSpace(ddl); ddl != "" {
			if _, err := db.Exec(ddl); err != nil {
				return err
			}
		}
	}

	return nil
}

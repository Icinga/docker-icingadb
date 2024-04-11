// Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

package main

import (
	"compress/bzip2"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"syscall"

	icingaSql "github.com/Icinga/go-libs/sql"
	icingadb "github.com/icinga/icingadb/pkg/config"
	database "github.com/icinga/icingadb/pkg/icingadb"
	driver "github.com/icinga/icingadb/pkg/icingadb"
	"github.com/icinga/icingadb/pkg/logging"
	"go.uber.org/zap"
)

const configDir = "/etc/icingadb"

var config = path.Join(configDir, "icingadb.yml")
var myEnv = regexp.MustCompile(`(?s)\AICINGADB_(\w+?)_(\w+)=(.*)\z`)

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

			option := strings.ToLower(match[2])
			rawValue := match[3]
			var value interface{}

			var unJSONed interface{}
			if json.Unmarshal([]byte(rawValue), &unJSONed) == nil {
				switch unJSONed.(type) {
				case bool, float64:
					value = unJSONed
				}
			}

			if value == nil {
				if strings.Contains(rawValue, "-----BEGIN") && strings.Count(rawValue, "\n") > 1 {
					file := path.Join(configDir, fmt.Sprintf("%s_%s.pem", section, option))
					log.Debug(
						"writing env var to file",
						zap.String("var", fmt.Sprintf("ICINGADB_%s_%s", match[1], match[2])), zap.String("file", file),
					)

					if err := ioutil.WriteFile(file, []byte(rawValue), 0600); err != nil {
						return err
					}

					value = file
				} else {
					value = rawValue
				}
			}

			sectionCfg[option] = value
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

	hasSchema, errHS := dbHasSchema(idb, map[string]string{
		driver.MySQL:      cfg.Database.Database,
		driver.PostgreSQL: "public",
	}[idb.DriverName()])

	if errHS != nil {
		return errHS
	}

	if !hasSchema {
		log.Debug("importing schema into SQL database")

		file, errOp := os.Open(map[string]string{
			driver.MySQL:      "/mysql.schema.sql.bz2",
			driver.PostgreSQL: "/pgsql.schema.sql.bz2",
		}[idb.DriverName()])

		if errOp != nil {
			return errOp
		}

		defer file.Close()

		schema, errRA := ioutil.ReadAll(bzip2.NewReader(file))
		if errRA != nil {
			return errRA
		}

		if idb.DriverName() == driver.MySQL {
			for _, ddl := range icingaSql.MysqlSplitStatements(string(schema)) {
				if _, errEx := db.Exec(ddl); errEx != nil {
					return errEx
				}
			}
		} else {
			if _, errEx := db.Exec(string(schema)); errEx != nil {
				return errEx
			}
		}
	}

	return nil
}

func dbHasSchema(db *database.DB, dbName string) (bool, error) {
	type one struct {
		One uint8
	}

	rows, errQr := db.Query(
		db.Rebind("SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=? AND TABLE_NAME='icingadb_schema'"),
		dbName,
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

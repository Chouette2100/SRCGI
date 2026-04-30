package main

import (
	"fmt"
	"log"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/srdblib/v3"
)

func newDB(configFile string) (dbnew *sql.DB, dbmapnew *gorp.DbMap, err error) {

	var dbconfig *srdblib.DBConfig
	dbnew, dbconfig, err = srdblib.OpenDb(configFile)
	if err != nil {
		log.Printf("Database error. err = %v\n", err)
		err = fmt.Errorf("srdblib.OpenDb(): %w", err)
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}

	dbnew.SetMaxOpenConns(8)
	dbnew.SetMaxIdleConns(12)

	dbnew.SetConnMaxLifetime(time.Minute * 5)
	dbnew.SetConnMaxIdleTime(time.Minute * 5)

	// defer dbnew.Close()
	log.Printf("%+v\n", dbconfig)

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	dbmapnew = &gorp.DbMap{Db: dbnew,
		Dialect:         dial,
		ExpandSliceArgs: true, //スライス引数展開オプションを有効化する
	}

	srdblib.AddTableWithName(dbmapnew)

	return
}

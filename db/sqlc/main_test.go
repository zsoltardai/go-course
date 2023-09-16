package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/zsoltardai/simple_bank/util"
)

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../../")

	if err != nil {
		log.Fatal("Couldn't load envioment variables!", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}

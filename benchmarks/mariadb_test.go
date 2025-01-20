package benchmarks

import (
	"database/sql"
	"runtime"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/rah-0/hyperion/utils"
	"github.com/rah-0/hyperion/utils/testutil"
)

var (
	db *sql.DB
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() error {
			if err := dbConnect(); err != nil {
				return err
			}
			if err := dbTableSampleACreate(); err != nil {
				return err
			}

			return nil
		},
		UnloadResources: func() error {
			if err := dbTableSampleADrop(); err != nil {
				return err
			}
			if err := dbDisconnect(); err != nil {
				return err
			}

			return nil
		},
	})
}

func BenchmarkMariaDBSingleInsertFixedData(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := dbTableSampleAInsert("9xKf3QpLm2Ry7UbHt6NwEjVg8As5OcIy4B")
		if err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkMariaDBSingleInsertRandomData(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := dbTableSampleAInsert(uuid.NewString())
		if err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}

func dbTableSampleAInsert(fieldAValue string) error {
	stmt, err := db.Prepare(`
		INSERT INTO sample_a (FieldA)
		VALUES (?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement
	_, err = stmt.Exec(fieldAValue)
	if err != nil {
		return err
	}

	return nil
}

func dbConnect() (err error) {
	usr := utils.GetEnvVariable("MariaDB_Usr")
	pwd := utils.GetEnvVariable("MariaDB_Pwd")
	ip := utils.GetEnvVariable("MariaDB_IP")
	port := utils.GetEnvVariable("MariaDB_Port")
	dbName := utils.GetEnvVariable("MariaDB_Name")

	dsn := usr + ":" + pwd + "@tcp(" + ip + ":" + port + ")/" + dbName
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(runtime.NumCPU())
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetConnMaxIdleTime(time.Minute * 1)

	return nil
}

func dbDisconnect() error {
	return db.Close()
}

func dbTableSampleACreate() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS sample_a (
		Id INT(11) NOT NULL AUTO_INCREMENT,
		FirstInsert DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
		LastUpdate DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
		FieldA VARCHAR(36) NOT NULL DEFAULT '',
		PRIMARY KEY (Id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPACT`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func dbTableSampleADrop() error {
	dropTableQuery := `DROP TABLE IF EXISTS sample_a;`

	_, err := db.Exec(dropTableQuery)
	if err != nil {
		return err
	}
	return nil
}

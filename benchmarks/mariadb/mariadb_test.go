package mariadb

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/rah-0/hyperion/utils/testutil"
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
	defer testutil.RecoverBenchHandler(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := dbTableSampleAInsert("9xKf3QpLm2Ry7UbHt6NwEjVg8As5OcIy4B")
		if err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkMariaDBSingleInsertRandomData(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := dbTableSampleAInsert(uuid.NewString())
		if err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}

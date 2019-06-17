// +build docker

package odbc

import (
	"testing"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/sql/sqltest"
	"github.com/cayleygraph/cayley/internal/dock"
)

/*
func makeInformixVersion(image string) sqltest.DatabaseFunc {
	return func(t testing.TB) (string, graph.Options, func()) {
		var conf dock.Config

		conf.Image = image
		conf.Tty = true
		conf.Env = []string{
			`MYSQL_ROOT_PASSWORD=changeme`,
			`MYSQL_DATABASE=testdb`,
		}

		addr, closer := dock.RunAndWait(t, conf, "60000", nil)
		addr = `informix:informix@tcp(` + addr + `)/testdb`
		return addr, nil, func() {
			closer()
		}
	}
}

func TestInformix(t *testing.T) {
	sqltest.TestAll(t, Type, makeInformixVersion(mysqlImage), nil)
}

func TestMariaDB(t *testing.T) {
	sqltest.TestAll(t, Type, makeMysqlVersion(mariadbImage), nil)
}

func BenchmarkMysql(t *testing.B) {
	sqltest.BenchmarkAll(t, Type, makeMysqlVersion(mysqlImage), nil)
}

func BenchmarkMariadb(t *testing.B) {
	sqltest.BenchmarkAll(t, Type, makeMysqlVersion(mariadbImage), nil)
}
/*

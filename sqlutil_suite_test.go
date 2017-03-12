package sqlutil_test

import (
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	db     *sql.DB
	dbfile string
)

func TestSqlutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sqlutil Suite")
}

var _ = BeforeSuite(func() {
	tmpfile, err := ioutil.TempFile("", "sqlutil")
	Expect(err).To(BeNil())

	dbfile = tmpfile.Name()

	db, err = sql.Open("sqlite3", dbfile)
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	Expect(db.Close()).To(Succeed())
	Expect(os.Remove(dbfile)).To(Succeed())
})

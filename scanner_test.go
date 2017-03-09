package sqlutil_test

import (
	"database/sql"

	"github.com/phogolabs/sqlutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scanner", func() {
	type student struct {
		ID   string `sql:"id,varchar(50),pk"`
		Name string `sql:"name,text"`
	}

	var db *sql.DB

	BeforeEach(func() {
		var err error
		db, err = sql.Open("sqlite3", "sqlutil.db")
		Expect(err).To(BeNil())

		Expect(sqlutil.CreateTable(db, &student{})).To(Succeed())

		_, err = db.Exec("INSERT INTO student (id,name) VALUES ('e73sg9','hello')")
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		_, err := db.Exec("drop table student")
		Expect(err).To(BeNil())
		Expect(db.Close()).To(Succeed())
	})

	It("reads the content correctly", func() {
		rows, err := db.Query("SELECT id,name,name as full FROM student")
		Expect(err).To(BeNil())
		defer rows.Close()

		Expect(rows.Next()).To(BeTrue())

		record := student{}

		Expect(sqlutil.Scan(rows, &record)).To(Succeed())
		Expect(record.ID).To(Equal("e73sg9"))
		Expect(record.Name).To(Equal("hello"))
	})
})

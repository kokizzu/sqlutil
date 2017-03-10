package sqlutil_test

import (
	"database/sql"

	"github.com/phogolabs/sqlutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crud", func() {
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
	})

	AfterEach(func() {
		_, err := db.Exec("drop table student")
		Expect(err).To(BeNil())
		Expect(db.Close()).To(Succeed())
	})

	It("inserts user correctly", func() {
		Expect(sqlutil.Insert(db, &student{
			ID:   "1234",
			Name: "Jack",
		})).To(Succeed())

		rows, err := db.Query("SELECT id,name FROM student")
		Expect(err).To(BeNil())
		defer rows.Close()

		Expect(rows.Next()).To(BeTrue())

		record := student{}

		Expect(sqlutil.Scan(rows, &record)).To(Succeed())
		Expect(record.ID).To(Equal("1234"))
		Expect(record.Name).To(Equal("Jack"))
	})

	It("updates row correctly", func() {
		Expect(sqlutil.Insert(db, &student{
			ID:   "1234",
			Name: "Jack",
		})).To(Succeed())

		Expect(sqlutil.Update(db, &student{
			ID:   "1234",
			Name: "John",
		})).To(Succeed())

		rows, err := db.Query("SELECT id,name FROM student")
		Expect(err).To(BeNil())
		defer rows.Close()

		Expect(rows.Next()).To(BeTrue())

		record := student{}

		Expect(sqlutil.Scan(rows, &record)).To(Succeed())
		Expect(record.ID).To(Equal("1234"))
		Expect(record.Name).To(Equal("John"))
	})

	It("deletes row correctly", func() {
		Expect(sqlutil.Insert(db, &student{
			ID:   "1234",
			Name: "Jack",
		})).To(Succeed())

		Expect(sqlutil.Delete(db, &student{
			ID: "1234",
		})).To(Succeed())

		rows, err := db.Query("SELECT id,name FROM student")
		Expect(err).To(BeNil())
		defer rows.Close()

		Expect(rows.Next()).To(BeFalse())
	})
})

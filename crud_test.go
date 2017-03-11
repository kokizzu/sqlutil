package sqlutil_test

import (
	"github.com/phogolabs/sqlutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crud", func() {
	type student struct {
		ID   string `sql:"id,varchar(50),pk"`
		Name string `sql:"name,text"`
	}

	BeforeEach(func() {
		Expect(sqlutil.CreateTable(db, &student{})).To(Succeed())
	})

	AfterEach(func() {
		_, err := db.Exec("drop table student")
		Expect(err).To(BeNil())
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

	Context("when the provided type is not a pointer", func() {
		It("insert operation returns an error", func() {
			Expect(sqlutil.Insert(db, student{})).To(MatchError("Must be pointer to struct; got student"))
		})

		It("update operation returns an error", func() {
			Expect(sqlutil.Update(db, student{})).To(MatchError("Must be pointer to struct; got student"))
		})

		It("delete operation returns an error", func() {
			Expect(sqlutil.Delete(db, student{})).To(MatchError("Must be pointer to struct; got student"))
		})
	})
})

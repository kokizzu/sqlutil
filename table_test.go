package sqlutil_test

import (
	"github.com/phogolabs/sqlutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TableEntity", func() {
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

	It("queries row by primary key successfully", func() {
		_, err := db.Exec("INSERT INTO student (id, name) VALUES('1', 'Jack')")
		Expect(err).To(BeNil())

		s := &student{ID: "1"}
		Expect(sqlutil.Entity(s).QueryRow(db)).To(Succeed())
		Expect(s.ID).To(Equal("1"))
		Expect(s.Name).To(Equal("Jack"))
	})

	It("inserts user correctly", func() {
		cnt, err := sqlutil.Entity(&student{
			ID:   "1234",
			Name: "Jack",
		}).Insert(db)

		Expect(cnt).To(Equal(int64(1)))
		Expect(err).To(BeNil())

		rows, err := db.Query("SELECT id,name FROM student")
		Expect(err).To(BeNil())
		defer rows.Close()

		Expect(rows.Next()).To(BeTrue())

		record := student{}

		Expect(sqlutil.Entity(&record).Scan(rows)).To(Succeed())
		Expect(record.ID).To(Equal("1234"))
		Expect(record.Name).To(Equal("Jack"))
	})

	It("updates row correctly", func() {
		cnt, err := sqlutil.Entity(&student{
			ID:   "1234",
			Name: "Jack",
		}).Insert(db)

		Expect(cnt).To(Equal(int64(1)))
		Expect(err).To(BeNil())

		cnt, err = sqlutil.Entity(&student{
			ID:   "1234",
			Name: "John",
		}).Update(db)

		Expect(cnt).To(Equal(int64(1)))
		Expect(err).To(BeNil())

		rows, err := db.Query("SELECT id,name FROM student")
		Expect(err).To(BeNil())
		defer rows.Close()

		Expect(rows.Next()).To(BeTrue())

		record := student{}

		Expect(sqlutil.Entity(&record).Scan(rows)).To(Succeed())
		Expect(record.ID).To(Equal("1234"))
		Expect(record.Name).To(Equal("John"))
	})

	It("deletes row correctly", func() {
		cnt, err := sqlutil.Entity(&student{
			ID:   "1234",
			Name: "Jack",
		}).Insert(db)

		Expect(cnt).To(Equal(int64(1)))
		Expect(err).To(BeNil())

		cnt, err = sqlutil.Entity(&student{
			ID: "1234",
		}).Delete(db)

		Expect(cnt).To(Equal(int64(1)))
		Expect(err).To(BeNil())

		rows, err := db.Query("SELECT id,name FROM student")
		Expect(err).To(BeNil())
		defer rows.Close()

		Expect(rows.Next()).To(BeFalse())
	})

	Context("when the provided type is not a pointer", func() {
		It("should panic", func() {
			Expect(func() { sqlutil.Entity(student{}) }).To(Panic())
		})
	})
})

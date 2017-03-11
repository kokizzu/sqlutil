package sqlutil_test

import (
	"github.com/phogolabs/sqlutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scanner", func() {
	type student struct {
		ID   string `sql:"id,varchar(50),pk"`
		Name string `sql:"name,text"`
	}

	BeforeEach(func() {
		Expect(sqlutil.CreateTable(db, &student{})).To(Succeed())
		_, err := db.Exec("INSERT INTO student (id,name) VALUES ('e73sg9','hello')")
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		_, err := db.Exec("drop table student")
		Expect(err).To(BeNil())
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

	Context("when the provided type is not a pointer", func() {
		It("read operation returns an error", func() {
			rows, err := db.Query("SELECT id,name,name as full FROM student")
			Expect(err).To(BeNil())
			defer rows.Close()

			Expect(rows.Next()).To(BeTrue())

			Expect(sqlutil.Scan(rows, student{})).To(MatchError("Must be pointer to struct; got student"))
		})
	})
})

package sqlutil_test

import (
	"time"

	"github.com/phogolabs/sqlutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Entity", func() {
	type student struct {
		ID        string    `sql:"id,varchar(50),pk"`
		Name      string    `sql:"name,text"`
		CreatedAt time.Time `sql:"created_at,timestamp,not_null"`
		UpdatedAt time.Time `sql:"updated_at,timestamp,not_null"`
	}

	BeforeEach(func() {
		Expect(sqlutil.CreateTable(db, &student{})).To(Succeed())
	})

	AfterEach(func() {
		_, err := db.Exec("drop table student")
		Expect(err).To(BeNil())
	})

	It("queries row by primary key successfully", func() {
		_, err := db.Exec("INSERT INTO student (id, name, created_at, updated_at) VALUES('1', 'Jack', datetime(), datetime())")
		Expect(err).To(BeNil())

		s := &student{ID: "1"}
		Expect(sqlutil.Entity(s).QueryRow(db)).To(Succeed())
		Expect(s.ID).To(Equal("1"))
		Expect(s.Name).To(Equal("Jack"))
	})

	It("inserts user correctly", func() {
		s := &student{
			ID:   "1234",
			Name: "Jack",
		}

		cnt, err := sqlutil.Entity(s).Insert(db)
		Expect(cnt).To(Equal(int64(1)))
		Expect(err).To(BeNil())
		Expect(s.CreatedAt).NotTo(Equal(time.Time{}))
		Expect(s.UpdatedAt).NotTo(Equal(time.Time{}))
		Expect(s.UpdatedAt).To(BeTemporally("==", s.CreatedAt))

		rows, err := db.Query("SELECT id,name,created_at,updated_at FROM student")
		Expect(err).To(BeNil())
		defer func() {
			Expect(rows.Close()).To(Succeed())
		}()

		Expect(rows.Next()).To(BeTrue())

		record := student{}

		Expect(sqlutil.Entity(&record).Scan(rows)).To(Succeed())
		Expect(record.ID).To(Equal("1234"))
		Expect(record.Name).To(Equal("Jack"))
		Expect(record.CreatedAt).NotTo(Equal(time.Time{}))
		Expect(record.CreatedAt).To(BeTemporally("<=", time.Now()))
		Expect(record.UpdatedAt).To(BeTemporally("==", record.CreatedAt))
	})

	It("updates row correctly", func() {
		s := &student{
			ID:   "1234",
			Name: "Jack",
		}
		cnt, err := sqlutil.Entity(s).Insert(db)
		Expect(cnt).To(Equal(int64(1)))
		Expect(err).To(BeNil())
		Expect(s.CreatedAt).NotTo(Equal(time.Time{}))
		Expect(s.UpdatedAt).To(BeTemporally("==", s.CreatedAt))

		s.Name = "John"
		cnt, err = sqlutil.Entity(s).Update(db)
		Expect(cnt).To(Equal(int64(1)))
		Expect(err).To(BeNil())
		Expect(s.CreatedAt).NotTo(Equal(time.Time{}))
		Expect(s.UpdatedAt).To(BeTemporally(">", s.CreatedAt))

		rows, err := db.Query("SELECT id,name FROM student")
		Expect(err).To(BeNil())
		defer func() {
			Expect(rows.Close()).To(Succeed())
		}()

		Expect(rows.Next()).To(BeTrue())

		record := student{}

		Expect(sqlutil.Entity(&record).Scan(rows)).To(Succeed())
		Expect(record.ID).To(Equal("1234"))
		Expect(record.Name).To(Equal("John"))
	})

	Context("when the update fields are provided", func() {
		It("updates only the provided fields", func() {
			cnt, err := sqlutil.Entity(&student{
				ID:   "1234",
				Name: "Jack",
			}).Insert(db)

			Expect(cnt).To(Equal(int64(1)))
			Expect(err).To(BeNil())

			cnt, err = sqlutil.Entity(&student{
				ID:   "1234",
				Name: "Peter",
			}).Update(db, sqlutil.Fields{
				"name": "Smith",
			})

			Expect(cnt).To(Equal(int64(1)))
			Expect(err).To(BeNil())

			rows, err := db.Query("SELECT id,name FROM student")
			Expect(err).To(BeNil())
			defer func() {
				Expect(rows.Close()).To(Succeed())
			}()

			Expect(rows.Next()).To(BeTrue())

			record := student{}

			Expect(sqlutil.Entity(&record).Scan(rows)).To(Succeed())
			Expect(record.ID).To(Equal("1234"))
			Expect(record.Name).To(Equal("Smith"))
		})
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
		defer func() {
			Expect(rows.Close()).To(Succeed())
		}()

		Expect(rows.Next()).To(BeFalse())
	})

	Context("when the provided type is not a pointer", func() {
		It("should panic", func() {
			Expect(func() { sqlutil.Entity(student{}) }).To(Panic())
		})
	})
})

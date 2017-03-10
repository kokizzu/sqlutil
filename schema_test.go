package sqlutil_test

import (
	"time"

	"github.com/phogolabs/sqlutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schema", func() {
	AfterEach(func() {
		_, err := db.Exec("drop table m")
		Expect(err).To(BeNil())
	})

	It("creates a table successfully", func() {
		type m struct {
			ID        string    `sql:"id,varchar(50),pk"`
			Name      string    `sql:"name,text,not_null,unique,index"`
			CreatedAt time.Time `sql:"created_at,timestamp,null"`
		}

		Expect(sqlutil.CreateTable(db, &m{})).To(Succeed())

		rows, err := db.Query("pragma table_info(m)")
		Expect(err).To(BeNil())
		defer rows.Close()

		var (
			no           int
			name         string
			dataType     string
			notNull      int
			defaultValue interface{}
			isPK         int
		)

		Expect(rows.Next()).To(BeTrue())
		Expect(rows.Scan(&no, &name, &dataType, &notNull, &defaultValue, &isPK)).To(Succeed())
		Expect(name).To(Equal("id"))
		Expect(dataType).To(Equal("varchar(50)"))
		Expect(notNull).To(Equal(0))
		Expect(isPK).To(Equal(1))

		Expect(rows.Next()).To(BeTrue())
		Expect(rows.Scan(&no, &name, &dataType, &notNull, &defaultValue, &isPK)).To(Succeed())
		Expect(name).To(Equal("name"))
		Expect(dataType).To(Equal("text"))
		Expect(notNull).To(Equal(1))
		Expect(isPK).To(Equal(0))

		Expect(rows.Next()).To(BeTrue())
		Expect(rows.Scan(&no, &name, &dataType, &notNull, &defaultValue, &isPK)).To(Succeed())
		Expect(name).To(Equal("created_at"))
		Expect(dataType).To(Equal("timestamp"))
		Expect(notNull).To(Equal(0))
		Expect(isPK).To(Equal(0))
	})
})

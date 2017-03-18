package sqlutil_test

import (
	"reflect"
	"time"

	"github.com/phogolabs/sqlutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata", func() {
	var metadata *sqlutil.Metadata

	BeforeEach(func() {
		metadata = &sqlutil.Metadata{}
	})

	It("returns constrains as a string", func() {
		f := sqlutil.ColumnConstraintNull | sqlutil.ColumnConstraintNotNull | sqlutil.ColumnConstraintUnique
		Expect(f.String()).To(Equal("UNIQUE NULL NOT NULL"))
	})

	It("retrieves the fields information", func() {
		type m struct {
			ID        string    `sql:"id,varchar(50),pk,not_null,unique" sqlindex:"search" sqlforeignkey:"table1(a)"`
			Name      string    `sql:"name,text,not_null,unique" sqlindex:"name_idx" sqlforeignkey:"table1(b)"`
			CreatedAt time.Time `sql:"created_at,timestamp,null" sqlforeignkey:"table2(c)"`
			RefId     int       `sql:"ref_id,integer" sqlindex:"search" sqlindex:"ref_id"`
			IgnoreMe  string    `sql:"-"`
		}

		t := reflect.ValueOf(m{}).Type()
		schema, err := metadata.Schema(t)
		Expect(err).To(BeNil())
		Expect(schema.Table).To(Equal("m"))

		columns := schema.Columns
		Expect(columns).To(HaveLen(4))

		zero := sqlutil.ColumnConstraint(0)

		Expect(columns[0].Name).To(Equal("id"))
		Expect(columns[0].DataType).To(Equal("varchar(50)"))
		Expect(columns[0].Constraint & sqlutil.ColumnConstraintNotNull).NotTo(Equal(zero))
		Expect(columns[0].Constraint & sqlutil.ColumnConstraintUnique).NotTo(Equal(zero))

		Expect(columns[1].Name).To(Equal("name"))
		Expect(columns[1].DataType).To(Equal("text"))
		Expect(columns[1].Constraint & sqlutil.ColumnConstraintNotNull).NotTo(Equal(zero))
		Expect(columns[1].Constraint & sqlutil.ColumnConstraintUnique).NotTo(Equal(zero))

		Expect(columns[2].Name).To(Equal("created_at"))
		Expect(columns[2].DataType).To(Equal("timestamp"))
		Expect(columns[2].Constraint & sqlutil.ColumnConstraintNull).NotTo(Equal(zero))

		Expect(columns[3].Name).To(Equal("ref_id"))
		Expect(columns[3].DataType).To(Equal("integer"))

		indexes := schema.Indexes
		Expect(indexes).To(HaveLen(3))

		Expect(indexes[0].Name).To(Equal("search"))
		Expect(indexes[0].Columns).To(HaveLen(2))
		Expect(indexes[0].Columns).To(ContainElement("id"))
		Expect(indexes[0].Columns).To(ContainElement("ref_id"))

		Expect(indexes[1].Name).To(Equal("name_idx"))
		Expect(indexes[1].Columns).To(HaveLen(1))
		Expect(indexes[1].Columns).To(ContainElement("name"))

		fk := schema.ForeignKeys
		Expect(fk).To(HaveLen(2))

		Expect(fk[0].Columns).To(HaveLen(2))
		Expect(fk[0].Columns).To(ContainElement("id"))
		Expect(fk[0].Columns).To(ContainElement("name"))
		Expect(fk[0].ReferenceTable).To(Equal("table1"))
		Expect(fk[0].ReferenceTableColumns).To(HaveLen(2))
		Expect(fk[0].ReferenceTableColumns).To(ContainElement("a"))
		Expect(fk[0].ReferenceTableColumns).To(ContainElement("b"))

		Expect(fk[1].Columns).To(HaveLen(1))
		Expect(fk[1].Columns).To(ContainElement("created_at"))
		Expect(fk[1].ReferenceTable).To(Equal("table2"))
		Expect(fk[1].ReferenceTableColumns).To(HaveLen(1))
		Expect(fk[1].ReferenceTableColumns).To(ContainElement("c"))

		Expect(indexes[2].Name).To(Equal("ref_id"))
		Expect(indexes[2].Columns).To(HaveLen(1))
		Expect(indexes[2].Columns).To(ContainElement("ref_id"))

	})

	Context("when a tag is not provided", func() {
		It("returns an error", func() {
			type m struct {
				ID string
			}

			t := reflect.ValueOf(m{}).Type()
			_, err := metadata.Schema(t)
			Expect(err).To(MatchError(`Type "m": Missing tag for field "ID"`))
		})
	})
})

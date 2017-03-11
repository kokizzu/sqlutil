package sqlutil_test

import (
	"github.com/phogolabs/sqlutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MetadataTag", func() {

	It("gets the requested tag", func() {
		t := sqlutil.Tag(`sql:"name,text,not_null,unique" sqlindex:"name"`)
		tags := t.Get("sqlindex")
		Expect(tags).To(HaveLen(1))
		Expect(tags[0]).To(Equal("name"))
	})

	It("looks up the requested tag", func() {
		t := sqlutil.Tag(`sql:"name,text,not_null,unique" sqlindex:"name"`)
		tags, ok := t.Lookup("sqlindex")
		Expect(ok).To(BeTrue())
		Expect(tags).To(HaveLen(1))
		Expect(tags[0]).To(Equal("name"))
	})

	Context("when the tag is presented more than once", func() {
		It("gets the requested tag", func() {
			t := sqlutil.Tag(`sqlindex:"first_name" crazy:"tag" sqlindex:"family_name"`)
			tags := t.Get("sqlindex")
			Expect(tags).To(HaveLen(2))
			Expect(tags[0]).To(Equal("first_name"))
			Expect(tags[1]).To(Equal("family_name"))
		})

		It("looks up the requested tag", func() {
			t := sqlutil.Tag(`sqlindex:"first_name" sqlindex:"family_name"`)
			tags, ok := t.Lookup("sqlindex")
			Expect(ok).To(BeTrue())
			Expect(tags).To(HaveLen(2))
			Expect(tags[0]).To(Equal("first_name"))
			Expect(tags[1]).To(Equal("family_name"))
		})
	})

	Context("when the tag is not presented", func() {
		It("gets the empty tag", func() {
			t := sqlutil.Tag(`sql:"name,text,not_null,unique"`)
			tags := t.Get("sqlindex")
			Expect(tags).To(HaveLen(0))
		})

		It("looks up returns false", func() {
			t := sqlutil.Tag(`sql:"name,text,not_null,unique"`)
			tags, ok := t.Lookup("sqlindex")
			Expect(ok).To(BeFalse())
			Expect(tags).To(HaveLen(0))
		})
	})
})

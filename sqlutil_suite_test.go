package sqlutil_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSqlutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sqlutil Suite")
}

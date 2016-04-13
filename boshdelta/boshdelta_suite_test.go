package boshdelta_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBoshdelta(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Boshdelta Suite")
}

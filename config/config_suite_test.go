package config_test

import (
	. "github.com/noexcs/redis-go/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("config", func() {
	Context("Parse config file", func() {
		Setup()
		It("Port", func() {
			Expect(Properties.Port).To(Equal(6399))
		})
		It("Bind", func() {
			Expect(Properties.Bind).To(Equal("127.0.0.1"))
		})
		It("Debug", func() {
			Expect(Properties.Debug).To(Equal(true))
		})
		It("Requirepass", func() {
			Expect(Properties.Requirepass).To(Equal("123456"))
		})
	})
})

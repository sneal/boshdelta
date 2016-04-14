package boshdelta_test

import (
	"os"
	"path/filepath"

	. "github.com/sneal/bosh-delta/boshdelta"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Boshdelta", func() {
	Context("Redis BOSH Release 12", func() {
		var (
			releasePath string
			release     *Release
		)

		BeforeEach(func() {
			dir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			releasePath = filepath.Join(dir, "../fixtures/redis-boshrelease-12.tgz")
			release, err = NewRelease(releasePath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can be loaded", func() {
			Expect(release).NotTo(BeNil())
		})
		It("has the correct path", func() {
			Expect(release.Path).To(Equal(releasePath))
		})
		It("has uncommitted changes", func() {
			Expect(release.UncommittedChanges).To(BeTrue())
		})
		It("is version 12", func() {
			Expect(release.Version).To(Equal("12"))
		})
		It("is named redis", func() {
			Expect(release.Name).To(Equal("redis"))
		})
		It("contains 2 jobs", func() {
			Expect(release.Jobs).To(HaveLen(2))
		})
		It("contains the acceptance-tests job", func() {
			job := release.Jobs[0]
			Expect(job.Name).To(Equal("acceptance-tests"))
			Expect(job.Version).To(Equal("d1b927ac3839a52f7caa29ac80878d8625fb4bcc"))
			Expect(job.Sha1).To(Equal("dee26bdad425c0276df120c922551c994777cb42"))
			Expect(job.Fingerprint).To(Equal("d1b927ac3839a52f7caa29ac80878d8625fb4bcc"))
		})
		It("contains the redis job", func() {
			job := release.Jobs[1]
			Expect(job.Name).To(Equal("redis"))
			Expect(job.Version).To(Equal("a52b9ad8bf105901dbd55e6dfb3af359b24a2d14"))
			Expect(job.Sha1).To(Equal("fbb1fe127a25429c84a120c8f676c3e7f26ff051"))
			Expect(job.Fingerprint).To(Equal("a52b9ad8bf105901dbd55e6dfb3af359b24a2d14"))
		})
	})
})

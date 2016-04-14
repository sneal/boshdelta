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
		Context("acceptance-tests job", func() {
			var job *Job
			BeforeEach(func() {
				job = release.FindJob("acceptance-tests")
				Expect(job).ToNot(BeNil())
			})
			It("has the correct name", func() {
				Expect(job.Name).To(Equal("acceptance-tests"))
			})
			It("has the correct SHA1", func() {
				Expect(job.Sha1).To(Equal("dee26bdad425c0276df120c922551c994777cb42"))
			})
			It("has the correct version", func() {
				Expect(job.Version).To(Equal("d1b927ac3839a52f7caa29ac80878d8625fb4bcc"))
			})
			It("has the correct finger print", func() {
				Expect(job.Fingerprint).To(Equal("d1b927ac3839a52f7caa29ac80878d8625fb4bcc"))
			})
			It("has the correct properties", func() {
				Expect(job.Properties).To(HaveKey("redis.port"))
				Expect(job.Properties).To(HaveKey("redis.password"))
				Expect(job.Properties).To(HaveKey("redis.master"))
				Expect(job.Properties).To(HaveKey("redis.slave"))
			})
		})
		Context("redis job", func() {
			var job *Job
			BeforeEach(func() {
				job = release.FindJob("redis")
				Expect(job).ToNot(BeNil())
			})
			It("has the correct name", func() {
				Expect(job.Name).To(Equal("redis"))
			})
			It("has the correct SHA1", func() {
				Expect(job.Sha1).To(Equal("fbb1fe127a25429c84a120c8f676c3e7f26ff051"))
			})
			It("has the correct version", func() {
				Expect(job.Version).To(Equal("a52b9ad8bf105901dbd55e6dfb3af359b24a2d14"))
			})
			It("has the correct finger print", func() {
				Expect(job.Fingerprint).To(Equal("a52b9ad8bf105901dbd55e6dfb3af359b24a2d14"))
			})
			It("has the correct properties", func() {
				Expect(job.Properties).To(HaveKey("redis.port"))
				Expect(job.Properties).To(HaveKey("redis.password"))
				Expect(job.Properties).To(HaveKey("redis.master"))
				Expect(job.Properties).To(HaveKey("consul.service.name"))
				Expect(job.Properties).To(HaveKey("health.interval"))
				Expect(job.Properties).To(HaveKey("health.disk.critical"))
				Expect(job.Properties).To(HaveKey("health.disk.warning"))
			})
		})
	})
})

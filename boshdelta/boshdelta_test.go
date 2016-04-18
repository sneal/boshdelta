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
		var release *Release

		BeforeEach(func() {
			release = loadBoshRelease("redis-boshrelease-12.tgz")
		})

		It("can be loaded", func() {
			Expect(release).NotTo(BeNil())
		})
		It("has the correct path", func() {
			Expect(release.Path).To(ContainSubstring("redis-boshrelease-12.tgz"))
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
		Context("UniqueProperties", func() {
			var props map[string]*Property
			BeforeEach(func() {
				props = release.UniqueProperties()
			})
			It("contains unique properties across all jobs", func() {
				Expect(props).To(HaveLen(8))
				Expect(props).To(HaveKey("redis.port"))
				Expect(props).To(HaveKey("redis.password"))
				Expect(props).To(HaveKey("redis.master"))
				Expect(props).To(HaveKey("redis.slave"))
				Expect(props).To(HaveKey("consul.service.name"))
				Expect(props).To(HaveKey("health.interval"))
				Expect(props).To(HaveKey("health.disk.critical"))
				Expect(props).To(HaveKey("health.disk.warning"))
			})
		})
	})
	Context("Compare Redis BOSH release 1 to release 12", func() {
		var (
			release1  *Release
			release12 *Release
			delta     *Delta
		)

		BeforeEach(func() {
			release1 = loadBoshRelease("redis-boshrelease-1.tgz")
			release12 = loadBoshRelease("redis-boshrelease-12.tgz")
			delta = CompareReleases(release1, release12)
		})

		It("contains new properties from all jobs", func() {
			Expect(delta.ContainsProperty("redis.master")).To(BeTrue())
			Expect(delta.ContainsProperty("redis.slave")).To(BeTrue())
			Expect(delta.ContainsProperty("consul.service.name")).To(BeTrue())
			Expect(delta.ContainsProperty("health.interval")).To(BeTrue())
			Expect(delta.ContainsProperty("health.disk.critical")).To(BeTrue())
			Expect(delta.ContainsProperty("health.disk.warning")).To(BeTrue())
		})

		It("does not contain existing properties", func() {
			Expect(delta.ContainsProperty("redis.port")).To(BeFalse())
			Expect(delta.ContainsProperty("redis.password")).To(BeFalse())
		})
	})
})

var _ = Describe("Pivnet release", func() {
	Context("Invalid .pivotal file", func() {
		var (
			release *PivnetRelease
			err     error
		)
		BeforeEach(func() {
			release, err = NewPivnetReleaseFromFile(fixturePath("redis-boshrelease-1.tgz"))
		})
		It("returns an error when the file doesn't end in .pivotal", func() {
			Expect(err).To(MatchError("Expected a .pivotal file, but instead got a .tgz file"))
		})
	})
	Context("Redis release 1.5", func() {
		var release *PivnetRelease

		BeforeEach(func() {
			// COPYFILE_DISABLE=1 tar czfv <src>/fixtures/p-redis-1.5.0.pivotal --exclude=".DS_Store" .
			release = loadPivnetRelease("p-redis-1.5.0.pivotal")
		})

		It("can be loaded", func() {
			Expect(release).NotTo(BeNil())
		})
		It("has the correct path", func() {
			Expect(release.Path).To(ContainSubstring("p-redis-1.5.0.pivotal"))
		})
		It("contains two releases", func() {
			Expect(release.Releases).To(HaveLen(2))
		})
		Context("redis BOSH release", func() {
			var redisRelease Release
			BeforeEach(func() {
				redisRelease = release.Releases[0]
			})
			It("isn't nil", func() {
				Expect(redisRelease).NotTo(BeNil())
			})
			It("has the correct relative path", func() {
				Expect(redisRelease.Path).To(Equal("./releases/redis-boshrelease-12.tgz"))
			})
		})
		Context("redis Xip release", func() {
			var xipRelease Release
			BeforeEach(func() {
				xipRelease = release.Releases[1]
			})
			It("isn't nil", func() {
				Expect(xipRelease).NotTo(BeNil())
			})
			It("has the correct relative path", func() {
				Expect(xipRelease.Path).To(Equal("./releases/xip-release-2.tgz"))
			})
		})
	})
})

func loadBoshRelease(releaseFileName string) *Release {
	release, err := NewReleaseFromFile(fixturePath(releaseFileName))
	Expect(err).NotTo(HaveOccurred())
	return release
}

func loadPivnetRelease(releaseFileName string) *PivnetRelease {
	release, err := NewPivnetReleaseFromFile(fixturePath(releaseFileName))
	Expect(err).NotTo(HaveOccurred())
	return release
}

func fixturePath(releaseFileName string) string {
	dir, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	return filepath.Join(dir, "../fixtures/", releaseFileName)
}

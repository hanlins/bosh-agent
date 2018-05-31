package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-agent/agent/action"
	"github.com/cloudfoundry/bosh-agent/platform/platformfakes"
	boshsettings "github.com/cloudfoundry/bosh-agent/settings"
	fakesettings "github.com/cloudfoundry/bosh-agent/settings/fakes"
	bosherrors "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var _ = Describe("ListDisk", func() {
	var (
		settingsService *fakesettings.FakeSettingsService
		platform        *platformfakes.FakePlatform
		logger          boshlog.Logger
		action          ListDiskAction
	)

	BeforeEach(func() {
		settingsService = &fakesettings.FakeSettingsService{}
		platform = &platformfakes.FakePlatform{}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		action = NewListDisk(settingsService, platform, logger)

		platform.IsPersistentDiskMountedStub = func(settings boshsettings.DiskSettings) (bool, error) {
			if settings.Path == "/dev/sdb" || settings.Path == "/dev/sdc" {
				return true, nil
			}

			return false, nil
		}
	})

	AssertActionIsSynchronousForVersion(action, 1)
	AssertActionIsSynchronousForVersion(action, 2)
	AssertActionIsAsynchronousForVersion(action, 3)

	AssertActionIsNotPersistent(action)
	AssertActionIsLoggable(action)

	AssertActionIsNotResumable(action)
	AssertActionIsNotCancelable(action)

	It("list disk run", func() {
		settingsService.Settings.Disks = boshsettings.Disks{
			Persistent: map[string]interface{}{
				"volume-1": "/dev/sda",
				"volume-2": "/dev/sdb",
				"volume-3": "/dev/sdc",
			},
		}

		value, err := action.Run()
		Expect(err).ToNot(HaveOccurred())
		values, ok := value.([]string)
		Expect(ok).To(BeTrue())
		Expect(values).To(ContainElement("volume-2"))
		Expect(values).To(ContainElement("volume-3"))
		Expect(len(values)).To(Equal(2))

		Expect(settingsService.SettingsWereLoaded).To(BeTrue())
	})

	Context("when unable to loadsettings", func() {
		BeforeEach(func() {
			settingsService.LoadSettingsError = bosherrors.Error("fake loadsettings error")
		})

		It("should return an error", func() {
			settingsService.Settings.Disks = boshsettings.Disks{
				Persistent: map[string]interface{}{
					"volume-1": "/dev/sda",
					"volume-2": "/dev/sdb",
					"volume-3": "/dev/sdc",
				},
			}

			_, err := action.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Refreshing the settings: fake loadsettings error"))
		})
	})
})

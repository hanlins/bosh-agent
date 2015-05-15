package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"

	"github.com/cloudfoundry/bosh-agent/agent/action"
	"github.com/cloudfoundry/bosh-agent/platform/cert/fakes"
	boshsettings "github.com/cloudfoundry/bosh-agent/settings"
	"github.com/cloudfoundry/bosh-utils/logger"
)

func init() {
	Describe("UpdateSettings", func() {
		var updateAction action.UpdateSettingsAction
		var certManager *fakes.FakeManager
		var log logger.Logger
		BeforeEach(func() {
			log = logger.NewLogger(logger.LevelNone)
			certManager = new(fakes.FakeManager)
			updateAction = action.NewUpdateSettings(certManager, log)
		})

		It("is synchronous", func() {
			Expect(updateAction.IsAsynchronous()).To(BeTrue())
		})

		It("is not persistent", func() {
			Expect(updateAction.IsPersistent()).To(BeFalse())
		})

		It("returns 'updated' on success", func() {
			newSettings := boshsettings.Settings{}
			result, err := updateAction.Run(newSettings)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal("updated"))
		})

		Context("When updating the certificates fails", func() {
			BeforeEach(func() {
				log = logger.NewLogger(logger.LevelNone)
				certManager = new(fakes.FakeManager)
				certManager.UpdateCertificatesReturns(errors.New("Error"))
				updateAction = action.NewUpdateSettings(certManager, log)
			})

			It("returns the error", func() {
				newSettings := boshsettings.Settings{}
				result, err := updateAction.Run(newSettings)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeEmpty())
			})
		})
	})
}

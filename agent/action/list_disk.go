package action

import (
	"errors"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshplatform "github.com/cloudfoundry/bosh-agent/platform"
	boshsettings "github.com/cloudfoundry/bosh-agent/settings"
)

type ListDiskAction struct {
	settingsService boshsettings.Service
	platform        boshplatform.Platform
	logger          boshlog.Logger
}

func NewListDisk(
	settingsService boshsettings.Service,
	platform boshplatform.Platform,
	logger boshlog.Logger,
) (action ListDiskAction) {
	action.settingsService = settingsService
	action.platform = platform
	action.logger = logger
	return
}

func (a ListDiskAction) IsAsynchronous() bool {
	return false
}

func (a ListDiskAction) IsPersistent() bool {
	return false
}

func (a ListDiskAction) Run() (interface{}, error) {
	settings := a.settingsService.GetSettings()
	volumeIDs := []string{}

	for volumeID, devicePath := range settings.Disks.Persistent {
		var isMounted bool

		isMounted, err := a.platform.IsPersistentDiskMounted(devicePath)
		if err != nil {
			return nil, bosherr.WrapErrorf(err, "Checking whether device %s is mounted", devicePath)
		}

		if isMounted {
			volumeIDs = append(volumeIDs, volumeID)
		} else {
			a.logger.Debug("list-disk-action", "Volume '%s' not mounted", volumeID)
		}
	}

	return volumeIDs, nil
}

func (a ListDiskAction) Resume() (interface{}, error) {
	return nil, errors.New("not supported")
}

func (a ListDiskAction) Cancel() error {
	return errors.New("not supported")
}

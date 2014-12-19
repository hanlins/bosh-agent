package cdrom

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshudev "github.com/cloudfoundry/bosh-agent/platform/udevdevice"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type LinuxCdrom struct {
	udev       boshudev.UdevDevice
	devicePath string
	runner     boshsys.CmdRunner
}

func NewLinuxCdrom(devicePath string, udev boshudev.UdevDevice, runner boshsys.CmdRunner) (cdrom LinuxCdrom) {
	cdrom = LinuxCdrom{
		udev:       udev,
		devicePath: devicePath,
		runner:     runner,
	}
	return
}

func (cdrom LinuxCdrom) WaitForMedia() (err error) {
	cdrom.udev.KickDevice("/dev/sr0")
	err = cdrom.udev.Settle()
	if err != nil {
		err = bosherr.WrapError(err, "Waiting for udev to settle")
		return
	}

	err = cdrom.udev.EnsureDeviceReadable(cdrom.devicePath)
	return
}

func (cdrom LinuxCdrom) Mount(mountPath string) (err error) {
	_, stderr, _, err := cdrom.runner.RunCommand("mount", cdrom.devicePath, mountPath)
	if err != nil {
		err = bosherr.WrapErrorf(err, "Mounting CDROM: %s", stderr)
	}
	return
}

func (cdrom LinuxCdrom) Unmount() (err error) {
	_, stderr, _, err := cdrom.runner.RunCommand("umount", cdrom.devicePath)
	if err != nil {
		err = bosherr.WrapErrorf(err, "Unmounting CDROM: %s", stderr)
	}
	return
}

func (cdrom LinuxCdrom) Eject() (err error) {
	_, stderr, _, err := cdrom.runner.RunCommand("eject", cdrom.devicePath)
	if err != nil {
		err = bosherr.WrapErrorf(err, "Ejecting CDROM: %s", stderr)
	}
	return
}

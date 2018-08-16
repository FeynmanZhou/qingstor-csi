package neonsan

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"os/exec"
	"strconv"
)

const (
	Int64Max        int64  = int64(^uint64(0) >> 1)
	PluginFolder    string = "/var/lib/kubelet/plugins/"
	DefaultPoolName string = "kube"
)

const (
	kib    int64 = 1024
	mib    int64 = kib * 1024
	gib    int64 = mib * 1024
	gib100 int64 = gib * 100
	tib    int64 = gib * 1024
	tib100 int64 = tib * 100
)

const (
	FileSystemExt3    string = "ext3"
	FileSystemExt4    string = "ext4"
	FileSystemXfs     string = "xfs"
	DefaultFileSystem string = FileSystemExt4
)

var (
	ConfigFilePath string = "/etc/neonsan/qbd.conf"
)

// ExecCommand
// Return cases:	normal output,	nil:	normal output
//					error logs,		error:	command execute error
func ExecCommand(command string, args []string) ([]byte, error) {
	glog.Infof("execCommand: command = \"%s\", args = \"%v\"", command, args)
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("code [%s]: message [%s]", err.Error(), output)
	}
	return output, nil
}

// ContainsVolumeCapability
// Does Array of VolumeCapability_AccessMode contain the volume capability of subCaps
func ContainsVolumeCapability(accessModes []*csi.VolumeCapability_AccessMode, subCaps *csi.VolumeCapability) bool {
	for _, cap := range accessModes {
		if cap.GetMode() == subCaps.GetAccessMode().GetMode() {
			return true
		}
	}
	return false
}

// ContainsVolumeCapabilities
// Does array of VolumeCapability_AccessMode contain volume capabilities of subCaps
func ContainsVolumeCapabilities(accessModes []*csi.VolumeCapability_AccessMode, subCaps []*csi.VolumeCapability) bool {
	for _, v := range subCaps {
		if !ContainsVolumeCapability(accessModes, v) {
			return false
		}
	}
	return true
}

// FormatVolumeSize convert volume size properly
func FormatVolumeSize(inputSize int64, step int64) int64 {
	if inputSize <= gib || step < 0 {
		return gib
	}
	remainder := inputSize % step
	if remainder != 0 {
		return inputSize - remainder + step
	}
	return inputSize
}

// Check file system type
// Support: ext3, ext4 and xfs
func IsValidFileSystemType(fs string) bool {
	switch fs {
	case FileSystemExt3:
		return true
	case FileSystemExt4:
		return true
	case FileSystemXfs:
		return true
	default:
		return false
	}
}

//	ParseIntToDec convert number string to decimal number string
func ParseIntToDec(hex string) (dec string) {
	i64, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatInt(i64, 10)
}
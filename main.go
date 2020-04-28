package main

import (
    "fmt"
    "time"
    "os"
    "github.com/pkg/errors"

    pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
    dpapi "github.com/intel/intel-device-plugins-for-kubernetes/pkg/deviceplugin"
)

const (
    devfsDriDirectory  = "/dev/tpmrm0"
    devTpmFsDirectory = "/dev/tpm0"

    namespace = "color.example.com"
)

type devicePlugin struct {
	devfsDir string
    tpmFsDir string
}

func newDevicePlugin(devfsDir string, tpmFsDir string) *devicePlugin {
	return &devicePlugin{
		devfsDir:         devfsDir,
        tpmFsDir:         tpmFsDir,
	}
}

func (dp *devicePlugin) Scan(notifier dpapi.Notifier) error {
	for {
		devTree, err := dp.scan()
		if err != nil {
			fmt.Errorf("Failed to scan: %s", err)
		}

		notifier.Notify(devTree)

		time.Sleep(5 * time.Second)
	}
}

func (dp *devicePlugin) scan() (dpapi.DeviceTree, error) {

    fmt.Printf("attempting to read %s\n", dp.devfsDir)
    
    _, err := os.Stat(dp.devfsDir)
    if err == nil {
        fmt.Printf("file %s exists\n", dp.devfsDir)
    } else if os.IsNotExist(err) {
        fmt.Printf("file %s not exists\n", dp.devfsDir)
        return nil, errors.Wrap(err, "Can't read /dev/mem")
    } else {
        fmt.Printf("file %s stat error\n: %v", dp.devfsDir, err)
        return nil, errors.Wrap(err, "Permission denied /dev/mem")
    }

    _, er := os.Stat(dp.tpmFsDir)
    if er == nil {
        fmt.Printf("file %s exists\n", dp.tpmFsDir)
    } else if os.IsNotExist(err) {
        fmt.Printf("file %s not exists\n", dp.tpmFsDir)
        return nil, errors.Wrap(er, "Can't read /dev/tpm")
    } else {
        fmt.Printf("file %s stat error\n: %v", dp.tpmFsDir, er)
        return nil, errors.Wrap(err, "Permission denied /dev/mem")
    }
    
    devTree := dpapi.NewDeviceTree()
    
    var nodes []pluginapi.DeviceSpec
    
    nodes = append(nodes, pluginapi.DeviceSpec{
        HostPath: dp.devfsDir,
        ContainerPath: dp.devfsDir,
        Permissions:   "rw",
    })

    var tpmNodes  []pluginapi.DeviceSpec

    tpmNodes = append(tpmNodes, pluginapi.DeviceSpec{
        HostPath: dp.tpmFsDir,
        ContainerPath: dp.tpmFsDir,
        Permissions:   "rw",
    })


    devID := "id1";
    fmt.Printf("adding %s with id %s", dp.devfsDir, devID)
    devTree.AddDevice("yellow", devID, dpapi.NewDeviceInfo(pluginapi.Healthy, nodes, nil, nil))

    fmt.Print("adding %s with id %s", dp.tpmFsDir, devID)
    devTree.AddDevice("red", "id2", dpapi.NewDeviceInfo(pluginapi.Healthy, tpmNodes, nil, nil))

    return devTree, nil
}


func main() {
    fmt.Println("my device plugin started")

    plugin := newDevicePlugin(devfsDriDirectory, devTpmFsDirectory)
    manager := dpapi.NewManager(namespace, plugin)
    manager.Run()
}

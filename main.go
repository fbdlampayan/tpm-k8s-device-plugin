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

    namespace = "color.example.com"
)

type devicePlugin struct {
	devfsDir string
}

func newDevicePlugin(devfsDir string) *devicePlugin {
	return &devicePlugin{
		devfsDir:         devfsDir,
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
    
    devTree := dpapi.NewDeviceTree()
    
    var nodes []pluginapi.DeviceSpec
    nodes = append(nodes, pluginapi.DeviceSpec{
        HostPath: dp.devfsDir,
        ContainerPath: dp.devfsDir,
        Permissions:   "rw",
    })

    devID := "id1";
    
    fmt.Printf("adding %s with id %s", dp.devfsDir, devID)
    devTree.AddDevice("yellow", devID, dpapi.NewDeviceInfo(pluginapi.Healthy, nodes, nil, nil))
    devTree.AddDevice("red", "id2", dpapi.NewDeviceInfo(pluginapi.Healthy, nodes, nil, nil))

    return devTree, nil
}


func main() {
    fmt.Println("my device plugin started")

    plugin := newDevicePlugin(devfsDriDirectory)
    manager := dpapi.NewManager(namespace, plugin)
    manager.Run()
}

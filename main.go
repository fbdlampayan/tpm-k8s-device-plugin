package main

import (
    "fmt"
    "time"
    "os"
    "path"
    "flag"
    "regexp"
    "io/ioutil"
    "github.com/pkg/errors"

    pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
    dpapi "github.com/intel/intel-device-plugins-for-kubernetes/pkg/deviceplugin"
)

const (
    devicesHostDirectory  = "/dev"
    tpmRmDeviceRegex = "mem" //`^tpmrm[0-9]*$`

    namespace = "fbdl.device.com"
    deviceType = "tpmrm"

    scanPeriod = 5 * time.Second
)

type devicePlugin struct {
    devHostDir string
    sharedCapacity int
    tpmrmDeviceReg *regexp.Regexp
    scanTicker *time.Ticker
    scanDone chan bool
}

func newDevicePlugin(devHostDirectory string, capacity int) *devicePlugin {
	return &devicePlugin{
        devHostDir:       devHostDirectory,
        sharedCapacity:   capacity,
        tpmrmDeviceReg:   regexp.MustCompile(tpmRmDeviceRegex),
        scanTicker:       time.NewTicker(scanPeriod),
        scanDone:         make(chan bool, 1),
	}
}

func (dp *devicePlugin) Scan(notifier dpapi.Notifier) error {
    defer dp.scanTicker.Stop()
    var previouslyFound int = -1

	for {
		devTree, err := dp.scan()
		if err != nil {
			fmt.Errorf("Failed to scan: %s", err)
		}

        found := len(devTree)
        if found != previouslyFound {
            fmt.Printf("TPM scan update: devices found: %d\n", found)
            previouslyFound = found
        }

		notifier.Notify(devTree)

        select {
        case <-dp.scanDone:
            return nil
        case <-dp.scanTicker.C:
        }
    }
}

func (dp *devicePlugin) scan() (dpapi.DeviceTree, error) {
    files, err := ioutil.ReadDir(dp.devHostDir)
    if err != nil {
        return nil, errors.Wrap(err, "Can't read dev directory")
    }

    devTree := dpapi.NewDeviceTree()
    for _, f := range files {
        var nodes []pluginapi.DeviceSpec

        if !dp.tpmrmDeviceReg.MatchString(f.Name()) {
            continue
        }

        devPath := path.Join(dp.devHostDir, f.Name())

        fmt.Printf("Adding %s to tmprm %s\n", devPath, f.Name())
        nodes = append(nodes, pluginapi.DeviceSpec{
            HostPath: devPath,
            ContainerPath: devPath,
            Permissions: "rw",
        })

        if len(nodes) > 0 {
            for i := 0; i < dp.sharedCapacity; i++ {
                devID := fmt.Sprintf("%s-%d", f.Name(), i)
                fmt.Printf("device ID: %s for device: %+v\n", devID, f)
                devTree.AddDevice(deviceType, devID, dpapi.NewDeviceInfo(pluginapi.Healthy, nodes, nil, nil, nil))
            }
        }
    }

    return devTree, nil
}


func main() {
    var capacityDesired int

    flag.IntVar(&capacityDesired, "capacity", 1, "number of pods sharing the same tpmrm device file")
    flag.Parse()

    if capacityDesired < 1 {
        fmt.Println("capacity must be greater than zero")
        os.Exit(1)
    }

    fmt.Println("tpm device plugin started")

    plugin := newDevicePlugin(devicesHostDirectory, capacityDesired)
    manager := dpapi.NewManager(namespace, plugin)
    manager.Run()
}

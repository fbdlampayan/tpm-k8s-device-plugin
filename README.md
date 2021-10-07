# tpm-k8s-device-plugin
custom kubernetes device plugin using Intel's library

# About this
This is a sample implementation of a kubernetes device plugin using Intel's SDK library [intel-device-plugin-for-kubernetes](https://github.com/intel/intel-device-plugins-for-kubernetes).

It demonstrates how to implement a custom k8s device plugin with it and be able to mount a device unto a pod within a kubernetes cluster. For instance, the Trusted Platform Module (TPM) chip.

# Usage

Example manifests on how to deploy and use are under the `examples` directory

- Deploying the plugin is demonstrated in `examples/sample-plugin-daemonset.yaml`, where it mounts the device under volumes and uses hostPath to locate its location from the host machine's /dev directory.

```
      volumes:
      - name: dev-mem
        hostPath:
          path: /dev/tpmrm0
```

- Once deployed the dependent pods should declare the appropriate namespace and device name being advertised by the plugin. As demonstrated in `examples/sample-plugin-usage.yaml`, where it declares under its container's required resources.

```
resources:
            limits:
                fbdl.device.com/tpmrm: 1
```

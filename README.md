# tpm-k8s-device-plugin
custom kubernetes device plugin using Intel's library

# About
This is a sample implementation of a kubernetes device plugin using Intel's SDK library [intel-device-plugin-for-kubernetes](https://github.com/intel/intel-device-plugins-for-kubernetes).

It demonstrates how to implement a custom k8s device plugin with it and be able to mount a device unto a pod within a kubernetes cluster. For instance, the Trusted Platform Module (TPM) chip.

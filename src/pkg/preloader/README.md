# The Preloader

This package is a wrapper around
[daisy](https://github.com/GoogleCloudPlatform/compute-daisy).
It pulls in config information from the `build.config` file that is generated by the `cos_customizer`
which is indirectly populated by the user provided `cloudbuild.yaml` config.
The preloader takes this information and generates a daisy config using this template here:

`src/data/build_image.wf.json`

Once daisy is called and all resources are created, cloud-init will execute the
[provisioner](https://cos.googlesource.com/cos/tools/+/refs/heads/master/src/cmd/provisioner)
which is the next step of the process. Once cloud-init is done and the provisioner runs to
completion, control will be handed back to daisy which will clean up any unecessary resources
and exit the build.

The provisioner is packaged in an image called `cidata` built in this directory's `BUILD.bazel`
file. The image is created, artifacts are copied into the image, and the image is then embedded in
`preload.go` where it is then later uploaded to GCS to create a disk to be mounted to the preload
VM.

A scratch image is also built in the same `BUILD.bazel` file and is used for the `install-gpu` step
so that the toolchain isn't installed onto the boot disk and is instead bindmounted and executed
from temporary storage. Much like cidata, the scratch image is also embedded in the `preload.go`
file and uploaded to GCS to allow daisy to create a temporary ext4 disk for gpu installation.
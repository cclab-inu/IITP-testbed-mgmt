# The Development of Darkweb Hidden Service Identification and Real IP Trace Technology Testbed Environment Based-on Kubernetes


* Requirements

  You can install Docker and Kubernetes on any Ubuntu platform (Recommend : 20.04 (Focal Fossa)).

* Prerequisites

  - 3 Virtual machines

  - You need to disable the swap partition in advance for Kubernetes setup.

  ```text
  $ sudo vi /etc/fstab
  (comment out the line for swap)
  $ sudo reboot
  ```

## Testing Guide
### 0. To install the testbed from basic to your Ubuntu run
``` text
  $ cd IITP-testbed-mgmt/scripts
  ~/IITP-testbed-mgmt/scripts$ ./install_devs.sh
  ~/IITP-testbed-mgmt/scripts$ ./install_docker.sh
  ~/IITP-testbed-mgmt/scripts$ ./install_k8s.sh
```
* The VM reboots when ./install_devs.sh runs.
* Please reconnect after some time ...

### 1. Compiling for execution
``` text
$ cd IITP-testbed-mgmt/src
~/IITP-testbed-mgmt/src$ go build main.go
~/IITP-testbed-mgmt/src$ ./main [Usage]
```
The output is as below, and you can enter the  "Usage".
``` text
Usage:
  create-cluster
  delete-cluster
  restart-cluster
  deploy-pods
  delete-pods
  restart-pods
  pull-image
  delete-image
  print-logs
```

### 2. Test
Currently, "**create-cluster**", "**delete-cluster**", "**deploy-pods**", and "**delete-pods**" are implemented.

* create-cluster
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./main create-cluster
  ```

* delete-cluster
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./main delete-cluster
  ```

* deploy-pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./main deploy-pods
    Image Version: [Image name] [Version]
  ```
  An example of deploying 'nginx:1.22' pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./main deploy-pods
  Image Version: nginx 1.22
  1.22: Pulling from library/nginx
  e9995326b091: Pull complete 
  6cc239fad459: Pull complete 
  55bbc49cb4de: Pull complete
  ---
  Status: Downloaded newer image for nginx:1.22
  docker.io/library/nginx:1.22

  deployment.apps/nginx-1.23-deploy created
  ```

* delete-pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./main delete-pods
  Type(all/choice): [Type]
  ```
  #### Type 'all'
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./main delete-pods
  Type (all/choice) : all
  deployment.apps "nginx-1.22-deploy" deleted
  deployment.apps "nginx-1.23-deploy" deleted
  deployment.apps "wordpress-6.1.0-deploy" deleted
  ---
  ```
    An example of Type 'all'
    ```text
    $ ~/IITP-testbed-mgmt/src$ ./main delete-pods
    Type (all/choice) : all
    deployment.apps "[image name]-[version]-deploy" deleted
    ----
    ```

  #### Type 'choice'
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./main delete-pods
  Type(all/choice): choice
  [image]-[version]-deploy
  [image]-[version]-deploy
  ---

  Image Version: [Image name] [Version]
  deployment.apps [Image name]-[Version]-deploy" deleted
  ```
  An example of Type 'choice'
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./main delete-pods
  Type (all/choice) : choice
  nginx-1.22-deploy
  nginx-1.23-deploy
  wordpress-6.1.0-deploy
  ---

  Image Version : nginx 1.22
  deployment.apps "nginx-1.22-deploy" deleted 
  ```

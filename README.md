# Testbed Management System for Kubernetes

* Requirements

  You can install Docker and Kubernetes on any Ubuntu platform (Recommend : 20.04 (Focal Fossa)).

* Prerequisites

  - 3 Virtual machines

  - Disable the swap partition in advance for Kubernetes setup.
  ```text
  $ sudo vi /etc/fstab
  (comment out the line for swap)
  $ sudo reboot
  ```
  - Install dependencies
  ```text
  $ cd IITP-testbed-mgmt/scripts
  ~/IITP-testbed-mgmt/scripts$ ./install_devs.sh
  ~/IITP-testbed-mgmt/scripts$ ./install_docker.sh
  ~/IITP-testbed-mgmt/scripts$ ./install_k8s.sh
  ```
  

## Test Guide

### 1. Compile and run the executable file
``` text
$ cd IITP-testbed-mgmt/src
~/IITP-testbed-mgmt/src$ go build -o testbed
~/IITP-testbed-mgmt/src$ ./testbed
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

### 2. Check the supported container images

```text
wordpress:6.1.0
nginx:1.23
nginx:1.22
httpd:2.4
httpd:2
```

### 3. Test the functions
"**create-cluster**", "**delete-cluster**", "**deploy-pods**", and "**delete-pods**" are available now

* create-cluster
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./testbed create-cluster
  ```

* delete-cluster
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./testbed delete-cluster
  ```

* deploy-pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./testbed deploy-pods
    Image Version: [Image name] [Version]
  ```
  An example of deploying 'nginx:1.22' pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./testbed deploy-pods
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
  $ ~/IITP-testbed-mgmt/src$ ./testbed delete-pods
  Type(all/choice): [Type]
  ```
  #### Type 'all'
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./testbed delete-pods
  Type (all/choice) : all
  deployment.apps "nginx-1.22-deploy" deleted
  deployment.apps "nginx-1.23-deploy" deleted
  deployment.apps "wordpress-6.1.0-deploy" deleted
  ---
  ```
    An example of Type 'all'
    ```text
    $ ~/IITP-testbed-mgmt/src$ ./testbed delete-pods
    Type (all/choice) : all
    deployment.apps "[image name]-[version]-deploy" deleted
    ----
    ```

  #### Type 'choice'
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./testbed delete-pods
  Type(all/choice): choice
  [image]-[version]-deploy
  [image]-[version]-deploy
  ---

  Image Version: [Image name] [Version]
  deployment.apps [Image name]-[Version]-deploy" deleted
  ```
  An example of Type 'choice'
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./testbed delete-pods
  Type (all/choice) : choice
  nginx-1.22-deploy
  nginx-1.23-deploy
  wordpress-6.1.0-deploy
  ---

  Image Version : nginx 1.22
  deployment.apps "nginx-1.22-deploy" deleted 
  ```

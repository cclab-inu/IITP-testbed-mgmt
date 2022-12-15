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
  - Change the conf/local.yaml to suit your environment.
  ```yaml
  project:
  home: /home/<user name>
  user : <user name>
  
  cluster:
    master: 127.0.0.1
    worker1: 
      ip: xxx.xxx.xxx.xxx
      id: yyyyy
      pw: zzzzz
    worker2: 
      ip: xxx.xxx.xxx.xxx
      id: yyyyy
      pw: zzzzz
  ...
  ```

## Test Guide

### 1. Compile and run the executable file
``` text
$ cd IITP-testbed-mgmt/src
~/IITP-testbed-mgmt/src$ go build -o mgmt
~/IITP-testbed-mgmt/src$ ./mgmt
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
nginx:1.22
nginx:1.23
httpd:2
httpd:2.4
mongo-express:0.54
mongo-express:0.49.0

joomla
drupal
wordpress
```

### 3. Test the functions
All available now

* create-cluster
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt create-cluster
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX", "message":Createing Cluster... "}
  [████████                                     ]  15%     15/100
  ```

* delete-cluster
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt delete-cluster
   {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX", "message":Deleting Cluster... "}
  [████████                                     ]  15%     15/100

  ```

* deploy-pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt deploy-pods [Image name] [Version]
  ```
  An example of deploying 'nginx:1.22' pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt deploy-pods nginx 1.22
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX", "message":"Deploying pod: nginx"}
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX", "message":deployment.apps/nginx-1.22 created\n"}
  ```

* delete-pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt delete-pods [Type]
  ```
  #### Type [all]
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt delete-pods all
  ```
    An example of Type 'all'
    ```text
    $ ~/IITP-testbed-mgmt/src$ ./testbed delete-pods all
    {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"Deleting pod: all"}
    {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"deployment.apps \"drupal\" deleted\ndeployment.apps \"joomla\" deleted\ndeployment.apps \"nginx-1.23\" deleted\n"}
    {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"release \"drupal\" uninstalled\nrelease \"joomla\" uninstalled\n"}
    ```
  #### Type [Image] [Version]
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt delete-pods [Image name] [Version]
  ```
  An example of Type [Image] [Version]
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./testbed delete-pods nginx 1.22
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"Deleting pod: nginx"}
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"deployment.apps \"nginx-1.22\" deleted\n"}
  ```

* restart-pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt restart-pods [Image name] [Version]
  ```
  An example of restarting 'nginx:1.22' pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt deploy-pods nginx 1.22
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"Restarting pod: nginx"}
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"deployment.apps/nginx-1.22 restarted\n"} 
  ```

* pull-image
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt pull-image [Image name] [Version]
  ```
  An example of restarting 'nginx:1.22' pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt pull-image nginx 1.22
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"1.22: Pulling from library/nginx\n025c56f98b67: Pulling fs layer\n8a9d2fc4eac8: 
  ...
  Status: Downloaded newer image for nginx:1.22\ndocker.io/library/nginx:1.22\n"}
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"Successfully saved Image: nginx-1.22.tar \n"}
  ```

* delete-image
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt delete-image [Image name] [Version]
  ```
  An example of restarting 'nginx:1.22' pods
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt delete-image nginx 1.22
  {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"Successfully deleted Image: nginx-1.22 \n"}
  ```

* print-logs
  #### Type [all]
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt print-logs all
  ```
    An example of Type 'all'
    ```text
    $ ~/IITP-testbed-mgmt/src$ ./testbed print-logs all
    {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"** Network Level Log **\n"}
    ...
    {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"** System Level Log **\n"}
    ...
    {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"** Application Level Log **\n"}
   ...
    ```
  
  #### Type [network]
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt print-logs network
  ```
    An example of Type 'network'
    ```text
    $ ~/IITP-testbed-mgmt/src$ ./mgmt print-logs network
    {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"** Network Level Log **"}
    ...

    ```
  
  #### Type [system]
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt print-logs system
  ```
    An example of Type 'system'
    ```text
    $ ~/IITP-testbed-mgmt/src$ ./mgmt print-logs system
    {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"** System Level Log **"}
    ...
    ```
  
  #### Type [app]
  ```text
  $ ~/IITP-testbed-mgmt/src$ ./mgmt print-logs app
  ```
    An example of Type 'app'
    ```text
    $ ~/IITP-testbed-mgmt/src$ ./mgmt print-logs app
    {"level":"info","time":"XXXX-XX-XXXXX:XX:XXX","message":"** Application Level Log **"}
    ...
    ```
  

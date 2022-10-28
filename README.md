# The Development of Darkweb Hidden Service Identification and Real IP Trace Technology Testbed Environment Based-on Kubernetes

* Requirements

  You can install Docker and Kubernetes on any Ubuntu platform (e.g., 18.04, 20.04).

* Prerequisites

  - 3 Virtual machines

  - You need to disable the swap partition in advance for Kubernetes setup.

  ```text
  $ sudo vi /etc/fstab
  (comment out the line for swap)
  $ sudo reboot
  ```
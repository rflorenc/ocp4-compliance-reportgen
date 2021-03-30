# Report publisher for ocp4 compliance-operator

## Purpose

This tool generates per scan oscap html pages derived from compliance-operator results and presents them in a single web page under the root of the webserver, by default /opt/nginx/html.

### Running

Once compliance-operator scans have finished.

`$ cd ./hack && generate-html-results.sh`


`generate-html-results.sh` creates a Deployment which spawns a reportgen pod per scan, mounts the pvc used for the scan and exposes the results via Openshift route. 

![example report root](screenshots/worker-results1.PNG?raw=true) 

![example report](screenshots/worker-results2.PNG?raw=true)


Kudos to @kharyam !

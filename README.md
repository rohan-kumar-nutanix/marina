# content-management-marina
Marina Service

# Welcome to your new service Marina for Content Management.
We've created an empty service structure to show you the setup and workflow with Canaveral.

### Directory Structure
The top level directory of your repository should be set up like this:
  1. `README.md`: this file contains a textual description of the repository.
  2. `.circleci/`: this directory contains CircleCI's `config.yml` file.
  3. `hooks/`: this directory, if present, can contain *ad hoc* scripts that customize your build.
  4. `package/`:  add your `Dockerfile` under `package/docker/` to build a docker image.  (Note:  You can refer to files and folders directly in your `Dockerfile` because all files and folders under `services/` will be copied into the same folder as the `Dockerfile` during build.)
  5. `services/`: this directory should have a subdirectory for each `service`, *e.g.* `services/my-service/`.  Each subdirectory (often there is only one) would contain the definition (source and tests) for the service.
  6. `blueprint.json`: this file, if present, contains instructions for Canaveral to deploy the service.
  7. `common/`: this directory contains common go util functions.
  8. `config/`: this directory holds the deployment configuration files, helm chart configuration etc.
  9. `error/`: this directory holds error.go, Marina service errors and Interface Types.
  10. `grpc/`: this directory has all the new entity servies which are served by gRPC RPC framework.
  11. `marina/`: this directory holds the Main service, init, metric functions.
  12. `protos/`: this directory holds all *.proto definitions for the entities.
  13. `proxy/`: this directory all proxy RPC work. Proxy all Catalog RPC's and delegate to Catalog service.
  14. `task/`: this directory hold the base task and task realted infrastructure for Marina service.
  15. `test/`: this directory contains test scripts etc.
  16. `utils/`: this directory contains util code.
  17. `vendor/`: This directory is GoLang vendor, contains service external dependencies.

### Build
Canaveral uses CircleCI for building, packaging, and alerting its Deployment Engine. Your repository should have been registered with CircleCI when it was provisioned.  Here are some additional steps you should follow to ensure proper builds:

### Dependencies
go version go1.17 or later
protoc: libprotoc 3.6.1
protoc-gen-go v1.5.2
protoc-gen-go-grpc v1.1.0
protoc-gen-grpc-gateway v2.4.0

### Build in local box (Mac, Linux, UBVM)
Make sure above the repo is checkout, above dependencies are installed.
set the following as per your user and env in .bashrc [shell rc file]
export GIT_TOKEN = <your github token>
export CANAVERAL_ARTIFACTORY_USERNAME = <your email id>
export CANAVERAL_ARTIFACTORY_PASSWORD = <your api key>
export GOPRIVATE="github.com/nutanix-core,github.com/nutanix"
export GONOSUMDB="github.com/nutanix-core,github.com/nutanix"
export GONOPROXY="github.com/nutanix-core,github.com/nutanix"

execute all Make cmds at repo Top level dir.

#### Building protos
`make build-protos` To build protos present in protos folder.

#### Building Server
`go mod vendor` to sync the vendor folder/ update dependencies.
`make server` To build Marina server binary. marina_server binary file will be created.


#### Running Server
1. Disable IPTables in PC. (sudo service iptables stop)
2. In your local box, add entry in /etc/hosts to resolve pcip to <your PC IP>
   ex: `10.33.112.233 pcip`
3. `./marina_server` to run server from the binary file.
    or
4. `go run marina/*.go` to build code and run the marina server.

#### Running Server with proxy support
1. In PC, stop catalog service. Create catalog.gflags under /home/nutanix/config/
2. add/set gflags "--catalog_port=9202"
3. `./marina_server --proxyrpc_port=2007 --legacycatalog_port=9202` to run server from the binary file.
   or
4. `go run marina/*.go --proxyrpc_port=2007 --legacycatalog_port=9202` to build code and run the marina server.

#### Building Docker Image for the Marina service
Install Docker tools in your box. `docker --version` `Docker version 20.10.12, build e91ed57` cmd should be working.
`docker build -t marina_service:v1 -f package/docker/Dockerfile` This will build a docker image.

#### Running Marina server in PC
`GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"  -o marina_server marina/*.go`
cmd will build and generate binary with PC arch type.
Copy the `marina_server` to PC and run it in PC.
`nohup ./marina_server_task_proxy --proxyrpc_port=2007 --legacycatalog_port=9202 &`

#### [Marina Deployment in PC using MSP platform] Running Marina server in PC using MSP (CMSP should be enabled in PC)
1. build the docker image from above cmd.
2. `docker save <image id>  > marina.tar` dump the image to tar file.
3. `scp marina.tar nutanix@pcip:` copy marina.tar to PC
4. `scp config/marina-deployment.yaml nutanix@pcip:`
5. Login to PC and perform below actions.
6. `mspctl cluster kubeconfig > k.cfg && export KUBECONFIG=/home/nutanix/k.cfg`
7. `kubectl get all` to make sure kubernetes is working fine.
8. `docker load < marina.tar`  To load the docker image to PC.
9. `docker images` to fetch the Image id of the loaded image.
10. `docker tag <image id> msp-registry.catalog.cluster.local:5001/marina-cms:v1` to tag the Image with name "marina-cms:v1"
11. `docker push msp-registry.catalog.cluster.local:5001/marina-cms:v1` to push the image to repo.
12. `mspctl application apply marina -c catalog -f marina-deployment.yaml` to deploy Marina service using MSP service.
13. `kubectl get all` to check if the pod got deployed for the Marina service
14. Marina service now runs as container in PC with above process.
15. `mspctl application delete marina -c catalog -f marina-deployment.yaml` to delete the Application.

#### Accessing Marina pod.
`kubectl exec --stdin --tty  <marina-podname>   -- "/bin/sh"` login to Marina container.

#### Slack Channel for collaboration
[cms-marina] (https://nutanix.slack.com/archives/C02PWJBD5LM)  
  
  
##### Ensure `.circleci/config.yml` has the correct variables (docker image only)
  1. Specify your preferred `CANAVERAL_BUILD_SYSTEM` (default is noop)
  2. Specify your preferred `CANAVERAL_PACKAGE_TOOLS` (use "docker" if deploying a docker image, use "noop" if no packaging is needed)
  3. **[OPTIONAL]** Specify the target `DOCKERFILE_NAME` to use  (default is Dockerfile)

You'll be able to monitor the build at [circleci.canaveral-corp.us-west-2.aws](https://circleci.canaveral-corp.us-west-2.aws/)

### Deployment
To use Canaveral for deployment, `blueprint.json` should be placed at the top level of the repo.  Spec for the blueprint can be found at [Canaveral Blueprint Spec](https://confluence.eng.nutanix.com:8443/x/5kbdBQ).

__Questions, issues or suggestions? Reach us at https://nutanix.slack.com/messages/canaveral-onboarding/.__

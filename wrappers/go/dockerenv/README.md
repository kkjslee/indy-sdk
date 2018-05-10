# Indy Getting Started with Go

## Running getting-started with docker-compose

### Prerequisites

`docker` and `docker-compose` should be installed.

### Run

`docker-compose up`

The command above will create `getting-started` and `indy_pool` (collection of the validator nodes) images if they hasn't been done yet, create containers and run them.
The validators run by default on IP `10.0.0.2`, this can be changed by changing `pool_ip` in the `docker-compose` file.

Attach to the docker container:

`docker exec -ti getting_started /bin/bash`

Navigate to the `test` folder:

`cd $GOPATH/src/github.com/hyperledger/indy-sdk-go/test`

Run the Alice demo:

`go test`

### Stop

`docker-compose down`

The command above will stop and delete created network and containers.

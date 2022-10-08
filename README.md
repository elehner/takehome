# takehome server
## Installation

Install [Golang](https://go.dev/dl/)

## Running the server
### Locally
To run the server, you can run clone the repo and then run:
```go build; ./takehomeserver```
from this directory

### Through Docker
First (if you haven't already done so), install [docker](https://docs.docker.com/get-docker/)

Once it's installed, clone the repo. From within the directory, run ```docker build --tag takehome-server .```

After it's built, use ```docker run -p 8080:8080 takehome-server``` to run the container and connect the server with port 8080 on your local machine.

At this point you can now test the app manually. See more on this below.

## Testing
### Go tests
You can run `go test ./...`

### Manual testing
For ease of use, there is an "http_requests" directory holding a sample JSON and http file, which can be run using the [REST Client VSCode plugin](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).
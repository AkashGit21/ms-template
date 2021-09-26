# ms-template

## Developing locally

### Prerequisites
1. Golang v1.16+
1. go protobuf tools (to generate code from proto files)

### Steps
1. Clone the repository to your desired location.
    ```sh
    git clone github.com/AkashGit21/ms-template
    ```

1. Make changes to the proto files as per your need and run the below command to generate Go files.
    ```sh
    make gen
    ```
1. Update the server and client with your business logic
1. To run the server, use the following command: 
    ```sh
    make run
    ```

## Installation
The ms-template can be installed by building your local Docker image or simply by installing form source using Go.

### Docker
The pre-requisite for this step is the requirement of Docker installed locally. Hence, install Docker if not already installed.
Then, run the following command:
  ```sh
  make docker-run
  ```

This will first build a docker image locally. Please keep calm as the process may take a few minutes. Then, it goes on to run the image in a container (in the background) and you shoul be able to check the container named `ms-server` using  ```docker ps```. 

To stop the container, run:
  ```sh 
  make docker-stop
  ```

### Local installation 

1. Clone the repository to your desired location.
    ```sh
    git clone github.com/AkashGit21/ms-template
    ```
1. Run the following command to make an executable file named `ms-server` in the same directory.
    ```sh
    make build
    ```
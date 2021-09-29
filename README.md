# ms-template

This repository tries to include most of the boilerplate code needed to develop gRPC APIs using Go and gRPC. Some of the important features are: REST/HTTP gateway(developed usig grpc-gateway), Logging, Authentication/Authorization, Rate Limiting, etc.

## Features
1. **gRPC API** alongwith the **REST/HTTP API** built using grpc-gateway.
2. **Services** 
    
    Currently, three services are supported:
    * ##### Identity Service:
        
        This enables the user to perform actions such as: creating User, deleting the User, updating his/her info. and fetch Users by Id or whole. The operations are authorized depending on User permissions and such.
        - [X] **POST** `/v1/users` Create a new user with unique Username, email and other such details.
        - [X] **GET** `/v1/users` Lists all the users present at the given time.
        - [ ] **PUT** `/v1/users/{username}` Update the User with specified *username*.
        - [X] **DELETE** `/v1/users/{username}` Deletes the User with specified *username*.
        - [X] **GET** `/v1/users/{username}` Fetch the details of User with specified *username*.


    * ##### Auth Service:
        
        A simple authentication service for User verification.
        - [X] **POST** `/v1/auth/login` Takes user input (Username & Password) to generate a token.

    * ##### Movie Service:
        
        This enables the user to perform following activities: inserting a Movie, deleting Movie, updating its info. and/or fetching Movies by Id or as a whole. The operations are authorized depending on User permissions and such.
        - [X] **POST** `/v1/movies` Adds a new Movie with unique Name, Tags, Summary and other such details. This returns the generated id for the inserted Movie.
        - [X] **GET** `/v1/movies` Lists all the Movies present at the given time.
        - [X] **PUT** `/v1/movies/{id}` Update an already present Movie with new values.
        - [X] **DELETE** `/v1/movies/{id}` Delete an existing Movie with given ID. 
        - [X] **GET** `/v1/movies/{id}` Fetches the movie with given ID.
1. **Unit Tests**

    *IN PROGRESS*
1. **Documentation**

    Visit [here](https://akashgit21.github.io/ms-template/docs) to view the documentation of Services. Search for the service you like to explore (e.g: IdentityService-0.1.yaml, AuthService-0.1.yaml, or MovieService-0.1.yaml )

    Documentation will be updated as per the modifications and changes in business logic. Also, the documentation is auto-generated. Hence, less focus may be there in some scenarios.

1. **Rate Limiting**

    The middleware to rate limit requests has been applied. By default, it is 2requests/min. To modify the limit, update `RefreshDuration` and `QueriesPerInterval` in this [file](./internal/server/interceptors/rate_limit.go).

1. **Pagination** 
    
    The concept of pagination has been supported for List calls such as ```GET v1/movies```.     
    **Query Parameters** such as *page_size* and *page_token* are supported for List calls. The response returns a *next_page_token* as well.


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
The pre-requisite for this process is the requirement of Docker (to be installed locally). Hence, install Docker if not already installed.

#### Steps

1. Clone the repository to your desired location.
    ```sh
    git clone github.com/AkashGit21/ms-template
    ```

1. Now, lets build the Docker image locally using:
    ```sh
    make docker-build
    ```

1. Then, to run the container use:
    ```sh
    make docker-run
    ```
Now, you should be able to access the REST API at `localhost:8081` and gRPC API at `localhost:8082`.

This will first build a docker image locally. Please keep calm as the process may take a few minutes. Then, it goes on to run the image in a container (in the background) and you shoul be able to check the container named `ms-server` using  ```docker ps```. 

To stop the container, use:
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

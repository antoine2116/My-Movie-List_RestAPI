# My Movie List - REST API

[![codecov](https://codecov.io/gh/antoine2116/My-Movie-List_RestAPI/branch/test-codecov/graph/badge.svg?token=BVPYGYDIY7)](https://codecov.io/gh/antoine2116/My-Movie-List_RestAPI)

This **Golang/Gin** REST API provides authentication using Basic and OAuth 2 with Google and GitHub, along with a watchlist management system. Authentication is done using JWT tokens. The application is built on top of a MongoDB database and follows clean architecture principles. 

Built for the [My Movie List App](https://github.com/antoine2116/My-Movie-List_App) project. A Next.js web application that allows users to search for movies and add them to their watchlist.

Testing is done using the standard `testing` package and the `testify` library. Code coverage is tracked using `codecov` via GitHub Actions.



## Getting Started

### Prequisite

- Go (1.3 or higher)
- MongoDB (3.6 or higher)
- Docker (for optional MongoDB containerization)


### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/antoine2116/My-Movie-List_RestAPI
    ```

2. Install dependencies
    ```sh
    go mod download
    ```

3. Create the `config/local.yml` following the template below : 
    ```yml
    server:
    port: 8081
    secret: "<your-secret>"
    token_duration: 24
    client:
      uri: "<my-movie-list-app-url>"
    database:
      uri: "mongodb://localhost:27017/"
      db: "<my-movie-list-db-name>"
    google_oauth:
      client_id: <your-client-id>
      client_secret: <your-client-secret>
      redirect_url: "http://localhost:8081/auth/google/callback"
    github_oauth:
      client_id: <your-client-id>
      client_secret: <your-client-secret>
      redirect_url: "http://localhost:8081/auth/github/callback"
    ```

4. Run the application
    ```sh
    go run cmd/main.go
    ```

### Dockerized MongoDB
To run a Docker container with MongoDB, execute the following command :
```sh
docker run -d -p 27017:27017 --name my-movie-list-mongodb mongo
```
## Endpoints

### User Registration

- **Endpoint**: `/auth/register`
- **Method**: POST
- **Description**: Register a new user with provided email and password.
- **Request Body**:
  ```json
  {
    "email": "somebody@mail.com",
    "password": "strongpassword"
  }
  ```
- **Responses**:
  - `201 Created` : Registration successful. Returns a JSON object with an authentication token.
  - `400 Bad Request` : Invalid input or validation errors.
  - `409 Conflict` : User with the provided email already exists.

### User Login

- **Endpoint**: `/auth/login`
- **Method**: POST
- **Description**: Authenticate a user using their email and password.
- **Request Body**:
  ```json
  {
    "email": "somebody@mail.com",
    "password": "strongpassword"
  }
  ```
- **Responses**:
  - `200 OK` : Login successful. Returns a JSON object with an authentication token.
  - `400 Bad Request` : Invalid input or validation errors.
  - `401 Unauthorized` : Incorrect email or password.

### Google OAuth2 Login
- **Endpoint**: `/auth/google/callback`
- **Method**: GET
- **Description**: Handle Google OAuth2 authentication callback.
- **Responses**:
  - `302 Found` : Redirects to the provided client URI with a token cookie.
  - `400 Bad Request` : Invalid token or missing query parameter.
  - `401 Unauthorized` : Cannot retrieve user information.

### GitHub OAuth2 Login
- **Endpoint**: `/auth/github/callback`
- **Method**: GET
- **Description**: Handle GitHub OAuth2 authentication callback.
- **Responses**:
  - `302 Found` : Redirects to the provided client URI with a token cookie.
  - `400 Bad Request` : Invalid token or missing query parameter.
  - `401 Unauthorized` : Cannot retrieve user information.

### User Profile
- **Endpoint**: `/api/profile`
- **Method**: GET
- **Description**: Retrieve the profile of the currently authenticated user.
- **Responses**:
  - `200 OK` : Returns the user's profile information.
  - `401 Unauthorized` : User is not authenticated.

## Todo
- Watchlist management system (in progress)
- CI/CD pipeline (in progress)
- Dockerfile
- Swagger documentation
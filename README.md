# BookQuest

This is a Go project that leverages HTMX for dynamic frontend interactions and PostgreSQL for the database.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites
Make sure you have the following installed on your system:
- [Go](https://golang.org/doc/install) 1.16+
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

Use Docker Compose file to build and run the application.

```sh
docker compose up -d
```

This command will build the Docker image for your application, start the PostgreSQL service, and then start your application service.

### Access the Application

Once the services are up and running, you can access the application at http://localhost:8090.

### Stopping the Application

To stop the running application and PostgreSQL services, use the following command:

```sh
docker compose down
```


## Configuration

Create a `.env` file in the root directory and add your database configuration:
```env
DATABASE=postgresql://myuser:mypassword@localhost:5432/mydatabase
OIDC_SERVER=my-domain.zitadel.cloud
OIDC_CLIENT_ID=<EXMPLE_CLIENT_ID>
SECRET__KEY=<EXAMPLE_SECRET_KEY> 
OIDC_REDIRECT_URI= http://localhost:8089/auth/callback
LISTEN_PORT=8090
```

## Setting Up OAuth

### Zitadel

1. Create a new project in Zitadel and set up an application within the project.
2. Configure the redirect URIs to point to your Go application endpoints, e.g., `http://localhost:8080/auth/callback`.
3. Note down the client ID and client secret provided by Zitadel.
4. Add the following configuration to your `.env` file:
    ```env
    ZITADEL_CLIENT_ID=your_zitadel_client_id
    ZITADEL_CLIENT_SECRET=your_zitadel_client_secret
    ZITADEL_REDIRECT_URI=http://localhost:8080/auth/callback
    ZITADEL_ISSUER=https://issuer.zitadel.ch
    ```

### Keycloak

1. Create a new realm in Keycloak.
2. Within the realm, create a new client and configure the redirect URIs to point to your Go application endpoints, e.g., `http://localhost:8080/auth/callback`.
3. Note down the client ID and client secret provided by Keycloak.
4. Add the following configuration to your `.env` file:
    ```env
    KEYCLOAK_CLIENT_ID=your_keycloak_client_id
    KEYCLOAK_CLIENT_SECRET=your_keycloak_client_secret
    KEYCLOAK_REDIRECT_URI=http://localhost:8080/auth/callback
    KEYCLOAK_ISSUER=http://localhost:8080/auth/realms/your_realm
    ```

## Running the Application
The server will run all database migrations on startup. If you wish to manually deploy the migrations see the Database section

```sh
bookquest server
```

## Database Management

```sh
bookquest db --help
```

### Migrate database

```sh
bookquest db migrate
```  

### Rollback the last migration group      

```sh
bookquest db rollback
```       

### Lock migrations

```sh
bookquest db lock
```    

### Unlock migrations       

```sh
bookquest db unlock
```      

### Migrations status   

```sh
bookquest db status
```         

### Mark migrations as applied without actually running them

```sh
bookquest db mark_applied
```   

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed Go 1.22. You can download it from [Go's official website](https://golang.org/dl/).
- You have a working environment for Go development. If not, follow the setup instructions provided by Go's documentation.

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/robrotheram/BookQuest.git
    cd BookQuest
    ```

2. Download the project dependencies:

    ```bash
    go mod tidy
    ```

## Building the Project

To build the project, follow these steps:

1. Navigate to the project directory (if not already there):

    ```bash
    cd BookQuest
    ```

2. Build the project:

    ```bash
    go build -o bookquest ./cmd/main.go
    ```

    This command will compile the Go code and create an executable named `bookquest` in the current directory.

## Running the Project

Once the project is built, you can run the executable:

```bash
./bookquest
```

## Contributing
We welcome contributions! Please see CONTRIBUTING.md for details.

## License
This project is licensed under the MIT License - see the LICENSE file for details.
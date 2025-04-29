# Swift Code Database

This project provides a RESTful API for managing a database of banks and their SWIFT codes. It allows users to retrieve, add, and delete bank information. The application is built using **Go** and utilizes **PostgreSQL** as its database. The entire application is containerized using Docker for easy setup and deployment.

## Documentation Sections

## Setup

Start by cloning the repository to your local computer

```sh
git clone https://github.com/Xenoneqq/swift-code-database
```

To set up the project you will require Docker installed. You can do so by (installation links here)

## Launching the Project Using Docker

The application can be launched in different modes depending on the Docker command used. The commands and modes are as follows:

### Default (Production-like)

This mode launches the app with a pre-determined database populated from the provided CSV file ``bank_data.csv`` and allows for data retrieval from requests to ``https://localhost:8080/api/v1/swift-codes``.

```sh
docker-compose up --build -d
```

### Debug

Launches the app in debug mode with a separate test database. No data is inserted into the database. The app also allows for HTTP requests and provides a safe environment to test edge cases or specific features.

```sh
docker-compose --env-file .env.debug up --build -d
```

### Testing

Launches the app in test mode and prepares it for automatic testing. The tests will run after the app and database (the same one used for debug mode) are up and running. Once the tests are complete, the console will display the test results. Afterward, the database will remain open for requests and will function similarly to the debug mode.

```sh
docker-compose --env-file .env.test up --build
```
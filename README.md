# Swift Code Database

This project provides a RESTful API for managing a database of banks and their SWIFT codes. It allows users to retrieve, add, and delete bank information. The application is built using **Go** and utilizes **PostgreSQL** as its database. The entire application is containerized using Docker for easy setup and deployment.

## Documentation Sections

- [Project Infrastructure](#project-infrastructure)
- [Features](#features)
- [Getting Started](#getting-started)
- [Running with Docker](#running-with-docker)
- [API Reference](#api-reference)
- [Testing](#testing)

## Project Infrastructure

The application is composed of several cooperating systems, each responsible for a different layer of functionality:

#### Main Database

This component stores and manages all project data. It supports adding, removing, and retrieving entries. The database is preloaded with mock data provided in the `/data` directory as CSV files.

#### Swift App

This is the main server application responsible for handling all incoming requests. It processes the backend logic, validates the data, and interacts with the database to ensure accurate responses or appropriate error handling.

The app uses verified ISO 3166-1 alpha-2 country codes and official country names obtained from [Restcountries v2.0](https://restcountries.com). This helps validate incoming data, whether from CSV files or API `CREATE` requests.

#### Test Database

A separate database used solely for testing purposes. It is isolated from production data and used exclusively to validate edge cases. This database is the target environment for automated tests.

#### Test App

This module is responsible for running automated tests against backend endpoints and logic. It identifies functional issues and provides detailed error reports, aiding in the validation of existing features.

## Features

- Written in Go with Fiber web framework
- PostgreSQL for persistent storage
- Automated testing with separate database
- Dockerized (dev/test/prod)
- Input validation based on official ISO and SWIFT formats

## Getting Started

To get started with this project, first clone the repository to your local machine:

```sh
git clone https://github.com/Xenoneqq/swift-code-database
```

Once cloning is complete, navigate into the project directory:

```sh
cd swift-code-database
```

### Prerequisites

Before setting up the project, ensure you have [Docker](https://www.docker.com/get-started) installed on your system. You can follow the installation instructions for your platform here:

- [Install Docker on Windows](https://docs.docker.com/desktop/install/windows-install/)
- [Install Docker on macOS](https://docs.docker.com/desktop/install/mac-install/)
- [Install Docker on Linux](https://docs.docker.com/engine/install/)

Once Docker is installed, you're ready to proceed with the setup steps.

## Running with Docker

The application supports multiple Docker run modes, each suited for a different environment or use case. Below are the available modes and their respective commands:

### 1. Default (Production-like)

This mode runs the application with a default database populated from the `bank_data.csv` file located in the project. It exposes an API endpoint for data access at:

```
http://localhost:8080/api/v1/swift-codes
```

To start the application in this mode:

```sh
docker-compose up --build -d
```

### 2. Debug Mode

This mode is intended for local debugging and feature testing. It runs the application with a separate test database, without any preloaded data. The environment accepts HTTP requests and is suitable for testing edge cases or development-specific scenarios.

To start in debug mode:

```sh
docker-compose --env-file .env.debug up --build -d
```

### 3. Testing Mode

This mode prepares the environment for automated testing. It starts the application and test database, then automatically runs the test suite. After the tests complete, the environment remains active, allowing further manual inspection or requests.

To run in test mode:

```sh
docker-compose --env-file .env.test up --build
```

## API Reference

All endpoints are served under the `https://localhost:8080/api/v1/swift-codes` base path unless stated otherwise.

You can interact with the API using tools like **Postman** *(recommended for easier testing)*, or by using command-line tools such as **curl**.

> Example request snippets below use curl for demonstration purposes and can be copied and executed directly in your terminal.

**There is a Postman collection available!**
Click the button below to quickly import the collection into Postman and start interacting with the API.

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://god.gw.postman.com/run-collection/40303085-0029cdfc-531e-4a77-b6d1-3f9081b44d81?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D40303085-0029cdfc-531e-4a77-b6d1-3f9081b44d81%26entityType%3Dcollection%26workspaceId%3D7502ffdb-055f-4b1d-9cbf-a0607f523d1e)

---

### `GET /api/v1/swift-codes`

Returns a list of all banks, sorted by bank name, headquarter status, country name, and SWIFT code.

#### Example use
```sh
curl http://localhost:8080/api/v1/swift-codes
```

#### Response
- `200 OK`: List of bank entries
- `400 Bad Request`: Failed to fetch data

> If there are no banks in the database, an empty array and a custom message will be included in the response.

---

### `GET /api/v1/swift-codes/:id`

Fetches a single bank using its SWIFT code.

#### Example use

1. When the bank is a **Headquarter**:
```sh
curl http://localhost:8080/api/v1/swift-codes/KCCPPLPWXXX
```
- **Explanation**:  The `XXX` at the end of the SWIFT code indicates a headquarter.

2. When the bank is a **Branch**:
```sh
curl http://localhost:8080/api/v1/swift-codes/KCCPPLPWASI
```
- **Explanation**: The `ASI` suffix in the SWIFT code indicates a branch of the bank.

3. When the bank might **not exist**:
```sh
curl http://localhost:8080/api/v1/swift-codes/GTBKPLWAAAA
```
- **Explanation**: If the SWIFT code doesn't exist, the response will return a `404 Not Found` status, indicating that no bank was found.

#### Response
- `200 OK`: Bank found; returns detailed info
- `404 Not Found`: Bank does not exist
- `400 Bad Request`: Multiple banks found with the same code or invalid query

> If the bank is a headquarter, branch information will be included in the response.

---

### `GET /api/v1/swift-codes/country/:country`

Retrieves all banks for a given country, based on its ISO2 code (e.g., "PL", "DE", "FR").

#### Example use
```sh
curl http://localhost:8080/api/v1/swift-codes/country/PL
```

#### Response
- `200 OK`: List of banks for the specified country
- `404 Not Found`: No banks found
- `400 Bad Request`: Failed to fetch data

---

### `POST /api/v1/swift-codes`

Creates a new bank entry using a JSON payload. The submitted data must meet the following validation rules:

- `swiftCode`, `countryName`, and `countryISO2` must be written entirely in uppercase letters.
- `swiftCode` must be exactly **11 characters** long.
- The first **six characters** of the `swiftCode` must contain only letters (no digits).
- If the bank is a headquarter (`isHeadquarter: true`), the `swiftCode` must end with `"XXX"`.
- The `countryISO2` must match the **5th and 6th characters** of the `swiftCode`.
- The `countryName` and `countryISO2` must correspond to real, valid country data.


#### Payload
```json
{
  "swiftCode": "ABCDCNSPXXX",
  "bankName": "Bank Name",
  "countryName": "COUNTRY NAME",
  "countryISO2": "XX",
  "isHeadquarter": true,
  "address": "Street 123"
}
```

#### Example use
```sh
curl -X POST http://localhost:8080/api/v1/swift-codes -H "Content-Type: application/json" -d "{\"swiftCode\":\"GTBKPLWAXXX\",\"bankName\":\"Generic Test Bank\",\"countryName\":\"POLAND\",\"countryISO2\":\"PL\",\"isHeadquarter\":true,\"address\":\"123 Bank St\"}"
```

#### Response
- `201 Created`: Entry created successfully
- `400 Bad Request`: Invalid input or bank already exists
- `422 Unprocessable Entity`: Malformed payload

---

### `DELETE /api/v1/swift-codes/:id`

Deletes a bank based on its SWIFT code.

#### Example use
```sh
curl -X DELETE http://localhost:8080/api/v1/swift-codes/GTBKPLWAXXX
```

*This example deletes the bank created in the previous section using the POST endpoint.*


#### Response
- `200 OK`: Entry deleted successfully
- `404 Not Found`: Bank not found
- `400 Bad Request`: Error during deletion

---

### `GET /api`

Returns the status of the API server.

#### Example use
```sh
curl http://localhost:8080/api
```

#### Response
- `200 OK`: Server is active

## Testing

This project includes an automated test suite that verifies the behavior of the API and core application logic. All tests are isolated from production data and run against a dedicated test database.

### What’s Covered

- API endpoint behavior (GET, POST, DELETE)
- SWIFT code validation rules
- Country name/code verification
- Handling of duplicate or malformed data

### Running the Tests

You can run the full test suite using Docker in testing mode:

```sh
docker-compose --env-file .env.test up --build
```

This will:

- Spin up the application and a dedicated test database
- Run all Go test files automatically
- Output detailed test results to the console

### Test Directory Structure

All test files are located in the `/tests` directory and use Go's built-in [`testing`](https://pkg.go.dev/testing) package. The tests are designed to run automatically during container startup when in test mode.

> ⚠️ After the tests finish, the application and test database will remain active — similar to launching in DEBUG mode.

# API Documentation

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [License](#license)

## Introduction

This project is a Go application that manages customers, services, and transactions. It uses the Gin web framework and PostgreSQL as the database.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go installed on your machine
- PostgreSQL database server
- Git (optional)

## Installation

1. Clone the repository:

   ```bash
   https://github.com/Damarwendha/challenge-goapi.git
   ```

2. Change into the project directory:

   ```bash
   cd challenge-goapi
   ```

3. Install dependencies:

   ```bash
   go get -u github.com/gin-gonic/gin
   go get -u github.com/lib/pq
   ```

4. Set up the database:

   - Create a PostgreSQL database.
   - Update the database connection details in the `config/database.go` file.

## Usage

To run the application, use the following command:

```bash
go run .
```

The application will start, and you can access it at [http://localhost:8080](http://localhost:8080).

## API Endpoints

- **Create Customer:**
  ```
  POST /customers
  ```
  Body Request Example (raw):
  ```
  {
    "id": 1234, // Should be a unique id
    "name": "John Doe",
    "phone": "085856209376",
  }
    ```

- **Update Customer:**
  ```
  PUT /customers/:id
  ```
    Body Request Example (raw):
  ```
  {
    "name": "John Doe",
    "phone": "085856209376",
  }
    ```

- **Delete Customer:**
  ```
  DELETE /customers/:id
  ```

- **Create Service:**
  ```
  POST /services
  ```
    Body Request Example (raw):
  ```
  {
    "id": 1234, // Should be a unique id
    "name": "John Doe",
    "price": 10000,
    "unit_type_id": 1 // either 1 or 2
  }
    ```

- **List Services:**
  ```
  GET /services
  ```

- **Update Service:**
  ```
  PUT /services/:id
  ```
  
     Body Request Example (raw):
  ```
  // EMPTY ATTRIBUTE NOT HANDLED YET, SO WHEN ONE OR MORE ATTRIBUTE IS EMPTY IT WILL STILL UPDATE THE ATTRIBUTE TO ZERO VALUE
  {
    "name": "John Doe",
    "price": 10000,
    "unit_type_id": 1 // either 1 or 2
  }
    ```

- **Delete Service:**
  ```
  DELETE /services/:id
  ```

- **Create Transaction:**
  ```
  POST /transactions
  ```
     Body Request Example (raw):
  ```
  {
    "id": 1234, // Should be a unique id
    "customer_id": 1, // Existing customer id
    "service_id": 1, // Existing service id
    "quantity": 10000,
  }
    ```

- **List Transactions:**
  ```
  GET /transactions
  ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```

Replace the placeholder URLs, usernames, and other details with your actual project information. This template provides a basic structure, and you can expand or modify it based on your project's specific requirements.

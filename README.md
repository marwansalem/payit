# Table of Contents

- [Table of Contents](#table-of-contents)
- [PayIt API Documentation](#payit-api-documentation)
  - [Endpoints](#endpoints)
    - [Accounts](#accounts)
      - [Get Account by ID](#get-account-by-id)
      - [Get All Accounts](#get-all-accounts)
      - [Create Account](#create-account)
      - [Update Account by ID](#update-account-by-id)
    - [Transfers](#transfers)
      - [Get Transfer by ID](#get-transfer-by-id)
      - [Get All Transfers](#get-all-transfers)
      - [Create Transfer](#create-transfer)
- [Unit Tests](#unit-tests)
- [In-Memory Data](#in-memory-data)
- [Concurrency](#concurrency)
# PayIt API Documentation

This document provides information about the endpoints available in the PayIt API, as well as the data models/entities used in the API.

## Endpoints

### Accounts

#### Get Account by ID

- **HTTP Method:** GET
- **Endpoint:** `/accounts/:id`
- **Description:** Retrieve an account by its ID.

  **Response (JSON):**

  ```json
  {
    "id": "812f2e44-6b34-4e26-a418-9affc416c032",
    "name": "John Doe",
    "balance": 1000.0
  }
  ```
    - **Note:** The `balance` field in the request body is of type string, but in the response, it will be of type float. (Why? Example data is supplied with Balance as a string, but to facilitate transfers it was changed to float)

#### Get All Accounts

- **HTTP Method:** GET
- **Endpoint:** `/accounts`
- **Description:** Retrieve a list of all accounts.

#### Create Account

- **HTTP Method:** POST
- **Endpoint:** `/accounts`
- **Description:** Create a new account.

  **Request Body (JSON):**
  
  ```json
  {
    "name": "John Doe",
    "balance": "1000.00"
  }
  ```

#### Update Account by ID

- **HTTP Method:** PUT
- **Endpoint:** `/accounts/:id`
- **Description:** Update an account by its ID.

  **Request Body (JSON):**
  
  ```json
  {
    "name": "Updated Name",
    "balance": "1500.00"
  }
  ```



  - **Note:** The `id` field in the request URL specifies the account to update.

### Transfers

#### Get Transfer by ID

- **HTTP Method:** GET
- **Endpoint:** `/transfers/:id`
- **Description:** Retrieve a transfer by its ID.

  **Response (JSON):**

  ```json
  {
    "id": "33f47cae-580f-11ee-8c99-0242ac120002",
    "senderID": "7d4b886d-dd4d-4133-9bc7-26f0c6233e8c",
    "receiverID": "53f30a2e-f9c8-4442-a278-b622a42e4054",
    "amount": 500.0,
    "succeeded": true
  }
  ```

#### Get All Transfers

- **HTTP Method:** GET
- **Endpoint:** `/transfers`
- **Description:** Retrieve a list of all transfers.

#### Create Transfer

- **HTTP Method:** POST
- **Endpoint:** `/transfers`
- **Description:** Make a new transfer.

  **Request Body (JSON):**
  
  ```json
  {
    "id": "33f47cae-580f-11ee-8c99-0242ac120002",
    "senderID": "7d4b886d-dd4d-4133-9bc7-26f0c6233e8c",
    "receiverID": "53f30a2e-f9c8-4442-a278-b622a42e4054",
    "amount": 500.0,
  }

## API Test
  * A collection of requests is provided.
  * For [Postman](./PayIT_postman.json)
  * For [ThunderClient](./PayIT_thunderclient.json)

# Running PayIt

To run the PayIt program, follow these steps:

1. Using source code

   ```bash
   go run main.go

2. Build and run executable
   ```bash
   go build -o payit
   ./payit
   ```
3. Server starts on port `8080`
  - **Note:** The application loads `accounts.json` data into the memory, if successful this message will be logged: `Ready to receive requests`.
  - **Note:** Optionally replace `accounts.json` with your own `accounts.json` file, but it should follow the same format, before running.
# Unit Tests

To run unit tests run:

   ```bash
   go test ./...
   ```
# In-Memory Data

Go's `sync.Map` is used as the main structure for holding Accounts & Transfers entities.

There is a separate for each kind of entity.

`sync.Map` facilitates concurrent read and write access to the data .

# Concurrency
The transfer must carefully modify the sender and receiver accounts. 

This is done by introducing a `mutex.Lock` for each account. This means that the locking granularity is on the **row** level.

Transfers include locking the sender's account, inspecting their balance, and then deducting the transfer amount, and finally and unlocking the sender's account.

As for the receiver, their account is locked until their balance is updated, to ensure that there are no write after write issues, and finally ublocking the receiver's account.

As for fetching and creating accounts and transfers, the concurrency is handled by the `sync.Map` 
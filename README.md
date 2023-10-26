# Go Payment App

### Description

This is a simple transfer app from customer to merchant.

### API Documentation

#### Create User

Request :

- Method : `POST`
- Endpoint : `/users`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "email": "albaras@gmail.com",
  "username": "albar",
  "password": "secret",
  "role": "user"
}
```

#### Login User

Request :

- Method : `POST`
- Endpoint : `/users/login`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

  ```json
  {
    "username": "albar",
    "password": "secret"
  }
  ```

#### List User

  <!-- Only user with role admin can access this route -->

Request :

- Method : `GET`
- Endpoint : `/users`
- Header :

  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token

#### Update User

Only user with role admin can access this route

Request :

- Method : `PUT`
- Endpoint : `/users`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token
- Body :

```json
{
  "id": "sdifjij39214",
  "email": "albaras@gmail.com",
  "username": "albar",
  "password": "secret",
  "role": "user"
}
```

#### Create Customer

Request :

- Method : `POST`
- Endpoint : `/customers`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token
- Body :

  ```json
  {
    "name": "Albar Adimas Suntoro"
  }
  ```

#### List Customer

Only user with role admin can access this route

Request :

- Method : `GET`
- Endpoint : `/customers`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token

#### Get Customer By Name

Request :

- Method : `GET`
- Endpoint : `/customers/:name`
- Header :

  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token

#### Delete Customer

Request :

- Method : `DELETE`
- Endpoint : `/customers/:id`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token

#### TOP-UP

Request :

- Method : `GET`
- Endpoint : `/transaction/top-up`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token
- Body :

```
{
    "amount" : 10000
}

```

#### Create Merchant

Request :

Only user with role admin can access this route

- Method : `POST`
- Endpoint : `/merchants`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token
- Body :

  ```json
  {
    "name": "XLEND",
    "description": "Perusahan Retail",
    "busines_type": "retail"
  }
  ```

#### Delete Customer

Only user with role admin can access this route

Request :

- Method : `DELETE`
- Endpoint : `/merchants/:id`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token

#### Get Merchant By Id

Request :

- Method : `GET`
- Endpoint : `/customers/:id`
- Header :

  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token

#### List Merchant

Request :

- Method : `GET`
- Endpoint : `/merchants`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token

#### Create Transaction

Request :

- Method : `POST`
- Endpoint : `/transactions`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token
- Body :

  ```json
  {
    "receiver_merchant_id": "659092c2-da66-42bf-b61c-0464dabb9a2e",
    "amount": 100
  }
  ```

#### List Transaction

Request :

- Method : `GET`
- Endpoint : `/transactions`
- Header :
  - Content-Type : application/json
  - Accept : application/json
  - Authorization : Bearer token

### How to run

- Clone this repository
- Copy `.env.example` to `.env`
- Change the value of `.env` to your configuration
- Download packages with `go mod tidy`
- Create database db_payments
- Copy `config/database/init.sql` to your database
- Run the app with `go run main.go`

### How to run with docker

- Clone this repository
- Copy `.env.example` to `.env`
- Change the value of `.env` to your configuration
- Run `docker-compose up -d`
- Enter the db container with `docker exec -it payments-app-db-1 sh`
- login to postgres with `psql -U postgres`
- Create database db_payments with `CREATE DATABASE db_payments;`
- Copy `config/database/init.sql` to your database

### How to deploy

- Clone this repository to your server
- Do the same as how to run with docker
- Setup nginx to your server
- Setup ssl to your server
- Setup reverse proxy to your server

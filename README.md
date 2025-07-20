# GO GET JOB

A backend project based on Golang and PostgreSQL that showcase CRUD operations using GraphQL. The application uses docker for containerization.

## 👨‍🚀 Postman Collection Documentation

[Postman Documentation URL](https://documenter.getpostman.com/view/40689865/2sB34kEKCS)

## ✨Features

- **Dependency Injection** for modularity
- **Pagination** to reduce load
- **Error handling** for debugging
- **Persistant storage** using PostgreSQL
- **Caching** to improve performance

## 📌 Tech Stack

- **Backend:** Golang, GraphQL
- **Database:** PostgreSQL
- **Cache:** Redis
- **Deployment:** Github (Code Repository)
- **Containerization:** Docker

![Project Screenshot](./public/images/Screenshot%202025-07-14%20130748.png)

## 🚀 Quick Start

### 1. Prerequisites

Make sure you have the `Docker Desktop` installed on your system:

### 2. Clone the Repository

```bash
git clone https://github.com/harshitrajsinha/go-get-job.git
cd go-get-job
```

### 3. Set up environment variable

Create .env file in root directory

```bash
POSTGRES_USER=postgres
POSTGRES_PASSWORD=yourstrongpassword
POSTGRES_DB=yourfavouritedbname

DB_USER=postgres
DB_NAME=yourfavouritedbname
DB_PASS=yourstrongpassword
DB_PORT=5432
DB_HOST=db
APP_PORT=8000
```

### 4. Run the application

Run the following command in your bash terminal

```bash
docker-compose up --build
```

### Caveat

In case you run into any problems while running docker-compose such as password authentication failure or role 'user' does not exists. Run the following commands

```bash
docker-compose down # remove existing container
docker volume rm go-get-job_db_data # remove stored data from <folder>_db_data
docker-compose up --build # rebuild the container with update
```

![Project Screenshot](./public/images/Screenshot%202025-07-14%20130601.png)

![Project Screenshot](./public/images/Screenshot%202025-07-14%20130633.png)

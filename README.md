# Fullstack Developer Test - Starter Project (Complete)

This repository contains a starter implementation for the Fullstack Developer Test Challenge:

- **product-service** (NestJS)
- **order-service** (Go + Fiber)
- RabbitMQ, Redis, PostgreSQL via Docker Compose

---

## 📦 Requirements

- Docker & Docker Compose installed
- Node.js & npm (for local NestJS development, optional)
- Go (for order-service, optional)

---

## 🚀 Running the Stack Locally

From the project root:

```bash
docker-compose up --build
```

## **Services**

| Service         | URL / Connection Details                  |
|-----------------|------------------------------------------|
| RabbitMQ UI     | http://localhost:15672 (guest/guest)     |
| Redis           | redis://localhost:6379                    |
| PostgreSQL      | localhost:5432 (user/pass)               |
| Order Service   | http://localhost:4000                     |
| Product Service | http://localhost:3000                     |


### API Collection

You can import the Postman collection for this project [here](./MrScrapper.postman_collection.json).

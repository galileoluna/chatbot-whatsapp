# WhatsApp Bot Project

This project is a WhatsApp bot built using the Go programming language and the WhatsMeow library. It uses PostgreSQL as its database, and it's containerized using Docker.

## Prerequisites

- Go (version 1.16 or later)
- Docker and Docker Compose
- PostgreSQL (version 13 or later)

## How to Run

1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Build the Docker images using Docker Compose with the following command:

```bash
docker-compose up --build
```

This command will start the PostgreSQL database and the application.

## How It Works
The bot connects to WhatsApp using the WhatsMeow library. If it's the first time running, it will generate a QR code that you need to scan using your phone to authenticate.
Once authenticated, the bot listens for incoming messages.
When a message is received, it sends a reply back to the sender.
The bot uses a PostgreSQL database to store device information.
The database connection details are configured in the docker-compose.yml file.

## Stopping the Application

To stop the application, you can use the following command:

```bash
docker-compose down
```
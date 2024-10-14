# SwiftChat Backend

[![Go 1.23x](https://img.shields.io/badge/Go-1.23.x-blue.svg)](https://go.dev/) [![Go Build](https://github.com/kwa0x2/SwiftChat-Backend/actions/workflows/go.yml/badge.svg)](https://github.com/kwa0x2//SwiftChat-Backend/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kwa0x2/swiftchat-backend?style=flat-square)](https://goreportcard.com/report/github.com/kwa0x2/swiftchat-backend) [![Website](https://img.shields.io/badge/Website-chat.nettasec.com-red.svg)](https://chat.nettasec.com/)

SwiftChat backend is built with the Go Gin framework. It uses PostgreSQL for database management, JWT for authentication, Redis for sessions, and S3 for profile pictures. Real-time communication is handled by Socket.IO, and the entire application is containerized with Docker.

## Installation:

**1. Build and Run the API with Docker Compose:**

```bash
docker-compose up -d --build
```

## Diagram:

![diagram](https://i.hizliresim.com/pnlzrcu.png)

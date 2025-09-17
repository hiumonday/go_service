# Go Service API

A robust REST API service built with Go, Gin framework, PostgreSQL, Redis, and Kafka for team and asset management.

## 🏃‍♂️ Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Git

### Running with Docker Compose

```bash
# Clone the repository
git clone https://github.com/hiumonday/go_service.git
cd seta-go_service/go_service

# Start all services
docker compose up --build

# The API will be available at http://localhost:8080
```

## 🚀 Features

- **Team Management**: Create teams, manage members and managers
- **Asset Management**: Folders and notes with sharing capabilities
- **User Import**: Bulk user import functionality
- **Authentication**: JWT-based authentication middleware
- **Caching**: Redis integration for improved performance
- **Message Queue**: Kafka integration for async processing
- **Containerized**: Docker support for easy deployment

## 🛠 Tech Stack

- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Message Queue**: Apache Kafka
- **Authentication**: JWT
- **Containerization**: Docker & Docker Compose

## 📚 API Documentation

### Base URL

- **Local**: `http://localhost:8080/api/v1`
- **Production**: `https://go-service-app.onrender.com/api/v1`

### Public Endpoints

| Method | Endpoint       | Description           |
| ------ | -------------- | --------------------- |
| GET    | `/public/ping` | Health check endpoint |

### Protected Endpoints (Requires Authentication)

#### Team Management

| Method | Endpoint                             | Description              |
| ------ | ------------------------------------ | ------------------------ |
| POST   | `/teams`                             | Create a new team        |
| POST   | `/teams/:teamId/members`             | Add member to team       |
| GET    | `/teams/:teamId/members`             | Get team members         |
| DELETE | `/teams/:teamId/members/:memberId`   | Remove member from team  |
| DELETE | `/teams/:teamId/managers/:managerId` | Remove manager from team |

#### Asset Management

**Folders**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/assets/folders` | Create a new folder |
| GET | `/assets/folders/:folderId` | Get folder details |
| PUT | `/assets/folders/:folderId` | Update folder |
| DELETE | `/assets/folders/:folderId` | Delete folder |

**Notes**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/assets/folders/:folderId/notes` | Create note in folder |
| GET | `/assets/notes/:noteId` | Get note details |
| PUT | `/assets/notes/:noteId` | Update note |
| DELETE | `/assets/notes/:noteId` | Delete note |

**Sharing**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/assets/folders/:folderId/shares` | Share folder with user |
| DELETE | `/assets/folders/:folderId/shares/:userId` | Revoke folder share |
| POST | `/assets/notes/:noteId/shares` | Share note with user |
| DELETE | `/assets/notes/:noteId/shares/:userId` | Revoke note share |

#### Manager Operations

| Method | Endpoint                        | Description                    |
| ------ | ------------------------------- | ------------------------------ |
| GET    | `/manager/teams/:teamId/assets` | Get team assets (manager only) |
| GET    | `/manager/users/:userId/assets` | Get user assets (manager only) |

#### Import Operations

| Method | Endpoint        | Description       |
| ------ | --------------- | ----------------- |
| POST   | `/import-users` | Bulk import users |

## 🔐 Authentication

All protected endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## 🌐 Deployment

### Production Environment

The application is deployed on **Render** and accessible at:
**https://go-service-app.onrender.com/**

### CI/CD Pipeline

This project uses **GitHub Actions** for continuous integration and deployment:

- **Trigger**: Automatic deployment on every push to the main branch
- **Process**:
  1. Code is pushed to GitHub repository
  2. GitHub Actions workflow is triggered
  3. Docker image is built automatically
  4. Image is deployed to Render platform
  5. Service is updated with zero downtime

### Environment Variables

Required environment variables for deployment:

```env
# JWT configuration
ACCESS_TOKEN_SECRET=abc
REFRESH_TOKEN_SECRET=abc

#User_service url

USER_SERVICE_URL=abc

# Server configuration
PORT=8080

BOOTSTRAP_HOST=abc
KAFKA_USERNAME=abc
KAFKA_PASSWORD=abc

REDIS_ADDRESS=abc
REDIS_PASSWORD=abc
```

## 🧪 Development

### Local Development Setup

1. **Install dependencies**:

   ```bash
   go mod download
   ```

2. **Set up environment variables**:

   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

---

**Built with ❤️ by Híu**

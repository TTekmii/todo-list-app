# 📝 Todo List API & Client

Production-ready Todo List application built with **Go (Clean Architecture)** and **Vue 3**. Features secure JWT authentication, structured logging, graceful shutdown, and a responsive dark-mode UI.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Vue.js](https://img.shields.io/badge/Vue.js-4FC08D?style=for-the-badge&logo=vuedotjs&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)


## ✨ Features

-   **Clean Architecture**: Strict separation of concerns (Domain, App, Infrastructure, Transport) for high testability and maintainability.
-   **Secure Auth**: JWT-based stateless authentication with Bcrypt password hashing.
-   **Structured Logging**: Context-aware logging using `log/slog` (JSON format for production, Pretty for development).
-   **Graceful Shutdown**: Proper handling of OS signals (SIGINT/SIGTERM) to complete active requests before exiting.
-   **Responsive UI**: Modern, lightweight frontend built with Vue 3 (Composition API) and Tailwind CSS, featuring Dark Mode support.
-   **Configuration Management**: Flexible config loading via Viper (supports YAML files and Environment Variables).
-   **API Documentation**: Auto-generated Swagger/OpenAPI docs available at `/swagger/index.html`.
-   **Docker Support**: Easy deployment with Docker Compose for PostgreSQL database.

## 🏗 Architecture

The project follows a layered architecture to ensure scalability and clean code organization:

```
internal/
├── domain/          # Business entities (Models) and Interface contracts (Repositories)
├── app/             # Business logic implementation (Services)
│   ├── auth/        # Authentication service (JWT, Bcrypt)
│   └── todo/        # Todo list and item services
├── infrastructure/  # External dependencies implementation
│   ├── database/    # PostgreSQL connection setup
│   └── repository/  # SQLX-based repository implementations
└── transport/       # HTTP delivery layer
    ├── http/
    │   ├── handler/     # Request parsing, validation, and response formatting
    │   ├── middleware/  # Auth checks, logging, recovery
    │   ├── dto/         # Data Transfer Objects for API isolation
    │   └── router.go    # Route definitions
    └── server.go    # HTTP server configuration and lifecycle
```
## 🛠 Tech Stack

### Backend

- **Language:** Go 1.21+
- **Web Framework:** [Gin Gonic](https://gin-gonic.com/ "gin gonic")
- **Database:** PostgreSQL
- **SQL Toolkit:** [SQLX](https://github.com/jmoiron/sqlx "sqlx")
- **Authentication:** JWT (`golang-jwt/jwt/v5`) & Bcrypt
- **Logging:** Standard library `log/slog`
- **Config:** [Viper](https://github.com/spf13/viper "viper")
- **Docs:** [Swagger](https://swagger.io/ "swagger")

### Frontend

- **Framework:** Vue 3 (Composition API)
- **Styling:** Tailwind CSS
- **HTTP Client:** Axios

## 🚀 Getting Started

#### Prerequisites
- Go 1.21 or higher
- Docker & Docker Compose (for running PostgreSQL)
- Git

#### 1. Clone the Repository
```bash
git clone https://github.com/TTekmii/todo-list-app.git
cd todo-list-app
```

#### 2. Configure Environment

Create a `.env` file in the root directory. You can copy `.env.example` if available.

```env
# Server Configuration
PORT=8000

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret_password
DB_NAME=todo_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your_super_duper_secret_key
JWT_TTL=12h

# Logger Configuration
LOG_LEVEL=debug      # Use 'info' or 'warn' for production
LOG_FORMAT=pretty    # Use 'json' for production
APP_ENV=development
```

> *For `JWT_SECRET` generation run this command in the Linux, Mac or Git Bash terminal:*

```bash
openssl rand -base64 48
```
>*...or in PowerShell:*

```powershell
$bytes = New-Object Byte[] 48
[System.Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
[Convert]::ToBase64String($bytes)
```

#### 3. Start Database

Run PostgreSQL using Docker Compose:

```bash
docker-compose up -d
```

#### 4. Run Migrations (Optional)
If you have migration files, apply them now. Otherwise, ensure your database schema matches the structs in `internal/domain/model`.

#### 5. Run the Application

```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/api/main.go
```

The API will start on `http://localhost:8000`.

## 📦 Building for Windows

To build the application with a custom icon and version information on Windows:

#### 1. Ensure you have `goversioninfo` installed:
 ```bash
 go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
```

#### 2. Run the build command using Make:

```bash
make build-windows
```

This will generate `todo-app.exe` with the application icon and metadata embedded.

## 📖 API Documentation

Once the server is running, interactive Swagger documentation is available at:
👉 http://localhost:8000/swagger/index.html

#### Key Endpoints

Method | Endpoint | Description | Auth Required
:------|:---------|:------------|:------:
POST | `/auth/sign-up` | Register a new user | No
POST | `/auth/sign-in` | Login and get JWT token | No
GET | `/api/lists` | Get all user lists | Yes
POST | `/api/lists` | Create a new list | Yes
DELETE | `/api/lists/:id` | Delete a list | Yes
GET | `/api/lists/:id/items` | Get items in a list | Yes
POST | `/api/lists/:id/items` | Add item to list | Yes
DELETE | `/api/lists/:listId/items/:itemId` | Delete an item | Yes

## 🎨 Frontend Client

A simple, responsive single-page application is included in the `web/` directory.

**1.** Ensure the backend API is running on `http://localhost:8000`.
**2.** Open the file `web/index.html` in any modern web browser (Chrome, Firefox, Edge).
**3.** Register a new account or log in to start managing your tasks.
**4.** Use the toggle button in the top-right corner to switch between Dark and Light modes.

> ***Note:*** *Since this is a static HTML file, no build process or Node.js server is required for the frontend. It communicates directly with the Go API via Axios.*

## 🛑 Stopping the Application

#### 1. Stop the Go Server
If you are running the API in a terminal (`go run ...` or `./todo-app.exe`):
- Press **`Ctrl + C`** in the terminal window.
- The server will perform a **Graceful Shutdown**, completing active requests and closing connections safely.

#### 2. Stop the Database
The PostgreSQL database runs in a Docker container and continues to run even after the server is stopped. To shut it down and free up system resources:

```bash
docker-compose down
```

>***Note:*** *Avoid killing the process forcefully (e.g., closing the terminal window directly) as it may interrupt active database transactions.*
>***Tip:*** *You can also use `docker-compose stop` if you plan to start working again soon (this pauses the containers without removing them).* 

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
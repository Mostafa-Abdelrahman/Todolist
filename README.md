# Todo List Full-Stack Application

A modern todo list application built with Vue.js frontend and Go backend.

## Features

- **Frontend (Vue.js 3 + Tailwind CSS)**
  - Modern, responsive UI
  - Dashboard view with all lists and tasks
  - Today's tasks view
  - Upcoming tasks view
  - Collapsible sidebar navigation
  - Real-time API integration

- **Backend (Go + Gin + GORM + SQLite)**
  - RESTful API endpoints
  - CRUD operations for todos and lists
  - SQLite database with automatic migrations
  - CORS enabled for frontend integration
  - Clean architecture with services and handlers

## Screenshots

### Dashboard View
![Dashboard](screens/Screenshot%202025-07-14%20142824.png)

### Create/Edit List
![Create/Edit List](screens/Screenshot%202025-07-14%20142837.png)

### Today's Tasks & Upcoming Tasks
![Today/Upcoming](screens/Screenshot%202025-07-14%20142910.png)

## Project Structure

```
TodoList/
├── UI/                 # Vue.js frontend
│   ├── src/
│   │   ├── components/ # Vue components
│   │   ├── views/      # Page components
│   │   ├── services/   # API service
│   │   └── router/     # Vue Router
│   └── package.json
├── Server/             # Go backend
│   ├── cmd/server/     # Main application
│   ├── internal/
│   │   ├── api/        # HTTP handlers and routes
│   │   ├── models/     # Database models
│   │   ├── services/   # Business logic
│   │   ├── config/     # Configuration
│   │   └── database/   # Database connection
│   └── go.mod
├── screens/            # Screenshots for documentation
└── run.sh              # Script to run both servers
```

## Quick Start

### Prerequisites

- **Go** (1.19 or later)
- **Node.js** (16 or later)
- **npm** or **yarn**

### Installation

1. **Clone the repository** (if not already done)
2. **Install frontend dependencies:**
   ```bash
   cd TodoList/UI
   npm install
   ```

3. **Install backend dependencies:**
   ```bash
   cd TodoList/Server
   go mod download
   ```

### Running the Application

#### Option 1: Run both servers with one command
```bash
cd TodoList
./run.sh
```

#### Option 2: Run servers separately

**Backend (Terminal 1):**
```bash
cd TodoList/Server
go run cmd/server/main.go
```

**Frontend (Terminal 2):**
```bash
cd TodoList/UI
npm run dev
```

### Access the Application

- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080
- **API Health Check:** http://localhost:8080/api/health

## API Endpoints

### Todos
- `GET /api/todos` - Get all todos
- `GET /api/todos/:id` - Get specific todo
- `POST /api/todos` - Create new todo
- `PUT /api/todos/:id` - Update todo
- `DELETE /api/todos/:id` - Delete todo

### Lists
- `GET /api/lists` - Get all lists
- `GET /api/lists/:id` - Get specific list
- `POST /api/lists` - Create new list
- `PUT /api/lists/:id` - Update list
- `DELETE /api/lists/:id` - Delete list

## Development

### Frontend Development
- Located in `TodoList/UI/`
- Uses Vue 3 with Composition API
- Tailwind CSS for styling
- Vite for fast development server

### Backend Development
- Located in `TodoList/Server/`
- Go with Gin web framework
- GORM for database operations
- SQLite database (stored as `todos.db`)

### Database
- SQLite database file: `TodoList/Server/todos.db`
- Automatic migrations on startup
- Models: `Todo` and `List`

## Testing

### Backend Tests
```bash
cd TodoList/Server
go test ./...
```

### API Testing with curl
```bash
# Get all todos
curl http://localhost:8080/api/todos

# Create a todo
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Todo","listId":1,"due":"2024-12-25","done":false}'

# Get all lists
curl http://localhost:8080/api/lists

# Create a list
curl -X POST http://localhost:8080/api/lists \
  -H "Content-Type: application/json" \
  -d '{"name":"My List"}'
```

## Configuration

The backend uses environment variables for configuration:

- `SERVER_PORT` - Server port (default: 8080)
- `DATABASE_URL` - Database connection string (default: "todos.db")
- `ENVIRONMENT` - Environment (default: "development")

## Troubleshooting

1. **Frontend can't connect to backend:**
   - Ensure backend is running on port 8080
   - Check CORS configuration in backend
   - Verify API_BASE_URL in frontend service

2. **Database issues:**
   - Delete `todos.db` file and restart server
   - Check database migrations

3. **Port conflicts:**
   - Change `SERVER_PORT` environment variable
   - Update `API_BASE_URL` in frontend service

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

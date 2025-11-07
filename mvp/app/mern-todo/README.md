# MERN Todo App

A simple full-stack todo application built with MongoDB, Express, React, and Node.js, designed to test NebulaBox's containerization and deployment capabilities.

## 🏗️ Architecture

```
┌─────────────────┐
│   React (Nginx) │  Port 80 (mapped to 3000)
│   Frontend      │
└────────┬────────┘
         │
         │ HTTP API calls
         │
┌────────▼────────┐
│   Express API   │  Port 5000
│   Backend       │
└────────┬────────┘
         │
         │ MongoDB queries
         │
┌────────▼────────┐
│   MongoDB        │  Port 27017
│   Database       │
└─────────────────┘
```

## 📁 Project Structure

```
mern-todo/
├── backend/
│   ├── server.js          # Express API server
│   ├── package.json       # Backend dependencies
│   └── Dockerfile         # Backend container (optional)
├── frontend/
│   ├── src/
│   │   ├── App.js         # Main React component
│   │   ├── App.css        # Styles
│   │   ├── index.js       # React entry point
│   │   └── index.css      # Global styles
│   ├── public/
│   │   └── index.html     # HTML template
│   ├── package.json       # Frontend dependencies
│   └── Dockerfile         # Frontend container (optional)
├── scripts/
│   └── start.sh           # Startup script (MongoDB + Backend + Nginx)
├── buildspec.json         # NebulaBox build specification
├── README.md              # This file
└── TESTING_GUIDE.md       # Detailed testing instructions
```

## 🚀 Quick Start

### Prerequisites
- NebulaBox API server running (port 8081)
- NebulaBox Dashboard running (port 3001)
- Docker/containerd (for container runtime)

### Build & Run

1. **Build the image** (via Dashboard):
   - Go to http://localhost:3001/buildspec
   - Load `buildspec.json`
   - Click "Validate" then "Build"

2. **Run the container** (via Dashboard):
   - Go to http://localhost:3001/containers
   - Create container with:
     - Image: `mern-todo:latest`
     - Ports: `3000:80` and `5000:5000`

3. **Access the app**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:5000/api

## 🧪 Testing

See [TESTING_GUIDE.md](./TESTING_GUIDE.md) for detailed testing instructions.

## 📡 API Endpoints

### Health Check
```
GET /api/health
Returns: { status: "ok", message: "..." }
```

### Get All Todos
```
GET /api/todos
Returns: Array of todo objects
```

### Create Todo
```
POST /api/todos
Body: { text: "Todo text" }
Returns: Created todo object
```

### Update Todo
```
PUT /api/todos/:id
Body: { text: "Updated text", completed: true }
Returns: Updated todo object
```

### Delete Todo
```
DELETE /api/todos/:id
Returns: { message: "Todo deleted successfully", todo: {...} }
```

## 🔧 Configuration

### Environment Variables

**Backend:**
- `PORT`: Backend server port (default: 5000)
- `MONGODB_URI`: MongoDB connection string (default: mongodb://localhost:27017/todos)

**Frontend:**
- `REACT_APP_API_URL`: Backend API URL (default: http://localhost:5000)

### Ports

- **3000**: Frontend (mapped from container port 80)
- **5000**: Backend API
- **27017**: MongoDB (internal)

## 📦 Technologies

- **MongoDB**: NoSQL database
- **Express.js**: Web framework for Node.js
- **React**: Frontend library
- **Node.js**: JavaScript runtime
- **Nginx**: Web server (for serving React build)
- **NebulaBox**: Container platform

## 🎯 Features

- ✅ Add todos
- ✅ Mark todos as complete
- ✅ Delete todos
- ✅ Persistent storage (MongoDB)
- ✅ RESTful API
- ✅ Modern React UI
- ✅ Health checks
- ✅ Containerized deployment

## 🐛 Troubleshooting

See [TESTING_GUIDE.md](./TESTING_GUIDE.md) for troubleshooting tips.

## 📝 License

This is a test application for NebulaBox.


#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting Todo List Application...${NC}"

# Function to cleanup background processes
cleanup() {
    echo -e "\n${RED}Stopping servers...${NC}"
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Start backend server
echo -e "${GREEN}Starting Go backend server...${NC}"
cd Server
go run cmd/server/main.go &
BACKEND_PID=$!
cd ..

# Wait a moment for backend to start
sleep 2

# Start frontend server
echo -e "${GREEN}Starting Vue.js frontend server...${NC}"
cd UI
npm run dev &
FRONTEND_PID=$!
cd ..

echo -e "${BLUE}Both servers are running!${NC}"
echo -e "${GREEN}Backend:${NC} http://localhost:8080"
echo -e "${GREEN}Frontend:${NC} http://localhost:5173"
echo -e "${BLUE}Press Ctrl+C to stop both servers${NC}"

# Wait for both processes
wait 
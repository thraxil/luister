.PHONY: all build frontend backend clean help

# Default target
all: build

# Build everything
build: frontend backend

# Build the frontend (Vue.js app)
frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

# Build the backend (Statically linked Go binary)
backend:
	@echo "Building backend (statically linked)..."
	CGO_ENABLED=0 go build -o luister .

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f luister
	rm -rf frontend/dist
	rm -rf frontend/node_modules

# Help target
help:
	@echo "Available targets:"
	@echo "  all       - Build both frontend and backend (default)"
	@echo "  build     - Alias for all"
	@echo "  frontend  - Install dependencies and build the frontend"
	@echo "  backend   - Build the statically linked Go binary"
	@echo "  clean     - Remove binary and build artifacts"
	@echo "  help      - Show this help message"

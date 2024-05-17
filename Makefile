# Build project
code.build:
	go build -v -o launch ./cmd

# Start application from launch script
app.run:
	go run ./cmd/main.go

# BSC Transaction Listener Service

This repository contains a Go service that listens to new blocks on the Binance Smart Chain (BSC) network and processes transactions involving specific contract addresses and watched wallet addresses. The service is built using Go, the go-ethereum client, and exposes a simple HTTP API using Fiber.

## Features

- Connects to a BSC node via WebSocket.
- Monitors new blocks in real-time.
- Detects and processes transactions to watched addresses and specified contract addresses (e.g., USDT, USDC, BUSD).
- Decodes contract transfer data for supported tokens.
- Easily add or remove watched addresses.
- Exposes a `/metrics` endpoint for monitoring.

## Getting Started

1. **Install dependencies:**  
   Run `go mod tidy` to install all required Go modules.

2. **Configure contracts and addresses:**  
   Edit `main.go` to specify contract addresses and watched wallet addresses.

3. **Run the service:**  
   ```
   go run main.go
   ```

4. **Access metrics:**  
   Visit `http://localhost:3002/metrics` for service metrics.

## Project Structure

- `main.go`: Application entry point and HTTP server.
- `bsc/bsc.go`: BSC listener logic and transaction processing.

## Requirements

- Go 1.23+
- Access to a BSC node with WebSocket support.

## License

MIT
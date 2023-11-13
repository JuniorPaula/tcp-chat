# TCP Chat Server

This README provides an overview of a Go programming project that implements a simple TCP chat server. The chat server allows multiple clients to connect, exchange messages, and handles basic security features such as banning clients based on their behavior.

## Project Structure

- The main program is defined in the `main` package.
- The project is built using Go (Golang).
- The code listens for incoming TCP connections on a specified port.
- It creates a separate goroutine for each connected client, enabling concurrent communication.
- It includes basic security features such as banning misbehaving clients.

## Project Components

1. **Constants**

    The project defines various constants:
    - `Port`: The TCP port on which the server listens.
    - `SafeMode`: A boolean that, when set to `true`, redacts sensitive information in log messages.
    - `MessageRate`: The rate at which clients can send messages.
    - `BanLimit`: The duration for which clients are banned if they exceed the strike limit.
    - `StrikeLimit`: The number of strikes (message rate violations) before a client is banned.

2. **Message Types and Structures**

    The project defines several message types and structures:
    - `MessageType`: An enumeration for message types, including `ClientConnected`, `ClientDisconnected`, and `NewMessage`.
    - `Message`: A structure containing the message type, connection, and text.

3. **Server Logic**

    The `server` function handles the core logic of the chat server. It maintains a map of connected clients, as well as a map of banned clients with their ban expiration times.

    - Upon receiving a `ClientConnected` message, the server checks if the client is banned. If not, it adds the client to the list of connected clients.
    - If a client is banned, it sends a ban message and closes the connection.
    - When a `ClientDisconnected` message is received, the server removes the client from the list of connected clients.
    - For a `NewMessage` message, the server checks the message rate for the sender. If the sender has violated the message rate limit, they receive strikes. If strikes exceed the limit, the client is banned.

4. **Client Logic**

    The `client` function handles communication with individual clients. It reads messages from the client and forwards them to the `server`.

5. **Main Function**

    The `main` function initializes the server by listening on the specified port and starts a goroutine for the `server`. It also accepts incoming client connections, creating a new goroutine for each connected client.

## Getting Started

To run the chat server, follow these steps:

1. Make sure you have Go (Golang) installed on your system.

2. Copy the code to a Go source file (e.g., `main.go`).

3. Build the project using the `go build` command:

    ```sh
    go build main.go
    ```

4. Run the resulting executable:

    ```sh
    ./main
    ```

The chat server will start listening on the specified port, and clients can connect to it.

## Usage

- Clients can connect to the server using a TCP client, such as `telnet` or custom client software.
- The server will log connections, disconnections, and messages exchanged between clients.
- It implements basic security features to ban misbehaving clients who exceed the message rate limit.

## Contributors

This project was created by [Your Name] and is open to contributions and improvements.

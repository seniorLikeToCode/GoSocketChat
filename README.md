# Simple Chat Server

This project is a simple chat server implemented in Go. It allows multiple clients to connect, choose nicknames, join chat rooms, send messages, and see available rooms.

## Features

-   **Nicknames**: Clients can set their nicknames.
-   **Join Rooms**: Clients can join existing chat rooms or create new ones.
-   **Message Rooms**: Clients can send messages to all members of a room.
-   **List Rooms**: Clients can list all available rooms.
-   **Quit**: Clients can quit the server gracefully.

## Getting Started

### Prerequisites

-   Go (version 1.16+)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/seniorLikeToCode/GoSocketChat.git
    cd GoSocketChat
    ```

2. Build the server:

    ```sh
    go build -o main.go
    ```

3. Run the server:
    ```sh
    ./main
    ```

The server will start and listen for connections on port 5000.

## Usage

### Connecting to the Server

You can use `telnet` or `netcat` to connect to the server:

```sh
telnet localhost 5000
```

or

```sh
nc localhost 5000
```

### Commands

Once connected, you can use the following commands:

-   **/nick NAME**: Set your nickname.

    ```sh
    /nick mynickname
    ```

-   **/join ROOM_NAME**: Join a room or create a new one if it doesn't exist.

    ```sh
    /join room1
    ```

-   **/msg MSG**: Send a message to the current room.

    ```sh
    /msg Hello everyone!
    ```

-   **/rooms**: List all available rooms.

    ```sh
    /rooms
    ```

-   **/quit**: Quit the server.
    ```sh
    /quit
    ```

### Example Session

```sh
telnet localhost 5000
```

```
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
/nick Alice
all right, I will call you Alice
/join room1
Welcome to room1
/msg Hello, room1!
Alice: Hello, room1!
/rooms
available rooms: room1
/quit
sad to see you go :(
Connection closed by foreign host.
```

## Logging

The server logs all connections, disconnections, and errors to the standard output.

---

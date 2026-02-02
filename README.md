# Peril: A RabbitMQ Pub/Sub Strategy Game

A multiplayer strategy game demonstrating publish-subscribe architecture using RabbitMQ. Players spawn military units, move them across continents, and engage in battles when units occupy the same location. This project was built as part of the [Boot.dev Learn Pub/Sub course](https://www.boot.dev/courses/learn-pub-sub-rabbitmq-golang).

## Architecture Overview

This project implements a distributed game system where:
- **Server**: Manages game state (pause/resume) and aggregates game logs
- **Clients**: Individual players who spawn units, move armies, and receive war notifications
- **RabbitMQ**: Message broker handling all inter-client communication

### Message Flow

```
Client A (Move) → RabbitMQ Topic Exchange → All Clients
                                          ↓
                                   Client B detects conflict
                                          ↓
                                   War Recognition published
                                          ↓
                                   All Clients notified
```

### RabbitMQ Exchanges

1. **peril_direct** (Direct Exchange)
   - Pause/resume game state commands from server to all clients

2. **peril_topic** (Topic Exchange)
   - Army movements: `army_moves.*`
   - War recognitions: `war.*`
   - Game logs: `game_logs.*`

## Game Mechanics

### Unit Types (Ranks)
- **Infantry**: Basic ground units
- **Cavalry**: Mobile combat units
- **Artillery**: Heavy firepower units

### Locations (Continents)
- Americas
- Europe
- Africa
- Asia
- Australia
- Antarctica

### Combat System
When two players have units in the same location, a war is triggered and both players are notified via the `war.*` routing key.

## Prerequisites

- **Go 1.22.1+**
- **Docker** (for RabbitMQ)
- **RabbitMQ 3.13** (via Docker)

## Installation & Setup

### 1. Clone and Install Dependencies

```bash
cd rabbitmq
go mod download
```

### 2. Start RabbitMQ

The project includes a helper script to manage RabbitMQ via Docker:

```bash
# Start RabbitMQ container (creates if doesn't exist)
./rabbit.sh start

# Stop RabbitMQ container
./rabbit.sh stop

# View RabbitMQ logs
./rabbit.sh logs
```

RabbitMQ will be available at:
- **AMQP Port**: `localhost:5672`
- **Management UI**: `http://localhost:15672` (username: `guest`, password: `guest`)

### 3. Verify RabbitMQ is Running

```bash
docker ps | grep peril_rabbitmq
```

You should see the container running on ports 5672 and 15672.

## Usage

### Starting the Game Server

The server manages global game state and logs all player actions:

```bash
go run ./cmd/server
```

**Server Commands:**
- `pause` - Pause the game (prevents all client movements)
- `resume` - Resume the game
- `quit` - Shut down the server
- `help` - Display available commands

### Starting Game Clients

Each client represents a player. You can run multiple clients in separate terminals:

```bash
# Terminal 1
go run ./cmd/client

# Terminal 2
go run ./cmd/client

# Terminal 3
go run ./cmd/client
```

When starting, each client will prompt you for a username.

### Multiple Servers (Optional)

You can run multiple server instances to demonstrate fanout behavior:

```bash
# Start 3 server instances
./multiserver.sh 3
```

Press `Ctrl+C` to stop all server instances.

## Client Commands

Once a client is running, you can issue the following commands:

### Spawn Units
Create a new military unit at a location:

```bash
spawn <location> <rank>
```

**Examples:**
```bash
spawn europe infantry
spawn asia cavalry
spawn africa artillery
```

### Move Units
Move one or more units to a new location:

```bash
move <location> <unitID> [unitID] [unitID]...
```

**Examples:**
```bash
move asia 1          # Move unit 1 to Asia
move europe 2 3 4    # Move units 2, 3, and 4 to Europe
```

**Notes:**
- Units are assigned sequential IDs starting from 1
- You can only move your own units
- Moving to a location occupied by another player triggers war
- Game must not be paused to move units

### Check Status
View your current units and game state:

```bash
status
```

Output shows:
- Whether the game is paused
- Your username
- Number of units you control
- Details of each unit (ID, location, rank)

### Spam Game Logs
Send multiple malicious/test log messages (for testing):

```bash
spam <count>
```

**Example:**
```bash
spam 5    # Sends 5 random war quotes to game logs
```

### Other Commands
```bash
help    # Display command reference
quit    # Exit the client
```

## Typical Gameplay Session

### Example: Two-Player Battle

**Terminal 1 - Server:**
```bash
$ go run ./cmd/server
Peril game server connected to RabbitMQ!

Possible commands:
* pause
* resume
* quit
* help
>
```

**Terminal 2 - Player "Alice":**
```bash
$ go run ./cmd/client
Welcome to the Peril client!
Please enter your username:
> Alice
Welcome, Alice!

> spawn europe infantry
Spawned a(n) infantry in europe with id 1

> spawn europe cavalry
Spawned a(n) cavalry in europe with id 2

> move asia 1 2
Moved 2 units to asia
```

**Terminal 3 - Player "Bob":**
```bash
$ go run ./cmd/client
Welcome to the Peril client!
Please enter your username:
> Bob

> spawn africa artillery
Spawned a(n) artillery in africa with id 1

> move asia 1
Moved 1 units to asia

==== Move Detected ====
Alice is moving 2 unit(s) to asia
* infantry
* cavalry
You have units in asia! You are at war with Alice!
```

At this point:
- Bob receives notification that Alice moved units to Asia
- Bob's client detects his units are also in Asia
- War is declared between Alice and Bob
- Both players are notified of the conflict

## Project Structure

```
rabbitmq/
├── cmd/
│   ├── server/          # Game server implementation
│   │   ├── main.go      # Server entry point
│   │   └── logs.go      # Game log handler
│   └── client/          # Game client implementation
│       ├── main.go      # Client entry point
│       ├── pause.go     # Pause/resume handler
│       ├── move.go      # Army movement handler
│       └── war.go       # War recognition handler
├── internal/
│   ├── gamelogic/       # Core game logic
│   │   ├── gamedata.go  # Game data types (Player, Unit, Location)
│   │   ├── gamestate.go # Thread-safe game state management
│   │   ├── gamelogic.go # User input/output helpers
│   │   ├── move.go      # Movement and war detection logic
│   │   ├── spawn.go     # Unit spawning logic
│   │   ├── pause.go     # Pause handling logic
│   │   ├── war.go       # War declaration logic
│   │   └── logs.go      # Game logging
│   ├── routing/         # RabbitMQ routing configuration
│   │   ├── routing.go   # Exchange and routing key constants
│   │   └── models.go    # Message models (PlayingState, GameLog)
│   └── pubsub/          # RabbitMQ pub/sub abstractions
│       ├── subscribe.go # Generic subscription handler
│       ├── json.go      # JSON marshalling/publishing
│       ├── binding.go   # Queue declaration and binding
│       └── logs.go      # Logging utilities
├── rabbit.sh            # RabbitMQ container management script
├── multiserver.sh       # Multiple server launcher script
├── go.mod               # Go module definition
└── README.md            # This file
```

## Key Implementation Details

### Queue Types

The project uses two queue types:

1. **Transient Queues** (Client-specific)
   - Auto-deleted when client disconnects
   - Used for: pause notifications, army movements per client
   - Example: `pause.Alice`, `army_moves.Bob`

2. **Durable Queues** (Persistent)
   - Survive broker restarts
   - Used for: game logs, war notifications
   - Ensures important events aren't lost

### Message Serialization

- **JSON**: Used for most messages (PlayingState, ArmyMove, RecognitionOfWar)
- **Gob**: Used for GameLog messages (Go binary encoding)

### Thread Safety

The `GameState` type uses `sync.RWMutex` to ensure thread-safe access to player data across concurrent goroutines handling RabbitMQ messages.

### Routing Keys

| Pattern | Description | Example |
|---------|-------------|---------|
| `pause` | Direct routing for pause/resume | `pause` |
| `game_logs.*` | All game logs | `game_logs.Alice` |
| `army_moves.*` | All army movements | `army_moves.Bob` |
| `war.*` | War recognition notifications | `war.Alice.Bob` |

## Learning Objectives

This project demonstrates:

✅ **Publisher/Subscriber Pattern**
- Decoupled clients communicating via message broker
- Multiple subscribers receiving same messages

✅ **Topic-Based Routing**
- Wildcard routing keys (`army_moves.*`, `war.*`)
- Selective message consumption based on topics

✅ **Direct Exchange Routing**
- Broadcast messages (pause/resume to all clients)

✅ **Message Acknowledgements**
- Manual ACK/NACK with requeue logic
- Ensuring reliable message processing

✅ **Queue Declaration Strategies**
- Durable vs transient queues
- Auto-delete queues for temporary clients

✅ **Concurrent Message Handling**
- Goroutines processing messages asynchronously
- Thread-safe shared state management

## Troubleshooting

### Client can't connect to RabbitMQ
```
Error: could not connect to RabbitMQ
```
**Solution:** Ensure RabbitMQ is running:
```bash
./rabbit.sh start
docker ps | grep peril_rabbitmq
```

### Game is paused error
```
Error: the game is paused, you can not move units
```
**Solution:** On the server terminal, type `resume` and press Enter.

### Unit not found error
```
Error: unit with ID X not found
```
**Solution:** Use the `status` command to view your available units and their IDs.

### Invalid location/rank error
**Solution:** Valid locations are: `americas`, `europe`, `africa`, `asia`, `australia`, `antarctica`
Valid ranks are: `infantry`, `cavalry`, `artillery`

### Port already in use (5672)
```
Error: port is already allocated
```
**Solution:** Stop any existing RabbitMQ containers:
```bash
docker stop peril_rabbitmq
./rabbit.sh start
```

## Further Exploration

- **RabbitMQ Management UI**: Visit `http://localhost:15672` to view:
  - Active queues and their message counts
  - Exchange bindings and routing topology
  - Message rates and throughput graphs
  - Connection and channel details

- **Experiment with Exchanges**: Try modifying routing keys and bindings to change message flow patterns

- **Load Testing**: Use the `spam` command to generate high message volumes and observe RabbitMQ performance

- **Add Features**: Consider implementing:
  - Battle resolution logic (unit strength calculations)
  - Territory control tracking
  - Score/victory conditions
  - Alliance system between players

## Credits

Built by Anu Raghavan as part of the [Boot.dev Pub/Sub course](https://www.boot.dev/courses/learn-pub-sub-rabbitmq-golang).

## License

This is an educational project. Feel free to use and modify as needed for learning purposes.

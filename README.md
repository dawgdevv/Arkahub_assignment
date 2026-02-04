# EnergyGrid Mock API and Client Implientaiton

This is the mock backend server and client for the EnergyGrid Data Aggregator coding assignment.

## Prerequisites

- Node.js (v14 or higher)
- npm (Node Package Manager)
- Go()

## Setup and Run The Mock API Server

2.  **Install dependencies:**

    ```bash
    npm install
    ```

3.  **Start the server:**

    ```bash
    npm start
    ```

    Or directly:

    ```bash
    node server.js
    ```

4.  **Verify:**
    You should see the following output:
    ```
    ⚡ EnergyGrid Mock API running on port 3000
       Constraints: 1 req/sec, Max 10 items/batch
    ```
    The server is now listening at `http://localhost:3000`.

## API Details

- **Base URL:** `http://localhost:3000`
- **Endpoint:** `POST /device/real/query`
- **Auth Token:** `interview_token_123`

### Security Headers Required

Every request must include:

- `timestamp`: Current time in milliseconds.
- `signature`: `MD5( URL + Token + timestamp )`

### Constraints

- **Rate Limit:** 1 request per second.
- **Batch Size:** Max 10 serial numbers per request.

See `instructions.md` for full details.

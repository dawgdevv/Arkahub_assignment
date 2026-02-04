# EnergyGrid Mock API and Client Implimentaiton

This is the mock backend server and client for the EnergyGrid Data Aggregator coding assignment.

## Prerequisites

- Node.js (v14 or higher)
- npm (Node Package Manager)
- Go(1.22+)

## Client folder structure and modules:

1. main.go: Entry point; runs fetch, aggregation, and export.

2. config/config.go: Constants for API URL, token, batch size, total devices.
3. auth/signature.go: GenerateSignature() MD5 signing.
4. cleint/client.go: HTTP client, rate limiting, batching loop, retries.
5. utils/serialnumbers.go: Serial number generation and batching helpers.
6. models/types.go: Request/response models and device data types.
7. aggregator/aggregator.go: Aggregation stats output.
8. export/export.go: JSON + report export.

## Brief explanation of my approach (how I handled rate limiting and concurrency).

1. Rate limiting: fetchBatch() enforces a minimum 1‑second gap between requests by tracking lastReqTime and sleeping when needed, so only one request is in flight at a time.

2. Concurrency: there is no parallelism by design; batches are processed sequentially in FetchAllDevices(), which keeps concurrency at 1 and guarantees compliance with the strict 1 req/s limit.

Reason for this approach: the API enforces a strict global 1 req/sec limit, so a simple sequential pipeline maximizes compliance and avoids 429s without extra complexity.

## Alternative handling options and when to use them

| Option                                  | What it does                                   | When to use it                                               |
| --------------------------------------- | ---------------------------------------------- | ------------------------------------------------------------ |
| Token bucket limiter (global)           | Enforces a max request rate across all workers | Rate limit is strict and shared across all workers/instances |
| Leaky bucket limiter                    | Smooths bursts by releasing at a steady rate   | Upstream rejects bursts even if average rate is fine         |
| Worker pool + global limiter            | Bounded concurrency with a shared rate gate    | Rate limit is strict but you want parallel parsing/IO        |
| Semaphore (max in-flight)               | Caps simultaneous requests                     | Upstream allows high QPS but limited concurrent connections  |
| Retry with exponential backoff + jitter | Spreads retries to avoid thundering herd       | 429/5xx responses or flaky networks                          |
| Adaptive rate (auto-tune)               | Decrease on 429, slowly ramp up                | Rate limit is not fixed or varies over time                  |
| Circuit breaker                         | Stops calls after repeated failures            | Extended outages or consistent 5xx/429 errors                |
| Batch size tuning                       | Adjusts payload size per request               | Batch limits change or you need faster partial retries       |
| Priority queue                          | Serves critical devices first                  | SLA devices need faster updates                              |
| Timeout + context cancellation          | Cancels slow/stuck requests                    | Protects workers from hanging connections                    |

## Setup and Run The Mock API Server and Client

```bash
git clone <repo_url>

cd Arkahub_assignment
```

### Terminal 1: Start Mock API Server

```bash
npm install
npm start
```

**Verify:**
You should see the following output:
`     ⚡ EnergyGrid Mock API running on port 3000
       Constraints: 1 req/sec, Max 10 items/batch
    `
The server is now listening at `http://localhost:3000`.

### Terminal 2: Run Go Client

```bash
cd clientingo
go run main.go
```

### Output

The client will:

- Fetch data from 500 devices in approximately 50 seconds
- Display real-time progress
- Generate aggregation report
- Export 2 files:
  - `energygrid_devices_*.json`
  - `energygrid_report_*.txt`

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

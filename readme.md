# Bet Placement and Settlement HTTP API Service

This service handles bet placements and bet settlements using an in-memory cache to store user and event data.

## API Endpoints

### Health Check

- `GET /ping` - Service health check
  ```
  curl --location 'localhost:5001/ping' --data ''
  ```

### Account Management

- `GET /get-account` - Check user account balance
  ```
  curl --location --request GET 'localhost:5001/get-account' \
  --header 'Content-Type: application/json' \
  --data '{
      "user_id": "6209cf4b-92f0-4c0e-8f8e-e4a518cd2430"
  }'
  ```

### Bet Operations

- `POST /place-bet` - Place a new bet

  ```
  curl --location 'localhost:5001/place-bet' \
  --header 'Content-Type: application/json' \
  --data '{
      "user_id": "6209cf4b-92f0-4c0e-8f8e-e4a518cd2430",
      "event_id": "91bb5494-ca38-42e9-a20d-cfa9a07900e6",
      "odds": 150,
      "amount": 100
  }'
  ```

- `POST /settle-bets` - Settle all bets for an event
  ```
  curl --location 'localhost:5001/settle-bets' \
  --header 'Content-Type: application/json' \
  --data '{
      "event_id": "91bb5494-ca38-42e9-a20d-cfa9a07900e6",
      "result_status": "WIN"
  }'
  ```

## Test Data

### Events

| Event ID                             | Odds | Status   |
| ------------------------------------ | ---- | -------- |
| 91bb5494-ca38-42e9-a20d-cfa9a07900e6 | 150  | OPEN     |
| 5b1043d8-151f-441d-868b-d8227961d54f | 160  | OPEN     |
| 38c86515-3eb8-458b-9f11-5cb677dbbe6f | 170  | RESULTED |

### Users

| User ID                              | Balance |
| ------------------------------------ | ------- |
| 6209cf4b-92f0-4c0e-8f8e-e4a518cd2430 | 100000  |
| 774ab462-b536-458b-ab46-30903f45001f | 0       |

## Assumptions

- Service starts with pre-populated user and event data
- Users can place unlimited bets
- Bets are only placed for events with "OPEN" status
- All amounts are in cents (integer values)

## Setup

1. Clone the repository
2. Run the service:
   ```
   make run
   ```
   - Default HTTP port: `5001` (configurable via `HTTP_PORT` environment variable)

## Testing

1. Verify service is running:

   ```
   curl --location 'localhost:5001/ping' --data ''
   ```

   (Should return HTTP 200)

2. Place a bet using `/place-bet` endpoint

   - Successful bets will have "UNRESULTED" status
   - Account balance will be deducted immediately

3. Check account balance with `/get-account`

4. Settle bets with `/settle-bets`. All associated bets will be updated with the result.

## Notes

- This service uses in-memory storage - all data will be lost when the service stops

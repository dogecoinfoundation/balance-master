# Balance Master
This is a program that follows the Dogecoin core and tracks UTXOs that you have asked it to track.

## APIs
- POST /trackers - Post an address so the UTXOs for this address will be tracked.
- GET /utxos?address=<myaddress> - Retrieve spendable UTXOs for an address

## Config
- "doge-schema", "http", "Dogecoin schema"
- "doge-host", "localhost", "Dogecoin host"
- "doge-port", "22556", "Dogecoin port"
- "doge-user", "test", "Dogecoin user"
- "doge-password", "test", "Dogecoin password"
- "database-url", "sqlite://balance-master.db", "Database URL"
- "migrations-path", "file://../db/migrations", "Migrations path"
- "rpc-server-host", "localhost", "RPC server host"
- "rpc-server-port", "8899", "RPC server port"

## Running
`go run cmd/balance-master/main.go`
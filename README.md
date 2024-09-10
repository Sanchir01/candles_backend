# Backend app for online shop candles online

## Build & Run (Locally)
### Prerequisites
- go 1.23
- docker
- air
- golangci-lint (optional)

Create .env file in root directory and add follow values:

```dotenv
JWT_SECRET="sdw@#!@#Fxd"
CONFIG_PATH="./config/config.yaml"
PASSWORD_POSTGRES='postgres'
STORAGE_PATH="postgresql://postgres:postgres@localhost:5432/postgres?ssl=disabled"
```

use `air` and open browser localhost:5000


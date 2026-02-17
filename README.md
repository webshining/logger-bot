# Logging bot

## Getting started

### Download binary

\#:/> `curl -L -O https://github.com/webshining/logger-bot/releases/latest/download/app`

\#:/> `chmod +x ./app`

### Environment variables

| Variable name   | variable type | settings                  |
| --------------- | ------------- | ------------------------- |
| BOT_TOKEN       | string        | required                  |
| BOT_ADMIN       | int64         | required                  |
|                 |               |                           |
| DATABASE_DRIVER | string        | default(sqlite)           |
| DATABASE_URL    | string        | default(database.sqlite3) |

### Application start (local)

\#:/> `./app`

# SQLServe

A web server where SQLite is the source of truth. Routes, templates, and static files are all stored in SQLite, with a minimal Go server handling the HTTP layer.

## Features

- Routes defined as SQL rows
- Templates stored in SQLite
- Static files as BLOBs
- SQL-based routing and templating
- Frontend generated from tables
- Minimal Go server (50 lines)

## Setup

1. Install Go 1.21 or later
2. Install SQLite3
3. Clone this repository
4. Run:
   ```bash
   go mod download
   go run main.go
   ```

The server will start on port 8080 by default. Set the `PORT` environment variable to change this.

## Database Schema

The server uses the following tables:
- `routes`: HTTP routes and their handlers
- `templates`: HTML templates
- `static_files`: Static assets as BLOBs
- `functions`: SQL functions for handling requests

## Adding New Routes

To add a new route, insert a row into the `routes` table:

```sql
INSERT INTO routes (path, method, handler_function) 
VALUES ('/about', 'GET', 'handle_about');
```

Then create the corresponding handler function:

```sql
INSERT INTO functions (name, sql_code) 
VALUES ('handle_about', 'SELECT content FROM templates WHERE name = ''about'';');
```

## Adding Static Files

To add a static file:

```sql
INSERT INTO static_files (path, content, mime_type) 
VALUES ('/images/logo.png', ?, 'image/png');
```

## License

MIT 
-- Routes table to define HTTP endpoints
CREATE TABLE IF NOT EXISTS routes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    handler_function TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(path, method)
);

-- Templates table for storing HTML templates
CREATE TABLE IF NOT EXISTS templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Static files table for storing static assets
CREATE TABLE IF NOT EXISTS static_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL UNIQUE,
    content BLOB NOT NULL,
    mime_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Functions table for storing SQL functions
CREATE TABLE IF NOT EXISTS functions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    sql_code TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert some initial data
INSERT OR REPLACE INTO routes (path, method, handler_function) VALUES
    ('/', 'GET', 'handle_home'),
    ('/static/*', 'GET', 'handle_static'),
    ('/hello', 'GET', 'handle_hello');

-- Create a basic home page template
INSERT OR REPLACE INTO templates (name, content) VALUES
    ('home', '<!DOCTYPE html>
<html>
<head>
    <title>SQLServe</title>
</head>
<body>
    <h1>Welcome to SQLServe</h1>
    <p>A web server powered by SQLite!</p>
</body>
</html>'),
    ('hello', '<!DOCTYPE html>
<html>
<head>
    <title>Hello World - SQLServe</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 40px auto;
            padding: 20px;
            line-height: 1.6;
        }
        h1 { color: #2c3e50; }
    </style>
</head>
<body>
    <h1>Hello, World!</h1>
    <p>This page is being served directly from SQLite!</p>
    <p>Current time: ?</p>
</body>
</html>');

-- Create basic handler functions
INSERT OR REPLACE INTO functions (name, sql_code) VALUES
    ('handle_home', 'SELECT content FROM templates WHERE name = ''home'';'),
    ('handle_static', 'SELECT content, mime_type FROM static_files WHERE path = ?;'),
    ('handle_hello', 'SELECT replace(content, ''?'', datetime(''now'', ''localtime'')), ''text/html'' FROM templates WHERE name = ''hello'';'); 
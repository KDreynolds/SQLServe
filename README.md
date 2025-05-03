# SQLServe

A web server where SQLite is the source of truth. Because who needs a real programming language when you have SQL? 😏

## The Gimmick

- Routes? They're just rows in a table
- Templates? Stored in SQLite
- Static files? BLOBs, obviously
- Logic? Pure SQL, baby
- Frontend? Generated from tables
- The actual server? A 50-line Go shim that's just there to make HTTP requests look less suspicious

## The Magic

The server is actually the database file itself. No, really. It's a shell script glued to a SQLite database. When you run it:

```bash
./sqlserve.db
```

It:
1. Extracts itself (very meta)
2. Runs the database as a server
3. Makes your coworkers question their life choices

## Setup

1. Have Go and SQLite installed (or don't, we're not your parents)
2. Run:
   ```bash
   ./build.sh
   ```
3. Run the server:
   ```bash
   ./sqlserve.db
   ```

## How It Works (Oversimplified)

1. HTTP request comes in
2. Go asks SQLite: "What do?"
3. SQLite says: "Run this SQL"
4. Go runs SQL
5. Magic happens
6. Profit

## Adding Routes

Want a new route? Just add a row:

```sql
INSERT INTO routes (path, method, handler_function) 
VALUES ('/about', 'GET', 'handle_about');
```

Then write the handler in SQL:

```sql
INSERT INTO functions (name, sql_code) 
VALUES ('handle_about', 'SELECT content FROM templates WHERE name = ''about'';');
```

## Why?

- Because we can
- Because it makes people uncomfortable
- Because SQL is technically Turing complete
- Because Uncle Bob will probably hate it
- Because I do whatever I want

## License

MIT (because we're not monsters) 
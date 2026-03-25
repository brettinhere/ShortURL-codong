# codong-shorturl

A URL shortener service written in [Codong](https://codong.org) — 42 lines, no frameworks, no boilerplate.

**Live:** https://codong.org/short-url/

## The Code

```codong
// main.cod — the entire service

import { generate_id, is_valid_url } from "./lib/url_utils.cod"

db.connect("sqlite://shorturl.db")
db.exec("""
    CREATE TABLE IF NOT EXISTS urls (
        code TEXT PRIMARY KEY,
        long_url TEXT NOT NULL,
        hits INTEGER DEFAULT 0,
        created_at TEXT DEFAULT (datetime('now'))
    )
""")

server = web.serve(port: 8082)

server.post("/api/shorten", fn(req) {
    body = req.json()
    if !is_valid_url(body.url) {
        return {status: 400, body: {error: "invalid url"}}
    }
    code = generate_id(6)
    db.exec("INSERT INTO urls (code, long_url) VALUES (?, ?)", code, body.url)
    return {status: 200, body: {
        code: code,
        short_url: "https://codong.org/s/{code}",
    }}
})

server.get("/s/:code", fn(req) {
    row = db.find_one("SELECT long_url FROM urls WHERE code = ?", req.params.code)
    if row == null {
        return {status: 404, body: {error: "not found"}}
    }
    db.exec("UPDATE urls SET hits = hits + 1 WHERE code = ?", req.params.code)
    return {status: 301, headers: {"Location": row.long_url}}
})

server.get("/api/stats/:code", fn(req) {
    row = db.find_one("SELECT * FROM urls WHERE code = ?", req.params.code)
    if row == null {
        return {status: 404, body: {error: "not found"}}
    }
    return {status: 200, body: row}
})

server.listen()
```

## API

### Shorten a URL
```
POST /short-url/api/shorten
{"url": "https://example.com/very/long/url"}

→ {"code": "abc123", "short_url": "https://codong.org/s/abc123"}
```

### Redirect
```
GET /s/abc123  →  301 to original URL
```

### Stats
```
GET /short-url/api/stats/abc123

→ {"code": "abc123", "long_url": "...", "hits": 42, "created_at": "..."}
```

## Line Count Comparison

| Language | Lines |
|----------|-------|
| **Codong** | **42** |
| Python (Flask + SQLAlchemy) | ~90 |
| Go (stdlib) | ~110 |
| JavaScript (Express) | ~95 |

## License

MIT

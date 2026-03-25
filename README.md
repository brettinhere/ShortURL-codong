# ShortURL — Built with Codong

> A free, open-source URL shortener. Paste a long link, get a short one instantly. No account required.

**🔗 Live demo → [codong.org/short-url](https://codong.org/short-url/)**

![ShortURL screenshot](./screenshot.jpg)

---

## Use it now — no setup needed

Go to **[https://codong.org/short-url/](https://codong.org/short-url/)**, paste your URL, click **Shorten**. Done.

- ✅ Instant short link — `codong.org/s/xxxxxx`
- ✅ 301 redirect — fast, SEO-friendly
- ✅ Click tracking — see how many times your link was visited
- ✅ No login, no account, no rate limit
- ✅ English / 中文 — switch languages in one click

---

## Self-host in 3 steps

```bash
git clone https://github.com/brettinhere/ShortURL-codong
cd ShortURL-codong/cmd
go mod tidy && go build -o shorturl .
./shorturl
```

That's it. Service starts on `:8082`. SQLite database is created automatically.

### API

```bash
# Shorten a URL
curl -X POST http://localhost:8082/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://your-long-url.com/very/long/path"}'
# → {"code":"abc123","short_url":"https://codong.org/s/abc123"}

# Check stats
curl http://localhost:8082/api/stats/abc123
# → {"code":"abc123","long_url":"...","hits":42,"created_at":"2026-03-25 07:12:46"}
```

---

## SEO built-in

This project ships with a complete SEO setup out of the box — so if you self-host, your page is immediately search engine and AI crawler ready:

| Feature | Included |
|---|---|
| `<meta>` description, keywords, canonical | ✓ |
| Open Graph tags (social share preview) | ✓ |
| Twitter Card (large image) | ✓ |
| Schema.org `WebApplication` structured data | ✓ |
| Schema.org `SoftwareSourceCode` structured data | ✓ |
| Schema.org `FAQPage` (AI answer boxes) | ✓ |
| `robots.txt` with GPTBot, Claude, Perplexity allowed | ✓ |
| `sitemap.xml` with hreflang EN/ZH | ✓ |
| OG image (`og-image.svg`) | ✓ |

The FAQPage schema means your site can appear directly in AI-powered search results (ChatGPT search, Perplexity, Google AI Overviews) as a structured answer — no extra work needed.

---

## Why Codong?

The entire backend service is **42 lines of Codong**. No frameworks. No ORM. No package manager.

```codong
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

Compare that to other languages:

| Language | Lines | Dependencies |
|---|---|---|
| **Codong** | **42** | **0** |
| Python (Flask + SQLAlchemy) | ~90 | 2 |
| Go (stdlib) | ~110 | 1 |
| JavaScript (Express) | ~95 | 2 |

Codong has `db` and `web` built in. You write the logic. The language handles the rest.

→ **Learn more about Codong: [codong.org](https://codong.org)**

---

## License

MIT

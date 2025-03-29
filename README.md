# 🧠 Telepath

**Telepath** is a lightweight, general-purpose reverse proxy that dynamically routes incoming HTTP requests to different targets based on either a custom header or the first segment of the request path. It's ideal for use cases like Salesforce or webhook integrations, where a single public endpoint needs to forward to different local or internal environments.

---

## 🚀 Features

- 🧭 Route by custom HTTP header or first path segment
- 🔁 Hot-reload configuration via `SIGHUP` (no restarts!)
- 🧰 CLI for adding, removing, and listing routes
- 🔐 Built with Go – fast, lightweight, and deployable as a single binary

---


---

## ⚙️ Configuration (`telepath.json`)

```json
{
  "mode": "header",                // or "path"
  "headerName": "X-Runtime-Env",  // required if mode is "header"
  "routes": {
    "qa": "http://internal.qa:8080",
    "john": "https://john-ngrok.io"
  }
}
```


Routing Modes
header: Uses the value of the specified header to determine the route.

path: Uses the first path segment (e.g. /john/photos) and strips it before proxying.

---
Running the Server
```bash
go run telepath.go
```

This will start the proxy on `:8080`
To reload config without restarting:
```bash
pkill -SIGNUP telepath
```

---
💻 CLI Usage

Add a Route
```bash
telepath route add john https://john-ngrok.io
```

Remove a Route
```bash
telepath route remove john
```

List all Routes
```bash
telepath route list
```

Refresh Config (sends SIGHUP to running server)
```bash
telepath refresh
```


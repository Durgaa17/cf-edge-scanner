# CF Edge Scanner

Fast concurrent **ping scanner** for the **official Cloudflare IPv4 ranges**.

When you start the scanner you can choose **which ranges** to scan (one, several, or all).

Works on:
- **Termux** (Android)
- Linux
- Windows 10/11
- macOS / Chromebook (Linux)

---

## Available Versions

| Version     | File       | Best For                        | Speed     |
|-------------|------------|---------------------------------|-----------|
| **Go**      | `scan.go`  | Laptop / Chromebook / Desktop   | Very Fast |
| **Python**  | `scan.py`  | Easy & Termux                   | Good      |

**Recommendation:** Use the **Go** version on laptop or Chromebook — it is faster and lighter.

---

## Quick Start

### Go Version (Recommended)

#### Termux
```bash
pkg update && pkg upgrade -y
pkg install golang git -y
git clone https://github.com/Durgaa17/cf-edge-scanner.git
cd cf-edge-scanner
go run scan.go
```

#### Linux / macOS / Chromebook
```bash
git clone https://github.com/Durgaa17/cf-edge-scanner.git
cd cf-edge-scanner
go run scan.go
```

#### Build a single binary (optional)
```bash
go build -o cf-edge-scanner scan.go
./cf-edge-scanner
```

#### Windows
```powershell
git clone https://github.com/Durgaa17/cf-edge-scanner.git
cd cf-edge-scanner
go run scan.go
```

---

### Python Version

#### Termux
```bash
pkg update && pkg upgrade -y
pkg install python git -y
git clone https://github.com/Durgaa17/cf-edge-scanner.git
cd cf-edge-scanner
python scan.py
```

#### Linux / macOS / Windows
```bash
git clone https://github.com/Durgaa17/cf-edge-scanner.git
cd cf-edge-scanner
python3 scan.py          # or python on Windows
```

---

## How it works

1. Shows the list of official Cloudflare ranges with approximate host counts.
2. You select the range(s) you want (e.g. `3 7 12` or `0` for all).
3. Concurrent ping scan runs.
4. Alive IPs are printed live and saved to a timestamped file (`reachable_YYYYMMDD_HHMMSS.txt`).

---

## Official Ranges Included

```
173.245.48.0/20
103.21.244.0/22
103.22.200.0/22
103.31.4.0/22
141.101.64.0/18
108.162.192.0/18
190.93.240.0/20
188.114.96.0/20
197.234.240.0/22
198.41.128.0/17
162.158.0.0/15
104.16.0.0/13
104.24.0.0/14
172.64.0.0/13
131.0.72.0/22
```

Source: https://www.cloudflare.com/ips-v4/

---

## Configuration

### Go (`scan.go`)
```go
Workers = 128              // Concurrent pings (lower on Termux if needed)
Timeout = 800 * time.Millisecond
```

### Python (`scan.py`)
```python
THREADS = 64      # Lower this on Termux / weak devices (try 16–32)
TIMEOUT = 0.8     # Ping timeout in seconds
```

---

## Notes

- Large ranges (`104.16.0.0/13`, `172.64.0.0/13`, `162.158.0.0/15`) contain hundreds of thousands of IPs. Scanning all of them takes a long time.
- On Termux you may need to grant network permission and keep the screen on.
- Go version is significantly faster and uses less memory than Python.

---

## Related

Based on the multi-language scanner: [cf-ip-scanner](https://github.com/Durgaa17/cf-ip-scanner)

---

## License

Free to use and modify.

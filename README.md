# CF Edge Scanner — `ping-test` branch

**Low memory / Sequential version**

This branch uses a **sequential** approach:
- Ping one IP → show result if alive → move to next IP
- Very low memory usage (good for laptops)
- No big concurrent workers

Perfect when the main branch uses too much RAM or feels slow on your laptop.

---

## How to use this branch

### Windows 11

```powershell
git clone -b ping-test https://github.com/Durgaa17/cf-edge-scanner.git
cd cf-edge-scanner
go run scan.go
```

Or if you already cloned the repo:

```powershell
git checkout ping-test
go run scan.go
```

### Linux / macOS / Chromebook / Termux

```bash
git clone -b ping-test https://github.com/Durgaa17/cf-edge-scanner.git
cd cf-edge-scanner
go run scan.go
```

---

## Difference from main branch

| Feature              | main branch          | ping-test branch     |
|----------------------|----------------------|----------------------|
| Mode                 | Concurrent (fast)    | Sequential (1 by 1)  |
| Memory usage         | Higher               | **Very low**         |
| Speed                | Faster               | Slower but stable    |
| Best for             | Powerful PC          | Laptop / low RAM     |

---

## Official Ranges Included

Same as main branch (the 15 official Cloudflare ranges).

You can still choose which ranges to scan from the menu.

---

## Notes

- This version is intentionally simple and light.
- It will take longer than the concurrent version, but it should not freeze or use high memory.
- Results are still saved to `reachable_YYYYMMDD_HHMMSS.txt`.

---

## License

Free to use and modify.

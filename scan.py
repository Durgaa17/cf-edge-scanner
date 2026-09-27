#!/usr/bin/env python3
"""
CF Edge Scanner - Concurrent ping scanner for official Cloudflare ranges
Works on Termux (Android), Linux, Windows, and macOS
"""

import subprocess
import platform
from concurrent.futures import ThreadPoolExecutor, as_completed
import ipaddress
import sys
import os
from datetime import datetime

# ====================== CONFIG ======================
THREADS = 64                 # Concurrent pings (lower on Termux if needed)
TIMEOUT = 0.8                # Timeout in seconds
# ====================================================

# Official Cloudflare IPv4 ranges (from https://www.cloudflare.com/ips-v4/)
CF_RANGES = [
    "173.245.48.0/20",
    "103.21.244.0/22",
    "103.22.200.0/22",
    "103.31.4.0/22",
    "141.101.64.0/18",
    "108.162.192.0/18",
    "190.93.240.0/20",
    "188.114.96.0/20",
    "197.234.240.0/22",
    "198.41.128.0/17",
    "162.158.0.0/15",
    "104.16.0.0/13",
    "104.24.0.0/14",
    "172.64.0.0/13",
    "131.0.72.0/22",
]


def clear_screen():
    os.system("cls" if platform.system().lower() == "windows" else "clear")


def is_alive(ip: str) -> tuple:
    system = platform.system().lower()

    if system == "windows":
        cmd = ["ping", "-n", "1", "-w", str(int(TIMEOUT * 1000)), ip]
    else:
        # Linux / Termux / macOS
        cmd = ["ping", "-c", "1", "-W", str(int(TIMEOUT)), ip]

    try:
        result = subprocess.run(
            cmd,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            timeout=TIMEOUT + 1.5
        )
        return ip, result.returncode == 0
    except Exception:
        return ip, False


def expand_ranges(selected_cidrs: list) -> list:
    ips = []
    for cidr in selected_cidrs:
        try:
            net = ipaddress.ip_network(cidr, strict=False)
            # Use .hosts() for /24 and larger; for very large ranges we still expand
            ips.extend(str(ip) for ip in net.hosts())
        except ValueError as e:
            print(f"[!] Skipping invalid CIDR {cidr}: {e}")
    return ips


def show_menu() -> list:
    clear_screen()
    print("=" * 55)
    print("   CF Edge Scanner  |  Official Cloudflare Ranges")
    print("=" * 55)
    print()

    for i, cidr in enumerate(CF_RANGES, 1):
        try:
            net = ipaddress.ip_network(cidr, strict=False)
            hosts = net.num_addresses - 2 if net.num_addresses > 2 else net.num_addresses
            print(f"  {i:2d}. {cidr:<18}  (~{hosts:,} hosts)")
        except Exception:
            print(f"  {i:2d}. {cidr}")

    print()
    print("  0.  Scan ALL ranges (warning: very large)")
    print("  q.  Quit")
    print()
    print("=" * 55)

    while True:
        choice = input("Select range number(s) separated by space (e.g. 1 5 12) or 0 for all: ").strip().lower()

        if choice in ("q", "quit", "exit"):
            print("Bye!")
            sys.exit(0)

        if choice == "0":
            confirm = input("You selected ALL ranges (millions of IPs). Continue? [y/N]: ").strip().lower()
            if confirm == "y":
                return CF_RANGES[:]
            continue

        try:
            indices = [int(x) for x in choice.split()]
            selected = []
            for idx in indices:
                if 1 <= idx <= len(CF_RANGES):
                    selected.append(CF_RANGES[idx - 1])
                else:
                    print(f"[!] Invalid number: {idx}")
                    selected = []
                    break
            if selected:
                return selected
        except ValueError:
            print("[!] Please enter numbers only (example: 3 7 11)")


def main():
    selected = show_menu()

    print()
    print(f"Selected ranges ({len(selected)}):")
    for c in selected:
        print(f"  • {c}")
    print()

    print("Expanding CIDRs... this may take a moment for large ranges.")
    ips = expand_ranges(selected)
    total = len(ips)

    if total == 0:
        print("No IPs to scan.")
        return

    print(f"Total hosts to scan: {total:,}")
    print(f"Threads: {THREADS}  |  Timeout: {TIMEOUT}s")
    print("Starting scan...\n")

    alive = []
    scanned = 0

    with ThreadPoolExecutor(max_workers=THREADS) as executor:
        futures = {executor.submit(is_alive, ip): ip for ip in ips}

        for future in as_completed(futures):
            scanned += 1
            ip, ok = future.result()
            if ok:
                alive.append(ip)
                print(f"[+] {ip}")

            # Progress every 500 IPs
            if scanned % 500 == 0 or scanned == total:
                print(f"    ... {scanned}/{total} scanned  |  {len(alive)} alive", end="\r")

    print()

    # Save results
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    filename = f"reachable_{timestamp}.txt"

    with open(filename, "w") as f:
        f.write("\n".join(sorted(alive, key=lambda x: ipaddress.ip_address(x))))

    print("\n" + "=" * 50)
    print(f"Done! {len(alive):,}/{total:,} hosts are reachable")
    print(f"Results saved to: {filename}")
    print("=" * 50)

    if alive:
        print("\nReachable IPs (first 30):")
        for ip in sorted(alive, key=lambda x: ipaddress.ip_address(x))[:30]:
            print(f"  {ip}")
        if len(alive) > 30:
            print(f"  ... and {len(alive) - 30} more (see file)")
    else:
        print("\nNo reachable hosts found.")


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n[!] Interrupted by user. Partial results may not be saved.")
        sys.exit(1)

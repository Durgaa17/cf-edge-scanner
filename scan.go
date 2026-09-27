package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ====================== CONFIG ======================
const (
	Timeout = 600 * time.Millisecond // Ping timeout
)

// ====================================================

// Official Cloudflare IPv4 ranges
var cfRanges = []string{
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
}

func isAlive(ip string) bool {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", fmt.Sprintf("%d", Timeout.Milliseconds()), ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", fmt.Sprintf("%d", int(Timeout.Seconds())), ip)
	}

	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

func showMenu() []string {
	fmt.Println(strings.Repeat("=", 55))
	fmt.Println("   CF Edge Scanner (Low Memory)  |  Sequential Ping")
	fmt.Println(strings.Repeat("=", 55))
	fmt.Println()

	for i, cidr := range cfRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			fmt.Printf("  %2d. %s\n", i+1, cidr)
			continue
		}
		ones, bits := ipNet.Mask.Size()
		hosts := 1 << (bits - ones)
		if hosts > 2 {
			hosts -= 2
		}
		fmt.Printf("  %2d. %-18s  (~%s hosts)\n", i+1, cidr, formatNumber(hosts))
	}

	fmt.Println()
	fmt.Println("  0.  Scan ALL ranges (warning: very large)")
	fmt.Println("  q.  Quit")
	fmt.Println()
	fmt.Println(strings.Repeat("=", 55))

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Select range number(s) separated by space (e.g. 1 5 12) or 0 for all: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "q" || input == "quit" || input == "exit" {
			fmt.Println("Bye!")
			os.Exit(0)
		}

		if input == "0" {
			fmt.Print("You selected ALL ranges. Continue? [y/N]: ")
			confirm, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(confirm)) == "y" {
				return append([]string{}, cfRanges...)
			}
			continue
		}

		parts := strings.Fields(input)
		var selected []string
		valid := true

		for _, p := range parts {
			idx, err := strconv.Atoi(p)
			if err != nil || idx < 1 || idx > len(cfRanges) {
				fmt.Printf("[!] Invalid number: %s\n", p)
				valid = false
				break
			}
			selected = append(selected, cfRanges[idx-1])
		}

		if valid && len(selected) > 0 {
			return selected
		}
	}
}

func formatNumber(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	var result []byte
	for i, c := range reverse(s) {
		if i > 0 && i%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return reverse(string(result))
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main() {
	selected := showMenu()

	fmt.Println()
	fmt.Printf("Selected ranges (%d):\n", len(selected))
	for _, c := range selected {
		fmt.Printf("  • %s\n", c)
	}
	fmt.Println()
	fmt.Println("Mode: Sequential (1 IP at a time) — Low memory")
	fmt.Println("Starting scan...\n")

	var alive []string
	scanned := 0

	for _, cidr := range selected {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			fmt.Printf("[!] Skipping %s: %v\n", cidr, err)
			continue
		}

		fmt.Printf("--- Scanning %s ---\n", cidr)

		for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
			if ip.Equal(ipNet.IP) || isBroadcast(ip, ipNet) {
				continue
			}

			ipStr := ip.String()
			scanned++

			if isAlive(ipStr) {
				alive = append(alive, ipStr)
				fmt.Printf("[+] %s\n", ipStr)
			}

			// Progress every 200 IPs
			if scanned%200 == 0 {
				fmt.Printf("    ... %d scanned  |  %d alive\r", scanned, len(alive))
			}
		}
		fmt.Println()
	}

	// Save results
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("reachable_%s.txt", timestamp)

	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	for _, ip := range alive {
		f.WriteString(ip + "\n")
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Done! %d alive out of %d scanned\n", len(alive), scanned)
	fmt.Printf("Results saved to: %s\n", filename)
	fmt.Println(strings.Repeat("=", 50))

	if len(alive) > 0 {
		fmt.Println("\nReachable IPs (first 30):")
		limit := 30
		if len(alive) < limit {
			limit = len(alive)
		}
		for i := 0; i < limit; i++ {
			fmt.Printf("  %s\n", alive[i])
		}
		if len(alive) > 30 {
			fmt.Printf("  ... and %d more (see file)\n", len(alive)-30)
		}
	} else {
		fmt.Println("\nNo reachable hosts found.")
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func isBroadcast(ip net.IP, network *net.IPNet) bool {
	broadcast := make(net.IP, len(network.IP))
	for i := range network.IP {
		broadcast[i] = network.IP[i] | ^network.Mask[i]
	}
	return ip.Equal(broadcast)
}

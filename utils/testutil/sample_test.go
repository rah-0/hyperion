package testutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestX(t *testing.T) {
	data := `BenchmarkEncodeBytes1Small-8            331346038              218.2 ns/op            24 B/op          1 allocs/op
BenchmarkEncodeBytes100Small-8          12789297              5647 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeBytes10000Small-8          126255            569372 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeBytes1000000Small-8          1198          60117976 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeBytes1Unreal-8           45447225              1493 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeBytes10Unreals-8          5216320             13784 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeBytes100Unreals-8          522309            139595 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeBytes1000Unreals-8          48220           1513438 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeJson1Small-8             302926762              235.1 ns/op            24 B/op          1 allocs/op
BenchmarkEncodeJson100Small-8            5796138             12506 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeJson10000Small-8            56889           1268709 ns/op              97 B/op          1 allocs/op
BenchmarkEncodeJson1000000Small-8            561         128433139 ns/op              26 B/op          1 allocs/op
BenchmarkEncodeJson1Unreal-8            15375618              4718 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeJson10Unreals-8           1542039             46899 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeJson100Unreals-8           150490            476419 ns/op              30 B/op          1 allocs/op
BenchmarkEncodeJson1000Unreals-8           14821           4888348 ns/op              24 B/op          1 allocs/op
BenchmarkEncodeProto1Small-8            675272739              106.9 ns/op            80 B/op          1 allocs/op
BenchmarkEncodeProto100Small-8           6710169             10709 ns/op            8000 B/op        100 allocs/op
BenchmarkEncodeProto10000Small-8           66267           1080003 ns/op          800000 B/op      10000 allocs/op
BenchmarkEncodeProto1000000Small-8           658         110364613 ns/op        80000002 B/op    1000000 allocs/op
BenchmarkEncodeProto1Unreal-8           26083810              2767 ns/op            4096 B/op          1 allocs/op
BenchmarkEncodeProto10Unreals-8          2594926             27592 ns/op           40960 B/op         10 allocs/op
BenchmarkEncodeProto100Unreals-8          257083            277134 ns/op          409600 B/op        100 allocs/op
BenchmarkEncodeProto1000Unreals-8          25994           2790997 ns/op         4096000 B/op       1000 allocs/op`

	fmt.Println(ParseBenchOutput(data))
}

// ParseBenchOutput parses raw benchmark data and converts it into a markdown table with human-readable units.
func ParseBenchOutput(raw string) string {
	// Regular expression to match benchmark output
	re := regexp.MustCompile(`(?m)^(Benchmark\w+)-\d+\s+(\d+)\s+([\d.]+) ns/op\s+([\d.]+) B/op\s+(\d+) allocs/op$`)

	// Header for the markdown table
	table := "| Benchmark                 | Iterations | Time per Op       | Memory per Op  | Allocations per Op |\n"
	table += "|---------------------------|------------|-------------------|----------------|--------------------|\n"

	// Iterate over each match
	matches := re.FindAllStringSubmatch(raw, -1)
	for _, match := range matches {
		// Extract fields
		benchmark := strings.TrimPrefix(match[1], "Benchmark")
		iterations, _ := strconv.Atoi(match[2])
		timePerOp := humanizeTime(match[3])
		memoryPerOp := humanizeMemory(match[4])
		allocsPerOp := match[5]

		// Add row to the table
		table += fmt.Sprintf("| %-25s | %10d | %-17s | %-14s | %-18s |\n",
			benchmark, iterations, timePerOp, memoryPerOp, allocsPerOp)
	}

	return table
}

// humanizeTime converts nanoseconds into a more human-readable time unit.
func humanizeTime(ns string) string {
	val, _ := strconv.ParseFloat(ns, 64)
	if val < 1e3 {
		return fmt.Sprintf("%.2f ns", val)
	} else if val < 1e6 {
		return fmt.Sprintf("%.2f µs", val/1e3)
	} else if val < 1e9 {
		return fmt.Sprintf("%.2f ms", val/1e6)
	}
	return fmt.Sprintf("%.2f s", val/1e9)
}

// humanizeMemory converts bytes into a more human-readable memory unit.
func humanizeMemory(bytes string) string {
	val, _ := strconv.ParseFloat(bytes, 64)
	if val < 1e3 {
		return fmt.Sprintf("%.2f B", val)
	} else if val < 1e6 {
		return fmt.Sprintf("%.2f KB", val/1e3)
	} else if val < 1e9 {
		return fmt.Sprintf("%.2f MB", val/1e6)
	}
	return fmt.Sprintf("%.2f GB", val/1e9)
}

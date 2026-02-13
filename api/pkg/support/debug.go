package support

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Dump prints variables with type information for debugging
func Dump(values ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	// Print location
	color.New(color.FgHiBlack).Printf("// %s:%d in %s\n", file, line, fn.Name())

	for i, v := range values {
		dumpValue(v, i)
	}
	fmt.Println()
}

// DD dumps and dies (prints and exits)
func DD(values ...interface{}) {
	Dump(values...)
	os.Exit(1)
}

// DumpJSON prints variables as formatted JSON
func DumpJSON(values ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	color.New(color.FgHiBlack).Printf("// %s:%d in %s\n", file, line, fn.Name())

	for _, v := range values {
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			color.Red("Error marshaling: %v", err)
			continue
		}
		fmt.Println(string(data))
	}
	fmt.Println()
}

func dumpValue(v interface{}, index int) {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	typeStr := "nil"
	if t != nil {
		typeStr = t.String()
	}

	color.New(color.FgCyan).Printf("[%d] ", index)
	color.New(color.FgYellow).Printf("(%s) ", typeStr)

	if v == nil {
		color.New(color.FgHiBlack).Println("nil")
		return
	}

	switch val.Kind() {
	case reflect.Struct:
		dumpStruct(val, 0)
	case reflect.Map:
		dumpMap(val, 0)
	case reflect.Slice, reflect.Array:
		dumpSlice(val, 0)
	case reflect.Ptr:
		if val.IsNil() {
			color.New(color.FgHiBlack).Println("nil")
		} else {
			fmt.Print("&")
			dumpValue(val.Elem().Interface(), index)
		}
	default:
		fmt.Printf("%+v\n", v)
	}
}

func dumpStruct(val reflect.Value, indent int) {
	t := val.Type()
	prefix := strings.Repeat("  ", indent)

	fmt.Println("{")
	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanInterface() {
			continue
		}

		color.New(color.FgGreen).Printf("%s  %s: ", prefix, field.Name)

		switch fieldVal.Kind() {
		case reflect.Struct:
			if field.Type == reflect.TypeOf(time.Time{}) {
				fmt.Printf("%v\n", fieldVal.Interface())
			} else {
				dumpStruct(fieldVal, indent+1)
			}
		case reflect.Map:
			dumpMap(fieldVal, indent+1)
		case reflect.Slice, reflect.Array:
			dumpSlice(fieldVal, indent+1)
		default:
			fmt.Printf("%+v\n", fieldVal.Interface())
		}
	}
	fmt.Printf("%s}\n", prefix)
}

func dumpMap(val reflect.Value, indent int) {
	prefix := strings.Repeat("  ", indent)

	if val.Len() == 0 {
		fmt.Println("{}")
		return
	}

	fmt.Println("{")
	for _, key := range val.MapKeys() {
		color.New(color.FgGreen).Printf("%s  %v: ", prefix, key.Interface())
		mapVal := val.MapIndex(key)
		fmt.Printf("%+v\n", mapVal.Interface())
	}
	fmt.Printf("%s}\n", prefix)
}

func dumpSlice(val reflect.Value, indent int) {
	prefix := strings.Repeat("  ", indent)

	if val.Len() == 0 {
		fmt.Println("[]")
		return
	}

	fmt.Println("[")
	for i := 0; i < val.Len(); i++ {
		fmt.Printf("%s  [%d] ", prefix, i)
		elem := val.Index(i)
		if elem.Kind() == reflect.Struct {
			dumpStruct(elem, indent+1)
		} else {
			fmt.Printf("%+v\n", elem.Interface())
		}
	}
	fmt.Printf("%s]\n", prefix)
}

// Timer provides a simple timing utility
type Timer struct {
	name  string
	start time.Time
}

// StartTimer starts a new timer with a name
func StartTimer(name string) *Timer {
	return &Timer{
		name:  name,
		start: time.Now(),
	}
}

// Stop stops the timer and prints the elapsed time
func (t *Timer) Stop() time.Duration {
	elapsed := time.Since(t.start)
	color.New(color.FgCyan).Printf("⏱ %s: %v\n", t.name, elapsed)
	return elapsed
}

// Lap prints the current elapsed time without stopping
func (t *Timer) Lap(label string) time.Duration {
	elapsed := time.Since(t.start)
	color.New(color.FgCyan).Printf("⏱ %s [%s]: %v\n", t.name, label, elapsed)
	return elapsed
}

// Benchmark runs a function n times and reports statistics
func Benchmark(name string, n int, fn func()) {
	var total time.Duration
	var min, max time.Duration

	for i := 0; i < n; i++ {
		start := time.Now()
		fn()
		elapsed := time.Since(start)

		total += elapsed
		if i == 0 || elapsed < min {
			min = elapsed
		}
		if elapsed > max {
			max = elapsed
		}
	}

	avg := total / time.Duration(n)

	color.New(color.FgCyan).Printf("📊 Benchmark: %s (%d iterations)\n", name, n)
	fmt.Printf("   Total: %v\n", total)
	fmt.Printf("   Avg:   %v\n", avg)
	fmt.Printf("   Min:   %v\n", min)
	fmt.Printf("   Max:   %v\n", max)
}

// MemUsage returns current memory usage information
func MemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return fmt.Sprintf("Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v",
		bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}

// PrintMemUsage prints current memory usage
func PrintMemUsage() {
	color.New(color.FgMagenta).Printf("💾 %s\n", MemUsage())
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// Stack prints the current stack trace
func Stack() {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	color.New(color.FgHiBlack).Printf("Stack trace:\n%s\n", buf[:n])
}

// Trace logs function entry and exit with timing
func Trace(name string) func() {
	start := time.Now()
	color.New(color.FgBlue).Printf("→ Entering %s\n", name)

	return func() {
		color.New(color.FgBlue).Printf("← Exiting %s (%v)\n", name, time.Since(start))
	}
}

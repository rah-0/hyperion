package profiler

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/rah-0/nabu"
)

// Start initializes and starts the pprof HTTP server
func Start(ip string, port int) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	mux := http.NewServeMux()

	// Register pprof handlers with our mux
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Start HTTP server in a goroutine
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			nabu.FromError(err).WithMessage("Error starting profiler").Log()
			return
		}
	}()

	nabu.FromMessage("Profiler started on " + addr).Log()
}

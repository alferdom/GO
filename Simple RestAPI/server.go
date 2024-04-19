package main

import (
	handlers "Simple_RestAPI/Handlers"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Parse Arguments
//
// -p to define port number on which server will listen
//
// -f to define HTML template filepPath
//
// returns tuple of pointers to parsed strings
func parseArgs() (*string, *string) {
	flag.Usage = func() {
		fmt.Println("Supported Flags:")
		flag.PrintDefaults()
	}
	portNumber := flag.Uint64("p", 8000, "Define port number on which server will listen")
	templPath := flag.String("f", "./threat.html.tmpl", "Define HTML template filepPath")
	flag.Parse()
	unparsedArgs := flag.Args()
	if len(unparsedArgs) != 0 {
		fmt.Println("Ignoring non-flag arguments:", unparsedArgs)
	}
	portNumStr := strconv.FormatUint(*portNumber, 10)
	return &portNumStr, templPath
}

// Parse HTML Template from file path
//
// returns parsed template
func parseTemplate(aTemplPath *string) *template.Template {
	tmpl, err := template.ParseFiles(*aTemplPath)
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}
	log.Printf("Parsed template file %q successfully\n", tmpl.Name())
	return tmpl
}

// Start HTTP server on local address with defined portnumber and handlers
//
// Server closes gracefully on SIGINT || SIGTERM signals and defers sync.WaitGroup
func startServer(aPort *string, aHandler *handlers.Handler, aWaitGroup *sync.WaitGroup) {
	server := &http.Server{Addr: ":" + *aPort, Handler: handlers.NewRouter(*aHandler)}
	log.Println("Server listening on port", *aPort, "...")

	go func() {
		defer aWaitGroup.Done()

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("Error while listening:", err.Error())
		}
		log.Println("Stopped listening and serving.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // wait for receiving signal
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}

func main() {
	port, templPath := parseArgs()
	template := parseTemplate(templPath)
	handler := handlers.NewHandler(template)
	snc := &sync.WaitGroup{}
	snc.Add(1)
	startServer(port, &handler, snc)
	snc.Wait()
}

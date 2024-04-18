package main

import (
	handlers "Simple_RestAPI/Handlers"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
	portNumber := flag.Uint64("p", 10533, "Define port number on which server will listen")
	templPath := flag.String("f", "./threat.html.tmpl", "Define HTML template filepPath")
	flag.Parse()
	unparsedArgs := flag.Args()
	if len(unparsedArgs) != 0 {
		fmt.Println("Ignoring non-flag arguments:", unparsedArgs)
	}
	portNumStr := strconv.FormatUint(*portNumber, 10)
	return &portNumStr, templPath
}

func parseTemplate(templPath *string) *template.Template {
	tmpl, err := template.ParseFiles(*templPath)
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}
	log.Printf("Parsed %q successfully\n", tmpl.Name())
	return tmpl
}

func main() {
	port, templPath := parseArgs()
	fmt.Println("Port number:", *port, "Template path:", *templPath)
	template := parseTemplate(templPath)
	handler := handlers.NewHandler(template)
	log.Printf("Server listening on port %s ...\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, handlers.NewRouter(handler)))
}

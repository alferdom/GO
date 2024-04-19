package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const submitPage = `<!DOCTYPE html>
<html>
<body>
<form action="/render" method="POST">
<label for="json_input">JSON input:</label><br/>
<textarea rows="4" cols="50" id="json_input"
			 name="json_input">
</textarea><br/>
<input type="submit" value="Submit"/>
</form>
</body>
</html>
`

type Variant struct {
	Name      string `json:"name"`
	DateAdded string `json:"dateAdded"`
}
type ThreatData struct {
	ThreatName    string    `json:"threatName"`
	Category      string    `json:"category"`
	Size          int64     `json:"size"`
	DetectionDate string    `json:"detectionDate"`
	Variants      []Variant `json:"variants"`
}

type handler struct {
	tmpl *template.Template
}

type Handler interface {
	HandlerRootGet(w http.ResponseWriter, r *http.Request)
	HandlerRenderPost(w http.ResponseWriter, r *http.Request)
}

// Init new Handler with template member
func NewHandler(aTmpl *template.Template) Handler {
	return &handler{tmpl: aTmpl}
}

// Init new ServeMux router
func NewRouter(aHandler Handler) *http.ServeMux {
	mx := http.NewServeMux()
	mx.HandleFunc("/", aHandler.HandlerRootGet)
	mx.HandleFunc("/render", aHandler.HandlerRenderPost)
	return mx
}

// Handle GET request to root "/"
//
// Respond with simple HTML web in format `text/html; charset=UTF-8`
func (h *handler) HandlerRootGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		writeStatus(w, r, http.StatusBadRequest, "400 Bad Request")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, submitPage)
}

// Handle POST request to "/render"
//
// Decode JSON from request Body and respond HTML temple filled with JSON decoded data
func (h *handler) HandlerRenderPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeStatus(w, r, http.StatusBadRequest, "400 Bad Request")
		return
	}
	threadData := ThreatData{}
	jsonDec := json.NewDecoder(r.Body)
	jsonDec.DisallowUnknownFields() // check for garbage values
	if err := jsonDec.Decode(&threadData); err != nil {
		writeStatus(w, r, http.StatusBadRequest, "400 Bad Request, Error: "+err.Error())
		return
	}
	log.Printf("Decoded JSON %+v", threadData)
	if err := h.tmpl.Execute(w, threadData); err != nil {
		writeStatus(w, r, http.StatusConflict, "409 Conflict, Error: "+err.Error())
		return
	}
}

// Write Status to reponse and log for better debug
func writeStatus(w http.ResponseWriter, r *http.Request, aStatusCode int, aMessage string) {
	w.WriteHeader(aStatusCode)
	fmt.Fprint(w, aMessage)
	log.Println(aMessage, "Request URL:", r.URL, "Method:", r.Method)
}

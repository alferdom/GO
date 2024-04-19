package main

import (
	handlers "Simple_RestAPI/Handlers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type handlerFnc func(w http.ResponseWriter, r *http.Request)

func expectedStatusOK(t *testing.T, fnc handlerFnc, w *httptest.ResponseRecorder, r *http.Request, aErrorMessage string) {
	fnc(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Error(aErrorMessage)
	}
}
func expectedStatusNotOK(t *testing.T, fnc handlerFnc, w *httptest.ResponseRecorder, r *http.Request, aErrorMessage string) {
	fnc(w, r)
	if w.Result().StatusCode == http.StatusOK {
		t.Error(aErrorMessage)
	}
}

func isContentTypeEqualHTML(t *testing.T, w *httptest.ResponseRecorder) {
	contentType := w.Result().Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Error("Content-Type is", contentType, "Expected `text/html; charset=utf-8`")
	}
}

func TestGET(t *testing.T) {
	// TEST 1 GET to correct address
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler := handlers.NewHandler(nil)
	expectedStatusOK(t, handler.HandlerRootGet, w, req, "GET failed, expected success")
	isContentTypeEqualHTML(t, w)

	// TEST 2 GET to wrong address
	req = httptest.NewRequest(http.MethodGet, "/asdff", nil)
	w = httptest.NewRecorder()
	expectedStatusNotOK(t, handler.HandlerRootGet, w, req, "GET success, expected 400 Bad Request")

	// TEST 3 GET to only POST address
	req = httptest.NewRequest(http.MethodGet, "/render", nil)
	w = httptest.NewRecorder()
	expectedStatusNotOK(t, handler.HandlerRenderPost, w, req, "GET success, expected 400 Bad Request")
}

func TestPOST(t *testing.T) {
	// TEST 1 POST to wrong address
	req := httptest.NewRequest(http.MethodPost, "/render/", nil)
	w := httptest.NewRecorder()
	handler := handlers.NewHandler(nil)
	expectedStatusNotOK(t, handler.HandlerRenderPost, w, req, "POST success, expected 400 Bad Request")

	// TEST 2 POST to only GET address
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	w = httptest.NewRecorder()
	expectedStatusNotOK(t, handler.HandlerRootGet, w, req, "POST success, expected 400 Bad Request")

	// TESTS WITH JSON REQUESTS IN POST BODY
	tmpDefaultPath := "./threat.html.tmpl"
	tmpl := parseTemplate(&tmpDefaultPath)
	if tmpl == nil {
		t.Fatal("Template did not parsed, exptected parsed on path", tmpDefaultPath)
	}
	handler = handlers.NewHandler(tmpl)

	// TEST 3 POST with partial JSON in body request
	threatData := handlers.ThreatData{ThreatName: "NEW THREAT", DetectionDate: "1.1.2024"}
	encodedBytes, err := json.Marshal(threatData)
	if err != nil {
		t.Fatal("Cannot marshal.", err)
	}
	req = httptest.NewRequest(http.MethodPost, "/render", bytes.NewBuffer(encodedBytes))
	w = httptest.NewRecorder()
	expectedStatusOK(t, handler.HandlerRenderPost, w, req, "POST unsuccessful, expected success")
	isContentTypeEqualHTML(t, w)

	// TEST 4 POST with empty JSON in body request
	threatData = handlers.ThreatData{}
	encodedBytes, err = json.Marshal(threatData)
	if err != nil {
		t.Fatal("Cannot marshal.", err)
	}
	req = httptest.NewRequest(http.MethodPost, "/render", bytes.NewBuffer(encodedBytes))
	w = httptest.NewRecorder()
	expectedStatusOK(t, handler.HandlerRenderPost, w, req, "POST unsuccessful, expected success")
	isContentTypeEqualHTML(t, w)

	// TEST 5 POST with unknown JSON keys
	encodedBytes, err = json.Marshal(map[string]any{"threatName": "THREAT_NAME", "size": 45, "UnkownKey": "VALUE", "UnkownKey2": "VALUE2"})
	if err != nil {
		t.Fatal("Cannot marshal.", err)
	}
	req = httptest.NewRequest(http.MethodPost, "/render", bytes.NewBuffer(encodedBytes))
	w = httptest.NewRecorder()
	expectedStatusOK(t, handler.HandlerRenderPost, w, req, "POST unsuccessful, expected success")
	isContentTypeEqualHTML(t, w)
}

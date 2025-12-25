package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type LogEntry struct {
	Time       string `json:"time,omitempty"`
	RemoteAddr string `json:"remote_addr"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	Status     string `json:"status"`
	ContentLen int64  `json:"content_len"`
	Error      string `json:"error"`
	Blocked    bool   `json:"blocked"`
}

var logFile *os.File

func InitLogger() {
	lf, err := os.OpenFile(Config.Log_filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("ERROR: could not open log file: %s\n", err.Error())
	}
	logFile = lf
}

func CloseLogger() {
	if logFile != nil {
		logFile.Sync()
		logFile.Close()
	}
}

func logAppend(entry LogEntry) error {
	ljs, err := json.Marshal(entry)
	if err != nil {
		log.Printf("ERROR: could not marshal log entry: %s\n", err.Error())
		return err
	}
	_, err = logFile.WriteString(string(ljs) + "\n")
	if err != nil {
		log.Printf("ERROR: could not write logs file: %s\n", err.Error())
		return err
	}
	return nil
}

func LogOk(req *http.Request, res *http.Response) {
	l := LogEntry{
		Time:       time.Now().Format(time.RFC3339),
		RemoteAddr: req.RemoteAddr,
		Method:     req.Method,
		URL:        req.URL.String(),
		Status:     res.Status,
		ContentLen: res.ContentLength,
	}

	logAppend(l)
}

func LogError(req *http.Request, err error) {
	l := LogEntry{
		Time:       time.Now().Format(time.RFC3339),
		RemoteAddr: req.RemoteAddr,
		Method:     req.Method,
		URL:        req.URL.String(),
		Error:      err.Error(),
	}

	logAppend(l)
}

func LogBlocked(req *http.Request) {
	l := LogEntry{
		Time:       time.Now().Format(time.RFC3339),
		RemoteAddr: req.RemoteAddr,
		Method:     req.Method,
		URL:        req.URL.String(),
		Blocked:    true,
	}

	logAppend(l)
}

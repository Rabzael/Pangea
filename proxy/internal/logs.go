package internal

import (
	"log"
	"net/http"
)

func logBase(req *http.Request) string {
	return "| " + req.RemoteAddr + " -> " + req.Method + " " + req.URL.String()
}

func LogOk(req *http.Request, res *http.Response) {
	log.Printf("%s\t:\t%s\n",
		logBase(req),
		res.Status,
	)
}

func LogError(req *http.Request, err error) {
	log.Printf("%s\t:\tERROR: %s\n",
		logBase(req),
		err.Error(),
	)
}

func LogBlocked(req *http.Request) {
	log.Printf("%s\t->\tBLOCKED\n",
		logBase(req),
	)
}

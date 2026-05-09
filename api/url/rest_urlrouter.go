package url

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

type RestUrlRouter struct {
	mux      *mux.Router
	hostname string
}

func (p *RestUrlRouter) Init(path string) *RestUrlRouter {
	if p.mux == nil {
		p.mux = mux.NewRouter().StrictSlash(false)
	}
	if len(p.hostname) == 0 {
		p.hostname, _ = os.Hostname()
	}

	if strings.Index(path, "/") == 0 && strings.LastIndex(path, "/")+1 == len(path) {
		p.mux = p.mux.PathPrefix(path).Subrouter()
	}
	return p
}

func (p *RestUrlRouter) AddObject(path string) *mux.Router {
	return p.mux.PathPrefix(path).Subrouter()
}

func (p *RestUrlRouter) GetRouter() *mux.Router {
	return p.mux
}

func (p *RestUrlRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			panic(err)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Requst-Host-Trace", p.hostname)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	p.mux.NotFoundHandler = p.UrlUnMatchHandler()
	p.mux.ServeHTTP(w, r)
	//http.NotFound(w, r)
}

var UrlUmatchErrStr = "{\"requestId\":\"%s\",\"error\":{\"code\":%d,\"status\":\"NOT_FOUND\",\"message\":\"url not found.\",\"detail\":[]}}"

func UrlUnFound(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	requestId := ""
	ctx := r.Context()
	if ctx.Value("trace_id") != nil {
		requestId = ctx.Value("trace_id").(string)
	}
	w.Write([]byte(fmt.Sprintf(UrlUmatchErrStr, requestId, http.StatusNotFound)))
}
func (p *RestUrlRouter) DefaultFunc(w http.ResponseWriter, r *http.Request) {
	h, v := http.DefaultServeMux.Handler(r)
	if v == "" {
		UrlUnFound(w, r)
	} else {
		h.ServeHTTP(w, r)
	}
}

func (p *RestUrlRouter) UrlUnMatchHandler() http.Handler {
	return http.HandlerFunc(p.DefaultFunc)
}

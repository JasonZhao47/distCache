package distCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const defaultBasePath = "/_distcache/"

// HTTPPool PeerPicker implementation
type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// note2self: interface{} = any
func (p *HTTPPool) Log(format string, v ...any) {
	log.Printf("[Server %s %s] %s", p.self, time.Now().String(), fmt.Sprintf(format, v))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check if path is basePath
	url := r.URL.String()
	if !strings.HasPrefix(url, p.basePath) {
		// problem
		p.Log("Request url base path error %s\n", url)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Url error!"))
		return
	}
	p.Log("Request %s on URL %s\n", r.Method, r.URL.String())
	// get key from path specified
	// /base/group/key
	path := strings.SplitN(url, "/", 3)
	groupName, key := path[1], path[2]
	g, err := GetGroup(groupName)
	if err != nil {
		p.Log("Get group error %s\n", url)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Internal Error"))
		return
	}
	bv, err := g.Get(groupName, key)
	if err != nil {
		p.Log("Get ByteView error %s %s\n", groupName, key)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Internal Error"))
		return
	}
	// set http headers
	w.Header().Set("Content-Type", "application/octet-stream")
	// write back byteView (View only)
	w.Write(bv.AsSlice())
}

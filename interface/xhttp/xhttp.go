package xhttp

import (
	"cacheme/cache"
	"cacheme/store"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Conf struct {
	Host string
	Port int
}

type HTTPCacheServer struct {
	conf *Conf
	c    cache.Cache
}

var server HTTPCacheServer

func CacheRoute(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		HandleGetKey(w, req)
	case http.MethodPost:
		HandleSetKey(w, req)
	case http.MethodDelete:
		HandleDelKey(w, req)
	default:
		HandleNotFound(w, req)
	}
}

func HandleGetKey(w http.ResponseWriter, req *http.Request) {
	pathList := strings.Split(req.URL.EscapedPath(), "/")
	key := pathList[len(pathList)-1]
	v, err := server.c.Get(key)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key + " found value: "))
	w.Write(v)
}

func HandleSetKey(w http.ResponseWriter, req *http.Request) {
	type SetKeyParams struct {
		Key   string
		Value string
	}
	var skp SetKeyParams
	err := json.NewDecoder(req.Body).Decode(&skp)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(skp.Key) == 0 || len(skp.Key) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = server.c.Set(skp.Key, []byte(skp.Value))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(skp.Key + " has been saved"))
}

func HandleDelKey(w http.ResponseWriter, req *http.Request) {
	pathList := strings.Split(req.URL.EscapedPath(), "/")
	key := pathList[len(pathList)-1]
	exist, err := server.c.Del(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if exist {
		w.Write([]byte(key + " has been removed"))
	} else {
		w.Write([]byte(key + " doesn't exist"))
	}
}

func HandleGetInfo(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello"))
}

func HandleNotFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func BoostrapHTTPServer(conf *Conf) {
	var err error
	var c cache.Cache
	c, err = store.NewStore("./badger")
	if err != nil {
		log.Println("can't open store directory")
		panic(err)
	}
	server = HTTPCacheServer{
		// c: mem.NewMemCache(),
		conf: conf,
		c:    c,
	}
	http.HandleFunc("/key/", CacheRoute)
	http.HandleFunc("/stat", HandleGetInfo)
	http.HandleFunc("*", HandleNotFound)
	log.Println("http server is listenning at port:", conf.Port)
	if err := http.ListenAndServe(net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)), nil); err != nil {
		log.Fatal("can't bootstrap http server: ", err.Error())
	}
}

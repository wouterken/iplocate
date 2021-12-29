package main

import (
  "fmt"
  "net"
  "log"
  "flag"
  "encoding/json"
  "net/http"
  "strings"
  "os/signal"
  "context"
  "os"
  "github.com/oschwald/maxminddb-golang"
)


func ipHandler(db *maxminddb.Reader) http.Handler {

  var record struct {
    Country struct {
      ISOCode string `maxminddb:"iso_code"`
    } `maxminddb:"country"`
    City struct {
      Confidence uint8             `maxminddb:"confidence"`
      GeoNameID  uint              `maxminddb:"geoname_id"`
      Names      map[string]string `maxminddb:"names"`
    } `maxminddb:"city"`
  } // Or any appropriate struct

  fn := func(w http.ResponseWriter, req *http.Request) {

    ipstr := strings.Split(req.RemoteAddr, ":")[0]
    userIP := net.ParseIP(ipstr)

    if userIP == nil {
        //return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
        log.Fatal("Couldn't detect user IP", req.RemoteAddr)
        fmt.Fprintf(w, "userip: %q is not IP:port", req.RemoteAddr)
        return
    }

    err := db.Lookup(userIP, &record)

    if err != nil {
      log.Fatal(err)
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(record)
  }
  return http.HandlerFunc(fn)
}

func main() {
  mux             := http.NewServeMux()
  dbLocation      := "GeoLite2-City.mmdb"
  db, err         := maxminddb.Open(dbLocation)
  ip              := ipHandler(db)
  idleConnsClosed := make(chan struct{})

  if err != nil {
    log.Fatal(err)
  }

  mux.Handle("/", ip)

  port := flag.Int("port", 3000, "Port to run server on")
  flag.Parse()

  srv := &http.Server{Addr: fmt.Sprintf(":%d", *port), Handler: mux}

  go func() {
    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)
    <-sigint
    log.Printf("Shutting down")
    srv.Shutdown(context.Background())
    close(idleConnsClosed)
  }()

  log.Printf("Listening on %d", *port)
  if err := srv.ListenAndServe(); err != http.ErrServerClosed {
      // unexpected error. port in use?
      log.Fatalf("ListenAndServe(): %v", err)
  }
}




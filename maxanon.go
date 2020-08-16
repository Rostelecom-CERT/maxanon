package maxanon

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"

	"github.com/Rostelecom-CERT/maxanon/storage"
	"github.com/julienschmidt/httprouter"
)

const bulkSize = 500000
const collectionName = "records"

// Reputation is structure which return in response
type Reputation struct {
	Anonymous         bool `json:"anonymous"`
	AnonymousVPN      bool `json:"anonymous_vpn"`
	IsHostingProvider bool `json:"is_hosting_provider"`
	IsPublicProxy     bool `json:"is_public_proxy"`
	IsTorExitNode     bool `json:"is_tor_exit_node"`
}

// App is base struct
type App struct {
	DB      storage.Database
	csvPath string
	Meta    struct {
		count  uint64
		buffer []interface{}
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func isIpv4Net(network string) bool {
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}/[0-9]{1,2}`)
	return re.MatchString(network)
}

func isIpv4Host(host string) bool {
	return net.ParseIP(host) != nil
}

func convert(boolString string) bool {
	if boolString == "1" {
		return true
	}
	return false
}

func (a *App) readFile(csvPath string) error {
	csvfile, err := os.Open(csvPath)
	if err != nil {
		return errors.New("Couldn't open the csv file")
	}

	r := csv.NewReader(csvfile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if isIpv4Net(record[0]) {
			ip, ipnet, err := net.ParseCIDR(record[0])
			if err != nil {
				return err
			}

			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
				if a.Meta.count == bulkSize {
					err := a.DB.InsertBulk(a.Meta.buffer)
					if err != nil {
						return err
					}
					a.Meta.count = 0
					a.Meta.buffer = []interface{}{}
				}
				a.Meta.buffer = append(a.Meta.buffer, storage.Data{
					IP:                ip.String(),
					Anonymous:         convert(record[1]),
					AnonymousVPN:      convert(record[2]),
					IsHostingProvider: convert(record[3]),
					IsPublicProxy:     convert(record[4]),
					IsTorExitNode:     convert(record[5]),
				})
				a.Meta.count++
			}
		}

	}
	return err
}
func (a *App) fillDB(csvPath string) error {
	err := a.readFile(csvPath)
	if err != nil {
		return err
	}
	return nil
}

// New return base struct
func New(csvPath string, dbType string, dbURL string) (*App, error) {
	var db storage.Database
	switch dbType {
	case "mongo":
		mongo := &storage.MongoDB{}
		err := mongo.Open(dbURL)
		if err != nil {
			return nil, err
		}
		db = mongo
	case "redis":
		redis := &storage.Redis{}
		err := redis.Open(dbURL)
		if err != nil {
			return nil, err
		}
		db = redis
	}
	return &App{
		DB:      db,
		csvPath: csvPath,
	}, nil
}

// Run starting main function
func (a *App) Run(port string) error {
	st, err := a.DB.Exist(collectionName)
	if err != nil {
		return err
	}
	if !st {
		err := a.fillDB(a.csvPath)
		if err != nil {
			return err
		}
	}
	err = a.Serve(port)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) infoIP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if ps.ByName("ip") == "" {
		http.Error(w, "No data provided", http.StatusBadRequest)
		return
	} else if !isIpv4Host(ps.ByName("ip")) {
		http.Error(w, "Wrong IP address type", http.StatusBadRequest)
		return
	}

	data, err := a.DB.Get(ps.ByName("ip"))
	if err != nil {
		http.Error(w, "Error", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Reputation{
		Anonymous:         data.Anonymous,
		AnonymousVPN:      data.AnonymousVPN,
		IsHostingProvider: data.IsHostingProvider,
		IsPublicProxy:     data.IsPublicProxy,
		IsTorExitNode:     data.IsTorExitNode,
	})
	fmt.Fprintln(w)
}

// Serve starting http service
func (a *App) Serve(addr string) error {
	router := httprouter.New()
	router.GET("/api/v1/info/:ip", a.infoIP)
	log.Printf("Listening on %s", addr)
	return http.ListenAndServe(addr, router)
}

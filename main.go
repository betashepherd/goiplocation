package main

import (
	//http://localhost:8181/debug/pprof/ to see the prof log
    _ "net/http/pprof"
	"runtime"
    "net/http"
    "fmt"
    "os"
    "bufio"
    "encoding/csv"
    "strconv"
    "strings"
    "reflect"
)

type IpCityType struct {
	ipstart int64
	ipend int64
	cname string
	ccode int64
}

func (s IpCityType) IsEmpty() bool {
  return reflect.DeepEqual(s,IpCityType{})
}

const (
	IP_CITY       = "./conf/ip_city.data"
)

var IpCityList []IpCityType

func _IPtoInt(ip string) int64 {
	bits := strings.Split(ip, ".")
    b0, _ := strconv.Atoi(bits[0])
    b1, _ := strconv.Atoi(bits[1])
    b2, _ := strconv.Atoi(bits[2])
    b3, _ := strconv.Atoi(bits[3])
    
    var sum int64
    sum += int64(b0) << 24
    sum += int64(b1) << 16
    sum += int64(b2) << 8
    sum += int64(b3)
    return sum
}

func _initConfig() {
	file, _ := os.Open(IP_CITY)
	defer file.Close()
	r := csv.NewReader(bufio.NewReader(file))
	data, _ := r.ReadAll()
	
	var IpCity IpCityType
	
	for _, d := range data {
		IpCity.ipstart = _IPtoInt(d[0])
		IpCity.ipend = _IPtoInt(d[1])
		IpCity.cname = d[2]
		code, _ := strconv.ParseInt(d[3], 10, 64)
		IpCity.ccode = code
		IpCityList = append(IpCityList, IpCity)
	}
}

func _GetLoc(ip string) IpCityType {

	left := 0
	right := len(IpCityList) - 1
	
	locip := _IPtoInt(ip)
	var IpInfo IpCityType
	
	for left <= right {
		mid := (left + right) / 2
		if IpCityList[mid].ipstart <= locip && IpCityList[mid].ipend >= locip {
			IpInfo = IpCityList[mid]
			break
		} else if IpCityList[mid].ipstart > locip {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	
	return IpInfo
}

// Default Request Handler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	
	//fmt.Println(r.URL.String())
	
	//ip := r.RemoteAddr
	ip := r.URL.Query().Get("ip")
	var IpInfo IpCityType
	 
	if ip != "" {
		IpInfo = _GetLoc(ip)
	}
	
	if IpInfo.IsEmpty() == false {
		fmt.Fprintf(w, "city:%s,code:%d", IpInfo.cname, IpInfo.ccode)
	} else {
		fmt.Fprintf(w, "loc failed!")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	_initConfig()
    http.HandleFunc("/", indexHandler)
    http.ListenAndServe(":8181", nil)
}

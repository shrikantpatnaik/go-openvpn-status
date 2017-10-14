package openvpnStatus

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

// Client Struct to store Client Information
type Client struct {
	CommonName     string
	RealAddress    string
	BytesReceived  string
	BytesSent      string
	ConnectedSince time.Time
}

// Routing Struct to store Routing Information
type Routing struct {
	VirtualAddress string
	CommonName     string
	RealAddress    string
	LastRef        time.Time
}

// GlobalStats struct to store Global stats
type GlobalStats struct {
	MaxBcastMcastQueueLen int
}

// Status struct to store the status
type Status struct {
	ClientList   []Client
	RoutingTable []Routing
	GlobalStats  GlobalStats
	UpdatedAt    time.Time
	IsUp         bool
}

type parseError struct {
	s string
}

func (e *parseError) Error() string {
	return e.s
}

var clientListHeaderColumns = [5]string{
	"Common Name",
	"Real Address",
	"Bytes Received",
	"Bytes Sent",
	"Connected Since",
}

var routingTableHeaderColumns = [4]string{
	"Virtual Address",
	"Common Name",
	"Real Address",
	"Last Ref",
}

const (
	clientListHeaders = 1 << iota
	routingTableHeaders
	globalStatsHeaders
)

func checkHeaders(headers []string) int {
	if checkClientListHeaders(headers) {
		return clientListHeaders
	} else if checkRoutingTableHeaders(headers) {
		return routingTableHeaders
	} else {
		return 0
	}
}

func checkClientListHeaders(headers []string) bool {
	for i, v := range headers {
		if v != clientListHeaderColumns[i] {
			return false
		}
	}
	return true
}

func checkRoutingTableHeaders(headers []string) bool {
	for i, v := range headers {
		if v != routingTableHeaderColumns[i] {
			return false
		}
	}
	return true
}

const dateLayout = "Mon Jan 2 15:04:05 2006"

// ParseFile parses OpenVPN Status file ad returns a Status struct
func ParseFile(file string) (*Status, error) {
	conn, err := os.Open(file)
	defer conn.Close()
	if err != nil {
		return &Status{IsUp: false}, err
	}

	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	var lastUpdatedAt string
	var clients []Client
	var routingTable []Routing
	var maxBcastMcastQueueLen int
	nextFieldType := 0
	isEmpty := true
	for scanner.Scan() {
		isEmpty = false
		var ct time.Time
		var rtt time.Time
		fields := strings.Split(scanner.Text(), ",")
		if fields[0] == "END" && len(fields) == 1 {
			// Stats footer.
		} else if fields[0] == "OpenVPN CLIENT LIST" {
			// Header
		} else if fields[0] == "ROUTING TABLE" {
			// Routing table header
		} else if fields[0] == "GLOBAL STATS" {
			nextFieldType = globalStatsHeaders
		} else if fields[0] == "Updated" && len(fields) == 2 {
			lastUpdatedAt = fields[1]
		} else if checkHeaders(fields) == clientListHeaders {
			nextFieldType = clientListHeaders
		} else if checkHeaders(fields) == routingTableHeaders {
			nextFieldType = routingTableHeaders
		} else if nextFieldType == clientListHeaders && len(fields) == 5 {
			ct, _ = time.Parse(dateLayout, fields[4])
			clients = append(clients, Client{fields[0], fields[1], fields[2], fields[3], ct})
		} else if nextFieldType == routingTableHeaders && len(fields) == 4 {
			rtt, _ = time.Parse(dateLayout, fields[3])
			routingTable = append(routingTable, Routing{fields[0], fields[1], fields[2], rtt})
		} else if nextFieldType == globalStatsHeaders && len(fields) == 2 {
			if fields[0] == "Max bcast/mcast queue length" {
				i, err := strconv.Atoi(fields[1])
				if err == nil {
					maxBcastMcastQueueLen = i
				}
			}
		} else {
			return &Status{IsUp: false}, &parseError{"Unable to Parse Status file"}
		}
	}
	if isEmpty {
		return &Status{IsUp: false}, &parseError{"Status File is empty"}
	}
	t, _ := time.Parse(dateLayout, lastUpdatedAt)
	return &Status{
		ClientList:   clients,
		RoutingTable: routingTable,
		GlobalStats:  GlobalStats{maxBcastMcastQueueLen},
		UpdatedAt:    t,
		IsUp:         true,
	}, nil
}

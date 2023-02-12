package bgptools

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

var errInvalidFormat = fmt.Errorf("invalid return format")

type Response struct {
	AS        uint32
	IP        string
	Prefix    string
	Country   string
	Registry  string
	Allocated string
	ASName    string
}

// parse parses the | delimited bgp.tools whois output
func parse(s string) (*Response, error) {
	var r Response

	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) != 2 {
		return nil, errInvalidFormat
	}

	d := strings.Split(lines[1], "|")
	as, err := strconv.Atoi(strings.TrimSpace(d[0]))
	if err != nil {
		return nil, errInvalidFormat
	}

	r.AS = uint32(as)
	r.IP = strings.TrimSpace(d[1])
	r.Prefix = strings.TrimSpace(d[2])
	r.Country = strings.TrimSpace(d[3])
	r.Registry = strings.TrimSpace(d[4])
	r.Allocated = strings.TrimSpace(d[5])
	r.ASName = strings.TrimSpace(d[6])

	return &r, nil
}

// query makes a whois query
func query(s string) ([]byte, error) {
	conn, err := net.Dial("tcp", "bgp.tools:43")
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.Write([]byte(s + "\r\n"))
	if err != nil {
		return nil, fmt.Errorf("whois: send to whois server failed: %w", err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	buffer, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

// Query makes a bgp.tools query
func Query(s string) (*Response, error) {
	whoisResponse, err := query(s)
	if err != nil {
		return nil, err
	}
	return parse(string(whoisResponse))
}

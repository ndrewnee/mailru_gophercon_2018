package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"strconv"
	"strings"
)

// глобальные переменные запрещены
// cgo запрещен

// User ...
type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Hits     []string `json:"hits"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

// Fast ...
func Fast(in io.Reader, out io.Writer, networks []string) {
	var (
		err   error
		user  User
		total int
		buf   = bytes.NewBuffer(nil)
	)

	parsedNets := parseNetworks(networks)

	for idx := 1; true; idx++ {
		err = json.NewDecoder(in).Decode(&user)
		if err != nil {
			if err == io.EOF {
				out.Write([]byte("Total: " + strconv.Itoa(total) + "\n"))
				out.Write(buf.Bytes())
				return
			}

			panic(err)
		}

		hasBrowser := hasValidBrowser(user.Browsers)
		if !hasBrowser {
			continue
		}

		hasIP := hasMoreThan3IPs(parsedNets, user.Hits)
		if !hasIP {
			continue
		}

		total++
		email := strings.Replace(user.Email, "@", " [at] ", -1)
		buf.WriteString("[" + strconv.Itoa(idx) + "] " + user.Name + " <" + email + ">\n")
	}
}

func parseNetworks(networks []string) []*net.IPNet {
	nets := make([]*net.IPNet, len(networks))
	for i := range networks {
		_, ipv4Net, err := net.ParseCIDR(networks[i])
		if err != nil {
			panic(err)
		}

		nets[i] = ipv4Net
	}

	return nets
}

func hasMoreThan3IPs(networks []*net.IPNet, hits []string) bool {
	var (
		netIP net.IP
		ok    bool
		count int
	)

	for i := range networks {
		for j := range hits {
			netIP = net.ParseIP(hits[j])
			ok = networks[i].Contains(netIP)
			if ok {
				count++
				if count >= 3 {
					return true
				}
				continue
			}
		}
	}

	return count >= 3
}

func hasValidBrowser(browsers []string) bool {
	var ok bool
	for i := range browsers {
		ok = strings.Contains(browsers[i], "Chrome/60.0.3112.90")
		if ok {
			return true
		}

		ok = strings.Contains(browsers[i], "Chrome/52.0.2743.116")
		if ok {
			return true
		}

		ok = strings.Contains(browsers[i], "Chrome/57.0.2987.133")
		if ok {
			return true
		}
	}

	return false
}

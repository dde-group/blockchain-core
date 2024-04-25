package network

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const HeaderXRequestID = "X-Request-ID"

var (
	regIpV4 *regexp.Regexp
	regIpV6 *regexp.Regexp
)

func init() {
	regIpV4, _ = regexp.Compile("^(25[0-5]|2[0-4]\\d|[0-1]?\\d?\\d)(\\.(25[0-5]|2[0-4]\\d|[0-1]?\\d?\\d)){3}(/([1-9]|1\\d|2\\d|3[0-2]))?$")
	regIpV6, _ = regexp.Compile("((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|" +
		"((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3})|:))|" +
		"(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)" +
		"(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|" +
		"((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|" +
		"(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)" +
		"(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|" +
		"2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|" +
		"((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|" +
		"((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}))|:)))(%.+)?")
}

func GetHttpRequestIP(r *http.Request) string {

	ipHeaders := []string{
		"x-forwarded-for",
		"Proxy-Client-IP",
		"WL-Proxy-Client-IP",
		"HTTP_CLIENT_IP",
		"HTTP_X_FORWARDED_FOR",
	}

	for _, ipHeader := range ipHeaders {
		ip := getRequestHeaderIP(r, ipHeader)
		if nil == CheckIp(ip) {
			return ip
		}
	}

	ip := r.RemoteAddr
	ip = strings.Trim(ip, " \b\t\r\n")
	ip = SplitIPPort(ip)

	if nil == CheckIp(ip) {
		return ip
	}

	return ""
}

func CheckIp(ip string) error {
	if len(ip) < 1 {
		return errors.New("short length")
	}

	if regIpV4 != nil && regIpV4.MatchString(ip) {
		return nil

	}

	if regIpV6 != nil && regIpV6.MatchString(ip) {
		return nil
	}

	return fmt.Errorf("invalid ip")
}

func SplitIPPort(addr string) string {

	_ipc := strings.Split(addr, ":")

	if len(_ipc) > 0 {
		addr = _ipc[0]
	}

	return addr
}

func getRequestHeaderIP(r *http.Request, headerName string) string {

	if r.Header != nil {
		ip := r.Header.Get(headerName)
		if len(ip) > 0 {
			ips := strings.Split(ip, ",")
			if len(ips) > 0 {
				for _, _ip := range ips {
					_ip = strings.Trim(_ip, " \b\t\r\n")
					_ip = SplitIPPort(_ip)
					if len(_ip) > 0 {
						return _ip
					}
				}
			}
		}
	}

	return ""
}

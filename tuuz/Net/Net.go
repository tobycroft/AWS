package Net

import (
	"crypto/tls"
	"fmt"
	"main.go/tuuz/Calc"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var transport *http.Transport

var dailer = &net.Dialer{
	Timeout:   5 * time.Second,
	KeepAlive: 0 * time.Second,
	DualStack: true,
}

func init() {
	transport = &http.Transport{
		DialContext:  dailer.DialContext,
		MaxIdleConns: 100,
	}
}

func Rpc(url string, postData map[string]interface{}, username, password string) (string, error) {
	req := Request()
	header := map[string]string{"Content-type": "application/json"}
	req.SetHeaders(header)
	req.SetBasicAuth(username, password)
	req.SetTimeout(5 * time.Second)
	req.DisableKeepAlives(true)
	req.SetTLSClient(&tls.Config{InsecureSkipVerify: true})
	req.Transport(transport)
	ret, err := req.Post(url, postData)
	body, err := ret.Content()
	if err != nil {
		return "", err
	} else {
		return body, err
	}
}

func Post(url string, queries map[string]interface{}, postData map[string]interface{}, headers map[string]string, cookies map[string]string) (string, error) {
	// 链式操作
	req := Request()
	req.SetHeaders(headers)
	req.SetCookies(cookies)
	req.SetTimeout(5 * time.Second)
	req.DisableKeepAlives(true)
	req.SetTLSClient(&tls.Config{InsecureSkipVerify: true})
	req.Transport(transport)
	ret, err := req.Post(url+"?"+Http_build_query(queries), postData)
	if err != nil {
		return "", err
	}
	return ret.Content()
}

func PostCookie(url string, queries map[string]interface{}, postData map[string]interface{}, headers map[string]string, cookies map[string]string) (string, map[string]interface{}, error) {
	req := Request()
	req.SetHeaders(headers)
	req.SetCookies(cookies)
	req.SetTimeout(5 * time.Second)
	req.DisableKeepAlives(true)
	req.SetTLSClient(&tls.Config{InsecureSkipVerify: true})
	req.Transport(transport)
	ret, err := req.Post(url+"?"+Http_build_query(queries), postData)
	body, err := ret.Content()

	cookie_arr := CookieHandler(ret.Cookies())
	//fmt.Println(cookie_arr)
	if err != nil {
		return "", cookie_arr, err
	} else {
		return body, cookie_arr, err
	}
}

func Get(url string, queries map[string]interface{}, headers map[string]string, cookies map[string]string) (string, error) {
	req := Request()
	req.SetHeaders(headers)
	req.SetCookies(cookies)
	req.SetTimeout(5 * time.Second)
	req.DisableKeepAlives(true)
	req.SetTLSClient(&tls.Config{InsecureSkipVerify: true})
	req.Transport(transport)
	ret, err := req.Get(url, queries)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	body, err := ret.Content()

	if err != nil {
		return "", err
	} else {
		return body, err
	}
}

func GetCookie(url string, queries map[string]interface{}, headers map[string]string, cookies map[string]string) (string, map[string]interface{}, error) {
	req := Request()
	req.SetHeaders(headers)
	req.SetCookies(cookies)
	req.SetTimeout(5 * time.Second)
	req.DisableKeepAlives(true)
	req.SetTLSClient(&tls.Config{InsecureSkipVerify: true})
	req.Transport(transport)
	ret, err := req.Get(url, queries)
	body, err := ret.Content()
	cookie_arr := CookieHandler(ret.Cookies())
	//fmt.Println(cookie_arr)
	if err != nil {
		return "", cookie_arr, err
	} else {
		return body, cookie_arr, err
	}
}

func CookieHandler(resp_headers []*http.Cookie) map[string]interface{} {
	cookie_arr := make(map[string]interface{})
	for _, resp_header := range resp_headers {
		cookie_arr[resp_header.Name] = resp_header.Value
	}
	return cookie_arr
}

func CookieHandler2(resp_header map[string]interface{}) map[string]interface{} {
	cookie := strings.Split(Calc.Any2String(resp_header["Set-Cookie"]), "; ")
	cookie_arr := make(map[string]interface{})
	for _, v := range cookie {
		split := strings.Split(v, "=")
		if CookieTagChecker(split[0]) == true {
			cookie_arr[split[0]] = split[1]
		}
	}
	return cookie_arr
}

func Http_build_query(querymap map[string]interface{}) string {
	query := make(url.Values)
	for k, v := range querymap {
		query.Add(k, Calc.Any2String(v))
	}
	//fmt.Println(query.Encode())
	return query.Encode()
}
package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type IPVersion uint8

const (
	v4 IPVersion = iota
	v6
)

// Get Public IP for the device
func MyPublicIP(version IPVersion) net.IP {
	urlList := [][]string{
		{ // IPv4 API
			"https://api.ipify.org?format=text",
			"https://myexternalip.com/raw",
			"https://v4.ident.me/",
		},
		{ // IPv6 API
			"https://api64.ipify.org?format=text",
			"https://v6.ident.me/",
		},
	}
	for _, url := range urlList[int(version)] {
		ip_timeout := make(chan string, 1)
		go func() {
			resp, err := http.Get(url)
			if err != nil {
				ip_timeout <- ""
			}
			ip, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				ip_timeout <- ""
			}
			ip_timeout <- string(ip)
		}()

		select {
		case ip_valid := <-ip_timeout:
			final_ip := net.ParseIP(ip_valid)
			if final_ip != nil {
				return final_ip
			}
		case <-time.After(1 * time.Second): // timeout after 1 second
			continue
		}
	}
	panic("Failed to fetch Public IP")
}

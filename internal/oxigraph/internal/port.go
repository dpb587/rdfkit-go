package internal

import (
	"fmt"
	"net"
	"time"
)

func IsPortOpen(host, port string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}

	conn.Close()

	return true
}

func WaitForPortOpen(host, port string, timeout time.Duration) error {
	stop := time.Now().Add(timeout)

	for !IsPortOpen(host, port, 3*time.Second) {
		if time.Now().After(stop) {
			return fmt.Errorf("timeout exceeded")
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

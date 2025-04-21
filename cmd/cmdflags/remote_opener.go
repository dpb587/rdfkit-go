package cmdflags

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RemoteOpenerFunc func(p string, originalErr error) (io.ReadCloser, http.Header, error)

func WebRemoteOpener(p string, originalErr error) (io.ReadCloser, http.Header, error) {
	if !strings.HasPrefix(p, "http:") && !strings.HasPrefix(p, "https:") {
		return nil, nil, originalErr
	}

	resp, err := http.DefaultClient.Get(p)
	if err != nil {
		if resp.Body != nil {
			resp.Body.Close()
		}

		return nil, nil, fmt.Errorf("remote: %v", err)
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.Body != nil {
			resp.Body.Close()
		}

		return nil, nil, fmt.Errorf("remote: unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, resp.Header, nil
}

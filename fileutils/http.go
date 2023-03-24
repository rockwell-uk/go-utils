package fileutils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// loose defaults.
var ConnectMaxWaitTime = 10 * time.Second
var RequestMaxWaitTime = 120 * time.Second

func Request(url, verb string, payload io.Reader, timeout *time.Duration) (*http.Response, context.CancelFunc, error) {
	if timeout == nil {
		timeout = &RequestMaxWaitTime
	}

	return sendWithContext(url, verb, payload, timeout)
}

func sendWithContext(url, verb string, payload io.Reader, timeout *time.Duration) (*http.Response, context.CancelFunc, error) {
	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: ConnectMaxWaitTime,
			}).DialContext,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)

	request, err := http.NewRequestWithContext(ctx, verb, url, payload)
	if err != nil {
		return nil, cancel, err
	}

	response, err := client.Do(request)

	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return nil, cancel, fmt.Errorf("request timeout: %v", err.Error())
	} else if err != nil {
		return nil, cancel, fmt.Errorf("request error: %v", err.Error())
	}

	return response, cancel, nil
}

func DownloadFileToPath(url, filePath string) error {
	var funcName string = "DownloadFileToPath"

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%v.%v: error creating file [%v], [%v]", packageName, funcName, filePath, err.Error())
	}
	defer out.Close()

	// Get the data
	resp, cancel, err := Request(url, http.MethodGet, nil, nil)
	defer cancel()
	if err != nil {
		return fmt.Errorf("%v.%v: error downloading file [%v], [%v]", packageName, funcName, url, err.Error())
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v.%v: bad return status: %s", packageName, funcName, resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("%v.%v: error writing file [%v], [%v]", packageName, funcName, filePath, err.Error())
	}

	return nil
}

func Get(url string) (*http.Response, context.CancelFunc, error) {
	return Request(url, http.MethodGet, nil, nil)
}

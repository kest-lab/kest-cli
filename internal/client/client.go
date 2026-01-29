package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
	Stream  bool
}

type Response struct {
	Status   int
	Headers  http.Header
	Body     []byte
	Duration time.Duration
}

func Execute(opt RequestOptions) (*Response, error) {
	client := &http.Client{
		Timeout: opt.Timeout,
	}

	req, err := http.NewRequest(opt.Method, opt.URL, bytes.NewBuffer(opt.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if opt.Stream {
		return handleStream(resp, duration)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:   resp.StatusCode,
		Headers:  resp.Header,
		Body:     body,
		Duration: duration,
	}, nil
}

func handleStream(resp *http.Response, duration time.Duration) (*Response, error) {
	fmt.Println("\n--- Stream Response ---")
	reader := bufio.NewReader(resp.Body)
	var fullBody bytes.Buffer

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if line != "" {
			fmt.Print(line)
			fullBody.WriteString(line)
		}
	}

	return &Response{
		Status:   resp.StatusCode,
		Headers:  resp.Header,
		Body:     fullBody.Bytes(),
		Duration: duration,
	}, nil
}

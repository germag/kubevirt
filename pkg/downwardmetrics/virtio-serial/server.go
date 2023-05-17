/*
 * This file is part of the kubevirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 *
 */

package virtio_serial

import (
	"bufio"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"net/url"
	"strings"
	"syscall"
	"time"

	"kubevirt.io/client-go/log"

	"kubevirt.io/kubevirt/pkg/downwardmetrics/vhostmd/api"
	diskutils "kubevirt.io/kubevirt/pkg/ephemeral-disk-utils"
	metricsScraper "kubevirt.io/kubevirt/pkg/monitoring/domainstats/downwardmetrics"
)

const (
	maxConnectAttempts = 6
	invalidRequest     = "INVALID REQUEST\n\n"
	emptyMetrics       = "<metrics><!-- host metrics not available --><!-- VM metrics not available --></metrics>"
)

func RunDownwardMetricsVirtioServer(context context.Context, nodeName, channelSocketPath, launcherSocketPath string) error {
	report, err := newMetricsReporter(nodeName, launcherSocketPath)
	if err != nil {
		return err
	}

	server := downwardMetricsServer{
		maxConnectAttempts: maxConnectAttempts,
		virtioSerialSocket: channelSocketPath,
		reportFn:           report,
	}
	return server.start(context)
}

type metricsReporter func() (*api.Metrics, error)

func newMetricsReporter(nodeName, launcherSocketPath string) (metricsReporter, error) {
	exists, err := diskutils.FileExists(launcherSocketPath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("virt-launcher socket not found")
	}

	scraper := metricsScraper.NewReporter(nodeName)

	return func() (*api.Metrics, error) {
		return scraper.Report(launcherSocketPath)
	}, nil
}

type downwardMetricsServer struct {
	maxConnectAttempts uint
	virtioSerialSocket string
	reportFn           metricsReporter
}

func (s *downwardMetricsServer) start(context context.Context) error {
	conn, err := connect(s.virtioSerialSocket, s.maxConnectAttempts)
	if err != nil {
		return err
	}

	go s.serve(context, conn)
	return nil
}

func (s *downwardMetricsServer) serve(context context.Context, conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Log.Reason(err).Warning("virtio-serial failed to close the connection")
		}
	}(conn)

	type reqResult struct {
		request string
		err     error
	}
	newRequest := make(chan reqResult)
	reader := bufio.NewReader(conn)

	for {
		// The virtio-serial vhostmd server implementation serves one request at a time,
		// so we make sure the client inside the guest cannot send another request before
		// the previous request was processed
		go func() {
			requestLine, err := waitForRequest(reader)
			if err != nil && errors.Is(err, net.ErrClosed) {
				// Closing the connection signals that the server should stop
				return
			}

			newRequest <- reqResult{request: requestLine, err: err}
		}()

		select {
		case res := <-newRequest:
			if res.err != nil {
				if errors.Is(res.err, io.EOF) {
					// qemu closed the connection
					return
				}
				log.Log.Reason(res.err).Warning("virtio-serial socket read failed")
				// Let's clean the remaining data in the connection to prevent the following
				// requests to fail due to left over bytes from the current invalid request
				reader.Reset(conn)
				replyError(conn)
				continue
			}

			response, err := s.handleRequest(res.request)
			if err != nil {
				log.Log.Reason(err).Error("failed to process the request")
				replyError(conn)
				continue
			}

			err = reply(conn, response)
			if err != nil {
				log.Log.Reason(err).Error("failed to send the metrics")
			}
		case <-context.Done():
			return
		}
	}
}

func (s *downwardMetricsServer) handleRequest(requestLine string) ([]byte, error) {
	err := parseRequest(requestLine)
	if err != nil {
		return nil, err
	}

	response, err := s.getXmlMetrics()
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *downwardMetricsServer) getXmlMetrics() ([]byte, error) {
	var xmlMetrics []byte
	metrics, err := s.reportFn()
	if err != nil {
		xmlMetrics = []byte(emptyMetrics)
		log.Log.Reason(err).Error("failed to collect the metrics")
	} else {
		xmlMetrics, err = xml.MarshalIndent(metrics, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to encode metrics: %v", err)
		}
	}

	// `vm-dump-metrics` expects `\n\n` as a termination symbol
	xmlMetrics = append(xmlMetrics, '\n', '\n')
	return xmlMetrics, nil
}

func connect(socketPath string, attempts uint) (net.Conn, error) {
	var conn net.Conn
	var err error

	multiplier := 1
	for i := uint(0); i < attempts; i++ {
		conn, err = net.Dial("unix", socketPath)
		if err == nil {
			break
		}

		// It is only tried again in case the socket doesn't exist or not one is
		// listening on the other end yet
		if !(errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ENOENT)) {
			break
		}

		if i == attempts {
			break // it was the last attempt
		}

		backoff := time.Duration(multiplier) * time.Second
		time.Sleep(backoff)
		multiplier *= 2
	}

	return conn, err
}

func waitForRequest(bufReader *bufio.Reader) (string, error) {
	reader := textproto.NewReader(bufReader)

	// First wait for an HTTP-like line, like GET /metrics/XML
	request, err := reader.ReadLine()
	if err != nil {
		return "", err
	}

	// Then wait for a blank line
	blankLine, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	if blankLine != "" {
		return "", errors.New("malformed request missing blank line")
	}

	return request, nil
}

func parseRequest(requestLine string) error {
	method, rawUri, ok := strings.Cut(requestLine, " ")
	if !ok {
		return fmt.Errorf("malformed request: %q", requestLine)
	}

	if method != "GET" {
		return fmt.Errorf("invalid method: %q", method)
	}

	requestUri, err := url.ParseRequestURI(rawUri)
	if err != nil {
		return err
	}

	// Currently this is the only valid request
	if requestUri.Path != "/metrics/XML" {
		return fmt.Errorf("invalid request: %q", requestUri.Path)
	}

	return nil
}

func replyError(conn net.Conn) {
	err := reply(conn, []byte(invalidRequest))
	if err != nil {
		log.Log.Reason(err).Error("reply error")
	}
}

func reply(conn net.Conn, response []byte) error {
	_, err := conn.Write(response)
	if err != nil {
		return err
	}
	return nil
}

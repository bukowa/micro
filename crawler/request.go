/*
Copyright © 2020 Mateusz Kurowski

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package crawler

import (
	"context"
	"io"
	"net/http"
)

// Request wraps http.Request.
type Request interface {
	Request() *http.Request
}

// NewRequest wraps http.NewRequest.
func NewRequest(method, url string, body io.Reader) (Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r := &BaseRequest{
		request: req,
	}
	return r, nil
}

// NewRequestWithContext wraps http.NewRequestWithContext.
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	r := &BaseRequest{
		request: req,
	}
	return r, nil
}

// BaseRequest implements Request.
type BaseRequest struct {
	request *http.Request
}

// Request returns http.Request instance.
func (r *BaseRequest) Request() *http.Request {
	return r.request
}

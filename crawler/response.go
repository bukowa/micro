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
	"net/http"
	"time"
)

// Response wraps http.Response.
type Response interface {
	Request
	Response() *http.Response
	Error() error
	Time() time.Duration
}

// NewResponse creates new Response.
var NewResponse = func(crawler *Crawler, took time.Duration, req Request, res *http.Response, err error) Response {
	r := &BaseResponse{
		xrequest:  req,
		xresponse: res,
		error:     err,
		took:      took,
	}
	return r
}

// BaseResponse is basic Response implementation.
type BaseResponse struct {
	xrequest  Request
	xresponse *http.Response
	error     error
	took      time.Duration
}

// Time returns time it took to complete request.
func (r *BaseResponse) Time() time.Duration {
	return r.took
}

// Request returns underlying http.Request.
func (r *BaseResponse) Request() *http.Request {
	return r.xrequest.Request()
}

// Response returns underlying http.Response.
func (r *BaseResponse) Response() *http.Response {
	return r.xresponse
}

// Error returns error as received from http.Client.
func (r *BaseResponse) Error() error {
	return r.error
}

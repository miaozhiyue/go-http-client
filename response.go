package client

import (
	"io/ioutil"
	"net/http"
	"unsafe"
)

type Response struct {
	http.Response
	Error    error
	request  *http.Request
	decoders DecoderGroup
}

func (r *Response) Raw() (*http.Response, error) {
	if r.Error != nil {
		return nil, r.Error
	}
	return (*http.Response)(unsafe.Pointer(r)), nil
}

// Read decodes response data via decoders.
func (r *Response) Read(out interface{}) (*http.Response, error) {
	defer r.Close()

	resp, err := r.Raw()
	if err != nil {
		return resp, err
	}

	return resp, r.decoders.Decode(resp, out)
}

// ReadBytes reads response body into a byte slice.
func (r *Response) ReadBytes() ([]byte, *http.Response, error) {
	defer r.Close()

	resp, err := r.Raw()
	if err != nil {
		return nil, resp, err
	}

	bytes, err := ioutil.ReadAll(r.Body)
	return bytes, resp, err
}

// Close closes response body.
func (r *Response) Close() {
	if r != nil && r.Body != nil {
		r.Body.Close()
	}
}

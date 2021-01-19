package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/kanodoe/dremio_scrapper/internal/entity"
)

// Middleware ...
type Middleware func(DremioClient) DremioClient

// DremioClient ...
type DremioClient interface {
	Login(context.Context) (*entity.DremioLoginResponse, error)
}

// DremioProxy ...
type DremioProxy struct {
	loginEndpoint               endpoint.Endpoint
	baseURI, username, password string
	timeout                     time.Duration
}

// NewDremioClient ...
func NewDremioClient(baseURI, username, password string, timeout time.Duration) DremioClient {
	endpoint := MakeDremioLoginClientEndpoint(baseURI, timeout)
	return &DremioProxy{
		loginEndpoint: endpoint,
		baseURI:       baseURI,
		username:      username,
		password:      password,
		timeout:       timeout,
	}
}

// Login ...
func (d *DremioProxy) Login(ctx context.Context) (*entity.DremioLoginResponse, error) {

	request := &entity.DremioLoginRequest{
		UserName: d.username,
		Password: d.password,
	}

	response, err := d.loginEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	res := response.(*entity.DremioLoginResponse)
	return res, nil
}

// MakeDremioLoginClientEndpoint ...
func MakeDremioLoginClientEndpoint(baseURI string, timeout time.Duration) endpoint.Endpoint {
	url, _ := url.Parse(baseURI + "/apiv2/login")
	client := http.DefaultClient
	client.Timeout = timeout

	return httptransport.NewClient(
		"POST",
		url,
		encodeDremioLoginRequest,
		decodeDremioLoginResponse,
		httptransport.SetClient(&http.Client{Timeout: timeout}),
	).Endpoint()
}

func encodeDremioLoginRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(*entity.DremioLoginRequest)
	r.Header.Add("Content-Type", "application/json")

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(req)

	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeDremioLoginResponse(_ context.Context, r *http.Response) (interface{}, error) {
	dremioLoginResponse := new(entity.DremioLoginResponse)
	if err := json.NewDecoder(r.Body).Decode(dremioLoginResponse); err != nil {
		return nil, err
	}

	return dremioLoginResponse, nil
}

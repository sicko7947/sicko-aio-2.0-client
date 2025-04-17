package psychoclient

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"time"

	http "github.com/zMrKrabz/fhttp"

	"github.com/gogf/gf/container/gmap"
	"github.com/google/uuid"
	"sicko-aio-2.0-client/models"
)

// Session interface allows operating a psycho client session
type Session interface {
	Do(id string, remove ...bool) (res *http.Response, respBody []byte, err *models.Error)
	RoundTrip(id string, remove ...bool) (res *http.Response, respBody []byte, err *models.Error)
	RemoveRequest(id string)
	Close()
	BuildRequest(req *RequestBuilder) (string, *models.Error)
	GetCookie(key string) (*http.Cookie, bool)
	GetAllCookies() map[string]*http.Cookie
	DeleteCookie(key string)
	SetCookies(cookeis map[string]*http.Cookie) bool
	ChangeRedirectBehavier(followRedirects bool)
}

type session struct {
	reqGmap *gmap.StrAnyMap
	cookies *gmap.StrAnyMap
	client  http.Client
}

// SessionBuilder : builder for making a new psycho client session
type SessionBuilder struct {
	UseDefaultClient bool
	FollowRedirects  bool
	Timeout          time.Duration
	Proxy            string
}

// RequestBuilder : builder for making a new psycho client session
type RequestBuilder struct {
	Endpoint string
	Method   string
	Headers  map[string]string
	Payload  io.Reader
}

type response struct {
	httpResponse *http.Response
	err          *models.Error
}

// New returns a new PsychoClient Request Session and an custom error for the configuration set by user
// Certificates, Private key, ClienthelloID and ClienthelloSpecs are passed to a new http request client
// Proxy also included on creating new client
func NewSession(b *SessionBuilder) (Session, *models.Error) {

	client, err := newClient(b.UseDefaultClient, b.FollowRedirects, defaultTimeout, b.Proxy)
	if err != nil {
		return nil, &models.Error{Error: err, Code: 500, Message: "Error Creating Request Session"}
	}

	return &session{
		reqGmap: gmap.NewStrAnyMap(true),
		cookies: gmap.NewStrAnyMap(true),
		client:  client,
	}, nil
}

func (a *session) Close() {
	a.cookies.Clear()
	a.reqGmap.Clear()
	a.client.CloseIdleConnections()
}

func (a *session) RemoveRequest(id string) {
	a.reqGmap.Remove(id)
}

// Build Request build an new http request and returning with an build request id that saved in a concurrent map
func (a *session) BuildRequest(builder *RequestBuilder) (string, *models.Error) {
	req, err := http.NewRequest(builder.Method, builder.Endpoint, builder.Payload) // setup request method & endpoint & payload'
	if err != nil {
		return "", &models.Error{Error: err, Code: 500, Message: "Error Builder Request"}
	}
	if len(builder.Headers) > 0 { // set the request headers
		for key, value := range builder.Headers {
			req.Header.Add(key, value)
		}
	}
	if a.cookies.Size() > 0 { // set the request cookie from the session
		a.cookies.Iterator(func(k string, v interface{}) bool {
			cookie := v.(*http.Cookie)
			if strings.Contains(cookie.Value, `"`) {
				return true
			}
			// if cookie.Domain == builder.Host {
			req.AddCookie(cookie)
			// }
			return true
		})
	}

	reqID := uuid.NewString()
	a.reqGmap.Set(reqID, req)
	return reqID, nil
}

// Do execute the http request by it's given request id
// reqID stands for the return id after created new request with request builder
// remove is that permists the session to remove request id after execution, false for not removing, default removes
func (a *session) Do(reqID string, remove ...bool) (res *http.Response, respBody []byte, err *models.Error) {
	channel := make(chan *response, 1) // setup request channel
	defer close(channel)

	go func() {

		var req interface{}
		if len(remove) > 0 && !remove[0] {
			req = a.reqGmap.Get(reqID)
		} else {
			req = a.reqGmap.Remove(reqID)
		}

		if req == nil {
			channel <- &response{
				httpResponse: nil,
				err:          &models.Error{Error: errors.New("ERROR_GETTING_REQUEST"), Code: 500, Message: "Internal Server Error"},
			}
			return
		}

		res, err := a.client.Do(req.(*http.Request))
		if err != nil || res == nil { // checking if empty response or err occured during the request
			channel <- &response{
				httpResponse: nil,
				err:          &models.Error{Error: err, Code: 500, Message: "Error Requesting"},
			}
			return
		}

		channel <- &response{ // puts the response back to channel
			httpResponse: res,
			err:          nil,
		}
	}()

	response := <-channel // waiting for channel to receive response

	res = response.httpResponse
	err = response.err

	if res != nil {
		// set new cookies to the session cookies
		newCookies := readSetCookies(res.Header)
		for k, v := range newCookies {
			a.cookies.Set(k, v)
		}

		// set response body
		if res.Body != nil {
			defer response.httpResponse.Body.Close()
			body, e := ioutil.ReadAll(res.Body)
			respBody = body
			if e != nil {
				err = &models.Error{Error: errors.New("ERROR_PARSING_RESPONSE_BODY"), Code: 500, Message: "Error Parsing Response Body"}
			}
		}
	}

	return res, respBody, err
}

// RoundTrip execute the http request by it's given request id
// reqID stands for the return id after created new request with request builder
// remove is that permists the session to remove request id after execution, false for not removing, default removes
func (a *session) RoundTrip(reqID string, remove ...bool) (res *http.Response, respBody []byte, err *models.Error) {
	channel := make(chan *response, 1) // setup request channel
	defer close(channel)

	go func() {

		var req interface{}
		if len(remove) > 0 && !remove[0] {
			req = a.reqGmap.Get(reqID)
		} else {
			req = a.reqGmap.Remove(reqID)
		}

		if req == nil {
			channel <- &response{
				httpResponse: nil,
				err:          &models.Error{Error: errors.New("ERROR_GETTING_REQUEST"), Code: 500, Message: "Internal Server Error"},
			}
			return
		}

		res, err := a.client.Transport.RoundTrip(req.(*http.Request))
		if err != nil || res == nil { // checking if empty response or err occured during the request
			channel <- &response{
				httpResponse: nil,
				err:          &models.Error{Error: err, Code: 500, Message: "Error Requesting"},
			}
			return
		}

		channel <- &response{ // puts the response back to channel
			httpResponse: res,
			err:          nil,
		}
	}()

	response := <-channel // waiting for channel to receive response

	res = response.httpResponse
	err = response.err

	if res != nil {
		// set new cookies to the session cookies
		newCookies := readSetCookies(res.Header)
		for k, v := range newCookies {
			a.cookies.Set(k, v)
		}

		// set response body
		if res.Body != nil {
			defer response.httpResponse.Body.Close()
			body, e := ioutil.ReadAll(res.Body)
			respBody = body
			if e != nil {
				err = &models.Error{Error: errors.New("ERROR_PARSING_RESPONSE_BODY"), Code: 500, Message: "Error Parsing Response Body"}
			}
		}
	}

	return res, respBody, err
}

// GetCookie function gets the cookie from the session by it's given key
// return nil on failure
func (a *session) GetCookie(key string) (*http.Cookie, bool) {
	if cookie := a.cookies.Get(key); cookie != nil {
		return cookie.(*http.Cookie), true
	}
	return nil, false
}

// GetAllCookies function gets all the cookie from the session
// return nil on failure
func (a *session) GetAllCookies() map[string]*http.Cookie {
	if a.cookies.Size() == 0 {
		return make(map[string]*http.Cookie)
	}

	cookies := make(map[string]*http.Cookie)
	a.cookies.Iterator(func(k string, v interface{}) bool {
		cookies[k] = v.(*http.Cookie)
		return true
	})
	return cookies
}

// SetCookies function sets the session cookies with the given cookie map
func (a *session) SetCookies(cookies map[string]*http.Cookie) bool {
	for k, v := range cookies {
		if strings.Contains(v.Value, `"`) {
			continue
		}
		a.cookies.Set(k, v)
	}
	return true
}

// DeleteCookie function removes the session cookie by key
func (a *session) DeleteCookie(key string) {
	a.cookies.Remove(key)
}

// ChangeRedirectBehavier function resets the redirect behavior
func (a *session) ChangeRedirectBehavier(followRedirects bool) {
	if followRedirects {
		a.client.CheckRedirect = nil
	} else {
		a.client.CheckRedirect = noRedirects
	}
}

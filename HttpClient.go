package gohttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	hlpr "github.com/rawansuww/go-httpwrapper/helpers"
	"github.com/rawansuww/go-httpwrapper/interfaces"
	"github.com/rawansuww/go-httpwrapper/types"
)

type httpClient struct {
	client  http.Client
	tracer  interfaces.ITracer
	headers map[string]string
	tid     string
	pid     string
	flg     string
}

func (httpClient *httpClient) Get(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.action("GET", headers, endpoint, qParam, params...)
}

func (httpClient *httpClient) Post(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.action("POST", headers, endpoint, qParam, params...)
}

func (httpClient *httpClient) Patch(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.action("PATCH", headers, endpoint, qParam, params...)
}

func (httpClient *httpClient) Put(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.action("PUT", headers, endpoint, qParam, params...)
}

func (httpClient *httpClient) Delete(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.action("DELETE", headers, endpoint, qParam, params...)
}

func (httpClient *httpClient) PostBody(headers map[string]string, body interface{}, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.actionBody("POST", headers, body, endpoint, qParam, params...)
}

func (httpClient *httpClient) PatchBody(headers map[string]string, body interface{}, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.actionBody("PATCH", headers, body, endpoint, qParam, params...)
}

func (httpClient *httpClient) PutBody(headers map[string]string, body interface{}, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.actionBody("PUT", headers, body, endpoint, qParam, params...)
}

func (httpClient *httpClient) DeleteBody(headers map[string]string, body interface{}, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.actionBody("DELETE", headers, body, endpoint, qParam, params...)
}

func (httpClient *httpClient) PostForm(headers map[string]string, form url.Values, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.actionForm("POST", headers, form, endpoint, qParam, params...)
}

func (httpClient *httpClient) PatchForm(headers map[string]string, form url.Values, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.actionForm("PATCH", headers, form, endpoint, qParam, params...)
}

func (httpClient *httpClient) PutForm(headers map[string]string, form url.Values, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.actionForm("PUT", headers, form, endpoint, qParam, params...)
}

func (httpClient *httpClient) DeleteForm(headers map[string]string, form url.Values, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error) {
	return httpClient.actionForm("DELETE", headers, form, endpoint, qParam, params...)
}

func (httpClient *httpClient) action(method string, headers map[string]string, endpoint string, qParam map[string][]string, pthParms ...string) (*types.Response, error) {
	endpoint, err := formatEp(endpoint, qParam, pthParms...)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	httpClient.formHeaders(req, headers)

	resp, err := httpClient.runRequest(req)
	if err != nil {
		return nil, err
	}
	respObj := types.Response(*resp)
	return &respObj, nil
}

func (httpClient *httpClient) actionBody(method string, headers map[string]string, body interface{}, endpoint string, qParam map[string][]string, pthParms ...string) (*types.Response, error) {
	endpoint, err := formatEp(endpoint, qParam, pthParms...)
	if err != nil {
		return nil, err
	}

	byts, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewReader(byts))
	if err != nil {
		return nil, err
	}

	httpClient.formHeaders(req, headers)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.runRequest(req)
	if err != nil {
		return nil, err
	}
	respObj := types.Response(*resp)
	return &respObj, nil
}

func (httpClient *httpClient) actionForm(method string, headers map[string]string, form url.Values, endpoint string, qParam map[string][]string, pthParms ...string) (*types.Response, error) {
	endpoint, err := formatEp(endpoint, qParam, pthParms...)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		method,
		endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	httpClient.formHeaders(req, headers)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.runRequest(req)
	if err != nil {
		return nil, err
	}
	respObj := types.Response(*resp)
	return &respObj, nil
}

func (httpClient *httpClient) formHeaders(req *http.Request, headers map[string]string) {
	for key, value := range httpClient.headers {
		req.Header.Add(key, value)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func (httpClient *httpClient) runRequest(req *http.Request) (*http.Response, error) {
	sid, err := hlpr.GenerateParentId()
	if err == nil {
		req.Header.Add(
			"traceparent",
			fmt.Sprintf("00-%s-%s-%s", httpClient.tid, sid, httpClient.flg),
		)
	} else {
		req.Header.Add(
			"traceparent",
			fmt.Sprintf("00-%s-%s-%s", httpClient.tid, httpClient.pid, httpClient.flg),
		)
	}
	start := time.Now()
	resp, err := httpClient.client.Do(req)
	end := time.Now()

	if err != nil {
		httpClient.tracer.TraceDependency(
			sid,
			"http",
			req.URL.Hostname(),
			fmt.Sprintf("%s %s", req.Method, req.URL.RequestURI()),
			false,
			start,
			end,
			types.NewField("method", req.Method),
			types.NewField("error", err.Error()),
		)
		return nil, err
	}
	httpClient.tracer.TraceDependency(
		sid,
		"http",
		req.URL.Hostname(),
		fmt.Sprintf("%s %s", req.Method, req.URL.RequestURI()),
		resp.StatusCode > 199 && resp.StatusCode < 300,
		start,
		end,
		types.NewField("method", req.Method),
		types.NewField("statusCode", strconv.Itoa(resp.StatusCode)),
	)
	return resp, err
}

// TODO Pre calculating length and allocating might improve performance
func formatEp(format string, qParam url.Values, pthParms ...string) (string, error) {
	end := len(format)
	prmCnt := len(pthParms)
	pthNum := 0
	var buffer []byte
	i := 0
	prev := 0
	for i < end {
		for i < end && format[i] != '{' {
			i++
		}
		if i == end {
			break
		}
		if format[i+1] != '}' {
			return "", fmt.Errorf("illegal character/Invalid format in url")
		}
		if pthNum >= prmCnt {
			return "", fmt.Errorf("not enough parameters provided")
		}
		// TODO Maybe can be done in one go?
		escaped := url.QueryEscape(pthParms[pthNum])
		buffer = append(buffer, format[prev:i]...)
		buffer = append(buffer, escaped...)
		pthNum++
		i += 2
		prev = i
	}
	if pthNum != prmCnt {
		return "", fmt.Errorf("too many parameters provided")
	}
	if prev < end {
		buffer = append(buffer, format[prev:end]...)
	}

	// TODO could be done in parallel, performance needs to be tested
	// TODO found out that url.Values has an Encode funtion that does this,
	//      need to test
	qryBuf := []byte("?")

	for key, vals := range qParam {
		esKey := url.QueryEscape(key)
		for _, val := range vals {
			// TODO just a spike, need to experiment
			query := esKey + "=" + url.QueryEscape(val) + "&"
			qryBuf = append(qryBuf, query...)
		}
	}
	buffer = append(buffer, qryBuf...)
	return string(buffer), nil
}

var _ interfaces.HttpClient = (*httpClient)(nil)

// - "Constructors"
func NewHttpClientProvider(tracer interfaces.ITracer, headers map[string]string, tid string, pid string, flg string) interfaces.HttpClient {
	return &httpClient{
		client:  *http.DefaultClient,
		tracer:  tracer,
		headers: headers,
		tid:     tid,
		pid:     pid,
		flg:     flg,
	}

}

//-------

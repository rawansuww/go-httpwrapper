package interfaces

import (
	"net/url"
	"time"

	"github.com/rawansuww/go-httpwrapper/types"
)

type HttpClient interface {
	Get(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	Post(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	Put(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	Patch(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	Delete(headers map[string]string, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	PostBody(headers map[string]string, body interface{}, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	PatchBody(headers map[string]string, body interface{}, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	PutBody(headers map[string]string, body interface{}, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	DeleteBody(headers map[string]string, body interface{}, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	PostForm(headers map[string]string, form url.Values, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	PatchForm(headers map[string]string, form url.Values, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	PutForm(headers map[string]string, form url.Values, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
	DeleteForm(headers map[string]string, form url.Values, endpoint string, qParam map[string][]string, params ...string) (*types.Response, error)
}

type ITracer interface {
	TraceRequest(isParent bool, method string, path string, query string, statusCode int, bodySize int, ip string, userAgent string, startTimestamp time.Time, eventTimestamp time.Time, fields ...types.Field)
	TraceDependency(spanId string, dependencyType string, serviceName string, commandName string, success bool, startTimestamp time.Time, eventTimestamp time.Time, fields ...types.Field)
}

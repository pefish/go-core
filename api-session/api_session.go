package api_session

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	_interface "github.com/pefish/go-core/api-session/interface"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type ApiMethod string

const (
	ApiMethod_Post   ApiMethod = `POST`
	ApiMethod_Get    ApiMethod = `GET`
	ApiMethod_Option ApiMethod = `OPTIONS`
	ApiMethod_All    ApiMethod = `ALL`
)

type StatusCode int

const (
	StatusCode_Continue           StatusCode = 100 // RFC 7231, 6.2.1
	StatusCode_SwitchingProtocols StatusCode = 101 // RFC 7231, 6.2.2
	StatusCode_Processing         StatusCode = 102 // RFC 2518, 10.1

	StatusCode_OK                   StatusCode = 200 // RFC 7231, 6.3.1
	StatusCode_Created              StatusCode = 201 // RFC 7231, 6.3.2
	StatusCode_Accepted             StatusCode = 202 // RFC 7231, 6.3.3
	StatusCode_NonAuthoritativeInfo StatusCode = 203 // RFC 7231, 6.3.4
	StatusCode_NoContent            StatusCode = 204 // RFC 7231, 6.3.5
	StatusCode_ResetContent         StatusCode = 205 // RFC 7231, 6.3.6
	StatusCode_PartialContent       StatusCode = 206 // RFC 7233, 4.1
	StatusCode_MultiStatus          StatusCode = 207 // RFC 4918, 11.1
	StatusCode_AlreadyReported      StatusCode = 208 // RFC 5842, 7.1
	StatusCode_IMUsed               StatusCode = 226 // RFC 3229, 10.4.1

	StatusCode_MultipleChoices  StatusCode = 300 // RFC 7231, 6.4.1
	StatusCode_MovedPermanently StatusCode = 301 // RFC 7231, 6.4.2
	StatusCode_Found            StatusCode = 302 // RFC 7231, 6.4.3
	StatusCode_SeeOther         StatusCode = 303 // RFC 7231, 6.4.4
	StatusCode_NotModified      StatusCode = 304 // RFC 7232, 4.1
	StatusCode_UseProxy         StatusCode = 305 // RFC 7231, 6.4.5

	StatusCode_TemporaryRedirect StatusCode = 307 // RFC 7231, 6.4.7
	StatusCode_PermanentRedirect StatusCode = 308 // RFC 7538, 3

	StatusCode_BadRequest                   StatusCode = 400 // RFC 7231, 6.5.1
	StatusCode_Unauthorized                 StatusCode = 401 // RFC 7235, 3.1
	StatusCode_PaymentRequired              StatusCode = 402 // RFC 7231, 6.5.2
	StatusCode_Forbidden                    StatusCode = 403 // RFC 7231, 6.5.3
	StatusCode_NotFound                     StatusCode = 404 // RFC 7231, 6.5.4
	StatusCode_MethodNotAllowed             StatusCode = 405 // RFC 7231, 6.5.5
	StatusCode_NotAcceptable                StatusCode = 406 // RFC 7231, 6.5.6
	StatusCode_ProxyAuthRequired            StatusCode = 407 // RFC 7235, 3.2
	StatusCode_RequestTimeout               StatusCode = 408 // RFC 7231, 6.5.7
	StatusCode_Conflict                     StatusCode = 409 // RFC 7231, 6.5.8
	StatusCode_Gone                         StatusCode = 410 // RFC 7231, 6.5.9
	StatusCode_LengthRequired               StatusCode = 411 // RFC 7231, 6.5.10
	StatusCode_PreconditionFailed           StatusCode = 412 // RFC 7232, 4.2
	StatusCode_RequestEntityTooLarge        StatusCode = 413 // RFC 7231, 6.5.11
	StatusCode_RequestURITooLong            StatusCode = 414 // RFC 7231, 6.5.12
	StatusCode_UnsupportedMediaType         StatusCode = 415 // RFC 7231, 6.5.13
	StatusCode_RequestedRangeNotSatisfiable StatusCode = 416 // RFC 7233, 4.4
	StatusCode_ExpectationFailed            StatusCode = 417 // RFC 7231, 6.5.14
	StatusCode_Teapot                       StatusCode = 418 // RFC 7168, 2.3.3
	StatusCode_MisdirectedRequest           StatusCode = 421 // RFC 7540, 9.1.2
	StatusCode_UnprocessableEntity          StatusCode = 422 // RFC 4918, 11.2
	StatusCode_Locked                       StatusCode = 423 // RFC 4918, 11.3
	StatusCode_FailedDependency             StatusCode = 424 // RFC 4918, 11.4
	StatusCode_TooEarly                     StatusCode = 425 // RFC 8470, 5.2.
	StatusCode_UpgradeRequired              StatusCode = 426 // RFC 7231, 6.5.15
	StatusCode_PreconditionRequired         StatusCode = 428 // RFC 6585, 3
	StatusCode_TooManyRequests              StatusCode = 429 // RFC 6585, 4
	StatusCode_RequestHeaderFieldsTooLarge  StatusCode = 431 // RFC 6585, 5
	StatusCode_UnavailableForLegalReasons   StatusCode = 451 // RFC 7725, 3

	StatusCode_InternalServerError           StatusCode = 500 // RFC 7231, 6.6.1
	StatusCode_NotImplemented                StatusCode = 501 // RFC 7231, 6.6.2
	StatusCode_BadGateway                    StatusCode = 502 // RFC 7231, 6.6.3
	StatusCode_ServiceUnavailable            StatusCode = 503 // RFC 7231, 6.6.4
	StatusCode_GatewayTimeout                StatusCode = 504 // RFC 7231, 6.6.5
	StatusCode_HTTPVersionNotSupported       StatusCode = 505 // RFC 7231, 6.6.6
	StatusCode_VariantAlsoNegotiates         StatusCode = 506 // RFC 2295, 8.1
	StatusCode_InsufficientStorage           StatusCode = 507 // RFC 4918, 11.5
	StatusCode_LoopDetected                  StatusCode = 508 // RFC 5842, 7.2
	StatusCode_NotExtended                   StatusCode = 510 // RFC 2774, 7
	StatusCode_NetworkAuthenticationRequired StatusCode = 511 // RFC 6585, 6
)

type HeaderName string

const (
	// ContentTypeHeaderKey is the header key of "Content-Type".
	HeaderName_ContentType HeaderName = "Content-Type"

	// LastModifiedHeaderKey is the header key of "Last-Modified".
	HeaderName_LastModified HeaderName = "Last-Modified"
	// IfModifiedSinceHeaderKey is the header key of "If-Modified-Since".
	HeaderName_IfModifiedSince HeaderName = "If-Modified-Since"
	// CacheControlHeaderKey is the header key of "Cache-Control".
	HeaderName_CacheControl HeaderName = "Cache-Control"
	// ETagHeaderKey is the header key of "ETag".
	HeaderName_ETag HeaderName = "ETag"

	// ContentDispositionHeaderKey is the header key of "Content-Disposition".
	HeaderName_ContentDisposition HeaderName = "Content-Disposition"
	// ContentLengthHeaderKey is the header key of "Content-Length"
	HeaderName_ContentLength HeaderName = "Content-Length"
	// ContentEncodingHeaderKey is the header key of "Content-Encoding".
	HeaderName_ContentEncoding HeaderName = "Content-Encoding"
	// GzipHeaderValue is the header value of "gzip".
	HeaderName_Gzip HeaderName = "gzip"
	// AcceptEncodingHeaderKey is the header key of "Accept-Encoding".
	HeaderName_AcceptEncoding HeaderName = "Accept-Encoding"
	// VaryHeaderKey is the header key of "Vary".
	HeaderName_Vary HeaderName = "Vary"
)

type ContentTypeValue string

const (
	// ContentBinaryHeaderValue header value for binary data.
	ContentTypeValue_Binary ContentTypeValue = "application/octet-stream"
	// ContentHTMLHeaderValue is the  string of text/html response header's content type value.
	ContentTypeValue_HTML ContentTypeValue = "text/html; charset=UTF-8"
	// ContentJSONHeaderValue header value for JSON data.
	ContentTypeValue_JSON ContentTypeValue = "application/json; charset=UTF-8"
	// ContentJavascriptHeaderValue header value for JSONP & Javascript data.
	ContentTypeValue_Javascript ContentTypeValue = "application/javascript; charset=UTF-8"
	// ContentTextHeaderValue header value for Text data.
	ContentTypeValue_Text ContentTypeValue = "text/plain; charset=UTF-8"
	// ContentXMLHeaderValue header value for XML data.
	ContentTypeValue_XML ContentTypeValue = "text/xml; charset=UTF-8"
	// ContentMarkdownHeaderValue custom key/content type, the real is the text/html.
	ContentTypeValue_Markdown ContentTypeValue = "text/markdown; charset=UTF-8"
	// ContentYAMLHeaderValue header value for YAML data.
	ContentTypeValue_YAML ContentTypeValue = "application/x-yaml; charset=UTF-8"
)

type ApiSessionClass struct {
	statusCode StatusCode

	Api            _interface.InterfaceApi
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	JwtHeaderName string
	JwtBody       map[string]interface{}
	UserId        uint64

	Lang       string
	ClientType string // web、android、ios

	Datas map[string]interface{}

	OriginalParams map[string]interface{} // 客户端传过来的原始参数
	Params         map[string]interface{} // 经过前置处理器修饰过的参数

	Defers []func() // api结束后执行的函数
}

func NewApiSession() *ApiSessionClass {
	return &ApiSessionClass{
		Datas: map[string]interface{}{},
	}
}

func (apiSession *ApiSessionClass) ScanParams(dest interface{}) {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "json",
		Result:           &dest,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(apiSession.Params)
	if err != nil {
		panic(err)
	}
}

func (apiSession *ApiSessionClass) AddDefer(defer_ func()) {
	apiSession.Defers = append(apiSession.Defers, defer_)
}

func (apiSession *ApiSessionClass) WriteJson(data interface{}) error {
	apiSession.SetHeader(string(HeaderName_ContentType), string(ContentTypeValue_JSON))
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	apiSession.ResponseWriter.WriteHeader(int(apiSession.statusCode))
	_, err = apiSession.ResponseWriter.Write(result)
	if err != nil {
		return err
	}
	return nil
}

func (apiSession *ApiSessionClass) SetHeader(key string, value string) {
	apiSession.ResponseWriter.Header().Set(key, value)
}

func (apiSession *ApiSessionClass) WriteText(text string) error {
	apiSession.ResponseWriter.Header().Set(string(HeaderName_ContentType), string(ContentTypeValue_Text))
	apiSession.ResponseWriter.WriteHeader(int(apiSession.statusCode))
	_, err := apiSession.ResponseWriter.Write([]byte(text))
	if err != nil {
		return err
	}
	return nil
}

func (apiSession *ApiSessionClass) SetStatusCode(code StatusCode) {
	apiSession.statusCode = code
}

func (apiSession *ApiSessionClass) GetPath() string {
	return apiSession.Request.URL.Path
}

func (apiSession *ApiSessionClass) GetMethod() string {
	return apiSession.Request.Method
}

func (apiSession *ApiSessionClass) GetHeader(name string) string {
	return apiSession.Request.Header.Get(name)
}

func (apiSession *ApiSessionClass) GetRemoteAddress() string {
	remoteHeaders := map[string]bool{
		`X-Forwarded-For`: true,
	}

	for headerName, enabled := range remoteHeaders {
		if enabled {
			headerValue := apiSession.GetHeader(headerName)
			// exception needed for 'X-Forwarded-For' only , if enabled.
			if headerName == "X-Forwarded-For" {
				idx := strings.IndexByte(headerValue, ',')
				if idx >= 0 {
					headerValue = headerValue[0:idx]
				}
			}

			realIP := strings.TrimSpace(headerValue)
			if realIP != "" {
				return realIP
			}
		}
	}

	addr := strings.TrimSpace(apiSession.Request.RemoteAddr)
	if addr != "" {
		// if addr has port use the net.SplitHostPort otherwise(error occurs) take as it is
		if ip, _, err := net.SplitHostPort(addr); err == nil {
			return ip
		}
	}

	return addr
}

func (apiSession *ApiSessionClass) GetUrlParams() map[string]string {
	values := map[string]string{}

	q := apiSession.Request.URL.Query()
	if q != nil {
		for k, v := range q {
			values[k] = strings.Join(v, ",")
		}
	}

	return values
}

func (apiSession *ApiSessionClass) GetFormValues() (map[string][]string, error) {
	err := apiSession.Request.ParseMultipartForm(32 << 20) // 默认32M
	if err != nil {
		return nil, err
	}
	var form map[string][]string
	if form := apiSession.Request.Form; len(form) > 0 {
		return form, nil
	}

	if form := apiSession.Request.PostForm; len(form) > 0 {
		return form, nil
	}

	if m := apiSession.Request.MultipartForm; m != nil {
		if len(m.Value) > 0 {
			return m.Value, nil
		}
	}

	return form, nil
}

func (apiSession *ApiSessionClass) ReadJSON(jsonObject interface{}) error {
	if apiSession.Request.Body == nil {
		return errors.New("unmarshal: empty body")
	}

	rawData, err := ioutil.ReadAll(apiSession.Request.Body)
	if err != nil {
		return err
	}

	apiSession.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))

	return json.Unmarshal(rawData, jsonObject)
}

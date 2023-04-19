package api_session

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/pefish/go-core/driver/logger"
	go_logger "github.com/pefish/go-logger"

	"github.com/mitchellh/mapstructure"
	_interface "github.com/pefish/go-core-type/api"
	_type "github.com/pefish/go-core-type/api-session"
	"io"
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

const (
	StatusCode_Continue           _type.StatusCode = 100 // RFC 7231, 6.2.1
	StatusCode_SwitchingProtocols _type.StatusCode = 101 // RFC 7231, 6.2.2
	StatusCode_Processing         _type.StatusCode = 102 // RFC 2518, 10.1

	StatusCode_OK                   _type.StatusCode = 200 // RFC 7231, 6.3.1
	StatusCode_Created              _type.StatusCode = 201 // RFC 7231, 6.3.2
	StatusCode_Accepted             _type.StatusCode = 202 // RFC 7231, 6.3.3
	StatusCode_NonAuthoritativeInfo _type.StatusCode = 203 // RFC 7231, 6.3.4
	StatusCode_NoContent            _type.StatusCode = 204 // RFC 7231, 6.3.5
	StatusCode_ResetContent         _type.StatusCode = 205 // RFC 7231, 6.3.6
	StatusCode_PartialContent       _type.StatusCode = 206 // RFC 7233, 4.1
	StatusCode_MultiStatus          _type.StatusCode = 207 // RFC 4918, 11.1
	StatusCode_AlreadyReported      _type.StatusCode = 208 // RFC 5842, 7.1
	StatusCode_IMUsed               _type.StatusCode = 226 // RFC 3229, 10.4.1

	StatusCode_MultipleChoices  _type.StatusCode = 300 // RFC 7231, 6.4.1
	StatusCode_MovedPermanently _type.StatusCode = 301 // RFC 7231, 6.4.2
	StatusCode_Found            _type.StatusCode = 302 // RFC 7231, 6.4.3
	StatusCode_SeeOther         _type.StatusCode = 303 // RFC 7231, 6.4.4
	StatusCode_NotModified      _type.StatusCode = 304 // RFC 7232, 4.1
	StatusCode_UseProxy         _type.StatusCode = 305 // RFC 7231, 6.4.5

	StatusCode_TemporaryRedirect _type.StatusCode = 307 // RFC 7231, 6.4.7
	StatusCode_PermanentRedirect _type.StatusCode = 308 // RFC 7538, 3

	StatusCode_BadRequest                   _type.StatusCode = 400 // RFC 7231, 6.5.1
	StatusCode_Unauthorized                 _type.StatusCode = 401 // RFC 7235, 3.1
	StatusCode_PaymentRequired              _type.StatusCode = 402 // RFC 7231, 6.5.2
	StatusCode_Forbidden                    _type.StatusCode = 403 // RFC 7231, 6.5.3
	StatusCode_NotFound                     _type.StatusCode = 404 // RFC 7231, 6.5.4
	StatusCode_MethodNotAllowed             _type.StatusCode = 405 // RFC 7231, 6.5.5
	StatusCode_NotAcceptable                _type.StatusCode = 406 // RFC 7231, 6.5.6
	StatusCode_ProxyAuthRequired            _type.StatusCode = 407 // RFC 7235, 3.2
	StatusCode_RequestTimeout               _type.StatusCode = 408 // RFC 7231, 6.5.7
	StatusCode_Conflict                     _type.StatusCode = 409 // RFC 7231, 6.5.8
	StatusCode_Gone                         _type.StatusCode = 410 // RFC 7231, 6.5.9
	StatusCode_LengthRequired               _type.StatusCode = 411 // RFC 7231, 6.5.10
	StatusCode_PreconditionFailed           _type.StatusCode = 412 // RFC 7232, 4.2
	StatusCode_RequestEntityTooLarge        _type.StatusCode = 413 // RFC 7231, 6.5.11
	StatusCode_RequestURITooLong            _type.StatusCode = 414 // RFC 7231, 6.5.12
	StatusCode_UnsupportedMediaType         _type.StatusCode = 415 // RFC 7231, 6.5.13
	StatusCode_RequestedRangeNotSatisfiable _type.StatusCode = 416 // RFC 7233, 4.4
	StatusCode_ExpectationFailed            _type.StatusCode = 417 // RFC 7231, 6.5.14
	StatusCode_Teapot                       _type.StatusCode = 418 // RFC 7168, 2.3.3
	StatusCode_MisdirectedRequest           _type.StatusCode = 421 // RFC 7540, 9.1.2
	StatusCode_UnprocessableEntity          _type.StatusCode = 422 // RFC 4918, 11.2
	StatusCode_Locked                       _type.StatusCode = 423 // RFC 4918, 11.3
	StatusCode_FailedDependency             _type.StatusCode = 424 // RFC 4918, 11.4
	StatusCode_TooEarly                     _type.StatusCode = 425 // RFC 8470, 5.2.
	StatusCode_UpgradeRequired              _type.StatusCode = 426 // RFC 7231, 6.5.15
	StatusCode_PreconditionRequired         _type.StatusCode = 428 // RFC 6585, 3
	StatusCode_TooManyRequests              _type.StatusCode = 429 // RFC 6585, 4
	StatusCode_RequestHeaderFieldsTooLarge  _type.StatusCode = 431 // RFC 6585, 5
	StatusCode_UnavailableForLegalReasons   _type.StatusCode = 451 // RFC 7725, 3

	StatusCode_InternalServerError           _type.StatusCode = 500 // RFC 7231, 6.6.1
	StatusCode_NotImplemented                _type.StatusCode = 501 // RFC 7231, 6.6.2
	StatusCode_BadGateway                    _type.StatusCode = 502 // RFC 7231, 6.6.3
	StatusCode_ServiceUnavailable            _type.StatusCode = 503 // RFC 7231, 6.6.4
	StatusCode_GatewayTimeout                _type.StatusCode = 504 // RFC 7231, 6.6.5
	StatusCode_HTTPVersionNotSupported       _type.StatusCode = 505 // RFC 7231, 6.6.6
	StatusCode_VariantAlsoNegotiates         _type.StatusCode = 506 // RFC 2295, 8.1
	StatusCode_InsufficientStorage           _type.StatusCode = 507 // RFC 4918, 11.5
	StatusCode_LoopDetected                  _type.StatusCode = 508 // RFC 5842, 7.2
	StatusCode_NotExtended                   _type.StatusCode = 510 // RFC 2774, 7
	StatusCode_NetworkAuthenticationRequired _type.StatusCode = 511 // RFC 6585, 6
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
	statusCode _type.StatusCode

	api            _interface.IApi
	responseWriter http.ResponseWriter
	request        *http.Request

	jwtHeaderName string
	jwtBody       map[string]interface{}
	userId        uint64

	lang       string
	clientType string // web、android、ios

	data     map[string]interface{}
	pathVars map[string]string

	originalParams map[string]interface{} // 客户端传过来的原始参数
	params         map[string]interface{} // 经过前置处理器修饰过的参数

	defers []func() // api结束后执行的函数
}

func NewApiSession() *ApiSessionClass {
	return &ApiSessionClass{
		data:   make(map[string]interface{}, 5),
		defers: make([]func(), 0, 20),
	}
}

func (apiSession *ApiSessionClass) SetPathVars(vars map[string]string) {
	apiSession.pathVars = vars
}

func (apiSession *ApiSessionClass) PathVars() map[string]string {
	return apiSession.pathVars
}

func (apiSession *ApiSessionClass) SetJwtHeaderName(headerName string) {
	apiSession.jwtHeaderName = headerName
}

func (apiSession *ApiSessionClass) JwtHeaderName() string {
	return apiSession.jwtHeaderName
}

func (apiSession *ApiSessionClass) SetLang(lang string) {
	apiSession.lang = lang
}

func (apiSession *ApiSessionClass) Lang() string {
	return apiSession.lang
}

func (apiSession *ApiSessionClass) SetClientType(clientType string) {
	apiSession.clientType = clientType
}

func (apiSession *ApiSessionClass) ClientType() string {
	return apiSession.clientType
}

func (apiSession *ApiSessionClass) SetJwtBody(jwtBody map[string]interface{}) {
	apiSession.jwtBody = jwtBody
}

func (apiSession *ApiSessionClass) ResponseWriter() http.ResponseWriter {
	return apiSession.responseWriter
}

func (apiSession *ApiSessionClass) Request() *http.Request {
	return apiSession.request
}

func (apiSession *ApiSessionClass) SetResponseWriter(w http.ResponseWriter) {
	apiSession.responseWriter = w
}

func (apiSession *ApiSessionClass) SetRequest(r *http.Request) {
	apiSession.request = r
}

func (apiSession *ApiSessionClass) JwtBody() map[string]interface{} {
	return apiSession.jwtBody
}

func (apiSession *ApiSessionClass) Params() map[string]interface{} {
	return apiSession.params
}

func (apiSession *ApiSessionClass) OriginalParams() map[string]interface{} {
	return apiSession.originalParams
}

func (apiSession *ApiSessionClass) SetParams(params map[string]interface{}) {
	apiSession.params = params
}

func (apiSession *ApiSessionClass) SetOriginalParams(originalParams map[string]interface{}) {
	apiSession.originalParams = originalParams
}

func (apiSession *ApiSessionClass) Api() _interface.IApi {
	return apiSession.api
}

func (apiSession *ApiSessionClass) SetApi(api _interface.IApi) {
	apiSession.api = api
}

func (apiSession *ApiSessionClass) SetUserId(userId uint64) {
	apiSession.userId = userId
}

func (apiSession *ApiSessionClass) UserId() uint64 {
	return apiSession.userId
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

	err = decoder.Decode(apiSession.params)
	if err != nil {
		panic(err)
	}
}

// Add defer handler.
// Defer handlers will be executed by order at the end of api session.
func (apiSession *ApiSessionClass) AddDefer(defer_ func()) {
	apiSession.defers = append(apiSession.Defers(), defer_)
}

func (apiSession *ApiSessionClass) Defers() []func() {
	return apiSession.defers
}

func (apiSession *ApiSessionClass) SetData(key string, data interface{}) {
	apiSession.data[key] = data
}

func (apiSession *ApiSessionClass) Data(key string) interface{} {
	return apiSession.data[key]
}

func (apiSession *ApiSessionClass) Redirect(url string) {
	http.Redirect(apiSession.responseWriter, apiSession.request, url, http.StatusTemporaryRedirect)
}

// Response json body.
func (apiSession *ApiSessionClass) WriteJson(data interface{}) error {
	apiSession.SetHeader(string(HeaderName_ContentType), string(ContentTypeValue_JSON))
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	apiSession.responseWriter.WriteHeader(int(apiSession.statusCode))
	_, err = apiSession.responseWriter.Write(result)
	if err != nil {
		return err
	}
	return nil
}

// Set header of response.
func (apiSession *ApiSessionClass) SetHeader(key string, value string) {
	apiSession.responseWriter.Header().Set(key, value)
}

// Response text body.
func (apiSession *ApiSessionClass) WriteText(text string) error {
	apiSession.responseWriter.Header().Set(string(HeaderName_ContentType), string(ContentTypeValue_Text))
	apiSession.responseWriter.WriteHeader(int(apiSession.statusCode))
	_, err := apiSession.responseWriter.Write([]byte(text))
	if err != nil {
		return err
	}
	return nil
}

// Set status code of response.
func (apiSession *ApiSessionClass) SetStatusCode(code _type.StatusCode) {
	apiSession.statusCode = code
}

// Get request host.
func (apiSession *ApiSessionClass) Host() string {
	return apiSession.request.Host
}

// Get request path.
func (apiSession *ApiSessionClass) Path() string {
	return apiSession.request.URL.Path
}

// Get request body.
func (apiSession *ApiSessionClass) Body() io.ReadCloser {
	return apiSession.request.Body
}

// Get request method (GET, POST, PUT, etc.).
func (apiSession *ApiSessionClass) Method() string {
	return apiSession.request.Method
}

// Read header by key from request headers.
func (apiSession *ApiSessionClass) Header(name string) string {
	return apiSession.request.Header.Get(name)
}

// Read remote address from request headers.
func (apiSession *ApiSessionClass) RemoteAddress() string {
	remoteHeaders := map[string]bool{
		`X-Forwarded-For`: true,
	}

	for headerName, enabled := range remoteHeaders {
		if enabled {
			headerValue := apiSession.Header(headerName)
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

	addr := strings.TrimSpace(apiSession.request.RemoteAddr)
	if addr != "" {
		// if addr has port use the net.SplitHostPort otherwise(error occurs) take as it is
		if ip, _, err := net.SplitHostPort(addr); err == nil {
			return ip
		}
	}

	return addr
}

// Read url params from get request.
func (apiSession *ApiSessionClass) UrlParams() map[string]string {
	values := make(map[string]string, 10)

	q := apiSession.request.URL.Query()
	if q != nil {
		for k, v := range q {
			values[k] = strings.Join(v, ",")
		}
	}

	return values
}

// Read form data from request.
func (apiSession *ApiSessionClass) FormValues() (map[string][]string, error) {
	err := apiSession.request.ParseMultipartForm(32 << 20) // 默认32M
	if err != nil {
		return nil, err
	}
	var form map[string][]string
	if form := apiSession.request.Form; len(form) > 0 {
		return form, nil
	}

	if form := apiSession.request.PostForm; len(form) > 0 {
		return form, nil
	}

	if m := apiSession.request.MultipartForm; m != nil {
		if len(m.Value) > 0 {
			return m.Value, nil
		}
	}

	return form, nil
}

// Read json data from request body.
func (apiSession *ApiSessionClass) ReadJSON(jsonObject interface{}) error {
	if apiSession.request.Body == nil {
		return errors.New("unmarshal: empty body")
	}

	rawData, err := ioutil.ReadAll(apiSession.request.Body)
	if err != nil {
		return err
	}

	apiSession.request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))

	return json.Unmarshal(rawData, jsonObject)
}

func (apiSession *ApiSessionClass) Logger() go_logger.InterfaceLogger {
	return logger.LoggerDriverInstance.Logger
}

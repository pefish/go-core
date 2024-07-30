package api_session

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/pefish/go-core/driver/logger"

	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"

	i_core "github.com/pefish/go-interface/i-core"
	i_logger "github.com/pefish/go-interface/i-logger"
	t_core "github.com/pefish/go-interface/t-core"
)

type ApiMethod string

const (
	ApiMethod_Post   ApiMethod = `POST`
	ApiMethod_Get    ApiMethod = `GET`
	ApiMethod_Option ApiMethod = `OPTIONS`
	ApiMethod_All    ApiMethod = `ALL`
)

const (
	StatusCode_Continue           t_core.StatusCode = 100 // RFC 7231, 6.2.1
	StatusCode_SwitchingProtocols t_core.StatusCode = 101 // RFC 7231, 6.2.2
	StatusCode_Processing         t_core.StatusCode = 102 // RFC 2518, 10.1

	StatusCode_OK                   t_core.StatusCode = 200 // RFC 7231, 6.3.1
	StatusCode_Created              t_core.StatusCode = 201 // RFC 7231, 6.3.2
	StatusCode_Accepted             t_core.StatusCode = 202 // RFC 7231, 6.3.3
	StatusCode_NonAuthoritativeInfo t_core.StatusCode = 203 // RFC 7231, 6.3.4
	StatusCode_NoContent            t_core.StatusCode = 204 // RFC 7231, 6.3.5
	StatusCode_ResetContent         t_core.StatusCode = 205 // RFC 7231, 6.3.6
	StatusCode_PartialContent       t_core.StatusCode = 206 // RFC 7233, 4.1
	StatusCode_MultiStatus          t_core.StatusCode = 207 // RFC 4918, 11.1
	StatusCode_AlreadyReported      t_core.StatusCode = 208 // RFC 5842, 7.1
	StatusCode_IMUsed               t_core.StatusCode = 226 // RFC 3229, 10.4.1

	StatusCode_MultipleChoices  t_core.StatusCode = 300 // RFC 7231, 6.4.1
	StatusCode_MovedPermanently t_core.StatusCode = 301 // RFC 7231, 6.4.2
	StatusCode_Found            t_core.StatusCode = 302 // RFC 7231, 6.4.3
	StatusCode_SeeOther         t_core.StatusCode = 303 // RFC 7231, 6.4.4
	StatusCode_NotModified      t_core.StatusCode = 304 // RFC 7232, 4.1
	StatusCode_UseProxy         t_core.StatusCode = 305 // RFC 7231, 6.4.5

	StatusCode_TemporaryRedirect t_core.StatusCode = 307 // RFC 7231, 6.4.7
	StatusCode_PermanentRedirect t_core.StatusCode = 308 // RFC 7538, 3

	StatusCode_BadRequest                   t_core.StatusCode = 400 // RFC 7231, 6.5.1
	StatusCode_Unauthorized                 t_core.StatusCode = 401 // RFC 7235, 3.1
	StatusCode_PaymentRequired              t_core.StatusCode = 402 // RFC 7231, 6.5.2
	StatusCode_Forbidden                    t_core.StatusCode = 403 // RFC 7231, 6.5.3
	StatusCode_NotFound                     t_core.StatusCode = 404 // RFC 7231, 6.5.4
	StatusCode_MethodNotAllowed             t_core.StatusCode = 405 // RFC 7231, 6.5.5
	StatusCode_NotAcceptable                t_core.StatusCode = 406 // RFC 7231, 6.5.6
	StatusCode_ProxyAuthRequired            t_core.StatusCode = 407 // RFC 7235, 3.2
	StatusCode_RequestTimeout               t_core.StatusCode = 408 // RFC 7231, 6.5.7
	StatusCode_Conflict                     t_core.StatusCode = 409 // RFC 7231, 6.5.8
	StatusCode_Gone                         t_core.StatusCode = 410 // RFC 7231, 6.5.9
	StatusCode_LengthRequired               t_core.StatusCode = 411 // RFC 7231, 6.5.10
	StatusCode_PreconditionFailed           t_core.StatusCode = 412 // RFC 7232, 4.2
	StatusCode_RequestEntityTooLarge        t_core.StatusCode = 413 // RFC 7231, 6.5.11
	StatusCode_RequestURITooLong            t_core.StatusCode = 414 // RFC 7231, 6.5.12
	StatusCode_UnsupportedMediaType         t_core.StatusCode = 415 // RFC 7231, 6.5.13
	StatusCode_RequestedRangeNotSatisfiable t_core.StatusCode = 416 // RFC 7233, 4.4
	StatusCode_ExpectationFailed            t_core.StatusCode = 417 // RFC 7231, 6.5.14
	StatusCode_Teapot                       t_core.StatusCode = 418 // RFC 7168, 2.3.3
	StatusCode_MisdirectedRequest           t_core.StatusCode = 421 // RFC 7540, 9.1.2
	StatusCode_UnprocessableEntity          t_core.StatusCode = 422 // RFC 4918, 11.2
	StatusCode_Locked                       t_core.StatusCode = 423 // RFC 4918, 11.3
	StatusCode_FailedDependency             t_core.StatusCode = 424 // RFC 4918, 11.4
	StatusCode_TooEarly                     t_core.StatusCode = 425 // RFC 8470, 5.2.
	StatusCode_UpgradeRequired              t_core.StatusCode = 426 // RFC 7231, 6.5.15
	StatusCode_PreconditionRequired         t_core.StatusCode = 428 // RFC 6585, 3
	StatusCode_TooManyRequests              t_core.StatusCode = 429 // RFC 6585, 4
	StatusCode_RequestHeaderFieldsTooLarge  t_core.StatusCode = 431 // RFC 6585, 5
	StatusCode_UnavailableForLegalReasons   t_core.StatusCode = 451 // RFC 7725, 3

	StatusCode_InternalServerError           t_core.StatusCode = 500 // RFC 7231, 6.6.1
	StatusCode_NotImplemented                t_core.StatusCode = 501 // RFC 7231, 6.6.2
	StatusCode_BadGateway                    t_core.StatusCode = 502 // RFC 7231, 6.6.3
	StatusCode_ServiceUnavailable            t_core.StatusCode = 503 // RFC 7231, 6.6.4
	StatusCode_GatewayTimeout                t_core.StatusCode = 504 // RFC 7231, 6.6.5
	StatusCode_HTTPVersionNotSupported       t_core.StatusCode = 505 // RFC 7231, 6.6.6
	StatusCode_VariantAlsoNegotiates         t_core.StatusCode = 506 // RFC 2295, 8.1
	StatusCode_InsufficientStorage           t_core.StatusCode = 507 // RFC 4918, 11.5
	StatusCode_LoopDetected                  t_core.StatusCode = 508 // RFC 5842, 7.2
	StatusCode_NotExtended                   t_core.StatusCode = 510 // RFC 2774, 7
	StatusCode_NetworkAuthenticationRequired t_core.StatusCode = 511 // RFC 6585, 6
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

type ApiSessionType struct {
	statusCode t_core.StatusCode

	api            i_core.IApi
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

func NewApiSession() *ApiSessionType {
	return &ApiSessionType{
		data:   make(map[string]interface{}, 5),
		defers: make([]func(), 0, 20),
	}
}

func (apiSession *ApiSessionType) SetPathVars(vars map[string]string) {
	apiSession.pathVars = vars
}

func (apiSession *ApiSessionType) PathVars() map[string]string {
	return apiSession.pathVars
}

func (apiSession *ApiSessionType) SetJwtHeaderName(headerName string) {
	apiSession.jwtHeaderName = headerName
}

func (apiSession *ApiSessionType) JwtHeaderName() string {
	return apiSession.jwtHeaderName
}

func (apiSession *ApiSessionType) SetLang(lang string) {
	apiSession.lang = lang
}

func (apiSession *ApiSessionType) Lang() string {
	return apiSession.lang
}

func (apiSession *ApiSessionType) SetClientType(clientType string) {
	apiSession.clientType = clientType
}

func (apiSession *ApiSessionType) ClientType() string {
	return apiSession.clientType
}

func (apiSession *ApiSessionType) SetJwtBody(jwtBody map[string]interface{}) {
	apiSession.jwtBody = jwtBody
}

func (apiSession *ApiSessionType) ResponseWriter() http.ResponseWriter {
	return apiSession.responseWriter
}

func (apiSession *ApiSessionType) Request() *http.Request {
	return apiSession.request
}

func (apiSession *ApiSessionType) SetResponseWriter(w http.ResponseWriter) {
	apiSession.responseWriter = w
}

func (apiSession *ApiSessionType) SetRequest(r *http.Request) {
	apiSession.request = r
}

func (apiSession *ApiSessionType) JwtBody() map[string]interface{} {
	return apiSession.jwtBody
}

func (apiSession *ApiSessionType) Params() map[string]interface{} {
	return apiSession.params
}

func (apiSession *ApiSessionType) OriginalParams() map[string]interface{} {
	return apiSession.originalParams
}

func (apiSession *ApiSessionType) SetParams(params map[string]interface{}) {
	apiSession.params = params
}

func (apiSession *ApiSessionType) SetOriginalParams(originalParams map[string]interface{}) {
	apiSession.originalParams = originalParams
}

func (apiSession *ApiSessionType) Api() i_core.IApi {
	return apiSession.api
}

func (apiSession *ApiSessionType) SetApi(api i_core.IApi) {
	apiSession.api = api
}

func (apiSession *ApiSessionType) SetUserId(userId uint64) {
	apiSession.userId = userId
}

func (apiSession *ApiSessionType) UserId() uint64 {
	return apiSession.userId
}

func (apiSession *ApiSessionType) ScanParams(dest interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: false,
		TagName:          "json",
		Result:           &dest,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(apiSession.params)
	if err != nil {
		return err
	}
	return nil
}

func (apiSession *ApiSessionType) MustScanParams(dest interface{}) {
	err := apiSession.ScanParams(dest)
	if err != nil {
		panic(err)
	}
}

// Add defer handler.
// Defer handlers will be executed by order at the end of api session.
func (apiSession *ApiSessionType) AddDefer(defer_ func()) {
	apiSession.defers = append(apiSession.Defers(), defer_)
}

func (apiSession *ApiSessionType) Defers() []func() {
	return apiSession.defers
}

func (apiSession *ApiSessionType) SetData(key string, data interface{}) {
	apiSession.data[key] = data
}

func (apiSession *ApiSessionType) Data(key string) interface{} {
	return apiSession.data[key]
}

func (apiSession *ApiSessionType) Redirect(url string) {
	http.Redirect(apiSession.responseWriter, apiSession.request, url, http.StatusTemporaryRedirect)
}

// Response json body.
func (apiSession *ApiSessionType) WriteJson(data interface{}) error {
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
func (apiSession *ApiSessionType) SetHeader(key string, value string) {
	apiSession.responseWriter.Header().Set(key, value)
}

// Response text body.
func (apiSession *ApiSessionType) WriteText(text string) error {
	apiSession.responseWriter.Header().Set(string(HeaderName_ContentType), string(ContentTypeValue_Text))
	apiSession.responseWriter.WriteHeader(int(apiSession.statusCode))
	_, err := apiSession.responseWriter.Write([]byte(text))
	if err != nil {
		return err
	}
	return nil
}

// Set status code of response.
func (apiSession *ApiSessionType) SetStatusCode(code t_core.StatusCode) {
	apiSession.statusCode = code
}

// Get request host.
func (apiSession *ApiSessionType) Host() string {
	return apiSession.request.Host
}

// Get request path.
func (apiSession *ApiSessionType) Path() string {
	return apiSession.request.URL.Path
}

// Get request body.
func (apiSession *ApiSessionType) Body() io.ReadCloser {
	return apiSession.request.Body
}

// Get request method (GET, POST, PUT, etc.).
func (apiSession *ApiSessionType) Method() string {
	return apiSession.request.Method
}

// Read header by key from request headers.
func (apiSession *ApiSessionType) Header(name string) string {
	return apiSession.request.Header.Get(name)
}

// Read remote address from request headers.
func (apiSession *ApiSessionType) RemoteAddress() string {
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
func (apiSession *ApiSessionType) UrlParams() map[string]string {
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
func (apiSession *ApiSessionType) FormValues() (map[string][]string, error) {
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
func (apiSession *ApiSessionType) ReadJSON(jsonObject interface{}) error {
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

func (apiSession *ApiSessionType) Logger() i_logger.ILogger {
	return logger.LoggerDriverInstance.Logger
}

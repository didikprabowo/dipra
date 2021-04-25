package dipra

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type (
	Header         string
	MimeType       string
	AccessControll string

	Context struct {
		http.ResponseWriter
		*http.Request
		params *params
		Binding
		engine *Engine
		Writen ResponseWriter
	}
)

const (
	HeaderAccept              Header = "Accept"
	HeaderAcceptEncoding      Header = "Accept-Encoding"
	HeaderAllow               Header = "Allow"
	HeaderAuthorization       Header = "Authorization"
	HeaderContentDisposition  Header = "Content-Disposition"
	HeaderContentEncoding     Header = "Content-Encoding"
	HeaderContentLength       Header = "Content-Length"
	HeaderContentType         Header = "Content-Type"
	HeaderCookie              Header = "Cookie"
	HeaderSetCookie           Header = "Set-Cookie"
	HeaderIfModifiedSince     Header = "If-Modified-Since"
	HeaderLastModified        Header = "Last-Modified"
	HeaderLocation            Header = "Location"
	HeaderUpgrade             Header = "Upgrade"
	HeaderVary                Header = "Vary"
	HeaderWWWAuthenticate     Header = "WWW-Authenticate"
	HeaderXForwardedFor       Header = "X-Forwarded-For"
	HeaderXForwardedProto     Header = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  Header = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       Header = "X-Forwarded-Ssl"
	HeaderXUrlScheme          Header = "X-Url-Scheme"
	HeaderXHTTPMethodOverride Header = "X-HTTP-Method-Override"
	HeaderXRealIP             Header = "X-Real-IP"
	HeaderXRequestID          Header = "X-Request-ID"
	HeaderXRequestedWith      Header = "X-Requested-With"
	HeaderServer              Header = "Server"
	HeaderOrigin              Header = "Origin"
	HeaderStatus              Header = "Status"

	MIMEApplicationJSON       MimeType = "application/json"
	MIMEApplicationJavaScript MimeType = "application/javascript"
	MIMEApplicationXML        MimeType = "application/xml"
	MIMETextXML               MimeType = "text/xml"
	MIMEApplicationForm       MimeType = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf   MimeType = "application/protobuf"
	MIMEApplicationMsgpack    MimeType = "application/msgpack"
	MIMETextHTML              MimeType = "text/html"
	MIMEApplicationYAML       MimeType = "application/x-yaml"
	MIMETextYAML              MimeType = "text/yaml"
	MIMETextPlain             MimeType = "text/plain"
	MIMEMultipartForm         MimeType = "multipart/form-data"
	MIMEOctetStream           MimeType = "application/octet-stream"

	AccessControllOrigin        AccessControll = "Access-Control-Allow-Origin"
	ACcessControllMethod        AccessControll = "Access-Control-Allow-Methods"
	ACcessControllHeaders       AccessControll = "Access-Control-Allow-Headers"
	AccessControllCredential    AccessControll = "Access-Control-Allow-Credentials"
	AccessControllMaxAge        AccessControll = "Access-Control-Max-Age"
	AccessControllExposeHeaders AccessControll = "Access-Control-Expose-Headers"
	AccessControllReqMethod     AccessControll = "Access-Control-Request-Method"
	AccessControllReqHeaders    AccessControll = "Access-Control-Request-Headers"
)

// Reset Context and response
func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.Writen.Reset(w)
	c.ResponseWriter = w
	c.Request = r
	c.SetBind(r)
}

// Query by Key ?=key=value
func (c *Context) Query(param string) string {
	return c.getQuery(param, "")
}

// // Param by wlidcard /:id
func (c *Context) Param(param string) string {
	return c.params.getParam(param)
}

// GetQuery By param
func (c *Context) getQuery(param string, DefaultQuery string) string {
	q := c.URL.Query()
	paramValue := q.Get(param)
	if len(paramValue) == 0 {
		if len(DefaultQuery) == 0 {
			return ""
		}

		return DefaultQuery
	}

	return paramValue
}

// JSON response
func (c *Context) JSON(code int, body interface{}) error {
	out, err := json.Marshal(body)
	if err != nil {
		defaulterrorHTTP(c.ResponseWriter, http.StatusInternalServerError, err)
	}
	p := map[string]string{
		string(HeaderContentType): string(MIMEApplicationJSON),
	}
	c.Writen.WriteHeader(p)
	c.Writen.WriteStatus(code)
	c.Writen.Write(out)
	return nil
}

// JSONP Response
func (c *Context) JSONP(code int, body interface{}) error {
	out, err := json.Marshal(body)
	if err != nil {
		defaulterrorHTTP(c.ResponseWriter, http.StatusInternalServerError, err)
	}
	p := map[string]string{
		string(HeaderContentType): string(MIMEApplicationJavaScript),
	}
	c.Writen.WriteHeader(p)
	c.Writen.WriteStatus(code)
	c.Writen.Write(out)
	return nil
}

// String response
func (c *Context) String(code int, body string) (err error) {
	if reflect.ValueOf(body).Kind() != reflect.String {
		return http.ErrNotSupported
	}

	p := map[string]string{
		string(HeaderContentType): string(MIMETextPlain),
	}
	c.Writen.WriteHeader(p)

	c.Writen.WriteStatus(code)
	c.Writen.Write([]byte(body))
	return err
}

// XML response
func (c *Context) XML(code int, body interface{}) (err error) {

	xml, err := xml.MarshalIndent(body, "", "")
	if err != nil {
		return err
	}
	p := map[string]string{
		string(HeaderContentType): string(MIMEApplicationXML),
	}
	c.Writen.WriteHeader(p)
	c.Writen.WriteStatus(code)
	c.Writen.Write(xml)
	return err
}

// SetCookie by input
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.GetResponse(), cookie)
}

// GetCookies all
func (c *Context) GetCookies() []*http.Cookie {
	return c.Cookies()
}

// GetCookie is get cookie by name
func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	return c.Cookie(name)
}

// Redirect http url
func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.GetResponse(), c.GetRequest(), url, code)
}

// File is used get file type
func (c *Context) File(path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	serveFile := func(path string) {
		http.ServeFile(c.GetResponse(), c.GetRequest(), path)
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + path
		if _, err = os.Open(index); err != nil {
			return
		}
	}

	serveFile(path)
	return
}

// GetResponse by Reset() or another set http with returns http.ResponseWriter
func (c *Context) GetResponse() http.ResponseWriter {
	return c.ResponseWriter
}

// GetRequest returns *http.request
func (c *Context) GetRequest() *http.Request {
	return c.Request
}

// SetError for Get Error
func (c *Context) SetError(err error) {
	c.engine.HandlerError(err, c)
}

// GetHeader for get value header
func (c *Context) GetHeader(key string) interface{} {
	switch key {
	case "*":
		return c.GetRequest().Header
	case "":
		return ""
	default:
		return c.GetRequest().Header.Get(key)
	}
}

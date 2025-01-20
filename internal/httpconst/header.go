// Package httpconst определяет константы приложения.
package httpconst

// HTTP header keys.
const (
	HeaderContentType     = "Content-Type"
	HeaderContentEncoding = "Content-Encoding"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderAuthorization   = "Authorization"
	HeaderUserAgent       = "User-Agent"
	HeaderHashSHA256      = "HashSHA256"
)

// HTTP header values.
const (
	// Content-Type: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	ContentTypePlain = "text/plain"
	ContentTypeHTML  = "text/html"
	ContentTypeXML   = "application/xml"
	ContentTypeJSON  = "application/json"

	// Content-Encoding: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
	ContentEncodingGzip     = "gzip"
	ContentEncodingCompress = "compress"
	ContentEncodingDeflate  = "deflate"
	ContentEncodingBr       = "br"
	ContentEncodingZstd     = "zstd"

	// Accept-Encoding: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
	AcceptEncodingGzip     = "gzip"
	AcceptEncodingCompress = "compress"
	AcceptEncodingDeflate  = "deflate"
	AcceptEncodingBr       = "br"
	AcceptEncodingZstd     = "zstd"
	AcceptEncodingIdentity = "identity"
	AcceptEncodingAll      = "*"
)

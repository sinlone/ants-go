package http

import (
	"ants/util"
	"io"
	Http "net/http"
	"strconv"
)

/*
what a request would do
*	basic request
*	global unique key
*   cookie jar index
*	spider belong to
*	parse belong to
*/
type Request struct {
	GoRequest  *Http.Request
	CookieJar  int // if the web site store message in cookie which show diffrent result,you shoud use it
	UniqueName string
	SpiderName string
	ParserName string // crawler will send the response to named parseFunction
	NodeName   string
	Depth      int
}

func NewRequest(method, url, spiderName, parserName string, body io.Reader, cookieJar int) (*Request, error) {
	httpRequest, err := Http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request := &Request{
		GoRequest:  httpRequest,
		CookieJar:  cookieJar,
		SpiderName: spiderName,
		ParserName: parserName,
		Depth:      0,
	}
	request.makeUniqueName()
	return request, err
}

// unique sign of q request
func (this *Request) makeUniqueName() {
	baseString := this.SpiderName
	baseString += ":" + this.ParserName
	baseString += ":" + this.GoRequest.Method
	baseString += ":" + this.GoRequest.URL.String()
	baseString += ":" + this.GoRequest.Form.Encode()
	this.UniqueName = strconv.FormatUint(util.HashString(baseString), 10)
}

func (this *Request) SetNodeName(nodeName string) {
	this.NodeName = nodeName
}

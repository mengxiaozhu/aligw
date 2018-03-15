package aligw

import (
	"net/url"
	"github.com/cocotyty/httpclient"
	"crypto/md5"
	"encoding/base64"
	"strings"
	"crypto/hmac"
	"hash"
	"crypto/sha256"
	"time"
	"strconv"
)

type Request struct {
	appId   string
	secret  string
	method  string
	host    string
	path    string
	queries url.Values
}

func New(host string, path string, appId string, secret string) *Request {
	return &Request{
		method: "GET",
		host:   host,
		path:   path,
		appId:  appId,
		secret: secret,
	}
}

func (req *Request) Query(k, v string) *Request {
	req.queries.Add(k, v)
	return req
}

func (req *Request) Send() *httpclient.HttpResponse {

	accept := ""
	date := time.Now().Format(time.RFC1123)
	contentType := "application/x-www-form-urlencoded; charset=UTF-8"
	headers := url.Values{}

	headers.Add("X-Ca-Version", "1")
	headers.Add("X-Ca-Stage", "RELEASE")
	headers.Add("X-Ca-Key", req.appId)
	headers.Add("X-Ca-Request-Mode", "null")
	headers.Add("X-Ca-Timestamp", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))

	//TODO GET or POST
	builder := httpclient.Get("http://" + req.host + req.path)
	builder.Head("Accept", accept)
	builder.Head("X-Ca-Signature-Headers", "X-Ca-Request-Mode,X-Ca-Version,X-Ca-Stage,X-Ca-Key,X-Ca-Timestamp ")
	builder.Head("Date", date)
	builder.Head("Content-Type", contentType)
	for k, vs := range headers {
		builder.Head(k, vs[0])
		builder.Head("X-Ca-Signature", req.Sign(req.method, accept, contentType, date, req.path, "", headers, req.queries))
	}

	for k, vs := range req.queries {
		builder.Query(k, vs[0])
	}

	return builder.Send()

}
func (req *Request) Sign(method, accept, contentType, date string, path, content string, headers url.Values, queries url.Values) string {
	h := md5.New()
	h.Write([]byte(content))
	contentMD5 := ""
	if content != "" {
		contentMD5 = base64.StdEncoding.EncodeToString(h.Sum(nil))
	}
	headersString, _ := url.QueryUnescape(strings.Replace(strings.Replace(headers.Encode(), "=", ":", -1), "&", "\n", -1) + "\n")
	stringToSign := method + "\n" + accept + "\n" + contentMD5 + "\n" + contentType + "\n" + date + "\n" + headersString + path
	if len(queries) > 0 {
		stringToSign += "?" + queries.Encode()
	}
	hh := hmac.New(func() hash.Hash {
		return sha256.New()
	}, []byte(req.secret))
	hh.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(hh.Sum(nil))
}

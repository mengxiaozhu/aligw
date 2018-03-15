package aligw

type AliGateway struct {
	AppKey    string `sm:"(.appKey)"`
	AppSecret string `sm:"(.appSecret)"`
}

func (gw *AliGateway) Ready() {

}

func (gw *AliGateway) Get(host string, path string) *Request {
	return New(host, path, gw.AppKey, gw.AppSecret)
}

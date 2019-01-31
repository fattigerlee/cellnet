package tests

import (
	"encoding/json"
	"fmt"
	"github.com/fattigerlee/cellnet"
	"github.com/fattigerlee/cellnet/codec"
	_ "github.com/fattigerlee/cellnet/codec/httpform"
	_ "github.com/fattigerlee/cellnet/codec/httpjson"
	"github.com/fattigerlee/cellnet/peer"
	httppeer "github.com/fattigerlee/cellnet/peer/http"
	"github.com/fattigerlee/cellnet/proc"
	_ "github.com/fattigerlee/cellnet/proc/http"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const httpTestAddr = "127.0.0.1:8081"

func TestHttp(t *testing.T) {

	p := peer.NewGenericPeer("http.Acceptor", "httpserver", httpTestAddr, nil)

	proc.BindProcessorHandler(p, "http", func(raw cellnet.Event) {

		// 不依赖httpmeta
		if matcher, ok := raw.Session().(httppeer.RequestMatcher); ok {
			switch {
			case matcher.Match("POST", "/gm"):

				// 默认返回json
				raw.Session().Send(&httppeer.MessageRespond{
					Msg: &HttpEchoACK{
						Token: "ok",
					},
				})

			}
		}

		switch raw.Message().(type) {
		case *HttpEchoREQ:

			raw.Session().Send(&httppeer.MessageRespond{
				StatusCode: http.StatusOK,
				Msg: &HttpEchoACK{
					Status: 0,
					Token:  "ok",
				},
				CodecName: "httpjson",
			})

		}

	})

	p.Start()

	requestThenValid(t, "GET", "/hello_form", &HttpEchoREQ{
		UserName: "kitty_form",
	}, &HttpEchoACK{
		Token: "ok",
	})

	requestThenValid(t, "POST", "/hello_json", &HttpEchoREQ{
		UserName: "kitty_json",
	}, &HttpEchoACK{
		Token: "ok",
	})

	postCheckBody(t, "/gm", &HttpEchoACK{
		Token: "ok",
	})

	p.Stop()

	//validPage(t, "http://127.0.0.1:8081", "")
}

func requestThenValid(t *testing.T, method, path string, req, expectACK interface{}) {

	p := peer.NewGenericPeer("http.Connector", "httpclient", httpTestAddr, nil).(cellnet.HTTPConnector)

	ack, err := p.Request(method, path, req)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(ack, expectACK) {
		t.Log("unexpect token result", err)
		t.FailNow()
	}

}

func validPage(t *testing.T, url, expectAck string) {
	c := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := c.Get(url)
	if err != nil {
		t.Log("http req failed", err)
		t.FailNow()
	}

	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("http response failed", err)
		t.FailNow()
	}

	body := string(bodyData)

	if body != expectAck {
		t.Log("unexpect result", err, body)
		t.FailNow()
	}
}

func postCheckBody(t *testing.T, path string, expectMsg interface{}) {
	resp, err := http.Post("http://"+httpTestAddr+path, "application/json", nil)

	if err != nil {
		t.Log("http req failed", err)
		t.FailNow()
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("http response failed", err)
		t.FailNow()
	}

	msg := reflect.New(reflect.TypeOf(expectMsg).Elem()).Interface()

	if err := json.Unmarshal(body, msg); err != nil {
		t.Log("json unmarshal failed", err)
		t.FailNow()
	}

	if !reflect.DeepEqual(msg, expectMsg) {
		t.Log("unexpect token result", err)
		t.FailNow()
	}

}

type HttpEchoREQ struct {
	UserName string
}

type HttpEchoACK struct {
	Token  string
	Status int32
}

func (self *HttpEchoREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *HttpEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
		Path:         "/hello_form",
		Method:       "GET",
		RequestCodec: codec.MustGetCodec("httpform"),
		RequestType:  reflect.TypeOf((*HttpEchoREQ)(nil)).Elem(),

		// 请求方约束
		ResponseCodec: codec.MustGetCodec("httpjson"),
		ResponseType:  reflect.TypeOf((*HttpEchoACK)(nil)).Elem(),
	})

	cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
		Path:         "/hello_json",
		Method:       "POST",
		RequestCodec: codec.MustGetCodec("httpjson"),
		RequestType:  reflect.TypeOf((*HttpEchoREQ)(nil)).Elem(),

		ResponseCodec: codec.MustGetCodec("httpjson"),
		ResponseType:  reflect.TypeOf((*HttpEchoACK)(nil)).Elem(),
	})

}

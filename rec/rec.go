package rec

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/rpc"
	"github.com/piupuer/go-helper/pkg/utils"
	pp "github.com/piupuer/pd-pipeline/api/build/pd-pipeline"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Rec struct {
	ops    Options
	uri    string
	path   string
	client *rpc.Grpc
	Error  error
}

type result struct {
	Key    []string `json:"key"`
	Value  []string `json:"value"`
	ErrNo  int      `json:"err_no"`
	ErrMsg string   `json:"err_msg"`
}

func New(options ...func(*Options)) (rec *Rec) {
	ops := getOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	rec = &Rec{
		ops:  *ops,
		path: fmt.Sprintf("%s/%s", ops.name, ops.method),
		uri:  fmt.Sprintf("%s/%s/%s", ops.http, ops.name, ops.method),
		client: rpc.NewGrpc(
			ops.grpc,
			rpc.WithGrpcCtx(ops.ctx),
		),
	}
	return
}

func (rec *Rec) Image(filename string) (rp string) {
	if rec.Error != nil {
		return
	}
	bs, _ := ioutil.ReadFile(filename)
	b64 := base64.StdEncoding.EncodeToString(bs)
	if len(b64) == 0 {
		rec.AddError(fmt.Errorf("invalid file"))
		return
	}
	if rec.ops.grpc != "" {
		rp = rec.grpc(b64)
	} else {
		rp = rec.http(b64)
	}
	return
}

func (rec *Rec) http(b64 string) (rp string) {
	body := map[string][]string{
		"key":   {"image"},
		"value": {b64},
	}
	r, _ := http.NewRequest(http.MethodPost, rec.uri, bytes.NewReader([]byte(utils.Struct2Json(body))))
	client := &http.Client{
		Timeout: time.Duration(rec.ops.timeout) * time.Second,
	}

	start := carbon.Now().TimestampMilli()
	res, err := client.Do(r)
	if err != nil {
		log.WithContext(rec.ops.ctx).Error("network error, %v", err)
		rec.AddError(errors.Errorf("network error"))
		return
	}
	bs, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	end := carbon.Now().TimestampMilli()
	var v result
	utils.Json2Struct(string(bs), &v)
	if v.ErrNo != 0 || len(v.Value) != 1 {
		if v.ErrNo != 0 {
			log.WithContext(rec.ops.ctx).Error("rec failed, %v", v.ErrMsg)
		} else {
			log.WithContext(rec.ops.ctx).Error("rec failed, invalid uri: %s", rec.uri)
		}
		rec.AddError(errors.Errorf("rec failed"))
		return
	}
	log.WithContext(rec.ops.ctx).Info("rec success, latency: %dms", end-start)
	rp = v.Value[0]
	return
}

func (rec *Rec) grpc(b64 string) (rp string) {
	client := pp.NewPipelineServiceClient(rec.client.Conn)
	start := carbon.Now().TimestampMilli()
	v, err := client.Inference(rec.ops.ctx, &pp.Request{
		Key:    []string{"image"},
		Value:  []string{b64},
		Name:   rec.ops.name,
		Method: rec.ops.method,
	})
	if err != nil {
		log.WithContext(rec.ops.ctx).Error("network error, %v", err)
		rec.AddError(errors.Errorf("network error"))
		return
	}
	end := carbon.Now().TimestampMilli()
	if v.ErrNo != 0 || len(v.Value) != 1 {
		if v.ErrNo != 0 {
			log.WithContext(rec.ops.ctx).Error("rec failed, %v", v.ErrMsg)
		} else {
			log.WithContext(rec.ops.ctx).Error("rec failed, invalid uri: %s", rec.path)
		}
		rec.AddError(errors.Errorf("rec failed"))
		return
	}
	log.WithContext(rec.ops.ctx).Info("rec success, latency: %dms", end-start)
	rp = v.Value[0]
	return
}

func (rec *Rec) AddError(err error) error {
	if rec.Error == nil {
		rec.Error = err
	} else if err != nil && !errors.Is(err, rec.Error) {
		rec.Error = fmt.Errorf("%v; %w", rec.Error, err)
	}
	return rec.Error
}

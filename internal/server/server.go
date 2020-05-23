package server

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/golang/snappy"

	produceTools "gitlab.tubecorporate.com/push/kafka-producer/internal/produce-tools"

	"github.com/valyala/fasthttp"
)

var producer *produceTools.Producer

func mainRouter(r *fasthttprouter.Router) {
	r.GET("/config", pathConfig)
	// POST send message in body to kafka topic
	r.POST("/:topic", pathRoot)
}

func pathRoot(ctx *fasthttp.RequestCtx) {
	if bytes.Equal(ctx.Method(), []byte(fasthttp.MethodPost)) &&
		len(ctx.PostBody()) > 1 {
		var v interface{}
		if err := json.Unmarshal(ctx.PostBody(), &v); err != nil {
			fmt.Println("ERROR", err, string(ctx.PostBody()))
			producer.Stat.IncFail(1)
			ctx.SetStatusCode(502)
			return
		}
		encoded := snappy.Encode(nil, ctx.PostBody())
		go producer.Push(encoded, ctx.UserValue("topic").(string))
		ctx.SetStatusCode(fasthttp.StatusCreated)
		return
	}
	producer.Stat.IncFail(1)
	ctx.NotModified()
}

func pathConfig(ctx *fasthttp.RequestCtx) {
	cfg := producer.GetConfig()
	b, _ := json.Marshal(cfg)
	_, _ = ctx.Write(b)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	return
}

func LaunchFastHTTPServer(port string, kafkaProducer *produceTools.Producer) error {
	producer = kafkaProducer
	producer.RunTimer()
	router := fasthttprouter.New()
	mainRouter(router)
	return fasthttp.ListenAndServe(port, router.Handler)
}

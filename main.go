package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

// 定义需要监控 Gauge 类型对象
var (
	queueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "order_service_order_queue_size",
		Help: "The size of order queue",
	}, []string{"type"})
)

type OrderQueue struct {
	queue chan string
}

func newOrderQueue() *OrderQueue {
	return &OrderQueue{
		queue: make(chan string, 100),
	}
}

// 产生订单消息
func (q *OrderQueue) produceOrder() {
	// 产生订单消息
	// 队列个数加1
	queueSize.WithLabelValues("make_order").Inc() // 下单队列
	// queueSize.WithLabelValues("cancel_order").Inc() // 取消订单队列
}

// 消费订单消息
func (q *OrderQueue) consumeOrder() {
	// 消费订单消息
	// 队列个数减1
	queueSize.WithLabelValues("make_order").Dec()
}

func ProduceOrder() {
	go func() {
		q := newOrderQueue()
		for {
			q.produceOrder()
			time.Sleep(2 * time.Second)
		}
	}()
}

func ConsumeOrder() {
	go func() {
		q := newOrderQueue()
		for {
			q.produceOrder()
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {
	ProduceOrder()
	ConsumeOrder()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * Run Prometheus Metrics and add more metrics for Service Observability.
 *
 */

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	SERVICE_PREFIX = "marina_"
)

type Metrics struct {
	counterMap map[string]prometheus.Counter
	gaugeMap   map[string]prometheus.Gauge
	addLock    sync.Mutex
}

func Init(prometheusPort int) *Metrics {
	if prometheusPort <= 0 {
		glog.Errorf("PrometheusPort is not specified, so listener is not started!!")
	} else {
		// start prometheus endpoint, used for some metrics
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(":"+strconv.Itoa(prometheusPort), nil)
			if err != nil {
				glog.Fatalf("Start prometheus endpoint encountered error = %v", err)
			}
		}()
	}

	return &Metrics{
		counterMap: map[string]prometheus.Counter{},
		gaugeMap:   map[string]prometheus.Gauge{},
	}
}

// If counter with same name was already registered, function will throw an error
func (m *Metrics) RegisterNewCounters(names []string) error {
	m.addLock.Lock()
	defer m.addLock.Unlock()
	for _, name := range names {
		_, found := m.counterMap[name]
		if found {
			return errors.New(fmt.Sprintf("Counter with name `%s` was already registered", name))
		}
		// will panic if duplicate metric is added
		m.counterMap[name] = promauto.NewCounter(prometheus.CounterOpts{Name: SERVICE_PREFIX + "counter_" + name,
			ConstLabels: prometheus.Labels{
				"enable_pulse4k_collection": "true",
			}})
	}
	return nil
}

// If gauge with same name was already registered, function will throw an error
func (m *Metrics) RegisterNewGauges(names []string) error {
	m.addLock.Lock()
	defer m.addLock.Unlock()
	for _, name := range names {
		_, found := m.gaugeMap[name]
		if found {
			return errors.New(fmt.Sprintf("Gauge with name `%s` was already registered", name))
		}
		// will panic if duplicate metric is added
		m.gaugeMap[name] = promauto.NewGauge(prometheus.GaugeOpts{Name: SERVICE_PREFIX + "gauge_" + name,
			ConstLabels: prometheus.Labels{
				"enable_pulse4k_collection": "true",
			}})
	}
	return nil
}

// Requires metric name to be registered in advance using RegisterNewCounters
// return TRUE if metric was recorded, otherwise FALSE
func (m *Metrics) Mark(name string, value int64) bool {
	_, found := m.counterMap[name]
	if found {
		m.counterMap[name].Add(float64(value))
		return true
	}
	return false
}

// Requires metric name to be registered in advance using RegisterNewGauges
// return TRUE if metric was recorded, otherwise FALSE
func (m *Metrics) Gauge(name string, value int64) bool {
	_, found := m.gaugeMap[name]
	if found {
		m.gaugeMap[name].Set(float64(value))
		return true
	}
	return false
}

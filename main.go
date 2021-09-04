package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/gogap/logrus"
	"github.com/julienschmidt/httprouter"
	v1 "k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

func main() {
	router := httprouter.New()
	router.POST("/filter", filter)
	log.Fatal(http.ListenAndServe(":8888", router))
}

func filter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Info("start filter")
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderArgs schedulerapi.ExtenderArgs
	var result schedulerapi.ExtenderFilterResult
	defer func() {
		if err := recover(); err != nil {
			result.Error = fmt.Sprintf("%v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		if response, err := json.Marshal(&result); err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			log.Info("response: ", string(response))
			w.Write(response)
		}
	}()
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		result = schedulerapi.ExtenderFilterResult{
			Error: err.Error(),
		}
		return
	}

	result = filterNodes(&extenderArgs, fitNodes)
	log.WithFields(log.Fields{
		"result": result}).Debug("Filter done")

}

func filterNodes(args *schedulerapi.ExtenderArgs, f func(*v1.Pod, []string) ([]string, error)) schedulerapi.ExtenderFilterResult {
	var nodes []string
	if args.NodeNames != nil && len(*args.NodeNames) > 0 {
		nodes = *args.NodeNames
	} else if args.Nodes != nil && len(args.Nodes.Items) > 0 {
		for _, node := range args.Nodes.Items {
			log.Info(node.Name)
			nodes = append(nodes, node.GetName())
		}
	}

	var result schedulerapi.ExtenderFilterResult
	fitnodes, err := f(args.Pod, nodes)
	if err != nil {
		result.Error = err.Error()
	}
	// 1. 跟自定义scheduler的策略有关，如果 nodeCacheCapable: true,则响应可省略 result.Nodes。
	// 2. 如果响应需要result.Nodes，设置nodeCacheCapable: false。
	// 第二种参考：https://www.qikqiak.com/post/custom-kube-scheduler/
	result.NodeNames = &fitnodes
	return result
}

func fitNodes(pod *v1.Pod, nodes []string) (out []string, err error) {
	if pod.Labels["want"] == "you" {
		for _, node := range nodes {
			if node == "10-9-101-66" {
				out = append(out, node)
			}
		}
	}
	return out, nil
}

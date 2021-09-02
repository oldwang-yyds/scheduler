package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/gogap/logrus"
	"github.com/julienschmidt/httprouter"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

func main() {
	router := httprouter.New()
	//router.GET("/", Index)
	router.POST("/scheduler/filter", filter)
	//	router.POST("/prioritize", Prioritize)
	log.Fatal(http.ListenAndServe(":8888", router))
}

func filter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Info("=================================")
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
			w.Write(response)
		}
	}()
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		result = schedulerapi.ExtenderFilterResult{
			Error: err.Error(),
		}
		return
	}

	result = filterNodes(&extenderArgs)
	log.WithFields(log.Fields{
		"result": result}).Debug("Filter done")

}

func filterNodes(args *schedulerapi.ExtenderArgs) schedulerapi.ExtenderFilterResult {
	log.Info("=================================")
	log.Info(args)
	var nodes []string
	if args.NodeNames != nil && len(*args.NodeNames) > 0 {
		nodes = *args.NodeNames
	} else if args.Nodes != nil && len(args.Nodes.Items) > 0 {
		for _, node := range args.Nodes.Items {
			nodes = append(nodes, node.GetName())
		}
	}

	var result schedulerapi.ExtenderFilterResult

	return result
}

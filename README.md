# scheduler

package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)
func main() {
	   // 初始化引擎
	   engine := gin.Default()
	   // 注册一个路由和处理函数
	   engine.Any("/filter", filter)
	   // 绑定端口，然后启动应用
	   engine.Run(":8888")

}
func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(schedulerapi.FailedNodesMap)
	pod := args.Pod

	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

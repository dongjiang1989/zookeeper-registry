package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-basic/uuid"
	"github.com/go-zk/registry"
)

var (
	zkRegistry *registry.ZkRegistry
	err        error
	port       = ":8888"
)

func main() {
	zkRegistry, err = registry.NewZkRegistry(
		registry.Hosts([]string{os.Getenv("ZOOKEEPER_ADDR")}),
		registry.Prefix(os.Getenv("ROOT_PATH")),
		registry.Timeout(15),
	)
	if err != nil {
		panic(err.Error())
	}
	defer zkRegistry.Close()

	node := &registry.Node{
		Id:      uuid.New(),
		Address: os.Getenv("SERVICE_NODE"),
		Port:    rand.Intn(8080),
	}
	zkRegistry.Register(os.Getenv("SERVICE_NAME"), node)

	//获取一个服务下的所有节点
	http.HandleFunc("/nodes", nodes)
	log.Println("start running : " + os.Getenv("SERVICE_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVICE_PORT"), nil))
}

// 随机获取一个服务节点
func getServiceNode(writer http.ResponseWriter, request *http.Request) {
	get := request.URL.Query()
	svc := get.Get("service")
	if svc == "" {
		writer.Write([]byte("请输入服务名称"))
		return
	}
	node, err := zkRegistry.GetServerNode(svc)
	if err != nil {
		writer.Write([]byte(err.Error()))
	} else {
		writer.Write([]byte(fmt.Sprintf("%+v", node)))
	}
}

// 获取一个服务下的所有节点
func nodes(writer http.ResponseWriter, request *http.Request) {
	get := request.URL.Query()
	svc := get.Get("service")
	if svc == "" {
		svc = os.Getenv("SERVICE_NAME")
	}
	all := zkRegistry.GetAllNode(svc)
	for _, item := range all {
		writer.Write([]byte(fmt.Sprintf("%+v", item) + "\n"))
	}
}

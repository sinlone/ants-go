package main

import (
	"ants/action"
	AHttp "ants/action/http"
	"ants/action/rpc"
	"ants/action/watcher"
	"ants/crawler"
	"ants/http"
	"ants/node"
	"ants/util"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	CONF_FILE = "/../conf/conf.json"
)

// init command line param
func initFlag(settings *util.Settings) {
	flag.IntVar(&settings.TcpPort, "tcp", settings.TcpPort, "tcp port")
	flag.IntVar(&settings.HttpPort, "http", settings.HttpPort, "http port")
}

// load setting from json file and input
func MakeSettings() *util.Settings {
	pwd, _ := os.Getwd()
	settings := util.LoadSettingFromFile(pwd + CONF_FILE)
	initFlag(settings)
	flag.Parse()
	return settings
}

// try to join the cluster,
// if there is no cluster,make itself master
func initCluster(settings *util.Settings, rpcClient action.RpcClientAnts, node *node.Node) {
	node.Join()
	isClusterExist := false
	if len(settings.NodeList) > 0 {
		for _, nodeInfo := range settings.NodeList {
			nodeSettings := strings.Split(nodeInfo, ":")
			ip := nodeSettings[0]
			port, _ := strconv.Atoi(nodeSettings[1])
			if ip == node.NodeInfo.Ip && port == node.NodeInfo.Port {
				continue
			}
			err := rpcClient.LetMeIn(ip, port)
			if err == nil {
				isClusterExist = true
			}
		}
	}
	if !isClusterExist {
		node.MakeMasterNode(node.NodeInfo.Name)
	}
	node.Ready()
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("let us go shipping")
	var wg sync.WaitGroup
	wg.Add(1)
	setting := MakeSettings()
	resultQuene := crawler.NewResultQuene()
	Node := node.NewNode(setting, resultQuene)
	var rpcClient action.RpcClientAnts = rpc.NewRpcClient(Node)
	var reporter, distributer action.Watcher = watcher.NewReporter(Node, rpcClient, resultQuene), watcher.NewDistributer(Node, rpcClient)
	var rpcServer action.RpcServerAnts = rpc.NewRpcServer(Node, setting.TcpPort, rpcClient, reporter, distributer)
	router := AHttp.NewRouter(Node, reporter, distributer)
	httpServer := http.NewHttpServer(setting, router)
	rpcServer.Start()
	httpServer.Start(wg)
	initCluster(setting, rpcClient, Node)
	wg.Wait()
}

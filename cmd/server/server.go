package main

import "fmt"
import "flag"
import "os"
import "github.com/monfron/mapago/ctrl/serverProtos"
import "github.com/monfron/mapago/ctrl/shared"


var CTRL_PORT = 64321
var DEF_BUFFER_SIZE = 8096 * 8

func main() {
	portPtr := flag.Int("port", CTRL_PORT, "port for interacting with control channel")
	callSizePtr := flag.Int("call-size", DEF_BUFFER_SIZE, "application buffer in bytes")


	flag.Parse()

	fmt.Println("mapago(c) - 2018")
	fmt.Println("Server side")
	fmt.Println("Port:", *portPtr)
	fmt.Println("Call-Size:", *callSizePtr)

	runServer(*portPtr, *callSizePtr)
}

func runServer(port int, callSize int) {
	var repDataObj *shared.DataObj
	ch := make(chan shared.ChResult)
	fmt.Println(ch)

	tcpObj := serverProtos.NewTcpObj("TcpConn1", port, callSize)
	tcpObj.Start(ch)

	/* WIP: disabled for reduced complexity
	udpObj := serverProtos.NewUdpObj("UdpConn1")
	udpObj.Start(ch)
	*/

	for {
		request := <- ch
		reqDataObj := shared.ConvJsonToDataStruct(request.Json)

		switch reqDataObj.Type {
		case shared.INFO_REQUEST:
			repDataObj = constructInfoReply(reqDataObj)
		case shared.MEASUREMENT_START_REQUESTS:
			fmt.Println("Construct MEASUREMENT_START_REP")

		case shared.MEASUREMENT_STOP_REQUEST:
			fmt.Println("Construct MEASUREMENT_STOP_REP")

		case shared.MEASUREMENT_INFO_REQUEST:
			fmt.Println("Construct MEASUREMENT_INFO_REP")

		case shared.TIME_DIFF_REQUEST:
			fmt.Println("Construct TIME_DIFF_REP")

		case shared.WARNING_ERR_MSG:
			fmt.Println("WARNING_ERR_MSG")

		default:
			fmt.Printf("Unknown type")
			os.Exit(1)
		}

		json := shared.ConvDataStructToJson(repDataObj)
		request.ConnObj.WriteAnswer(json)
	}
}

func constructInfoReply(reqDataObj *shared.DataObj) *shared.DataObj {
	fmt.Println("Constructing INFO_REP")

	// construct ID
	repDataObj := new(shared.DataObj)
	repDataObj.Type = shared.INFO_REPLY
	repDataObj.Id = shared.ConstructId()
	repDataObj.Seq = "0"
	repDataObj.Seq_rp = reqDataObj.Seq
	// maparo STD: timestamp replied untouched by server
	reqDataObj.Ts = reqDataObj.Ts
	// repDataObj.modules
	// repDataObj.Arch
	// repDataObj.Os
	repDataObj.Info = "fancyInfo"
	return repDataObj
}

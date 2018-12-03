package managementPlane

import "fmt"
import "os"
import "math/rand"
import "strings"
import "strconv"
import "github.com/monfron/mapago/control-plane/ctrl/shared"
import "github.com/monfron/mapago/measurement-plane"

var msmtStorage map[string]chan shared.ChMgmt2Msmt
var mapInited = false

func HandleMsmtStartReq(ctrlCh chan<- shared.ChMsmt2Ctrl, msmtStartReq *shared.DataObj, cltAddr string) {
	switch msmtStartReq.Measurement.Name {
	case "tcp-throughput":
		msmtId := constructMsmtId(cltAddr)
		msmtCh := make(chan shared.ChMgmt2Msmt)

		if mapInited == false {
			msmtStorage = make(map[string]chan shared.ChMgmt2Msmt)
			mapInited = true
		}

		msmtStorage[msmtId] = msmtCh
		fmt.Println("\nmsmtStorage content: ", msmtStorage)

		/*
			POSSIBLE BLOCKING CAUSE
			we have to call it via goroutine asynchronously
			or we stay within the for loop and block on the channel
			and cannot receive anything else
		*/
		go measurementPlane.NewTcpMsmt(msmtCh, ctrlCh)

		/*
			POSSIBLE BLOCKING CAUSE
			send blocks until corresponding read is called
			PROBLEM: this function is blocked => the callee is blocked aswell
			=> connClose() and HandleConn() cannot be called => no further requests
		*/

		mgmtCmd := new(shared.ChMgmt2Msmt)
		mgmtCmd.Cmd = "Msmt_start"
		mgmtCmd.MsmtId = msmtId
		msmtCh <- *mgmtCmd

	case "udp-throughput":
		fmt.Println("\nStarting UDP throughput module")

	case "quic-throughput":
		fmt.Println("\nStarting QUIC throughput module")

	case "udp-ping":
		fmt.Println("\nStarting UDP ping module")

	default:
		fmt.Printf("Unknown measurement module")
		os.Exit(1)
	}
}

func constructMsmtId(cltAddr string) string {
	// cut the port from clt_addr
	spltCltAddr := strings.Split(cltAddr, ":")
	msmtId := spltCltAddr[0] + "=" + strconv.Itoa(int(rand.Int31()))

	return msmtId
}

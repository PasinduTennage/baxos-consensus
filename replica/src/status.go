package src

import (
	"fmt"
	"time"
)

/*
	Handler for status message
		1. Invoke bootstrap / start consensus or printlog depending on the operation type
		2. Send a response back to the sender
*/

func (rp *Replica) handleStatus(message *Status) {
	fmt.Print("Status  " + fmt.Sprintf("%v", message) + " \n")
	if message.Type == 1 {
		if rp.serverStarted == false {
			rp.serverStarted = true
			rp.ConnectBootStrap()
			time.Sleep(2 * time.Second)
		}
	} else if message.Type == 2 {
		if rp.logPrinted == false {
			rp.logPrinted = true
			rp.cancel <- true
			rp.cancel <- true
			// empty the incoming channel
			go func() {
				for true {
					_ = <-rp.incomingChan
				}
			}()

			rp.printPaxosLogConsensus() // this is for consensus testing purposes

			//fmt.Printf("num go routines: %v \n", runtime.NumGoroutine())
		}
	} else if message.Type == 3 {
		if rp.consensusStarted == false {
			rp.consensusStarted = true
			rp.lastProposedTime = time.Now()
			rp.sendDummyRequests(rp.cancel)
			rp.paxosConsensus.run()
			//rp.debug("started paxos consensus with initial prepare", 0)

		}
	}

	//rp.debug("Sending status reply ", 0)

	statusMessage := Status{
		Type: message.Type,
		Note: message.Note,
	}

	rpcPair := RPCPair{
		Code: rp.messageCodes.StatusRPC,
		Obj:  &statusMessage,
	}

	rp.sendMessage(int32(message.Sender), rpcPair)
	//rp.debug("Sent status ", 0)

}

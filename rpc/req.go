package rpc

import (
	"errors"
	"github.com/davyxu/cellnet"
	"sync"
	"sync/atomic"
)

var (
	rpcIDSeq        int64
	requestByCallID sync.Map
)

type request struct {
	id     int64
	onRecv func(interface{})
}

var ErrTimeout = errors.New("time out")

func (self *request) RecvFeedback(msg interface{}) {
	self.onRecv(msg)
}

func (self *request) Send(ses cellnet.Session, msg interface{}) {
	data, meta, err := cellnet.EncodeMessage(msg)

	if err != nil {
		log.Errorf("rpc request message encode error: %s", err)
		return
	}

	ses.Send(&RemoteCallREQ{
		MsgID:  meta.ID,
		Data:   data,
		CallID: self.id,
	})
}

func createRequest(onRecv func(interface{})) *request {

	self := &request{
		onRecv: onRecv,
	}

	self.id = atomic.AddInt64(&rpcIDSeq, 1)

	requestByCallID.Store(self.id, self)

	return self
}

func getRequest(callid int64) *request {

	if v, ok := requestByCallID.Load(callid); ok {

		requestByCallID.Delete(callid)
		return v.(*request)
	}

	return nil
}

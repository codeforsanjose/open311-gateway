package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/open311-gateway/engine/common"
	"github.com/open311-gateway/engine/router"
	"github.com/open311-gateway/engine/structs"
	"github.com/open311-gateway/engine/telemetry"

	log "github.com/jeffizhungry/logrus"
)

// refreshMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to Refresh a Report.
type refreshMgr struct {
	id    int64
	start time.Time

	reqType structs.NRequestType
	nreq    *structs.NServiceRequest

	valid common.Validation

	routes structs.NRoutes
	rpc    *router.RPCCallMgr
}

func refresh(area string) (reterr error) {
	tid := "SrvRrsh"
	rqstID := common.RequestID()
	mgr := refreshMgr{
		id:    rqstID,
		start: time.Now(),
		nreq: &structs.NServiceRequest{
			NRequestCommon: structs.NRequestCommon{
				ID: structs.NID{
					RqstID: rqstID,
				},
				Rtype: structs.NRTServicesAll,
			},
			Area: "all",
		},
	}

	telemetry.SendTelemetry(mgr.id, tid, "open")
	defer func() {
		if reterr != nil {
			telemetry.SendTelemetry(mgr.id, tid, "error")
		} else {
			telemetry.SendTelemetry(mgr.id, tid, "done")
		}
	}()

	area = strings.ToUpper(area)
	switch {
	case area == "ALL":
		mgr.reqType = structs.NRTServicesAll
	case len(area) >= 2:
		mgr.reqType = structs.NRTServicesArea
		return fmt.Errorf("Services Refresh currently only supports requests for ALL areas")
	default:
		return fmt.Errorf("Invalid Area %q specified for Services Refresh", area)
	}

	routes, err := router.RoutesAll()
	log.Debugf("router.RoutesAll(): %T: %#[1]v  len: %d", routes, len(routes))
	switch {
	case err != nil:
		return err
	case len(routes) == 0:
		return fmt.Errorf("unable to get route")
	}
	mgr.routes = routes

	log.Debug("Before callRPC: " + mgr.String())
	if err := mgr.callRPC(); err != nil {
		log.Error("processRefresh.callRPC() failed - " + err.Error())
		return err
	}

	return nil
}

// -------------------------------------------------------------------------------
//                        ROUTER.REQUESTER INTERFACE
// -------------------------------------------------------------------------------
func (r *refreshMgr) RType() structs.NRequestType {
	return r.reqType
}

func (r *refreshMgr) Routes() structs.NRoutes {
	return r.routes
}

func (r *refreshMgr) Data() interface{} {
	return r.nreq
}

func (r *refreshMgr) Processer() func(ndata interface{}) error {
	return r.processReply
}

// -------------------------------------------------------------------------------
//                        RPC
// -------------------------------------------------------------------------------

// callRPC runs the calls to the Adapter(s).
func (r *refreshMgr) callRPC() (err error) {
	r.rpc, err = router.NewRPCCallMgr(r)
	if err != nil {
		return err
	}

	log.Debug("Before RPC" + r.String())
	if err = r.rpc.Run(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (r *refreshMgr) processReply(ndata interface{}) error {
	return servicesData.merge(ndata)
}

// ------------------------------ String -------------------------------------------------

// String displays the contents of the SearchRequest custom type.
func (r refreshMgr) String() string {
	ls := new(common.LogString)
	ls.AddF("Services refreshMgr - %d\n", r.id)
	ls.AddF("Request type: %v\n", r.reqType.String())
	ls.AddS(r.routes.String())
	if r.rpc != nil {
		ls.AddS(r.rpc.String())
	} else {
		ls.AddS("*****RPC uninitialized*****\n")
	}
	if r.nreq != nil {
		ls.AddS(r.nreq.String())
	} else {
		ls.AddS("**** nreq nil ****")
	}
	if r.routes != nil {
		ls.AddS(r.routes.String())
	} else {
		ls.AddS("**** routes nil ****")
	}
	return ls.Box(120) + "\n\n"
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================

package display

import (
	"fmt"
	"time"

	"github.com/open311-gateway/monitor/logs"
	"github.com/open311-gateway/monitor/telemetry"

	ui "github.com/gizak/termui"
)

var (
	// systemList  systems
	displayList displays

	engStatuses *sortedData
	engRequests *sortedData
	engAdpCalls *sortedData

	adpStatuses *sortedData
	adpCalls01  *sortedData

	log = logs.Log
)

// ==============================================================================================================================
//                                      DISPLAYS
// ==============================================================================================================================

type displays struct {
	d    []*ui.List
	data []func() []string
	l    []ui.Bufferer
}

func (r *displays) init() error {
	r.d = make([]*ui.List, 0)
	r.l = make([]ui.Bufferer, 0)

	r.newList("Engine Status", 0, 0, 10, 80, engStatuses.display)
	r.newList("Adapter Status", 80, 0, 10, 80, adpStatuses.display)

	r.newList("Eng Requests", 0, 10, 15, 80, engRequests.display)
	r.newList("Eng Adapter Calls", 80, 10, 15, 80, engAdpCalls.display)

	r.newList("Adapter Calls", 0, 25, 15, 160, adpCalls01.display)

	debugList := func() []string {
		return telemetry.DebugListLast(18)
	}

	r.newList("DEBUG", 0, 40, 20, 240, debugList)

	for _, uiList := range r.d {
		r.l = append(r.l, uiList)
	}
	return nil
}

func (r *displays) run() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	draw := func(t int) {
		for i, d := range r.d {
			d.Items = r.data[i]()
		}
		ui.Render(r.l...)
	}

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/d", func(ui.Event) {
		telemetry.DebugClear()
	})
	ui.Handle("/sys/kbd/c", func(ui.Event) {
		clearAll()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		draw(int(t.Count))
	})
	ui.Loop()
}

func (r *displays) newList(caption string, x, y, height, width int, getData func() []string) error {
	l := ui.NewList()
	l.Items = getData()
	l.ItemFgColor = ui.ColorWhite
	l.BorderLabel = caption
	l.Height = height
	l.Width = width
	l.X = x
	l.Y = y
	r.d = append(r.d, l)
	r.data = append(r.data, getData)
	return nil

}

func clearAll() {
	engStatuses.clear()
	engRequests.clear()
	engAdpCalls.clear()

	adpStatuses.clear()
	adpCalls01.clear()
	telemetry.DebugClear()

}

// Start initializes and starts the Display processes.  It should be called AFTER the
// init() processes.
func Start() {

	// Process messages
	go func() {
		msgChan := telemetry.GetMsgChan()

		for msg := range msgChan {
			telemetry.DebugMsg("Message type [%s] - %#v\n", msg.Mtype(), msg.Data())
			switch msg.Mtype() {
			case telemetry.MsgTypeES:
				if err := engStatuses.update(msg); err != nil {
					log.Error(err.Error())
				}
			case telemetry.MsgTypeER:
				if err := engRequests.update(msg); err != nil {
					log.Error(err.Error())
				}
			case telemetry.MsgTypeERPC:
				if err := engAdpCalls.update(msg); err != nil {
					log.Error(err.Error())
				}

			case telemetry.MsgTypeAS:
				if err := adpStatuses.update(msg); err != nil {
					log.Error(err.Error())
				}
			case telemetry.MsgTypeARPC:
				if err := adpCalls01.update(msg); err != nil {
					log.Error(err.Error())
				}
			default:
				log.Errorf("Invalid message received - type: %q msg: [%v]\n", msg.Mtype(), msg.Data())
			}
		}
	}()

	// runTests()
	displayList.run()
}

func runTests() {
	if err := engStatuses.update(telemetry.NewMessageTest([]string{"ES", "Sys1", "active!", "CS1, CS2", "127.0.0.1/5081"})); err != nil {
		log.Error(err.Error())
	}
	if err := engStatuses.update(telemetry.NewMessageTest([]string{"ES", "Sys2", "active", "", "127.0.0.1/5082"})); err != nil {
		log.Error(err.Error())
	}

	if err := engRequests.update(telemetry.NewMessageTest([]string{"ER", "10001", "Create", "Active", time.Now().Format(time.RFC3339), "SJ"})); err != nil {
		log.Error(err.Error())
	}

	if err := engAdpCalls.update(telemetry.NewMessageTest([]string{"ERPC", "10001-1", "active", "CS1-SJ-1", time.Now().Format(time.RFC3339)})); err != nil {
		log.Error(err.Error())
	}
	if err := engAdpCalls.update(telemetry.NewMessageTest([]string{"ERPC", "10001-2", "active", "CS1-SC-1", time.Now().Format(time.RFC3339)})); err != nil {
		log.Error(err.Error())
	}

	go func() {
		var cnt int
		for {
			cnt++
			id := fmt.Sprintf("Sys%02d", cnt)
			if err := engStatuses.update(telemetry.NewMessageTest([]string{"ES", id, "CS1, CS2", "active", ""})); err != nil {
				log.Fatalf(err.Error())
			}
			for name, data := range engStatuses.data {
				if time.Since(data.getLastUpdate()).Seconds() > 10 {
					engStatuses.data[name].setStatus("INACTIVE")
				}
			}
			if cnt == 10 {
				if err := engStatuses.update(telemetry.NewMessageTest([]string{"ES", "Sys3", "active", "XXX", "127.0.0.1/5083"})); err != nil {
					log.Fatalf(err.Error())
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

func init() {
	engStatuses = newSortedData(telemetry.MsgTypeES, true)
	engRequests = newSortedData(telemetry.MsgTypeER, false)
	engAdpCalls = newSortedData(telemetry.MsgTypeERPC, false)

	adpStatuses = newSortedData(telemetry.MsgTypeAS, true)
	adpCalls01 = newSortedData(telemetry.MsgTypeARPC, false)

	displayList.init()
}

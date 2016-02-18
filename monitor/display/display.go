package display

import (
	"fmt"
	"time"

	"Gateway311/monitor/logs"
	"Gateway311/monitor/telemetry"

	ui "github.com/gizak/termui"
)

var (
	// systemList  systems
	displayList displays

	engStatuses *sortedData
	engRequests *sortedData
	engAdpCalls *sortedData
	adpRequests *sortedData

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

	r.newList("Engine List", 0, 0, 7, 80, engStatuses.display)
	r.newList("Eng01 Requests", 0, 7, 7, 80, engRequests.display)
	r.newList("Eng01 Adapter Calls", 0, 14, 10, 80, engAdpCalls.display)

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

// Start initializes and starts the Display processes.  It should be called AFTER the
// init() processes.
func Start() {

	// Process messages
	go func() {
		msgChan := telemetry.GetMsgChan()

		for msg := range msgChan {
			switch msg.Mtype() {
			case telemetry.MsgTypeES:
				log.Debug("Message type ES - %v\n", msg.Data())
			case telemetry.MsgTypeER:
				log.Debug("Message type ER - %v\n", msg.Data())
			case telemetry.MsgTypeERPC:
				log.Debug("Message type ARPC - %v\n", msg.Data())
				if err := engAdpCalls.update(msg); err != nil {
					log.Error(err.Error())
				}

			case telemetry.MsgTypeAS:
				log.Debug("Message type AS - %v\n", msg.Data())
			case telemetry.MsgTypeARPC:
				log.Debug("Message type ARPC - %v\n", msg.Data())
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

	displayList.init()
}

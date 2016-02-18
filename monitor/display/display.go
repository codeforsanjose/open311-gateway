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

// RunTest blah blah
func RunTest() {
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
	displayList.run()
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

/*
func RunTest1() {

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	if err := engStatuses.update([]string{"S", "Sys1", "active!", "127.0.0.1/5081"}); err != nil {
		log.Fatalf(err.Error())
	}
	if err := engStatuses.update([]string{"S", "Sys2", "active", "127.0.0.1/5082"}); err != nil {
		log.Fatalf(err.Error())
	}

	go func() {
		var cnt int
		for {
			cnt++
			if err := engStatuses.update([]string{"S", "Sys1", "active", ""}); err != nil {
				log.Fatalf(err.Error())
			}
			for name, system := range engStatuses.m {
				if time.Since(system.lastUpdate).Seconds() > 10 {
					engStatuses.Lock()
					engStatuses.m[name].status = "INACTIVE"
					engStatuses.Unlock()
				}
			}
			if cnt == 10 {
				if err := engStatuses.update([]string{"S", "Sys3", "active", "127.0.0.1/5083"}); err != nil {
					log.Fatalf(err.Error())
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()

	list := ui.NewList()
	log.Debug("list type: %T", list)
	list.Items = engStatuses.display()
	list.ItemFgColor = ui.ColorWhite
	list.BorderLabel = "System List"
	list.Height = 7
	list.Width = 80
	list.Y = 4

	draw := func(t int) {
		list.Items = engStatuses.display()
		ui.Render(list)
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
*/

/*
func RunOld() {
	log.Debug("Here we go...")

	go Display()
	a, err := net.ResolveUDPAddr("udp", monitorAddr)
	if err != nil {
		fmt.Printf("Error resolving address - %s", err.Error())
		stop()
		os.Exit(-1)
	}
	fmt.Printf("Address: %s\n", spew.Sdump(a))
	MonitorConn, _ := net.ListenUDP("udp", a)
	defer MonitorConn.Close()
	MonitorConn.SetReadBuffer(1048576)

	buf := make([]byte, 1024)

	fmt.Printf("MsgID           Source                Type                 Destination\n")
	for {
		n, _, err := MonitorConn.ReadFromUDP(buf)
		msgStr := strings.Split(string(buf[0:n]), "|")
		// fmt.Printf("%-15s %-20s  %-20s %s\n", msgStr[MIMsgID], fmt.Sprintf("[%s: %s]", msgStr[MISysID], msgStr[MIStatus]), msgStr[MIOp], msgStr[MIDest])
		msgID, err := strconv.ParseInt(msgStr[MIMsgID], 10, 64)
		if err != nil || msgID <= 0 {
			log.Error("Invalid message id: %q", msgStr[MIMsgID])
			continue
		}
		_, found := sMap[msgID]
		if !found {
			sMap[msgID] = &fullStatus{
				status: status{
					ID:        msgID,
					CreatedAt: time.Now(),
					SysID:     msgStr[MISysID],
					Op:        msgStr[MIOp],
					Dest:      msgStr[MIDest],
				},
			}
		}

		switch msgStr[MIStatus] {
		case "eng-send":
			sMap[msgID].engSent = true
		case "eng-recv":
			sMap[msgID].engRecv = true
		case "adp-recv":
			sMap[msgID].adpRecv = true
		case "adp-send":
			sMap[msgID].adpSent = true
		}

		displayItem(sMap[msgID])

	}
}
*/

/*

// ==============================================================================================================================
//                                      SYSTEMS
// ==============================================================================================================================

type systems struct {
	sorted bool

	sync.RWMutex
	m   map[string]*sysData
	ind []string
}

func (r *systems) display() []string {
	systemList.RLock()
	defer systemList.RUnlock()

	// Update the index
	r.index()
	var list []string
	for _, name := range r.ind {
		list = append(list, r.m[name].String())
	}
	return list
}

func (r *systems) index() error {
	r.ind = make([]string, 0)
	for k := range r.m {
		r.ind = append(r.ind, k)
	}
	sort.Strings(r.ind)
	r.sorted = true
	return nil
}

func (r *systems) update(data []string) (key string, err error) {
	if data[sesType] != telemetry.MsgTypeES {
		return fmt.Errorf("invalid message type: %q sent to System Update - message: %v", data[sesType], data)
	}

	name := data[sesName]
	status := data[sesStatus]
	addr := data[sesAddr]

	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[name]; !ok {
		r.m[name] = &sysData{
			name:       name,
			lastUpdate: time.Now(),
			status:     status,
			addr:       addr,
		}
		r.sorted = false
	} else {
		r.m[name].name = name
		r.m[name].lastUpdate = time.Now()
		r.m[name].status = status
		if addr > "" {
			r.m[name].addr = addr
		}
	}
	return nil
}

// -------------------------------------------- sysData -------------------------------------------------------------------------

type sysData struct {
	name       string
	lastUpdate time.Time
	status     string
	addr       string
}

func (r sysData) String() string {
	return fmt.Sprintf("%-10s  %10s %8.1f  %s", r.name, r.status, time.Since(r.lastUpdate).Seconds(), r.addr)
}

*/

/*
func Run1() {

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = "Text Box"
	p.BorderFg = ui.ColorCyan
	// p.Handle("/timer/1s", func(e ui.Event) {
	// 	cnt := e.Data.(ui.EvtTimer)
	// 	if cnt.Count%2 == 0 {
	// 		p.TextFgColor = ui.ColorRed
	// 	} else {
	// 		p.TextFgColor = ui.ColorWhite
	// 	}
	// })

	strs := []string{"[0] gizak/termui", "[1] editbox.go", "[2] iterrupt.go", "[3] keyboard.go", "[4] output.go", "[5] random_out.go", "[6] dashboard.go", "[7] nsf/termbox-go"}
	list := ui.NewList()
	list.Items = strs
	list.ItemFgColor = ui.ColorYellow
	list.BorderLabel = "List"
	list.Height = 7
	list.Width = 25
	list.Y = 4

	draw := func(t int) {
		list.Items = strs[t%9:]
		ui.Render(p, list)
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

func RunExample() {
	// smpText := [3]string{
	// 	"This is the first string.  FIRST.  FIRST.  Firsty first first.",
	// 	"This is the second string.  SSSSEEEECCCCCOOOONNNNDDD.",
	// 	"This is the third string.  THIRD.  THIRD.  Thirdy 3rd third.",
	// }

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = "Text Box"
	p.BorderFg = ui.ColorCyan
	p.Handle("/timer/1s", func(e ui.Event) {
		cnt := e.Data.(ui.EvtTimer)
		if cnt.Count%2 == 0 {
			p.TextFgColor = ui.ColorRed
		} else {
			p.TextFgColor = ui.ColorWhite
		}
	})

	strs := []string{"[0] gizak/termui", "[1] editbox.go", "[2] iterrupt.go", "[3] keyboard.go", "[4] output.go", "[5] random_out.go", "[6] dashboard.go", "[7] nsf/termbox-go"}
	list := ui.NewList()
	list.Items = strs
	list.ItemFgColor = ui.ColorYellow
	list.BorderLabel = "List"
	list.Height = 7
	list.Width = 25
	list.Y = 4

	g := ui.NewGauge()
	g.Percent = 50
	g.Width = 50
	g.Height = 3
	g.Y = 11
	g.BorderLabel = "Gauge"
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan

	spark := ui.Sparkline{}
	spark.Height = 1
	spark.Title = "srv 0:"
	spdata := []int{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6, 4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}
	spark.Data = spdata
	spark.LineColor = ui.ColorCyan
	spark.TitleColor = ui.ColorWhite

	spark1 := ui.Sparkline{}
	spark1.Height = 1
	spark1.Title = "srv 1:"
	spark1.Data = spdata
	spark1.TitleColor = ui.ColorWhite
	spark1.LineColor = ui.ColorRed

	sp := ui.NewSparklines(spark, spark1)
	sp.Width = 25
	sp.Height = 7
	sp.BorderLabel = "Sparkline"
	sp.Y = 4
	sp.X = 25

	sinps := (func() []float64 {
		n := 220
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()

	lc := ui.NewLineChart()
	lc.BorderLabel = "dot-mode Line Chart"
	lc.Data = sinps
	lc.Width = 50
	lc.Height = 11
	lc.X = 0
	lc.Y = 14
	lc.AxesColor = ui.ColorWhite
	lc.LineColor = ui.ColorRed | ui.AttrBold
	lc.Mode = "dot"

	bc := ui.NewBarChart()
	bcdata := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bclabels := []string{"S0", "S1", "S2", "S3", "S4", "S5"}
	bc.BorderLabel = "Bar Chart"
	bc.Width = 26
	bc.Height = 10
	bc.X = 51
	bc.Y = 0
	bc.DataLabels = bclabels
	bc.BarColor = ui.ColorGreen
	bc.NumColor = ui.ColorBlack

	lc1 := ui.NewLineChart()
	lc1.BorderLabel = "braille-mode Line Chart"
	lc1.Data = sinps
	lc1.Width = 26
	lc1.Height = 11
	lc1.X = 51
	lc1.Y = 14
	lc1.AxesColor = ui.ColorWhite
	lc1.LineColor = ui.ColorYellow | ui.AttrBold

	p1 := ui.NewPar("Hey!\nI am a borderless block!")
	p1.Border = false
	p1.Width = 26
	p1.Height = 2
	p1.TextFgColor = ui.ColorMagenta
	p1.X = 52
	p1.Y = 11

	draw := func(t int) {
		g.Percent = t % 101
		list.Items = strs[t%9:]
		sp.Lines[0].Data = spdata[:30+t%50]
		sp.Lines[1].Data = spdata[:35+t%50]
		lc.Data = sinps[t/2%220:]
		lc1.Data = sinps[2*t%220:]
		bc.Data = bcdata[t/2%10:]
		ui.Render(p, list, g, sp, lc, bc, lc1, p1)
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

*/

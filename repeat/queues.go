package repeat

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strconv"

	"github.com/empijei/Wapty/ui"
)

var done = make(chan struct{})

func RepeaterLoop() {
	for {
		select {
		case cmd := <-uiRepeater.DataChannel:
			switch action := parseRepeaterAction(cmd.Action); action {
			case CREATE:
				r := NewRepeater()
				status.Add(r)
			case GO:
				handleGo(&cmd)
			case GET:
				handleGet(&cmd)
			default:
				log.Println("Unknown repeater action: " + cmd.Action)
			}
		case <-done:
			return
		}
	}
}

func handleGo(cmd *ui.Command) {
	if len(cmd.Args) != 3 {
		//TODO
		log.Println("Wrong number of parameters")
		return
	}
	host := cmd.Args[0]
	tls := cmd.Args[1] == "true"
	ri, err := strconv.Atoi(cmd.Args[2])
	if err != nil {
		log.Println(err)
		return
	}
	body := bytes.NewBuffer(cmd.Payload)
	status.RLock()
	defer status.RUnlock()
	if len(status.Repeats) <= ri {
		log.Println("Repeater out of range")
		return
	}
	r := status.Repeats[ri]
	var res io.Reader
	if res, err = r.Repeat(body, host, tls); err != nil {
		log.Println(err)
		return
	}
	_ = res
	//TODO send response
	//BOOKMARK
	uiRepeater.Send(
		ui.Command{},
	)
}

func handleGet(cmd *ui.Command) {
	if len(cmd.Args) != 2 {
		//TODO
		log.Println("Wrong number of parameters")
		return
	}
	ri, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		log.Println(err)
		return
	}
	itemn, err := strconv.Atoi(cmd.Args[1])
	if err != nil {
		log.Println(err)
		return
	}
	status.RLock()
	defer status.RUnlock()
	if len(status.Repeats) <= ri {
		log.Println("Repeater out of range")
		return
	}
	r := status.Repeats[ri]
	if len(r.history) <= itemn {
		log.Println("Repeater item out of range")
		return
	}
	repitem, err := json.Marshal(r.history[itemn])
	if err != nil {
		log.Println("Error while marshaling repeat item")
		return
	}
	uiRepeater.Send(
		ui.Command{
			Action:  GET.String(),
			Args:    cmd.Args,
			Payload: repitem,
		},
	)
}
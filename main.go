package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

var (
	weboscketID, callID, callStatus, callFrom, callTo string
)

func startChrome() (string, context.CancelFunc) {
	procCtx, closeChrome := context.WithCancel(context.Background())

	cmd := exec.CommandContext(procCtx, `C:\Program Files\Google\Chrome\Application\chrome.exe`,
		"--remote-debugging-port=9222",
		"about:blank",
	)

	// cmd := exec.CommandContext(procCtx, `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
	// 	"--remote-debugging-port=9222",
	// 	"about:blank",
	// )

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	defer stderr.Close()
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	var url string
	wsURL, err := readOutput(stderr, nil, nil)
	if err != nil {
		fmt.Println(err)
	}

	if wsURL != "" {
		url = wsURL
	} else {
		url = "ws://[::1]:9222/devtools/browser"
	}

	return url, closeChrome
}

func detectPlatform(url, platform string, alreadyListenting []target.ID) (context.Context, []target.ID) {
	allocatorContext, _ := chromedp.NewRemoteAllocator(context.Background(), url)
	//defer cancel()
	ctx, _ := chromedp.NewContext(allocatorContext)
	//defer cancel()

	relatedTabs := []target.ID{}
	// get the list of the targets
	tabs, err := chromedp.Targets(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(tabs) == 0 {
		log.Println("no tabs available")
	}

	for _, tab := range tabs {
		url := tab.URL
		if strings.Contains(url, platform) && !contains(alreadyListenting, tab.TargetID) {
			relatedTabs = append(relatedTabs, tab.TargetID)
		}
	}

	return ctx, relatedTabs
}

func contains(s []target.ID, e target.ID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	url, _ := startChrome()

	stillActive := true
	var tabCtx context.Context

	alreadyListenting := []target.ID{}
	for stillActive {
		ctx, tabids := detectPlatform(url, "dialpad.com", alreadyListenting)
		// create context attached to the specified target ID.
		// this example just uses the first target,
		// you can search for the one you want.
		for _, tabid := range tabids {
			alreadyListenting = append(alreadyListenting, tabid)
			tabCtx, _ = chromedp.NewContext(ctx, chromedp.WithTargetID(tabid))
			//defer cancel()

			// Listen for messages if the websocket url corresponds to
			// a supported application
			chromedp.ListenTarget(tabCtx, func(ev interface{}) {
				if ev, ok := ev.(*network.EventWebSocketCreated); ok {
					if strings.Contains(ev.URL, "wss://uberconf.ubervoip.net/") {
						weboscketID = ev.RequestID.String()
						fmt.Println("New socket stablished: " + weboscketID)
					}
				}
			})

			chromedp.ListenTarget(tabCtx, func(ev interface{}) {
				if ev, ok := ev.(*network.EventWebSocketClosed); ok {
					if weboscketID == ev.RequestID.String() {
						fmt.Println("The socket has been closed...")
					} else {
						fmt.Println("This socket just closed: " + ev.RequestID.String())
					}
				}
			})

			chromedp.ListenTarget(tabCtx, func(ev interface{}) {
				if ev, ok := ev.(*network.EventWebSocketFrameReceived); ok {
					parseCall(ev.Response.PayloadData)
				}

				if ev, ok := ev.(*network.EventWebSocketFrameSent); ok {
					parseCall(ev.Response.PayloadData)
				}
			})

			if err := chromedp.Run(tabCtx, chromedp.Tasks{network.Enable()}); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func parseCall(payload string) {
	if strings.HasPrefix(payload, "INVITE") {
		data := strings.Split(payload, "\n")
		for _, line := range data {
			if callID == "" || callStatus == "" || callFrom == "" || callTo == "" {
				keyvalue := strings.Split(line, ":")
				if keyvalue[0] == "From" {
					fromData := strings.Split(keyvalue[2], "@")
					callFrom = fromData[0]
				}
				if keyvalue[0] == "To" {
					toData := strings.Split(keyvalue[2], "@")
					callTo = toData[0]
				}
				if keyvalue[0] == "Call-ID" {
					callID = keyvalue[1]
				}
			} else {
				break
			}
		}
		callStatus = "call_starting"
	}

	for strings.HasPrefix(payload, "SIP/2.0 200 OK") {
		if callStatus == "call_starting" {
			callStatus = "on_call"
		}
	}
}

// https://networktest.twilio.com/

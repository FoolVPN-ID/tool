package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FoolVPN-ID/tool/api"
	"github.com/FoolVPN-ID/tool/modules/telegram"
)

func main() {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(osSignals)

	for {
		ctx, cancel := context.WithCancel(context.Background())

		// Services
		go api.RunWithContext(ctx)
		go telegram.RunWithContext(ctx)

		// Notify services to gracefully shutdown
		for {
			osSignal := <-osSignals
			cancel()

			// Initial wait
			time.Sleep(3 * time.Second)

			if osSignal != syscall.SIGHUP {
				// Exit on terminatation
				fmt.Println("Good bye...")
				os.Exit(0)
			}

			break
		}
	}
}

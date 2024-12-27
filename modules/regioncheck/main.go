package regioncheck

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/FoolVPN-ID/tool/common"
	"github.com/FoolVPN-ID/tool/modules/config"
	box "github.com/sagernet/sing-box"
)

func MakeLibrary() LibraryStruct {
	return LibraryStruct{
		Runner: []func(http.Client) runnerResultStruct{
			YoutubeCDN,
			Netflix,
		},
	}
}

func (lib *LibraryStruct) Run(rawConfig string) error {
	boxConfig, err := config.BuildSingboxConfig(rawConfig)
	if err != nil {
		return err
	}

	boxInstance, err := box.New(box.Options{
		Context: context.Background(),
		Options: boxConfig,
	})

	if err != nil {
		return err
	}

	defer boxInstance.Close()
	boxInstance.Start()

	// Build http client
	listenPort := boxConfig.Inbounds[0].MixedOptions.ListenPort
	proxyClient, _ := url.Parse(fmt.Sprintf("socks5://0.0.0.0:%d", listenPort))
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyClient),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Go routine goes here 💪🏻
	var wg sync.WaitGroup
	for _, runner := range lib.Runner {
		wg.Add(1)
		go func() {
			defer common.RecoverFromPanic()
			defer wg.Done()
			result := runner(*httpClient)
			lib.Result = append(lib.Result, result)
		}()
	}
	wg.Wait()

	return nil
}

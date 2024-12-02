package main

import (
	"flag"
	"log"

	"github.com/tunnelchaos/hopger/pkg/config"
	"github.com/tunnelchaos/hopger/pkg/pretalxConverter"
	rssconverter "github.com/tunnelchaos/hopger/pkg/rssConverter"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "config", "config.toml", "path to the config file")
	flag.Parse()
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	for _, event := range conf.Events {
		log.Printf("Event: %s", event.Name)
		for _, info := range event.Infos {
			log.Printf("  Info: %s", info.Name)
			switch info.Type {
			case config.InfoTypeRSS:
				log.Printf("    Type: RSS")
				err := rssconverter.Convert(event.Name, info, conf.Server)
				if err != nil {
					log.Printf("    Error: %v", err)
				}
			case config.InfoPretalx:
				log.Printf("    Type: Pretalx")
				pretalxConverter.Convert(event.Name, info, conf.Server)
				if err != nil {
					log.Printf("    Error: %v", err)
				}
			}
		}
	}
}

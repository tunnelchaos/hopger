package main

import (
	"flag"
	"log"

	"github.com/tunnelchaos/hopger/pkg/assemblyConverter"
	"github.com/tunnelchaos/hopger/pkg/config"
	"github.com/tunnelchaos/hopger/pkg/converter"
	"github.com/tunnelchaos/hopger/pkg/pretalxConverter"
	"github.com/tunnelchaos/hopger/pkg/rssConverter"
)

var (
	configPath string
)

var converterRegistry = map[config.InfoType]converter.Converter{
	config.InfoTypeRSS:       &rssConverter.RSSConverter{},
	config.InfoPretalx:       &pretalxConverter.PretalxConverter{},
	config.InfoHubAssemblies: &assemblyConverter.AssemblyConverter{},
}

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
			converter, exists := converterRegistry[info.Type]
			if !exists {
				log.Println("   No converter found for this info type:", info.Type)
				continue
			}
			err := converter.Convert(event.Name, info, conf.Server)
			if err != nil {
				log.Printf("   Failed to convert info: %v", err)
			}
		}
	}
}

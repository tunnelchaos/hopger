package converter

import "github.com/tunnelchaos/go-packages/config"

type Converter interface {
	Convert(eventname string, info config.Info, server config.Server) error
}

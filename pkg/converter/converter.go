package converter

import "github.com/tunnelchaos/hopger/pkg/config"

type Converter interface {
	Convert(eventname string, info config.Info, server config.Server) error
}

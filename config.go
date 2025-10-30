
package gurobireceiver

import "go.opentelemetry.io/collector/component"

type Config struct {
    component.ReceiverSettings `mapstructure:",squash"`
    URI        string `mapstructure:"uri"`
    AccessID   string `mapstructure:"access_id"`
    SecretKey  string `mapstructure:"secret_key"`
    Interval   string `mapstructure:"interval"`
}

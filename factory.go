package gurobireceiver

import (
    "context"
    "go.opentelemetry.io/collector/component"
    "go.opentelemetry.io/collector/receiver/receiverhelper"
    "go.opentelemetry.io/collector/receiver"
)

const typeStr = "gurobi"

func NewFactory() receiver.Factory {
    return receiver.NewFactory(
        typeStr,
        createDefaultConfig,
        receiver.WithMetrics(createMetricsReceiver),
    )
}

func createDefaultConfig() component.Config {
    return &Config{
        ReceiverSettings: component.NewReceiverSettings(component.NewID(typeStr)),
        URI:              "",
        AccessID:         "",
        SecretKey:        "",
        Interval:         "30s",
    }
}
package gurobireceiver
import (
    "encoding/json"
    "net/http"
    "time"

    "go.opentelemetry.io/collector/consumer"
    "go.opentelemetry.io/collector/receiver"
    "go.opentelemetry.io/collector/pdata/pmetric"
)

func createMetricsReceiver(
    ctx context.Context,
    settings receiver.CreateSettings,
    cfg component.Config,
    nextConsumer consumer.Metrics,
) (receiver.Metrics, error) {
    rCfg := cfg.(*Config)

    ticker := time.NewTicker(time.Duration(30) * time.Second)

    go func() {
        for {
            select {
            case <-ticker.C:
                metrics := fetchGurobiMetrics(rCfg)
                _ = nextConsumer.ConsumeMetrics(ctx, metrics)
            case <-ctx.Done():
                ticker.Stop()
                return
            }
        }
    }()

    return receiverhelper.NewNopMetricsReceiver(), nil
}

func fetchGurobiMetrics(cfg *Config) pmetric.Metrics {
    req, _ := http.NewRequest("GET", cfg.URI, nil)
    req.Header.Set("X-GUROBI-ACCESS-ID", cfg.AccessID)
    req.Header.Set("X-GUROBI-SECRET-KEY", cfg.SecretKey)
    req.Header.Set("Accept", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return pmetric.NewMetrics()
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    metrics := pmetric.NewMetrics()
    rm := metrics.ResourceMetrics().AppendEmpty()
    sm := rm.ScopeMetrics().AppendEmpty()
    m := sm.Metrics().AppendEmpty()
    m.SetName("gurobi_job_count")
    m.SetEmptyGauge().DataPoints().AppendEmpty().SetIntValue(int64(len(result["jobs"].([]interface{}))))

    return metrics
}
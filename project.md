# Weather Observatory Project

## Objective

Build a Go service that:

- Collects weather data from public APIs
- Exposes Prometheus metrics
- Produces structured logs
- Emits OpenTelemetry traces
- Visualizes data in Grafana
- Uses Loki and Tempo for observability

---

# Phase 1 - Data Collector

## Task 1.1 - Create Project

- Initialize Go module
- Create repository structure

### Deliverables

- Go project initialized
- Basic README

---

## Task 1.2 - Select Weather API

### Options

- National Weather Service (recommended)
- OpenWeather
- WeatherAPI

### Deliverables

- Working API client
- Configuration support

---

## Task 1.3 - Build Weather Model

Create:

```go
type WeatherObservation struct {
    Temperature float64
    Humidity    float64
    WindSpeed   float64
    Timestamp   time.Time
}
```

### Deliverables

- Weather data model
- JSON parsing

---

## Task 1.4 - Poll Weather Data

- Fetch every 5 minutes
- Store latest observation

### Deliverables

- Scheduler
- Retry handling

---

# Phase 2 - Prometheus Metrics

## Task 2.1 - Install Prometheus Client

### Deliverables

- Metrics package integrated

---

## Task 2.2 - Create Gauges

Metrics:

- weather_temperature_f
- weather_humidity_percent
- weather_wind_speed_mph

---

## Task 2.3 - Create Counters

Metrics:

- weather_api_requests_total
- weather_api_errors_total

---

## Task 2.4 - Expose /metrics

### Deliverables

- Metrics endpoint
- Validation via curl

---

# Phase 3 - Prometheus Deployment

## Task 3.1 - Docker Compose

Services:

- Prometheus
- Grafana

---

## Task 3.2 - Configure Scraping

### Deliverables

- Prometheus scrape configuration

---

## Task 3.3 - Validate Metrics

Queries:

- weather_temperature_f
- weather_api_requests_total

---

# Phase 4 - Grafana Dashboards

## Task 4.1 - Configure Datasource

- Add Prometheus datasource

---

## Task 4.2 - Create Dashboard

Panels:

### Current Conditions

- Temperature
- Humidity
- Wind Speed

### Historical Trends

- Temperature over time
- Humidity over time

### API Health

- Request count
- Error count

---

## Task 4.3 - Add Variables

Variables:

- City
- Environment

---

## Task 4.4 - Add Thresholds

Temperature:

- Green < 75
- Yellow < 90
- Red >= 90

---

# Phase 5 - Logging

## Task 5.1 - Structured Logging

Suggested libraries:

- zap
- zerolog

---

## Task 5.2 - Deploy Loki

### Deliverables

- Loki container

---

## Task 5.3 - Deploy Promtail

### Deliverables

- Log collection

---

## Task 5.4 - Grafana Log Dashboard

Visualizations:

- API failures
- Collector activity

---

# Phase 6 - Tracing

## Task 6.1 - Install OpenTelemetry

### Deliverables

- OTel SDK integrated

---

## Task 6.2 - Create Spans

Trace flow:

- GetWeather
- HTTPCall
- ParseResponse

---

## Task 6.3 - Deploy Tempo

### Deliverables

- Tempo container
- OTel Collector

---

## Task 6.4 - Export Traces

Flow:

Go App -> OTel Collector -> Tempo

---

# Phase 7 - Correlation

## Task 7.1 - Link Metrics, Logs, and Traces

### Deliverables

- End-to-end observability workflow

---

## Task 7.2 - Add Trace IDs to Logs

### Deliverables

- Correlated troubleshooting

---

## Task 7.3 - Incident Dashboard

Dashboard Panels:

- Temperature
- API latency
- Error rate
- Recent failures
- Trace explorer

---

# Phase 8 - Enhancements

## Option A - Multi-City Monitoring

Support:

- St Paul
- Minneapolis
- Chicago
- Denver

---

## Option B - Forecast Accuracy

Compare:

- Forecasted temperature
- Actual temperature

---

## Option C - Severe Weather Alerts

Track:

- Tornado watches
- Storm warnings
- Flood alerts

---

## Option D - Air Quality

Track:

- AQI
- Smoke
- Pollen

---

# Stretch Goal

Route all telemetry through OpenTelemetry Collector:

Go App -> OpenTelemetry Collector -> Prometheus / Loki / Tempo -> Grafana

## Learning Outcomes

- Go development
- Metrics design
- Prometheus
- Grafana dashboards
- Loki logging
- Tempo tracing
- OpenTelemetry
- Production-style observability

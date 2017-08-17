mackerel-plugin-timeline
=====================

timeline custom metrics plugin for mackerel.io agent.

## mackerel-plugin-timeline

### Interface

```
type TimeLine interface {
	ToConut(line string) error
	ToMetrics(metricName string) map[string]interface{}
	ToGraph(metricName string) map[string]mackerelplugin.Graphs
	ParseTime(line string) time.Time
}
```

### Synopsis

```
command [options] file
  -m int
        minute interval (default 1)
  -metric string
        specify a to metric name
  -datetime string
        start datetime
  -location string
        datetime location name (default "Asia/Tokyo")
  -i    display info
  -h    this help
  -v    show version and exit
```
mackerel-plugin-timeline
=====================

timeline custom metrics plugin for mackerel.io agent.

## mackerel-plugin-timeline

### Interface

```
type TimeLine interface {
	ToConut(line string) error
	ToMetrics() map[string]interface{}
	ToGraph() map[string]mackerelplugin.Graphs
	ParseTime(line string) time.Time
}
```

### Synopsis

```
command [options] file
  -datetime string
        start datetime
  -h    this help
  -layout string
        datetime layout (default "2006-01-02T15:04:05Z07:00")
  -location string
        datetime location name (default "Asia/Tokyo")
  -m int
        time minutes (default 5)
  -v    show version and exit
```
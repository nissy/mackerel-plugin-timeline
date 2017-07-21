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
  -m int
        minute interval (default 1)
  -datetime string
        start datetime
  -location string
        datetime location name (default "Asia/Tokyo")
  -i    display info
  -h    this help
  -v    show version and exit
```
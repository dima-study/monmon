package providers

import (
	_ "github.com/dima-study/monmon/internal/stats/providers/cpuload" // init для cpuload
	_ "github.com/dima-study/monmon/internal/stats/providers/loadavg" // init для loadavg
)

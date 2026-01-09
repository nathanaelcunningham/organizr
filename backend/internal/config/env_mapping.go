package config

var envKeyMap = map[string]string{
	"qbittorrent.url":          "QBITTORRENT_URL",
	"qbittorrent.username":     "QBITTORRENT_USERNAME",
	"qbittorrent.password":     "QBITTORRENT_PASSWORD",
	"paths.destination":        "PATHS_DESTINATION",
	"paths.template":           "PATHS_TEMPLATE",
	"paths.no_series_template": "PATHS_NO_SERIES_TEMPLATE",
	"paths.operation":          "PATHS_OPERATION",
	"paths.local_mount":        "PATHS_LOCAL_MOUNT",
	"monitor.interval_seconds": "MONITOR_INTERVAL_SECONDS",
	"monitor.auto_organize":    "MONITOR_AUTO_ORGANIZE",
	"mam.baseurl":              "MAM_BASEURL",
	"mam.secret":               "MAM_SECRET",
}

func getEnvKey(dbKey string) string {
	return envKeyMap[dbKey]
}

package opksshplugingoogleworkspace

import "time"

const (
	DefaultConfigPath    = "/etc/opkssh-plugin-google-workspace/config.yaml"
	DefaultCachePath     = "/var/cache/opkssh-plugin-google-workspace/cache.json"
	DefaultLogPath       = "/var/log/opkssh-plugin-google-workspace.log"
	DefaultCacheDuration = time.Minute * 15
)

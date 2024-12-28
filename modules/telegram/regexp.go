package telegram

import "regexp"

var (
	PROXY_IP_REGEXP   = regexp.MustCompile(`^[\d\.]+:\d+`)
	CONFIG_VPN_REGEXP = regexp.MustCompile(`\w+:\/\/.+$`)
)

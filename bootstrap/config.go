package bootstrap

import (
	"flag"
	"fmt"
	"net"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/clicklord/lms/config"
	"github.com/clicklord/lms/log"
)

func LoadConfig(logger log.Logger) (*config.DmsConfig, error) {
	// default config
	var cfg = &config.DmsConfig{
		Path:             "",
		IfName:           "",
		Http:             ":1338",
		FriendlyName:     "",
		DeviceIcon:       "",
		DeviceIconSizes:  []string{"48,128"},
		LogHeaders:       false,
		FFprobeCachePath: config.GetDefaultFFprobeCachePath(),
		ForceTranscodeTo: "",
	}

	path := flag.String("path", cfg.Path, "browse root path")
	ifName := flag.String("ifname", cfg.IfName, "specific SSDP network interface")
	http := flag.String("http", cfg.Http, "http server port")
	friendlyName := flag.String("friendlyName", cfg.FriendlyName, "server friendly name")
	deviceIcon := flag.String("deviceIcon", cfg.DeviceIcon, "device defaultIcon")
	deviceIconSizes := flag.String("deviceIconSizes", strings.Join(cfg.DeviceIconSizes, ","), "comma separated list of icon sizes to advertise, eg 48,128,256. Use 48:512,128:512 format to force actual size.")
	logHeaders := flag.Bool("logHeaders", cfg.LogHeaders, "log HTTP headers")
	fFprobeCachePath := flag.String("fFprobeCachePath", cfg.FFprobeCachePath, "path to FFprobe cache file")
	configFilePath := flag.String("config", "", "json configuration file")
	allowedIps := flag.String("allowedIps", "", "allowed ip of clients, separated by comma")
	forceTranscodeTo := flag.String("forceTranscodeTo", cfg.ForceTranscodeTo, "force transcoding to certain format, supported: 'chromecast', 'vp8', 'web'")
	transcodeLogPattern := flag.String("transcodeLogPattern", "", "pattern where to write transcode logs to. The [tsname] placeholder is replaced with the name of the item currently being played. The default is $HOME/.dms/log/[tsname]")
	flag.BoolVar(&cfg.NoTranscode, "noTranscode", false, "disable transcoding")
	flag.BoolVar(&cfg.NoProbe, "noProbe", false, "disable media probing with ffprobe")
	flag.BoolVar(&cfg.StallEventSubscribe, "stallEventSubscribe", false, "workaround for some bad event subscribers")
	flag.DurationVar(&cfg.NotifyInterval, "notifyInterval", 30*time.Second, "interval between SSPD announces")
	flag.BoolVar(&cfg.IgnoreHidden, "ignoreHidden", false, "ignore hidden files and directories")
	flag.BoolVar(&cfg.IgnoreUnreadable, "ignoreUnreadable", false, "ignore unreadable files and directories")
	ignorePaths := flag.String("ignore", "", "comma separated list of directories to ignore (i.e. thumbnails,thumbs)")
	flag.BoolVar(&cfg.AllowDynamicStreams, "allowDynamicStreams", false, "activate support for dynamic streams described via .dms.json metadata files")

	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		return nil, fmt.Errorf("%s: %s\n", "unexpected positional arguments", flag.Args())
	}

	cfg.Path, _ = filepath.Abs(*path)
	cfg.IfName = *ifName
	cfg.Http = *http
	cfg.FriendlyName = *friendlyName
	cfg.DeviceIcon = *deviceIcon
	cfg.DeviceIconSizes = strings.Split(*deviceIconSizes, ",")

	cfg.LogHeaders = *logHeaders
	cfg.FFprobeCachePath = *fFprobeCachePath
	cfg.AllowedIpNets = makeIpNets(*allowedIps)
	cfg.ForceTranscodeTo = *forceTranscodeTo
	cfg.IgnorePaths = strings.Split(*ignorePaths, ",")
	cfg.TranscodeLogPattern = *transcodeLogPattern

	if cfg.TranscodeLogPattern == "" {
		u, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("unable to resolve current user: %v", err)
		}
		cfg.TranscodeLogPattern = filepath.Join(u.HomeDir, ".lms", "log", "[tsname]")
	}

	if len(*configFilePath) == 0 {
		configFilePath = config.GetDefaultConfigPath()
	}
	cfg.Load(*configFilePath)

	logger.Printf("device icon sizes are %s", cfg.DeviceIconSizes)
	logger.Printf("allowed ip nets are %s", cfg.AllowedIpNets)
	logger.Printf("serving folder %s", cfg.Path)
	if cfg.AllowDynamicStreams {
		logger.Printf("Dynamic streams ARE allowed")
	}

	return cfg, nil
}

func makeIpNets(s string) []*net.IPNet {
	var nets []*net.IPNet
	if len(s) < 1 {
		_, ipnet, _ := net.ParseCIDR("0.0.0.0/0")
		nets = append(nets, ipnet)
		_, ipnet, _ = net.ParseCIDR("::/0")
		nets = append(nets, ipnet)
	} else {
		for _, el := range strings.Split(s, ",") {
			ip := net.ParseIP(el)

			if ip == nil {
				_, ipnet, err := net.ParseCIDR(el)
				if err == nil {
					nets = append(nets, ipnet)
				} else {
					log.Printf("unable to parse expression %s", el)
				}

			} else {
				_, ipnet, err := net.ParseCIDR(el + "/32")
				if err == nil {
					nets = append(nets, ipnet)
				} else {
					log.Printf("unable to parse ip %s", el)
				}
			}
		}
	}
	return nets
}

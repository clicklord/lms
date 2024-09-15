package bootstrap

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/clicklord/lms/config"
	"github.com/clicklord/lms/dlna/dms"
	"github.com/clicklord/lms/log"
	"github.com/nfnt/resize"
)

func LoadDMS(
	cfg *config.DmsConfig,
	cache *config.FFprobeCache,
	defaultIcon []byte,
	logger log.Logger,
) *dms.Server {
	dmsServer := &dms.Server{
		Logger: logger,
		Interfaces: func(ifName string) (ifs []net.Interface) {
			var err error
			if ifName == "" {
				ifs, err = net.Interfaces()
			} else {
				var if_ *net.Interface
				if_, err = net.InterfaceByName(ifName)
				if if_ != nil {
					ifs = append(ifs, *if_)
				}
			}
			if err != nil {
				log.Fatal(err)
			}
			var tmp []net.Interface
			for _, if_ := range ifs {
				if if_.Flags&net.FlagUp == 0 || if_.MTU <= 0 {
					continue
				}
				tmp = append(tmp, if_)
			}
			ifs = tmp
			return
		}(cfg.IfName),
		HTTPConn: func() net.Listener {
			conn, err := net.Listen("tcp", cfg.Http)
			if err != nil {
				log.Fatal(err)
			}
			return conn
		}(),
		FriendlyName:        cfg.FriendlyName,
		RootObjectPath:      filepath.Clean(cfg.Path),
		FFProbeCache:        cache,
		LogHeaders:          cfg.LogHeaders,
		NoTranscode:         cfg.NoTranscode,
		AllowDynamicStreams: cfg.AllowDynamicStreams,
		ForceTranscodeTo:    cfg.ForceTranscodeTo,
		TranscodeLogPattern: cfg.TranscodeLogPattern,
		NoProbe:             cfg.NoProbe,
		Icons: func() []dms.Icon {
			var icons []dms.Icon
			for _, size := range cfg.DeviceIconSizes {
				s := strings.Split(size, ":")
				if len(s) != 1 && len(s) != 2 {
					log.Fatal("bad device icon size: ", size)
				}
				advertisedSize, err := strconv.Atoi(s[0])
				if err != nil {
					log.Fatal("bad device icon size: ", size)
				}
				actualSize := advertisedSize
				if len(s) == 2 {
					// Force actual icon size to be different from advertised
					actualSize, err = strconv.Atoi(s[1])
					if err != nil {
						log.Fatal("bad device icon size: ", size)
					}
				}
				icons = append(icons, dms.Icon{
					Width:    advertisedSize,
					Height:   advertisedSize,
					Depth:    8,
					Mimetype: "image/png",
					Bytes:    readIcon(cfg.DeviceIcon, uint(actualSize), defaultIcon),
				})
			}
			return icons
		}(),
		StallEventSubscribe: cfg.StallEventSubscribe,
		NotifyInterval:      cfg.NotifyInterval,
		IgnoreHidden:        cfg.IgnoreHidden,
		IgnoreUnreadable:    cfg.IgnoreUnreadable,
		IgnorePaths:         cfg.IgnorePaths,
		AllowedIpNets:       cfg.AllowedIpNets,
	}

	return dmsServer
}

func getIconReader(path string, defaultIcon []byte) (io.ReadCloser, error) {
	if path == "" {
		return io.NopCloser(bytes.NewReader(defaultIcon)), nil
	}
	return os.Open(path)
}

func readIcon(path string, size uint, defaultIcon []byte) []byte {
	r, err := getIconReader(path, defaultIcon)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	imageData, _, err := image.Decode(r)
	if err != nil {
		panic(err)
	}
	return resizeImage(imageData, size)
}

func resizeImage(imageData image.Image, size uint) []byte {
	img := resize.Resize(size, size, imageData, resize.Lanczos3)
	var buff bytes.Buffer
	png.Encode(&buff, img)
	return buff.Bytes()
}

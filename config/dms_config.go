package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/clicklord/lms/dlna/dms"
	"github.com/clicklord/lms/rrcache"
)

type DmsConfig struct {
	Path                string
	IfName              string
	Http                string
	FriendlyName        string
	DeviceIcon          string
	DeviceIconSizes     []string
	LogHeaders          bool
	FFprobeCachePath    string
	NoTranscode         bool
	ForceTranscodeTo    string
	NoProbe             bool
	StallEventSubscribe bool
	NotifyInterval      time.Duration
	IgnoreHidden        bool
	IgnoreUnreadable    bool
	IgnorePaths         []string
	AllowedIpNets       []*net.IPNet
	AllowDynamicStreams bool
	TranscodeLogPattern string
}

func (config *DmsConfig) Load(configPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("config error (config file: '%s'): %v\n", configPath, err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Printf("config error: %v\n", err)
		return
	}
}

func GetDefaultConfigPath() *string {
	path := ""
	_user, err := user.Current()
	if err != nil {
		log.Print(err)
		return &path
	}
	path = filepath.Join(_user.HomeDir, ".lms/config.json")
	return &path
}

func GetDefaultFFprobeCachePath() (path string) {
	_user, err := user.Current()
	if err != nil {
		log.Print(err)
		return
	}
	path = filepath.Join(_user.HomeDir, ".lms/ffprobe-cache")
	return
}

type FFprobeCache struct {
	Cache *rrcache.RRCache
	sync.Mutex
}

func (fc *FFprobeCache) Get(key interface{}) (value interface{}, ok bool) {
	fc.Lock()
	defer fc.Unlock()
	return fc.Cache.Get(key)
}

func (fc *FFprobeCache) Set(key interface{}, value interface{}) {
	fc.Lock()
	defer fc.Unlock()
	var size int64
	for _, v := range []interface{}{key, value} {
		b, err := json.Marshal(v)
		if err != nil {
			log.Printf("Could not marshal %v: %s", v, err)
			continue
		}
		size += int64(len(b))
	}
	fc.Cache.Set(key, value, size)
}

func (cache *FFprobeCache) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var items []dms.FfprobeCacheItem
	err = dec.Decode(&items)
	if err != nil {
		return err
	}
	for _, item := range items {
		cache.Set(item.Key, item.Value)
	}
	log.Printf("added %d items from cache", len(items))
	return nil
}

func (cache *FFprobeCache) Save(path string) error {
	cache.Lock()
	items := cache.Cache.Items()
	cache.Unlock()
	f, err := ioutil.TempFile(filepath.Dir(path), filepath.Base(path))
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	err = enc.Encode(items)
	f.Close()
	if err != nil {
		os.Remove(f.Name())
		return err
	}
	if runtime.GOOS == "windows" {
		err = os.Remove(path)
		if err == os.ErrNotExist {
			err = nil
		}
	}
	if err == nil {
		err = os.Rename(f.Name(), path)
	}
	if err == nil {
		log.Printf("saved cache with %d items", len(items))
	} else {
		os.Remove(f.Name())
	}
	return err
}

# LMS

Local Media Server. UPnP DLNA server that includes basic video transcoding. Tested on a LG television, several Android UPnP apps, and Chromecast.

Project based on https://github.com/anacrolix/dms

## LMS config params

|parameter  | description |
|--|--|
| ``allowDynamicStreams`` | turns on support for `.lms.json` files in the path |
| ``allowedIps string`` | allowed ip of clients, separated by comma |
| ``config string`` | json configuration file |
| ``deviceIcon string`` | device icon |
| ``deviceIconSizes string`` | device icon sizes, separated by comma |
| ``fFprobeCachePath string`` | path to FFprobe cache file (default ``$HOME/.lms/ffprobe-cache``) |
| ``forceTranscodeTo string`` | force transcoding to certain format, supported: 'chromecast', 'vp8' |
| ``friendlyName string`` | server friendly name |
| ``http string`` | http server port (default ":1338") |
| ``ignoreHidden`` | ignore hidden files and directories |
| ``ignoreUnreadable`` | ignore unreadable files and directories |
| ``ignore`` | ignore comma separated list of paths (i.e. -ignore thumbnails,thumbs) |
| ``logHeaders`` | log HTTP headers |
| ``noProbe`` | disable media probing with ffprobe |
| ``noTranscode`` | disable transcoding |
| ``notifyInterval duration`` | interval between SSPD announces (default 30s) |
| ``path string`` | browse root path |
| ``stallEventSubscribe`` | workaround for some bad event subscribers |
| ``transcodeLogPattern`` | pattern where to write transcode logs to. The ``[tsname]`` placeholder is replaced with the name of the item currently being played. The default is ``$HOME/.lms/log/[tsname]``. You may turn off transcode logging entirely by setting it to ``/dev/null``. You may log to stderr by setting ``/dev/stderr`` |

An example json configuration file:
```
{
  "path": "/path/to/media/files",
  "friendlyName": "lms",
  "noTranscode": true,
  "deviceIcon": "/path/to/icon.png",
  "deviceIconSizes": ["48:512","128:512"]
}
```
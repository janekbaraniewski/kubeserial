package managers


var Available = map[string]*Manager{
	"octoprint":	&Octoprint,
  "openwebrx":  &OpenWebRX,
}

var Octoprint = Manager{
	Image: 		"janekbaraniewski/octoprint:1.3.10",
	RunCmnd: 	"mkdir /root/.octoprint && cp /data/config.yaml /root/.octoprint/config.yaml && /OctoPrint-1.3.10/run --iknowwhatimdoing --port 80",
	Config: 	octoprintConfig,
	ConfigPath:	"/data/config.yaml",
}

var OpenWebRX = Manager{
  Image:      "janekbaraniewski/openwebrx:latest", // TODO: fix version
  RunCmnd:    "/src/openwebrx/openwebrx.py",
  Config:     "",
  ConfigPath: "/tmp/conf.yaml",
}

var octoprintConfig = `
accessControl:
  enabled: false
plugins:
  announcements:
    _config_version: 1
    channels:
      _blog:
        read_until: 1573642500
      _important:
        read_until: 1521111600
      _octopi:
        read_until: 1573722900
      _plugins:
        read_until: 1573862400
      _releases:
        read_until: 1574699400
  discovery:
    upnpUuid: ef35acc7-a859-4947-980d-d5edb10508e4
  softwareupdate:
    _config_version: 6
  tracking:
    enabled: false
deviceProfiles:
  default: _default
serial:
  additionalPorts:
  - /dev/device
  autoconnect: true
  baudrate: 0
  port: /dev/device
server:
  firstRun: false
  onlineCheck:
    enabled: true
  pluginBlacklist:
    enabled: false
  seenWizards:
    corewizard: 3
    cura: null
    tracking: null`

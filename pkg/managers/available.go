package managers


var Available = map[string]*Manager{
	"octoprint":	&Octoprint,
}

var Octoprint = Manager{
	Image: "janekbaraniewski/octoprint:1.3.10",
	RunCmnd: " mkdir /root/.octoprint && cp /data/config.yaml /root/.octoprint/config.yaml && /OctoPrint-1.3.10/run --iknowwhatimdoing --port 80",
}

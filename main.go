package main

import (
	"flag"
	"github.com/pkg/profile"
	"github.com/revfyawo/gogame/ecs"
	"github.com/revfyawo/gogame/engine"
	"github.com/revfyawo/gogame/systems"
	"log"
	"os"
)

type defaultScene struct{}

func (s *defaultScene) Setup(w *ecs.World) {
	w.AddSystem(&systems.Camera{})
	w.AddSystem(&systems.MouseChunk{})
	w.AddSystem(&systems.MouseZoom{})
	w.AddSystem(&systems.ChunkRender{})
	w.AddSystem(&systems.Grid{})
	w.AddSystem(&systems.LandscapeGen{})
	w.AddSystem(&systems.LandscapeRender{})
	w.AddSystem(&systems.ChunkGen{})
	w.AddSystem(&systems.GUI{})
	w.AddSystem(&systems.UpdatesCounter{})
}

var cpuprofile = flag.Bool("cpuprofile", false, "profile CPU usage")
var memprofile = flag.Bool("memprofile", false, "profile memory usage")

func main() {
	flag.Parse()
	if *cpuprofile && *memprofile {
		log.Print("can't use both cpu & memory profiling")
		os.Exit(1)
	} else if *cpuprofile {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	} else if *memprofile {
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	}
	engine.Run(&defaultScene{})
}

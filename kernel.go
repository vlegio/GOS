package kernel

import (
	"screen"
)

func Load() {
	screen.Init()
	screen.Clear()
	screen.PrintStr("Hello world!")
}

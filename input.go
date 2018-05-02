package main

import (
	"fmt"
	"github.com/RobLoach/go-libretro/libretro"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const numPlayers = 5
const DeviceIDJoypadMenuToggle uint32 = 16
const DeviceIDJoypadFullscreenToggle uint32 = 17

var keyBinds = map[glfw.Key]uint32{
	glfw.KeyX:         libretro.DeviceIDJoypadA,
	glfw.KeyZ:         libretro.DeviceIDJoypadB,
	glfw.KeyA:         libretro.DeviceIDJoypadY,
	glfw.KeyS:         libretro.DeviceIDJoypadX,
	glfw.KeyUp:        libretro.DeviceIDJoypadUp,
	glfw.KeyDown:      libretro.DeviceIDJoypadDown,
	glfw.KeyLeft:      libretro.DeviceIDJoypadLeft,
	glfw.KeyRight:     libretro.DeviceIDJoypadRight,
	glfw.KeyEnter:     libretro.DeviceIDJoypadStart,
	glfw.KeyBackspace: libretro.DeviceIDJoypadSelect,
}

var buttonBinds = map[byte]uint32{
	0:  libretro.DeviceIDJoypadUp,
	1:  libretro.DeviceIDJoypadDown,
	2:  libretro.DeviceIDJoypadLeft,
	3:  libretro.DeviceIDJoypadRight,
	4:  libretro.DeviceIDJoypadStart,
	5:  libretro.DeviceIDJoypadSelect,
	6:  libretro.DeviceIDJoypadL3,
	7:  libretro.DeviceIDJoypadR3,
	8:  libretro.DeviceIDJoypadL,
	9:  libretro.DeviceIDJoypadR,
	10: DeviceIDJoypadMenuToggle, // Special case
	11: libretro.DeviceIDJoypadB,
	12: libretro.DeviceIDJoypadA,
	13: libretro.DeviceIDJoypadY,
	14: libretro.DeviceIDJoypadX,
}

// Input state for all the players
var (
	newState [numPlayers][libretro.DeviceIDJoypadR3 + 3]bool
	oldState [numPlayers][libretro.DeviceIDJoypadR3 + 3]bool
	released [numPlayers][libretro.DeviceIDJoypadR3 + 3]bool
	pressed  [numPlayers][libretro.DeviceIDJoypadR3 + 3]bool
)

func joystickCallback(joy int, event int) {
	var message string
	switch event {
	case 262145:
		message = fmt.Sprintf("Joystick #%d plugged: %s.", joy, glfw.GetJoystickName(glfw.Joystick(joy)))
		break
	case 262146:
		message = fmt.Sprintf("Joystick #%d unplugged.", joy)
		break
	default:
		message = fmt.Sprintf("Joystick #%d unhandled event: %d.", joy, event)
	}
	fmt.Printf("[Input]: %s\n", message)
	notify(message, 120)
}

func inputInit() {
	glfw.SetJoystickCallback(joystickCallback)
}

func inputPoll() {
	// Reset all retropad buttons to false
	for p := range newState {
		for k := range newState[p] {
			newState[p][k] = false
		}
	}

	// Process joypads of all players
	for p := range newState {
		buttonState := glfw.GetJoystickButtons(glfw.Joystick(p))
		if len(buttonState) > 0 {
			for k, v := range buttonBinds {
				if glfw.Action(buttonState[k]) == glfw.Press {
					newState[p][v] = true
				}
			}
		}
	}

	// Process keyboard keys
	for k, v := range keyBinds {
		if window.GetKey(k) == glfw.Press {
			newState[0][v] = true
		}
	}

	// Close on escape
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	// Toggle menu when P is pressed
	if window.GetKey(glfw.KeyP) == glfw.Press {
		newState[0][DeviceIDJoypadMenuToggle] = true
	}

	if window.GetKey(glfw.KeyF) == glfw.Press {
		newState[0][DeviceIDJoypadFullscreenToggle] = true
	}

	// Compute the keys pressed or released during this frame
	for p := range newState {
		for k := range newState[p] {
			pressed[p][k] = newState[p][k] && !oldState[p][k]
			released[p][k] = !newState[p][k] && oldState[p][k]
		}
	}

	// Toggle the menu if DeviceIDJoypadMenuToggle is pressed
	if released[0][DeviceIDJoypadMenuToggle] {
		menuActive = !menuActive
	}

	// Toggle fullscreen if DeviceIDJoypadFullscreenToggle is pressed
	if released[0][DeviceIDJoypadFullscreenToggle] {
		toggleFullscreen()
	}

	// Store the old input state for comparisions
	oldState = newState
}

func inputState(port uint, device uint32, index uint, id uint) int16 {
	if id >= 255 || index > 0 || device != libretro.DeviceJoypad {
		return 0
	}

	if newState[port][id] {
		return 1
	}
	return 0
}

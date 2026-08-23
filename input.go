package canvas

import "github.com/gopxl/pixel/v2"

// Key represents a keyboard key or button.
type Key = pixel.Button

// MouseButton represents a mouse button.
type MouseButton = pixel.Button

// Vec represents a 2D vector / coordinate.
type Vec = pixel.Vec

// V returns a new 2D vector.
func V(x, y float64) Vec {
	return pixel.V(x, y)
}

// Mouse buttons
const (
	MouseLeft   MouseButton = pixel.MouseButtonLeft
	MouseRight  MouseButton = pixel.MouseButtonRight
	MouseMiddle MouseButton = pixel.MouseButtonMiddle
)

// Keyboard keys
const (
	KeySpace        Key = pixel.KeySpace
	KeyApostrophe   Key = pixel.KeyApostrophe
	KeyComma        Key = pixel.KeyComma
	KeyMinus        Key = pixel.KeyMinus
	KeyPeriod       Key = pixel.KeyPeriod
	KeySlash        Key = pixel.KeySlash
	Key0            Key = pixel.Key0
	Key1            Key = pixel.Key1
	Key2            Key = pixel.Key2
	Key3            Key = pixel.Key3
	Key4            Key = pixel.Key4
	Key5            Key = pixel.Key5
	Key6            Key = pixel.Key6
	Key7            Key = pixel.Key7
	Key8            Key = pixel.Key8
	Key9            Key = pixel.Key9
	KeySemicolon    Key = pixel.KeySemicolon
	KeyEqual        Key = pixel.KeyEqual
	KeyA            Key = pixel.KeyA
	KeyB            Key = pixel.KeyB
	KeyC            Key = pixel.KeyC
	KeyD            Key = pixel.KeyD
	KeyE            Key = pixel.KeyE
	KeyF            Key = pixel.KeyF
	KeyG            Key = pixel.KeyG
	KeyH            Key = pixel.KeyH
	KeyI            Key = pixel.KeyI
	KeyJ            Key = pixel.KeyJ
	KeyK            Key = pixel.KeyK
	KeyL            Key = pixel.KeyL
	KeyM            Key = pixel.KeyM
	KeyN            Key = pixel.KeyN
	KeyO            Key = pixel.KeyO
	KeyP            Key = pixel.KeyP
	KeyQ            Key = pixel.KeyQ
	KeyR            Key = pixel.KeyR
	KeyS            Key = pixel.KeyS
	KeyT            Key = pixel.KeyT
	KeyU            Key = pixel.KeyU
	KeyV            Key = pixel.KeyV
	KeyW            Key = pixel.KeyW
	KeyX            Key = pixel.KeyX
	KeyY            Key = pixel.KeyY
	KeyZ            Key = pixel.KeyZ
	KeyLeftBracket  Key = pixel.KeyLeftBracket
	KeyBackslash    Key = pixel.KeyBackslash
	KeyRightBracket Key = pixel.KeyRightBracket
	KeyGraveAccent  Key = pixel.KeyGraveAccent
	KeyEscape       Key = pixel.KeyEscape
	KeyEnter        Key = pixel.KeyEnter
	KeyTab          Key = pixel.KeyTab
	KeyBackspace    Key = pixel.KeyBackspace
	KeyInsert       Key = pixel.KeyInsert
	KeyDelete       Key = pixel.KeyDelete
	KeyRight        Key = pixel.KeyRight
	KeyLeft         Key = pixel.KeyLeft
	KeyDown         Key = pixel.KeyDown
	KeyUp           Key = pixel.KeyUp
	KeyPageUp       Key = pixel.KeyPageUp
	KeyPageDown     Key = pixel.KeyPageDown
	KeyHome         Key = pixel.KeyHome
	KeyEnd          Key = pixel.KeyEnd
	KeyCapsLock     Key = pixel.KeyCapsLock
	KeyScrollLock   Key = pixel.KeyScrollLock
	KeyNumLock      Key = pixel.KeyNumLock
	KeyPrintScreen  Key = pixel.KeyPrintScreen
	KeyPause        Key = pixel.KeyPause
	KeyF1           Key = pixel.KeyF1
	KeyF2           Key = pixel.KeyF2
	KeyF3           Key = pixel.KeyF3
	KeyF4           Key = pixel.KeyF4
	KeyF5           Key = pixel.KeyF5
	KeyF6           Key = pixel.KeyF6
	KeyF7           Key = pixel.KeyF7
	KeyF8           Key = pixel.KeyF8
	KeyF9           Key = pixel.KeyF9
	KeyF10          Key = pixel.KeyF10
	KeyF11          Key = pixel.KeyF11
	KeyF12          Key = pixel.KeyF12
	KeyLeftShift    Key = pixel.KeyLeftShift
	KeyLeftControl  Key = pixel.KeyLeftControl
	KeyLeftAlt      Key = pixel.KeyLeftAlt
	KeyLeftSuper    Key = pixel.KeyLeftSuper
	KeyRightShift   Key = pixel.KeyRightShift
	KeyRightControl Key = pixel.KeyRightControl
	KeyRightAlt     Key = pixel.KeyRightAlt
	KeyRightSuper   Key = pixel.KeyRightSuper
	KeyMenu         Key = pixel.KeyMenu
)

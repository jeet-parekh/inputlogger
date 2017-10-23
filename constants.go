package inputlogger

const (
	_SUCCESS string = "The operation completed successfully."

	_WH_KEYBOARD_LL int32 = 13
	_WH_MOUSE_LL int32 = 14

	WM_KEYDOWN       uintptr = 0x0100
	WM_KEYUP         uintptr = 0x0101
	WM_SYSKEYDOWN    uintptr = 0x0104
	WM_SYSKEYUP      uintptr = 0x0105

	WM_MOUSEMOVE       uintptr = 0x0200
	WM_LBUTTONDOWN     uintptr = 0x0201
	WM_LBUTTONUP       uintptr = 0x0202
	WM_RBUTTONDOWN     uintptr = 0x0204
	WM_RBUTTONUP       uintptr = 0x0205
	WM_MOUSEWHEEL      uintptr = 0x020A
	WM_MOUSEHWHEEL     uintptr = 0x020E
)

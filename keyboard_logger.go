package inputlogger

import (
	"unsafe"
	"golang.org/x/sys/windows"
	"github.com/jeet-parekh/winapi"
)

type (
	// KeyboardEvent contains information about a low-level keyboard input event.
	KeyboardEvent struct {
		Flags uint32
		MessageID uint32
		ScanCode uint32
		Time uint32
		VirtualKeyCode uint32
	}

	// KeyboardLogger provides an interface to receive KeyboardEvents.
	KeyboardLogger struct {
		callbackFunction func(int32, uintptr, uintptr) uintptr
		hook uintptr
		messages chan *KeyboardEvent
		targetMessages map[uintptr]bool
	}
)

// NewKeyboardLogger creates and returns a KeyboardLogger.
// First parameter is the buffer size of the channel through which the messages would be sent.
// Second parameter is a MessageID -> bool map. Receive = true, Ignore = false.
func NewKeyboardLogger(bufferSize uint, targetMessages map[uintptr]bool) *KeyboardLogger {
	// Create a KeyboardLogger, assign it target messages and a callback function.
	logger := &KeyboardLogger {
		messages : make(chan *KeyboardEvent, bufferSize),
		targetMessages : make(map[uintptr]bool),
	}

	logger.targetMessages[WM_KEYDOWN] = targetMessages[WM_KEYDOWN]
	logger.targetMessages[WM_KEYUP] = targetMessages[WM_KEYUP]
	logger.targetMessages[WM_SYSKEYDOWN] = targetMessages[WM_SYSKEYDOWN]
	logger.targetMessages[WM_SYSKEYUP] = targetMessages[WM_SYSKEYUP]

	callbackFunction := func(nCode int32, wParam uintptr, lParam uintptr) uintptr {
		if logger.targetMessages[wParam] {
			llMsg := (*winapi.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
			logger.messages <- &KeyboardEvent {
				Flags : llMsg.Flags,
				MessageID : uint32(wParam),
				ScanCode : llMsg.ScanCode,
				Time : llMsg.Time,
				VirtualKeyCode : llMsg.VkCode,
			}
		}

		r1, err := winapi.CallNextHookEx(0, nCode, wParam, lParam)
		if err.Error() != _SUCCESS { panic(err.Error()) }
		return r1
	}
	logger.callbackFunction = callbackFunction

	return logger
}

// GetMessageChannel returns the channel through which the messages would be sent.
func (logger *KeyboardLogger) GetMessageChannel() <-chan *KeyboardEvent {
	return logger.messages
}

// Start the KeyboardLogger.
func (logger *KeyboardLogger) Start() {
	// Set the low-level keyboard hook, and start a goroutine which never terminates.
	var err error
	logger.hook, err = winapi.SetWindowsHookEx(_WH_KEYBOARD_LL, windows.NewCallback(logger.callbackFunction), 0, 0)
	if err.Error() != _SUCCESS { panic(err.Error()) }

	go func() {
		for {
			winapi.GetMessage(nil, 0, 0, 0)
		}
	}()
}

// Stop the KeyboardLogger and close the message channel.
func (logger *KeyboardLogger) Stop() {
	// Remove the low-level hook, and close the message channel.
	_, err := winapi.UnhookWindowsHookEx(logger.hook)
	if err.Error() != _SUCCESS { panic(err.Error()) }
	close(logger.messages)
}

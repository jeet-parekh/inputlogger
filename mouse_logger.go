package inputlogger

import (
	"unsafe"
	"golang.org/x/sys/windows"
	"github.com/jeet-parekh/winapi"
)

type (
	// MouseEvent contains information about a low-level mouse input event.
	MouseEvent struct {
		Flags uint32
		MessageID uint32
		MouseData uint32
		Time uint32
		CursorX int32
		CursorY int32
	}

	// MouseLogger provides an interface to receive MouseEvents.
	MouseLogger struct {
		callbackFunction func(int32, uintptr, uintptr) uintptr
		hook uintptr
		messages chan *MouseEvent
		targetMessages map[uintptr]bool
	}
)

// NewMouseLogger creates and returns a MouseLogger.
// First parameter is the buffer size of the channel through which the messages would be sent.
// Second parameter is a MessageID -> bool map. Receive = true, Ignore = false.
func NewMouseLogger(bufferSize uint, targetMessages map[uintptr]bool) *MouseLogger {
	// Create a MouseLogger, assign it target messages and a callback function.
	logger := &MouseLogger {
		messages : make(chan *MouseEvent, bufferSize),
		targetMessages : make(map[uintptr]bool),
	}

	logger.targetMessages[WM_LBUTTONDOWN] = targetMessages[WM_LBUTTONDOWN]
	logger.targetMessages[WM_LBUTTONUP] = targetMessages[WM_LBUTTONUP]
	logger.targetMessages[WM_RBUTTONDOWN] = targetMessages[WM_RBUTTONDOWN]
	logger.targetMessages[WM_RBUTTONUP] = targetMessages[WM_RBUTTONUP]
	logger.targetMessages[WM_MOUSEMOVE] = targetMessages[WM_MOUSEMOVE]
	logger.targetMessages[WM_MOUSEWHEEL] = targetMessages[WM_MOUSEWHEEL]
	logger.targetMessages[WM_MOUSEHWHEEL] = targetMessages[WM_MOUSEHWHEEL]

	callbackFunction := func(nCode int32, wParam uintptr, lParam uintptr) uintptr {
		if logger.targetMessages[wParam] {
			llMsg := (*winapi.MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
			logger.messages <- &MouseEvent {
				Flags : llMsg.Flags,
				MessageID : uint32(wParam),
				MouseData : llMsg.MouseData,
				Time : llMsg.Time,
				CursorX : llMsg.Pt.X,
				CursorY : llMsg.Pt.Y,
			}
		}

		r1, err := winapi.CallNextHookEx(0, nCode, wParam, lParam)
		if err.Error() != _SUCCESS { panic(err) }
		return r1
	}
	logger.callbackFunction = callbackFunction

	return logger
}

// GetMessageChannel returns the channel through which the messages would be sent.
func (logger *MouseLogger) GetMessageChannel() <-chan *MouseEvent {
	return logger.messages
}

// Start the MouseLogger.
func (logger *MouseLogger) Start() {
	// Set the low-level mouse hook, and start a goroutine which never terminates.
	var err error
	logger.hook, err = winapi.SetWindowsHookEx(_WH_MOUSE_LL, windows.NewCallback(logger.callbackFunction), 0, 0)
	if err.Error() != _SUCCESS { panic(err) }

	go func() {
		for {
			winapi.GetMessage(nil, 0, 0, 0)
		}
	}()
}

// Stop the MouseLogger and close the message channel.
func (logger *MouseLogger) Stop() {
	// Remove the low-level hook, and close the message channel.
	_, err := winapi.UnhookWindowsHookEx(logger.hook)
	if err.Error() != _SUCCESS { panic(err) }
	close(logger.messages)
}

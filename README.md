## Low-level input logger for Go (for Windows)

---

### About inputlogger

- Use inputlogger to capture low-level input events from keyboard or mouse (for Windows).

### Example

- This program would listen for 20 key-down input events followed by 20 mouse-move events.

```go
package main

import (
	"fmt"
	"github.com/jeet-parekh/inputlogger"
)

func main() {
	kl := inputlogger.NewKeyboardLogger(4, map[uintptr]bool { inputlogger.WM_KEYDOWN:true })
	kch := kl.GetMessageChannel()
	kl.Start()

	for i := 0; i < 20; i++ {
		fmt.Printf("%+v\n", <- kch)
	}
	kl.Stop()

	ml := inputlogger.NewMouseLogger(4, map[uintptr]bool { inputlogger.WM_MOUSEMOVE:true })
	mch := ml.GetMessageChannel()
	ml.Start()

	for i := 0; i < 20; i++ {
		fmt.Printf("%+v\n", <- mch)
	}
	ml.Stop()
}
```

---

### Notes

- Using the provided constants, you could specify the input events you would like to receive.
- If you want to flush the console input buffer at the end of the program, use [flushconin](https://github.com/jeet-parekh/flushconin).

---

package util

// PopMsg removes the first message in the message array.
func PopMsg(msgList [][]string) (head []string, tail [][]string) {
	head = msgList[0]
	tail = msgList[1:]
	return
}

// PopStr removes the first element in the array.
func PopStr(list []string) (head string, tail []string) {
	head = list[0]
	tail = list[1:]
	return
}

// Unwrap pops a frame off front of message and returns it as 'head', if the
// next frame is empty, it pops that empty frame. Return remaining frames of
// message as 'tail'.
func Unwrap(msg []string) (head string, tail []string) {
	head = msg[0]
	if len(msg) > 1 && msg[1] == "" {
		tail = msg[2:]
	} else {
		tail = msg[1:]
	}
	return
}

// Contains returns `true` if `value` was found in `arr`, `false` otherwise.
func Contains(arr []string, value string) bool {
	for _, a := range arr {
		if a == value {
			return true
		}
	}
	return false
}

// Keys returns a slice of strings from the given map `kv`.
func Keys(kv map[string]string) []string {
	list := make([]string, 0, len(kv))
	for k := range kv {
		list = append(list, k)
	}
	return list
}

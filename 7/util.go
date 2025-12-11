package main

// unescapeData will return a new slice of the data after removing escaped characters from the input data
func unescapeData(b []byte) []byte {
	out := make([]byte, 0, len(b))

	for i := 0; i < len(b); i++ {
		if b[i] == '\\' && i+1 < len(b) {
			// skip escape marker, take the escaped byte
			i++
		}
		out = append(out, b[i])
	}
	return out
}

// escapeData will return a new slice of the data after escaping / and \ characters from the input data
// Where the DATA contains forward slash ("/") or backslash ("\") characters, the sender must escape the slashes by prepending them each with a single backslash character
func escapeData(b []byte) []byte {
	// First, count how many extra bytes we need
	var extra int
	for _, c := range b {
		if c == '/' || c == '\\' {
			extra++
		}
	}

	out := make([]byte, 0, len(b)+extra)
	for _, c := range b {
		if c == '/' || c == '\\' {
			out = append(out, '\\') // insert escape
		}
		out = append(out, c)
	}
	return out
}

// reverseBytes will return a new reversed byte slice of the input data
func reverseBytes(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)

	// assume newline
	end := len(out) - 1
	if out[end] != '\n' {
		end = len(out)
	}

	for i, j := 0, end-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

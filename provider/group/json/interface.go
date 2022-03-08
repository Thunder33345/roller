package json

import "io"

//truncateSeeker is an io.ReadWriter that can be seeked and truncated
//internally used for compat on JSON.Save where Truncate(0) and Seek(0,0) will be called before writing
type truncateSeeker interface {
	Truncate(size int64) error
	io.Seeker
	io.ReadWriter
}

//reseter is an io.ReadWriter that can be reset before being written to, used for JSON.Save
type reseter interface {
	Reset()
	io.ReadWriter
}

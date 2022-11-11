package hardlinkfilemap

import (
	"errors"
	"os"
	"strconv"
	"syscall"
)

func FileIdentifier(fi os.FileInfo) (string, error) {
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return "", errors.New("failed to get file identifier")
	}

	return strconv.FormatUint(sys.Dev, 10) + "|" + strconv.FormatUint(sys.Ino, 10), nil
}

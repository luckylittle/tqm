package hardlinkfilemap

import (
	"os"
	"reflect"
	"strconv"
)

func FileIdentifier(fi os.FileInfo) (string, error) {
	// This is extreme hackies to get os to load the vol, idxhi and idxlo fields
	ok := os.SameFile(fi, fi)
	if !ok {
		return "", errors.New("error while getting file identifier")
	}

	v := reflect.Indirect(reflect.ValueOf(fi))
	vol := v.FieldByName("vol").Uint()
	idxhi := v.FieldByName("idxhi").Uint()
	idxlo := v.FieldByName("idxlo").Uint()

	return strconv.FormatUint(vol, 10) + "|" + strconv.FormatUint(idxhi, 10) + "|" + strconv.FormatUint(idxlo, 10), nil
}

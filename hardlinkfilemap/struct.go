package hardlinkfilemap

import (
	"github.com/sirupsen/logrus"
)

type HardlinkFileMap struct {
	// hardlinkFileMap map[string]map[string]config.Torrent
	hardlinkFileMap    map[string][]string
	log                *logrus.Entry
	torrentPathMapping map[string]string
}

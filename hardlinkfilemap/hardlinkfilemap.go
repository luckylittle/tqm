package hardlinkfilemap

import (
	"os"
	"strings"

	"github.com/l3uddz/tqm/config"
	"github.com/l3uddz/tqm/logger"
	"github.com/l3uddz/tqm/sliceutils"
)

func New(torrents map[string]config.Torrent, torrentPathMapping map[string]string) *HardlinkFileMap {
	tfm := &HardlinkFileMap{
		hardlinkFileMap:    make(map[string][]string),
		log:                logger.GetLogger("hardlinkfilemap"),
		torrentPathMapping: torrentPathMapping,
	}

	for _, torrent := range torrents {
		tfm.AddByTorrent(torrent)
	}

	return tfm
}

func (t *HardlinkFileMap) ConsiderPathMapping(path string) string {
	for mapFrom, mapTo := range t.torrentPathMapping {
		if strings.HasPrefix(path, mapFrom) {
			return strings.Replace(path, mapFrom, mapTo, 1)
		}
	}

	return path
}

func (t *HardlinkFileMap) FileIdentifierByPath(path string) (string, bool) {
	stat, err1 := os.Stat(path)
	if err1 != nil {
		t.log.Warnf("Failed to stat file: %s - %s", path, err1)
		return "", false
	}

	id, err2 := FileIdentifier(stat)
	if err2 != nil {
		t.log.Warnf("Failed to get file identifier: %s - %s", path, err2)
		return "", false
	}

	return id, true
}

func (t *HardlinkFileMap) AddByTorrent(torrent config.Torrent) {
	for _, f := range torrent.Files {
		f = t.ConsiderPathMapping(f)

		id, ok := t.FileIdentifierByPath(f)

		if !ok {
			continue
		}

		if _, exists := t.hardlinkFileMap[id]; exists {
			// file id already associated with other paths
			t.hardlinkFileMap[id] = append(t.hardlinkFileMap[id], f)
			continue
		}

		// file id has not been seen before, create id entry
		t.hardlinkFileMap[id] = []string{f}
	}
}

func (t *HardlinkFileMap) RemoveByTorrent(torrent config.Torrent) {
	for _, f := range torrent.Files {
		f = t.ConsiderPathMapping(f)

		id, ok := t.FileIdentifierByPath(f)

		if !ok {
			continue
		}

		if _, exists := t.hardlinkFileMap[id]; exists {
			// remove this path from the id entry
			i := sliceutils.IndexOfString(t.hardlinkFileMap[id], f)
			if i != -1 {
				t.hardlinkFileMap[id] = sliceutils.FastDelete(t.hardlinkFileMap[id], i)
			}

			// remove id entry if no more paths
			if len(t.hardlinkFileMap[id]) == 0 {
				delete(t.hardlinkFileMap, id)
			}

			continue
		}
	}
}

func (t *HardlinkFileMap) IsTorrentUnique(torrent config.Torrent) bool {
	for _, f := range torrent.Files {
		f = t.ConsiderPathMapping(f)

		id, ok := t.FileIdentifierByPath(f)

		if !ok {
			return false
		}

		t.log.Infof("File: %s - ID: %s", f, id)
		// preview the file id entry
		t.log.Infof("File ID Entry: %v", t.hardlinkFileMap[id])

		if paths, exists := t.hardlinkFileMap[id]; exists && len(paths) > 1 {
			return false
		}
	}

	return true
}

func (t *HardlinkFileMap) NoInstances(torrent config.Torrent) bool {
	for _, f := range torrent.Files {
		f = t.ConsiderPathMapping(f)

		id, ok := t.FileIdentifierByPath(f)

		if !ok {
			return false
		}

		if paths, exists := t.hardlinkFileMap[id]; exists && len(paths) != 0 {
			return false
		}
	}

	return true
}

func (t *HardlinkFileMap) Length() int {
	return len(t.hardlinkFileMap)
}

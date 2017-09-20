package movetorrent

import (
	"errors"

	"github.com/tubbebubbe/transmission"
)

//MoveShowHandle stores the config and details for moving files
type MoveTorrentHandle struct {
	client transmission.TransmissionClient
}

//NewMoveTorrentHandle creates a new Handle to manage moving shows
func NewMoveTorrentHandle(homeTvDirectory string) *MoveTorrentHandle {
	client := transmission.New("http://127.0.0.1:9091", "", "")
	return &MoveTorrentHandle{client: client}
}

func (m *MoveTorrentHandle) MoveTorrent() {

	cmd, err := newGetTorrentsCmd()
	m.client.ExecuteCommand(cmd)
}

func (m *MoveTorrentHandle) findByName(name string) (int, error) {
	torrents, err := m.client.GetTorrents()
	if err != nil {
		return -1, err
	}

	for _, item := range torrents {
		if item.Name == name {
			return item.ID, nil
		}
	}

	return -1, errors.New("Torrent not found")
}

func newGetTorrentsCmd() (*transmission.Command, error) {
	cmd := &transmission.Command{}

	cmd.Method = "torrent-get"
	cmd.Arguments.

	return cmd, nil
}

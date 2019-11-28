package rpc

import (
	"fmt"
	"github.com/hekmon/transmissionrpc"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

type TransmissionConnection struct {
	*transmissionrpc.Client
}

func NewTransmissionConnection(hostname string, port int, username, password string, useHTTPS bool) (*TransmissionConnection, error) {
	c, err := transmissionrpc.New(hostname, username, password,
		&transmissionrpc.AdvancedConfig{
			HTTPS: useHTTPS,
			Port: uint16(port),
		})
	tc := TransmissionConnection{c}
	return &tc, err
}

func (tc *TransmissionConnection) GetTorrentList(onlyActive bool) ([]*transmissionrpc.Torrent, error) {
	return tc.TorrentGetAll()
}

func (tc *TransmissionConnection) AddTorrent(url string) (*transmissionrpc.Torrent, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "torrent-*.torrent")
	if err != nil {return nil, err}
	defer os.Remove(tmpFile.Name())
	log.Info(fmt.Sprintf("Created temporary torrent file at %s.\n", tmpFile.Name()))

	resp, err := http.Get(url)
	if err != nil {return nil, err}
	defer resp.Body.Close()

	body := make([]byte,0)
	_, err = resp.Body.Read(body)
	if err != nil {return nil, err}

	err = ioutil.WriteFile(tmpFile.Name(), body, os.FileMode(int(0644)))
	if err != nil {return nil, err}

	torrent, err := tc.TorrentAddFile(tmpFile.Name())
	return torrent, err

}

func (tc *TransmissionConnection) RemoveTorrent(id int, deleteData bool) error {
	ids := []int64{int64(id)}
	rp := transmissionrpc.TorrentRemovePayload{
		IDs:             ids,
		DeleteLocalData: deleteData,
	}

	err := tc.TorrentRemove(&rp)
		if err == nil {
			log.Info(fmt.Sprintf("Removed torrent ID %d. Delete data: %v\n", id, deleteData))
		}
	return err
}

func (tc *TransmissionConnection) PauseTorrent(id int) error {
	ids := []int64{int64(id)}
	err := tc.TorrentStopIDs(ids)
	if err == nil {
		log.Info(fmt.Sprintf("Stopped torrent ID %d.\n", id))
	}
	return err
}

func (tc *TransmissionConnection) StartTorrent(id int) error {
	ids := []int64{int64(id)}
	err := tc.TorrentStartIDs(ids)
	if err == nil {
		log.Info(fmt.Sprintf("Started torrent ID %d.\n", id))
	}
	return err
}

func (tc *TransmissionConnection) IsConnected() (bool, int64, int64, error) {
	ok, serverVersion, serverMinimumVersion, err := tc.RPCVersion()
	return ok, serverVersion, serverMinimumVersion, err
}
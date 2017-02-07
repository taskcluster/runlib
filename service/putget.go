package service

import (
	"os"

	"github.com/juju/errors"
	"github.com/taskcluster/runlib/contester_proto"
	"github.com/taskcluster/runlib/tools"
)

func (s *Contester) Put(request *contester_proto.FileBlob, response *contester_proto.FileStat) error {
	resolved, sandbox, err := resolvePath(s.Sandboxes, *request.Name, true)
	if err != nil {
		return err
	}

	if sandbox != nil {
		sandbox.Mutex.Lock()
		defer sandbox.Mutex.Unlock()
	}

	var destination *os.File

	for {
		destination, err = os.Create(resolved)
		loop, err := OnOsCreateError(err)

		if err != nil {
			return errors.Annotatef(err, "os.Create(%q)", resolved)
		}
		if !loop {
			break
		}
	}
	data, err := request.Data.Bytes()
	if err != nil {
		return err
	}
	_, err = destination.Write(data)
	destination.Close()
	if err != nil {
		return errors.Annotate(err, "destination.Write")
	}
	if sandbox != nil {
		return sandbox.Own(resolved)
	}

	stat, err := tools.StatFile(resolved, true)
	if err != nil {
		return err
	}

	*response = *stat
	return nil
}

func (s *Contester) Get(request *contester_proto.GetRequest, response *contester_proto.FileBlob) error {
	resolved, sandbox, err := resolvePath(s.Sandboxes, *request.Name, false)
	if err != nil {
		return err
	}

	if sandbox != nil {
		sandbox.Mutex.RLock()
		defer sandbox.Mutex.RUnlock()
	}

	source, err := os.Open(resolved)
	if err != nil {
		return errors.Annotatef(err, "os.Open(%q)", resolved)
	}
	defer source.Close()

	response.Name = &resolved
	response.Data, err = contester_proto.BlobFromStream(source)
	return err
}

package core

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/jhump/protoreflect/dynamic"

	"github.com/Aliceikkk/ViaGenshin/internal/mapper"
)

func (s *Session) ConvertPacket(from, to mapper.Protocol, fromCmd uint16, p []byte) ([]byte, error) {
	name := s.mapping.CommandNameMap[from][fromCmd]
	fromDesc := s.mapping.MessageDescMap[from][name]
	if fromDesc == nil {
		return p, fmt.Errorf("unknown from message %s in %s", name, to)
	}
	fromPacket := dynamic.NewMessage(fromDesc)
	if err := fromPacket.Unmarshal(p); err != nil {
		return p, err
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(fromPacket); err != nil {
		return p, err
	}
	fromGob := buf.Bytes()
	toGob, err := s.HandlePacket(from, to, name, fromGob)
	if err != nil {
		return p, err
	}
	buf.Reset()
	buf.Write(toGob)
	dec := gob.NewDecoder(&buf)
	toPacket := dynamic.NewMessage(s.mapping.MessageDescMap[to][name])
	if err := dec.Decode(toPacket); err != nil {
		return p, err
	}
	toBytes, err := toPacket.Marshal()
	if err != nil {
		return p, err
	}
	return toBytes, nil
}

func (s *Session) ConvertPacketByName(from, to mapper.Protocol, name string, p []byte) ([]byte, error) {
	fromDesc := s.mapping.MessageDescMap[from][name]
	if fromDesc == nil {
		return p, fmt.Errorf("unknown from message %s in %s", name, to)
	}
	fromPacket := dynamic.NewMessage(fromDesc)
	if err := fromPacket.Unmarshal(p); err != nil {
		return p, err
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(fromPacket); err != nil {
		return p, err
	}
	fromGob := buf.Bytes()
	toGob, err := s.HandlePacket(from, to, name, fromGob)
	if err != nil {
		return p, err
	}
	buf.Reset()
	buf.Write(toGob)
	dec := gob.NewDecoder(&buf)
	toPacket := dynamic.NewMessage(s.mapping.MessageDescMap[to][name])
	if err := dec.Decode(toPacket); err != nil {
		return p, err
	}
	toBytes, err := toPacket.Marshal()
	if err != nil {
		return p, err
	}
	return toBytes, nil
}

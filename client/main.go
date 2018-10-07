package main

import (
	"context"
	"encoding/binary"
	"io"
	"os"

	encodepb "github.com/fvdveen/mu2-proto/go/proto/encode"
	"github.com/micro/cli"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/sirupsen/logrus"
)

func main() {
	var loc, url, file string

	srv := grpc.NewService(
		micro.Flags(
			cli.StringFlag{
				Destination: &loc,
				Name:        "location",
				Usage:       "Search service location",
				Value:       "mu2.service.encode",
			},
			cli.StringFlag{
				Destination: &url,
				Name:        "encode-url",
				Usage:       "The url to encode",
				Value:       "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			},
			cli.StringFlag{
				Destination: &file,
				Name:        "encode-file",
				Usage:       "The output file",
				Value:       "out.opus",
			},
		),
	)

	srv.Init()

	cl := encodepb.NewEncodeService(loc, srv.Client())

	s, err := cl.Encode(context.Background(), &encodepb.EncodeRequest{
		Url: url,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	var ops [][]byte

	for {
		res, err := s.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			logrus.Error(err)
			return
		}

		ops = append(ops, res.Opus)
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		logrus.Error(err)
	}

	for _, opus := range ops {
		// write header
		opuslen := int16(len(opus))
		err = binary.Write(f, binary.LittleEndian, &opuslen)
		if err != nil {
			logrus.Error(err)
			return
		}

		// write opus data to f
		err = binary.Write(f, binary.LittleEndian, &opus)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	if err = f.Close(); err != nil {
		logrus.Error(err)
		return
	}
}

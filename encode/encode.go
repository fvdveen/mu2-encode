package encode

import (
	"context"
	"fmt"
	"io"

	"github.com/fvdveen/mu2-proto/go/proto/encode"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

// NewService creates a new EncodeServiceHandler
func NewService() encode.EncodeServiceHandler {
	return &service{}
}

type service struct {
}

func (s *service) Encode(ctx context.Context, req *encode.EncodeRequest, stream encode.EncodeService_EncodeStream) error {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	videoInfo, err := ytdl.GetVideoInfo(req.Url)
	if err != nil {
		return fmt.Errorf("get video info: %v", err)
	}

	format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	downloadURL, err := videoInfo.GetDownloadURL(format)
	if err != nil {
		return fmt.Errorf("get download URL: %v", err)
	}

	es, err := dca.EncodeFile(downloadURL.String(), options)
	if err != nil {
		return fmt.Errorf("create encode session: %v", err)
	}
	defer es.Cleanup()

	for {
		f, err := es.OpusFrame()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("get opus frame: %v", err)
		}

		err = stream.Send(&encode.EncodeResponse{
			Opus: f,
			Url:  req.Url,
		})
		if err != nil {
			return fmt.Errorf("send opus: %v", err)
		}
	}

	if err := stream.Close(); err != nil {
		return fmt.Errorf("close stream: %v", err)
	}

	return nil
}

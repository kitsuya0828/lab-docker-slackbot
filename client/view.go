package main

import (
	"fmt"
	"strings"

	pb "github.com/Kitsuya0828/lab-docker-slackbot/proto/stat"
	"github.com/dustin/go-humanize"
	"github.com/slack-go/slack"
)

const (
	barLength = 10
)

func HomeTabView(hs []string, fs []*pb.FsStat, ds []*pb.DockerStat) slack.HomeTabViewRequest {
	fsMsgs := make([]string, 0, len(fs))
	imgMsgs := make([]string, 0, len(ds))
	ctnMsgs := make([]string, 0, len(ds))
	for i := range hs {
		progressBar := createProgressBar(fs[i].Used, fs[i].Total, barLength)
		fsMsgs = append(fsMsgs, fmt.Sprintf(
			"%s\n\tUsed: *%s* / *%s* (%.f%%)",
			progressBar, humanize.Bytes(fs[i].Used), humanize.Bytes(fs[i].Total),
			float64(fs[i].Used)/float64(fs[i].Total)*100,
		))
		imgMsgs = append(imgMsgs, fmt.Sprintf(
			"\tActive: *%d*/%d (%.f%%) images, Reclaimable: *%s*/%s (%.f%%)",
			ds[i].Images.Active, ds[i].Images.TotalCount,
			float64(ds[i].Images.Active)/float64(ds[i].Images.TotalCount)*100,
			humanize.Bytes(uint64(ds[i].Images.Reclaimable)), humanize.Bytes(uint64(ds[i].Images.Size)),
			ds[i].Images.Reclaimable/ds[i].Images.Size*100,
		))
		ctnMsgs = append(ctnMsgs, fmt.Sprintf(
			"\tActive: *%d*/%d (%.f%%) containers, Reclaimable: *%s*/%s (%.f%%)",
			ds[i].Containers.Active, ds[i].Containers.TotalCount,
			float64(ds[i].Containers.Active)/float64(ds[i].Containers.TotalCount)*100,
			humanize.Bytes(uint64(ds[i].Containers.Reclaimable)), humanize.Bytes(uint64(ds[i].Containers.Size)),
			ds[i].Containers.Reclaimable/ds[i].Containers.Size*100,
		))
	}

	blockSet := make([]slack.Block, 0, len(fs)*5)
	for i, hostname := range hs {
		blockSet = append(blockSet, slack.NewHeaderBlock(
			slack.NewTextBlockObject(slack.PlainTextType, hostname, false, false),
		))
		blockSet = append(blockSet, slack.NewSectionBlock(
			slack.NewTextBlockObject(
				slack.MarkdownType,
				fmt.Sprintf(":file_folder: *File System* `/`\n%s", fsMsgs[i]),
				false,
				false,
			),
			nil,
			nil,
		))
		blockSet = append(blockSet, slack.NewSectionBlock(
			slack.NewTextBlockObject(
				slack.MarkdownType,
				fmt.Sprintf(":whale: *Docker Images*\n%s", imgMsgs[i]),
				false,
				false,
			),
			nil,
			nil,
		))
		blockSet = append(blockSet, slack.NewSectionBlock(
			slack.NewTextBlockObject(
				slack.MarkdownType,
				fmt.Sprintf(":whale2: *Docker Containers*\n%s", ctnMsgs[i]),
				false,
				false,
			),
			nil,
			nil,
		))
		if i != len(hs)-1 {
			blockSet = append(blockSet, slack.NewDividerBlock())
		}
	}

	return slack.HomeTabViewRequest{
		Type: slack.VTHomeTab,
		Blocks: slack.Blocks{
			BlockSet: blockSet,
		},
	}
}

func createProgressBar(used uint64, total uint64, barLength int) string {
	ratio := float32(used) / float32(total)
	filledLength := int(float32(barLength) * ratio)
	emptyLength := barLength - filledLength

	var filledBar string
	switch {
	case ratio < 0.5:
		filledBar = strings.Repeat(":large_green_square:", filledLength)
	case ratio < 0.8:
		filledBar = strings.Repeat(":large_yellow_square:", filledLength)
	default:
		filledBar = strings.Repeat(":large_red_square:", filledLength)
	}

	emptyBar := strings.Repeat(":white_large_square:", emptyLength)

	return fmt.Sprintf("%s%s", filledBar, emptyBar)
}

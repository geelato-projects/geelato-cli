package progress

import (
	"sync"

	"github.com/cheggaaa/pb/v3"
)

type Progress struct {
	mu      sync.Mutex
	bar     *pb.ProgressBar
	total   int64
	current int64
	desc    string
}

func NewBar(total int64, description string) *Progress {
	if total == 0 {
		return nil
	}

	p := &Progress{
		total: total,
		desc:  description,
	}

	p.bar = pb.New64(total)
	p.bar.SetWidth(80)
	p.bar.SetTemplateString(`{{string .Name "` + description + `"}} {{counters .}} {{bar .}} {{percent .}}`)

	return p
}

func (p *Progress) Start() {
	if p == nil || p.bar == nil {
		return
	}
	p.bar.Start()
}

func (p *Progress) Update(percent int) {
	if p == nil || p.bar == nil {
		return
	}
	target := int64(float64(p.total) * float64(percent) / 100.0)
	p.bar.SetCurrent(target)
}

func (p *Progress) Increment(n int64) {
	if p == nil || p.bar == nil {
		return
	}
	p.bar.Add(int(n))
}

func (p *Progress) Stop() {
	if p == nil || p.bar == nil {
		return
	}
	p.bar.Finish()
}

package plugin

type matcher struct {
	r       runner
	matches map[string]bool
}

func newMatcher(matches []string, r runner) *matcher {
	this := &matcher{r: r, matches: make(map[string]bool, len(matches))}
	for _, m := range matches {
		this.matches[m] = true
	}
	return this
}

func (this *matcher) InChan() chan *PipelinePack {
	return this.r.InChan()
}

func (this *matcher) Match(pack *PipelinePack) bool {
	return this.matches[pack.Ident]
}

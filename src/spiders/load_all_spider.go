package spiders

import (
	"ants/spiders"
)

func LoadAllSpiders() map[string]*spiders.Spider {
	spiderMap := make(map[string]*spiders.Spider)
	deadLoopTest := MakeDeadLoopSpider()
	spiderMap[deadLoopTest.Name] = deadLoopTest
	dumpTestSpider := MakeDumpTestSpider()
	spiderMap[dumpTestSpider.Name] = dumpTestSpider
	return spiderMap
}

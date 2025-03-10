package do

import "book-management/internal/pkg/common"

type BorrowStatistics struct {
	ChildrenStoryNum    int
	ScienceKnowledgeNum int
	ArtEnlightenmentNum int
}

func (s *BorrowStatistics) ToMap() map[string]int {
	return map[string]int{
		common.ChildrenStory:    s.ChildrenStoryNum,
		common.ScienceKnowledge: s.ScienceKnowledgeNum,
		common.ArtEnlightenment: s.ArtEnlightenmentNum,
	}
}

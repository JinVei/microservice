package service

type ReplyCommentSvcConfig struct {
	CacheDura      string `json:"cacheDura"`      // comment comment cache duration
	IndexCacheDura string `json:"indexCacheDura"` // comment index cache duration
}

func defaultReplyCommentSvcConfig() ReplyCommentSvcConfig {
	return ReplyCommentSvcConfig{
		CacheDura:      "160h", // 7 day
		IndexCacheDura: "336h", // 14 day
	}
}

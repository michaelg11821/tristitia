package main

type NarouCharCountResp struct {
	AllCount  int   `json:"allcount,omitempty"`
	CharCount int64 `json:"length,omitempty"`
}

func getCharCountEndpoint(ncode string) string {
	return "https://api.syosetu.com/novelapi/api/?out=json&of=l&ncode=" + ncode
}

type NarouChapterNumResp struct {
	AllCount      int `json:"allcount,omitempty"`
	LatestChapter int `json:"general_all_no,omitempty"`
}

func getLatestChapterNumEndpoint(ncode string) string {
	return "https://api.syosetu.com/novelapi/api/?out=json&of=ga&ncode=" + ncode
}

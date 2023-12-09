package container

type Statistics struct {
	Empty              int64 `json:"empty"`
	Partial            int64 `json:"partial"`
	Full               int64 `json:"full"`
	WaitingForApproval int64 `json:"waiting_for_approval"`
}

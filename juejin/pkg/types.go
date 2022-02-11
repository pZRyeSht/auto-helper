package pkg

type WechatMarkdown struct {
	MsgType  string    `json:"msgtype"`
	Markdown *Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

type JueJinSignInResp struct {
	ErrNo  int        `json:"err_no"`
	ErrMsg string     `json:"err_msg"`
	Data   SignInResp `json:"data"`
}

type SignInResp struct {
	IncrPoint int64 `json:"incr_point"`
	SumPoint  int64 `json:"sum_point"`
}

type JueJinLotteryResp struct {
	ErrNo  int         `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Data   LotteryResp `json:"data"`
}

type LotteryResp struct {
	Id           int    `json:"id"`
	LotteryId    string `json:"lottery_id"`
	LotteryName  string `json:"lottery_name"`
	LotteryType  int    `json:"lottery_type"`
	LotteryImage string `json:"lottery_image"`
	LotteryDesc  string `json:"lottery_desc"`
	HistoryId    string `json:"history_id"`
}

type DefineEvent struct {
	// test event define
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}
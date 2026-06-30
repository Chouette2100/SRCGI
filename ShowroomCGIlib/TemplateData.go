package ShowroomCGIlib

type GraphPageData struct {
	Filename string
	Eventid  string
	Maxpoint string
	Gscale   string
}

type SimpleMessagePageData struct {
	Function string
	Comment  string
}

type ErrorPageData struct {
	Msg001   string
	Msg002   string
	ReturnTo string
	Eventid  string
	Maxpoint string
	Gscale   string
}

type NewUserPageData struct {
	Event_ID   string
	Event_name string
	Period     string
	Roomid     string
	Roomname   string
	Longname   string
	Shortname  string
	Roomurlkey string
	Genre      string
	Rank       string
	Nrank      string
	Prank      string
	Level      string
	Followers  string
	Fans       string
	Fans_lst   string
	Submit     string
	Label      string
	Msg1       string
	Msg2       string
	Msg2color  string
}

type NewEventPageData struct {
	Eventid   string
	Eventname string
	Period    string
	Noroom    string
	Msgcolor  string
	Stm       string
	Sts       string
	Maxcmap   string
	Msg       string
	Submit    string
}

type EditUserPageData struct {
	Eventid   string
	Eventname string
	Period    string
	Maxpoint  string
	Gscale    string
}

type ListLastPageData struct {
	Eventid         string
	Ieventid        string
	Userno          string
	Detail          string
	Isover          string
	Limit           string
	Page            int
	Maxrooms        int
	NoRooms         int
	Roomid          int
	Scorelist       []CurrentScore
	UpdateTime      string
	NextTime        string
	ReloadTime      string
	SecondsToReload string
	EventName       string
	Period          string
	Maxpoint        string
	Gscale          string
}

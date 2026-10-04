package ShowroomCGIlib

import (
	"html/template"
	"net/http"
)

type GraphPageData struct {
	Filename         string
	Eventid          string
	Maxpoint         string
	Gscale           string
	TurnstileSiteKey string
	TurnstileError   string
	RequestID        string
}

func (h *GraphPageData) SetTurnstileInfo(siteKey string, errorMsg string) {
	h.TurnstileSiteKey = siteKey
	h.TurnstileError = errorMsg
}

func (h *GraphPageData) GetTemplatePath() string {
	return "templates/graph-total.gtpl"
}

func (h *GraphPageData) GetTemplateName() string {
	return "graph-total.gtpl"
}

func (h *GraphPageData) GetFuncMap() *template.FuncMap {
	return nil
}

type TurnstileChallengePageData struct {
	TurnstileSiteKey      string
	TurnstileError        string
	RequestID             string
	ChallengeAction       string
	ChallengeHiddenFields map[string]string
}

func (h *TurnstileChallengePageData) SetTurnstileInfo(siteKey string, errorMsg string) {
	h.TurnstileSiteKey = siteKey
	h.TurnstileError = errorMsg
}

func (h *TurnstileChallengePageData) GetTemplatePath() string {
	return "templates/turnstilechallenge.gtpl"
}

func (h *TurnstileChallengePageData) GetTemplateName() string {
	return "turnstilechallenge"
}

func (h *TurnstileChallengePageData) GetFuncMap() *template.FuncMap {
	return nil
}

func BuildTurnstileChallengePageData(r *http.Request) *TurnstileChallengePageData {
	page := &TurnstileChallengePageData{
		ChallengeAction:       r.URL.Path,
		ChallengeHiddenFields: map[string]string{},
	}
	if v := r.Context().Value("requestid"); v != nil {
		page.RequestID = v.(string)
	}
	if err := r.ParseForm(); err == nil {
		for key, values := range r.Form {
			if key == "cf-turnstile-response" || key == "requestid" {
				continue
			}
			if len(values) > 0 {
				page.ChallengeHiddenFields[key] = values[0]
			}
		}
	}
	return page
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
	Maxpoint   string
	Gscale     string
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
	Eventid          string
	Eventname        string
	Period           string
	Maxpoint         string
	Gscale           string
	TurnstileSiteKey string
	TurnstileError   string
	RequestID        string
}

func (h *EditUserPageData) SetTurnstileInfo(siteKey string, errorMsg string) {
	h.TurnstileSiteKey = siteKey
	h.TurnstileError = errorMsg
}

func (h *EditUserPageData) GetTemplatePath() string {
	return "templates/edit-user1.gtpl"
}

func (h *EditUserPageData) GetTemplateName() string {
	return "edit-user1.gtpl"
}

func (h *EditUserPageData) GetFuncMap() *template.FuncMap {
	return nil
}

type ListLastPageData struct {
	Eventid          string
	Ieventid         string
	Userno           string
	Detail           string
	Isover           string
	Limit            string
	Page             int
	Maxrooms         int
	NoRooms          int
	Roomid           int
	Scorelist        []CurrentScore
	UpdateTime       string
	NextTime         string
	ReloadTime       string
	SecondsToReload  string
	EventName        string
	Period           string
	Maxpoint         string
	Gscale           string
	TurnstileSiteKey string
	TurnstileError   string
	RequestID        string
}

func (h *ListLastPageData) SetTurnstileInfo(siteKey string, errorMsg string) {
	h.TurnstileSiteKey = siteKey
	h.TurnstileError = errorMsg
}

func (h *ListLastPageData) GetTemplatePath() string {
	return "templates/list-last.gtpl"
}

func (h *ListLastPageData) GetTemplateName() string {
	return "list-last.gtpl"
}

func (h *ListLastPageData) GetFuncMap() *template.FuncMap {
	return nil
}

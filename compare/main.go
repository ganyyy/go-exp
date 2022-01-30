package main

type FinalRank struct {
	DistinctId string `json:"#distinct_id"`
	Type       string `json:"#type"`
	Time       string `json:"#time"`
	EventName  string `json:"#event_name"`
	Properties struct {
		Lib              string `json:"#lib"`
		LibVersion       string `json:"#lib_version"`
		C3V3Group        int    `json:"C3V3Group"`
		C3V3MyCaptainAID string `json:"C3V3MyCaptainAID"`
		C3V3MyMemberAID  string `json:"C3V3MyMemberAID"`
		C3V3MyTeamID     string `json:"C3V3MyTeamID"`
		C3V3MyTeamName   string `json:"C3V3MyTeamName"`
		C3V3MyTeamRank   int    `json:"C3V3MyTeamRank"`
		CnTime           string `json:"cn_time"`
		Logtime          int64  `json:"logtime"`
		ServiceIp        string `json:"service_ip"`
		ServiceName      string `json:"service_name"`
		Timestamp        string `json:"timestamp"`
		TypeName         string `json:"type_name"`
	} `json:"properties"`
}

var playerList = map[string]struct{}{}
var finalList = map[string]struct{}{}

func main() {

}

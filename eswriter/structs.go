package eswriter

type IndexInfo struct {
	Entry   interface{}
	Index   string
	TimeKey string
	Type    string
}

//--------------------------------------------

type Process struct {
	DD    string          `json:"dd"`
	RD    string          `json:"rd"`
	UD    string          `json:"ud"`
	NI    string          `json:"ni"`
	NM    string          `json:"nm"`
	OD    string          `json:"od"`
	MC    string          `json:"mc"`
	SC    string          `json:"sc"`
	MI    string          `json:"mi"`
	TI    string          `json:"ti"`
	OI    string          `json:"oi"`
	Jdata []ProcessDetail `json:"jdata"`
}

type ProcessDetail struct {
	Name          string `json:"name"`
	PID           string
	Md5           string `json:"md5"`
	CreateTime    string
	ActiveTime    string
	DisableTime   string
	EndTime       string
	TopActiveTime string `json:"topActiveTime"`
	CompanyName   string `json:"companyName"`
	DesCrib       string `json:"desCrib"`
	Version       string
	Signed        string `json:"signed"`
	SignedName    string `json:"signedName"`
	Path          string `json:"path"`
}

type Net struct {
	DD    string      `json:"dd"`
	RD    string      `json:"rd"`
	UD    string      `json:"ud"`
	NI    string      `json:"ni"`
	NM    string      `json:"nm"`
	OD    string      `json:"od"`
	MC    string      `json:"mc"`
	SC    string      `json:"sc"`
	MI    string      `json:"mi"`
	TI    string      `json:"ti"`
	OI    string      `json:"oi"`
	Jdata []NetDetail `json:"jdata"`
}

type NetDetail struct {
	SrcMac        string `json:"srcMac"`
	Mac           string
	Srcip         string `json:"srcip"`
	Ip            string `json:"ip"`
	Time          string `json:"time"`
	Cmd           string `json:"cmd"`
	Uri           string `json:"uri"`
	Host          string `json:"host"`
	Referer       string `json:"referer"`
	ContentLength string `json:"content-length"`
	AccessTime    string `json:"access-time"`
	IsConnect     string
}

type Copy struct {
	DD    string       `json:"dd"`
	RD    string       `json:"rd"`
	UD    string       `json:"ud"`
	NI    string       `json:"ni"`
	NM    string       `json:"nm"`
	OD    string       `json:"od"`
	MC    string       `json:"mc"`
	SC    string       `json:"sc"`
	MI    string       `json:"mi"`
	TI    string       `json:"ti"`
	OI    string       `json:"oi"`
	Jdata []CopyDetail `json:"jdata"`
}

type CopyDetail struct {
	Time     string
	Mode     string
	DataType string
	EXE      string
	WinF     string
}

type StartShutdown struct {
	DD    string                `json:"dd"`
	RD    string                `json:"rd"`
	UD    string                `json:"ud"`
	NI    string                `json:"ni"`
	NM    string                `json:"nm"`
	OD    string                `json:"od"`
	MC    string                `json:"mc"`
	SC    string                `json:"sc"`
	MI    string                `json:"mi"`
	TI    string                `json:"ti"`
	OI    string                `json:"oi"`
	Jdata []StartShutdownDetail `json:"jdata"`
}

type StartShutdownDetail struct {
	Time string `json:"time"`
	Cmd  string `json:"cmd"`
}

package web

type Credentials struct {
	SecurityToken   string `json:"SecurityToken"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	Expiration      string `json:"Expiration"`
}

type StsTokenResp struct {
	RequestId   string      `json:"RequestId"`
	Credentials Credentials `json:"Credentials"`
}

type GraphResp struct {
	Success   bool      `json:"success"`
	GraphJson GraphJson `json:"data"`
}

type GraphJson struct {
	GraphJsonStr string `json:"graphJsonStr"`
}

type Graph struct {
	CaseSensitive bool               `json:"caseSensitive"`
	Processes     map[string]Process `json:"processes"`
	Connections   []Connection       `json:"connections"`
}

type Process struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Def Def `json:"def"`
}

type Def struct {
	Ports	[]Port	`json:"ports"`
}

type Port struct {
	UUID	string	`json:"uuid"`
	Description	Description	`json:"description"`
}

type Description struct {
	EN_US	string	`json:"en_US"`
	ZH_CN	string	`json:"zh_CN"`
}

type Connection struct {
	Src SrcTgt `json:"src"`
	Tgt SrcTgt `json:"tgt"`
}

type SrcTgt struct {
	Process string `json:"process"`
	Port    string `json:"port"`
}
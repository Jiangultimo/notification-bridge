package types

// NetAlertXPayload 是 NetAlertX webhook 的顶层结构
type NetAlertXPayload struct {
	Text        string           `json:"text"`
	Username    string           `json:"username"`
	Attachments []AttachmentItem `json:"attachments"`
}

// AttachmentItem 是 webhook 中的附件项
type AttachmentItem struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"` // JSON 字符串,需要二次解析
}

// NetAlertXData 是 attachments[0].text 中的实际数据
type NetAlertXData struct {
	NewDevices          []NewDevice       `json:"new_devices"`
	NewDevicesMeta      MetaInfo          `json:"new_devices_meta"`
	DownDevices         []DownDevice      `json:"down_devices"`
	DownDevicesMeta     MetaInfo          `json:"down_devices_meta"`
	DownReconnected     []DownReconnected `json:"down_reconnected"`
	DownReconnectedMeta MetaInfo          `json:"down_reconnected_meta"`
	Events              []Event           `json:"events"`
	EventsMeta          MetaInfo          `json:"events_meta"`
	Plugins             []Plugin          `json:"plugins"`
	PluginsMeta         MetaInfo          `json:"plugins_meta"`
}

// MetaInfo 包含列表的元数据
type MetaInfo struct {
	Title       string   `json:"title"`
	ColumnNames []string `json:"columnNames"`
}

// NewDevice 表示新发现的设备
type NewDevice struct {
	MAC        string `json:"MAC"`
	Datetime   string `json:"Datetime"`
	IP         string `json:"IP"`
	EventType  string `json:"Event Type"`
	DeviceName string `json:"Device name"`
	Comments   string `json:"Comments"`
}

// DownDevice 表示离线的设备
type DownDevice struct {
	DevName   string `json:"devName"`
	MAC       string `json:"eve_MAC"`
	Vendor    string `json:"devVendor"`
	IP        string `json:"eve_IP"`
	DateTime  string `json:"eve_DateTime"`
	EventType string `json:"eve_EventType"`
}

// DownReconnected 表示重新连接的设备
type DownReconnected struct {
	DevName   string `json:"devName"`
	MAC       string `json:"eve_MAC"`
	Vendor    string `json:"devVendor"`
	IP        string `json:"eve_IP"`
	DateTime  string `json:"eve_DateTime"`
	EventType string `json:"eve_EventType"`
}

// Event 表示网络事件
type Event struct {
	MAC        string      `json:"MAC"`
	Datetime   string      `json:"Datetime"`
	IP         string      `json:"IP"`
	EventType  string      `json:"Event Type"`
	DeviceName string      `json:"Device name"`
	Comments   interface{} `json:"Comments"` // 可能为 null
}

// Plugin 表示插件信息
type Plugin struct {
	Plugin            string `json:"Plugin"`
	ObjectPrimaryID   string `json:"Object_PrimaryID"`
	ObjectSecondaryID string `json:"Object_SecondaryID"`
	DateTimeChanged   string `json:"DateTimeChanged"`
	WatchedValue1     string `json:"Watched_Value1"`
	WatchedValue2     string `json:"Watched_Value2"`
	WatchedValue3     string `json:"Watched_Value3"`
	WatchedValue4     string `json:"Watched_Value4"`
	Status            string `json:"Status"`
}

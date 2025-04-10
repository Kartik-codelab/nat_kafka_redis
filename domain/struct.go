package domain

import (
	"time"
)

type PositionPacket struct {
	Position PositionData `json:"position"`
	Device   DeviceData   `json:"device"`
}

type PositionData struct {
	ID         int64                  `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
	DeviceID   int64                  `json:"deviceId"`
	Type       *string                `json:"type"`
	Protocol   string                 `json:"protocol"`
	ServerTime time.Time              `json:"serverTime"`
	DeviceTime time.Time              `json:"deviceTime"`
	FixTime    time.Time              `json:"fixTime"`
	Outdated   bool                   `json:"outdated"`
	Valid      bool                   `json:"valid"`
	Latitude   float64                `json:"latitude"`
	Longitude  float64                `json:"longitude"`
	Altitude   float64                `json:"altitude"`
	Speed      float64                `json:"speed"`
	Course     float64                `json:"course"`
	Address    *string                `json:"address"`
	Accuracy   float64                `json:"accuracy"`
	Network    NetworkData            `json:"network"`
}

type DeviceData struct {
	ID          int64                  `json:"id"`
	Attributes  map[string]interface{} `json:"attributes"`
	GroupID     int64                  `json:"groupId"`
	Name        string                 `json:"name"`
	UniqueID    string                 `json:"uniqueId"`
	Status      string                 `json:"status"`
	LastUpdate  time.Time              `json:"lastUpdate"`
	PositionID  int64                  `json:"positionId"`
	GeofenceIds []int64                `json:"geofenceIds"`
	Phone       string                 `json:"phone"`
	Model       string                 `json:"model"`
	Contact     string                 `json:"contact"`
	Category    string                 `json:"category"`
	Disabled    bool                   `json:"disabled"`
}

type NetworkData struct {
	RadioType  string      `json:"radioType"`
	ConsiderIP bool        `json:"considerIp"`
	CellTowers []CellTower `json:"cellTowers"`
}

type CellTower struct {
	CellID            int `json:"cellId"`
	LocationAreaCode  int `json:"locationAreaCode"`
	MobileCountryCode int `json:"mobileCountryCode"`
	MobileNetworkCode int `json:"mobileNetworkCode"`
}

package message

import (
	"github.com/go-sql-driver/mysql"
)

// DeviceGroup is used to select which table to retrieve
type DeviceGroup int

// Different groups of devices
const (
	TestUserDevice DeviceGroup = iota
	AllUserDevice
	PaidUserDevice
)

// Device represents a row from ios_device_token.
type Device struct {
	Token      string `db:"device_token"`
	DeviceType string `db:"device_type"`
}

const (
	DeviceTypePhone = "phone"
	DeviceTypePad   = "pad"
)

// InvalidDevice represents a device that is unreachable
type InvalidDevice struct {
	Token      string
	StatusCode int
	Reason     string
	InvalidUTC mysql.NullTime
}

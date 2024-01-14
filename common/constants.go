package common

import (
	"time"
)

var StartTime = time.Now().Unix() // unit: second

var ItemsPerPage = 10

// USER ROLE
/*
root 10
内部操作员 8
分销商 6
API 2
Owner 3
Admin 1
Guest 0
*/
const (
	RoleGuestUser       = 0
	RoleAdminUser       = 1
	RoleApiUser         = 2
	RoleOwnerUser       = 3
	RoleDistributorUser = 6
	RoleOperatorUser    = 8
	RoleRootUser        = 10
)

// All duration's unit is seconds
// Shouldn't larger then RateLimitKeyExpirationDuration
var (
	GlobalApiRateLimitNum            = 180
	GlobalApiRateLimitDuration int64 = 3 * 60

	GlobalWebRateLimitNum            = 90
	GlobalWebRateLimitDuration int64 = 3 * 60

	UploadRateLimitNum            = 10
	UploadRateLimitDuration int64 = 60

	DownloadRateLimitNum            = 10
	DownloadRateLimitDuration int64 = 60

	CriticalRateLimitNum            = 20
	CriticalRateLimitDuration int64 = 20 * 60
)

var RateLimitKeyExpirationDuration = 20 * time.Minute

const (
	UserStatusEnabled  = 1 // don't use 0, 0 is the default value!
	UserStatusDisabled = 2 // also don't use 0
)

const (
	TokenStatusEnabled   = 1 // don't use 0, 0 is the default value!
	TokenStatusDisabled  = 2 // also don't use 0
	TokenStatusExpired   = 3
	TokenStatusExhausted = 4
)

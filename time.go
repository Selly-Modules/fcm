package fcm

import "time"

//
// NOTE: due to unique timezone in server's code, all using time will be convert to HCM timezone (UTC +7)
// All functions generate time, must be call util functions here
// WARNING: don't accept call time.Now() directly
//

const timezoneHCM = "Asia/Ho_Chi_Minh"

// GetHCMLocation ...
func GetHCMLocation() *time.Location {
	l, _ := time.LoadLocation(timezoneHCM)
	return l
}


// Now ...
func Now() time.Time {
	return time.Now().In(GetHCMLocation())
}

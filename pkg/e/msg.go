package e

var MsgFlags = map[int]string{
	SUCCESS:                  "ok",
	ERROR:                    "fail",
	INVALID_PARAMS:           "INVALID_PARAMS",
	ERROR_EXIST_EVENT:        "ERROR EXIST EVENT NAME",
	ERROR_EXIST_COUPON:       "ERROR EXIST COUPON CODE",
	ERROR_NOT_EXIST_COUPON:   "ERROR NOT EXIST COUPON CODE",
	ERROR_COUPON_FAIL_EXPIRE: "ERROR COUPON FAIL EXPIRE",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

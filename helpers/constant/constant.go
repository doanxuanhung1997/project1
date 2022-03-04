package constant

const (
	MaxLenOTPCode        = 6
	DateFormat           = "2006-01-02"
	DateTimeFormat       = "2006-01-02 15:04:05"
	Used                 = "USED"
	Active               = "ACTIVE"
	Inactive             = "INACTIVE"
	Refunded             = "REFUNDED"
	Blocked              = "BLOCKED"
	Deleted              = "DELETED"
	Processing           = "PROCESSING"
	Verified             = "VERIFIED"
	Started              = "STARTED"
	Unconnected          = "UNCONNECTED"
	Talking              = "TALKING"
	Completed            = "COMPLETED"
	Canceled             = "CANCELED"
	ExpiresOTP           = 3
	OrderBookAppointment = "APPOINTMENT"
	OrderBookCallNow     = "CALL_NOW"
	CheckIn              = "CHECK_IN"
	CheckOut             = "CHECK_OUT"
	OnCall               = "ON_CALL"
	OffCall              = "OFF_CALL"
	JoinCall             = "JOIN_CALL"
	MissCall             = "MISS_CALL"
	NotifyCoupon         = 1
	NotifyAppointment    = 2
	NotifyOrder          = 3
	WSEventNotify        = "notify"
	WSEventStartCall     = "start_call"
	WSEventEndCall       = "end_call"
	WSEventReadyCall     = "ready_call"
	WSEventReconnectCall = "reconnect_call"
	WSEventExtendCall    = "extend_call"
	CouponCallNow        = "cp_call_now"
	CouponBookingCG      = "cp_booking_cg"
	CouponBookingCV      = "cp_booking_cv"
	UTC7                 = 7
)

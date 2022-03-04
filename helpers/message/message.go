package message

const (
	// Message Error Validate
	MessageErrorConvertInput     = "Dữ liệu đầu vào không đúng."
	MessageErrorPermissionDenied = "Không được phép truy cập."

	//Email
	MessageErrorEmailEmpty         = "Địa chỉ email không thể trống."
	MessageErrorEmailInvalidFormat = "Địa chi email không đúng định dạng."
	MessageErrorEmailLength        = "Độ dài địa chỉ email phải từ 10 đến 256 ký tự."
	MessageErrorEmailExist         = "Địa chỉ email đã tồn tại."

	//Password
	MessageErrorPasswordEmpty   = "Mật khẩu không thể trống."
	MessageErrorPasswordFormat  = `Mật khẩu phải dài ít nhất 8 ký tự. Phải bao gồm chữ hoa, chữ thường, số và các ký tự đặc biệt !@#$%&*()-_+=[]{}|;:<>?/., `
	MessageErrorPasswordSameOld = "Mật khẩu mới không thể giống mật khẩu hiện tại."
	MessageErrorChangePassword  = "Mật khẩu cũ không đúng."

	//Name
	MessageErrorFirstNameEmpty       = "Tên không thể trống."
	MessageErrorFirstNameLength      = "Độ dài tên phải từ 1 đến 36 ký tự."
	MessageErrorFirstNameSpecialChar = "Tên không thể chứa ký tự đặc biệt trong (!@#$%^&*(),._?:{+}|<>/-) ."
	MessageErrorLastNameEmpty        = "Họ không thể trống."
	MessageErrorLastNameLength       = "Họ dài tên phải từ 1 đến 36 ký tự."
	MessageErrorLastNameSpecialChar  = "Họ không thể chứa ký tự đặc biệt trong (!@#$%^&*(),._?:{+}|<>/-) ."
	MessageErrorNameEmpty            = "Tên không thể trống."
	MessageErrorNameLength           = "Độ dài tên phải từ 1 đến 36 ký tựg"
	MessageErrorNameCharSpecial      = "Tên không thể chứa ký tự đặc biệt trong (!@#$%^&*(),._?:{+}|<>/-) ."

	//PhoneNumber
	MessageErrorPhoneNumberEmpty    = "Số điện thoại không thể trống."
	MessageErrorPhoneNumberInvalid  = "Số điện thoại không hợp lệ."
	MessageErrorPhoneNumberNotExist = "Số điện thoại không tồn tại."
	MessageErrorPhoneNumberExist    = "Số điện thoại đã tồn tại."

	//Datetime
	MessageErrorBirthdayEmpty       = "Ngày sinh không thể trống."
	MessageErrorBirthdayInvalid     = "Ngày sinh không hợp lệ."
	MessageErrorDateFormat          = "Ngày sinh không đúng định dạng."
	MessageErrorDateEmpty           = "Ngày không thể trống."
	MessageErrorTimeSlotEmpty       = "Ca không thể trống."
	MessageErrorBookingTimeEmpty    = "Thời gian đặt lịch không thể trống."
	MessageErrorDateScheduleInvalid = "Vui lòng đăng ký ca làm sau ngày hôm nay."

	//ListenerId - EmployeeId
	MessageErrorListenerIdInvalid  = "Listener id không hợp lệ."
	MessageErrorListenerIdEmpty    = "Listener id không thể trống."
	MessageErrorExpertIdEmpty      = "Expert id không thể trống."
	MessageErrorListenerIdNotExist = "Listener id không tồn tại."
	MessageErrorEmployeeIdEmpty    = "Employee Id không thể trống."
	MessageErrorEmployeeIdNotExist = "Employee Id không tồn tại."
	MessageErrorEmployeeIdExist    = "Employee Id đã tồn tại."

	//OTP
	MessageErrorCodeEmpty     = "Mã OTP không thể trống."
	MessageErrorCodeExpired   = "Mã OTP đã hết hạn."
	MessageErrorCodeIncorrect = "Mã OTP không đúng."

	//Diamond
	MessageErrorDiamondInvalid   = "Số lượng kim cương không hợp lệ."
	MessageErrorDiamondNotEnough = "Số lượng kim cương không đủ."

	MessageErrorOrderIdInvalid = "Order id không hợp lệ."
	MessageErrorOrderIdEmpty   = "Order id không thể trống."
	MessageErrorCallIdInvalid  = "Call id không hợp lệ."
	MessageErrorCallIdEmpty    = "Call id không thể trống."

	MessageErrorAccountInacctive = "Tài khoản đã bị vô hiệu hóa."
	MessageErrorStatusInvalid    = "Trạng thái không hợp lệ."
	MessageErrorActionInvalid    = "Hành động không hợp lệ."
	MessageErrorRoleInvalid      = "Vai trò ko hợp lệ."
	MessageErrorContentEmpty     = "Nội dung không thể trống."
	MessageErrorReceiverNotFound = "Người nhận không tìm thấy."
	MessageErrorScheduleNotExist = "Listener's không có ca làm việc."
	MessageErrorOrderNotAllowed  = "Chỉ được phép đặt lịch trước tối đa 1 cuộc hẹn chuyên viên và 1 cuộc hẹn chuyên gia."
	MessageErrorAppointmentExist = "Lịch hẹn này đã có người đặt."
	MessageErrorAppointmentLate  = "Vui lòng đặt lịch hẹn trước 20 phút."
	MessageErrorOrderDateLimit   = "Lịch hẹn chỉ cho phép đặt trước trong tuần."

	//Message Handle Fail
	MessageErrorReadNotifyFail      = "Đọc thông báo thất bại."
	MessageErrorLoginFail           = "Số điện thoại hoặc mật khẩu không chính xác."
	MessageErrorForgotPassword      = "Yêu cầu tìm lại mật khẩu thất bại."
	MessageErrorVerifyResetPassword = "Xác thực mã OTP thất bại."
	MessageErrorResetNewPassword    = "Đặt lại mật khẩu mới thất bại."
	MessageErrorResendOTP           = "Gủi lại mã OTP thất bại."
	MessageErrorSendOTP             = "Gủi mã OTP thất bại."
	MessageErrorCouponNotExist      = "Mã giảm giá không tồn tại."
	MessageErrorCouponInvalid       = "Mã giảm giá không hợp lệ."
	MessageErrorTokenInvalid        = "Token không hợp lệ."
	MessageErrorListenerUnconnect   = "Không thể kết nối đến Listener."
	MessageErrorStartConversation   = "Không thể bắt đầu cuộc gọi ngay bây giờ."
	MessageErrorAppointmentComingUp = "Bạn có một cuộc hẹn sắp tới."
	MessageErrorUnableExtend        = "Không thể gia hạn cuộc gọi."
	MessageErrorReferralNotExist    = "Mã giới thiệu không tồn tại."

	//Message Success
	MessageSuccess             = "Thành công."
	MessageNotifyResetPassword = "Vùi lòng đặt lại mật khẩu."

	//Notification Message
	NotificationAppointmentUser     = "Còn {number} phút nữa cuộc gọi với chuyên viên {listener_code} sẽ diễn ra"
	NotificationAppointmentListener = "Còn {number} phút nữa cuộc gọi với KH {user_name} sẽ diễn ra"

	//
	FirebaseRegisterSuccess = "Token đã được đăng kí thành công"
	FirebaseRegisterInvalid = "Thông tin đăng kí chauw hợp lệ, vui lòng kiểm tra lại"
	FirebaseSuccess         = "Tin nhắn đã được gửi thành công"
	FirebaseInvalidFormat   = "Tin nhắn sai định dạng, vui lòng kiểm tra lại"
	FirebaseInvalidToken    = "Người nhận chưa chính xác"
	FirebaseError           = "Tin nhắn chưa được gửi thanh công, vui lòng thử lại sau ít phút"

	//SMS Content
	SmsOTP     = "Ma OTP tai khoan SandexCare cua QK la: {otp_code}"
	SmsRemind  = "QK co lich hen voi SandexCare vao {hour} {day_of_week} {date}. Vui long goi cho chung toi de khong bo lo cuoc hen nhe!"
	SmsBooking = "Ban co cuoc hen voi khach hang vao {hour} {day_of_week} {date}. Vui long sap xep thoi gian de khong bo lo cuoc hen."

	//
	Abnormal            = "Lỗi chưa xác định từ service"
	InvalidToken        = "Token chưa hợp lệ"
	InvalidRefreshToken = "Refresh-Token chưa hợp lệ"
)

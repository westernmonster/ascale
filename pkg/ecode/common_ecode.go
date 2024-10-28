package ecode

// All common ecode
var (
	OK = add(200) // 正确

	AppKeyInvalid      = add(1)
	AccessKeyErr       = add(2)
	SignCheckErr       = add(3)
	MethodNoPermission = add(4)
	NoLogin            = add(101)
	UserDisabled       = add(102)
	CaptchaErr         = add(105)
	UserInactive       = add(106)
	AppDenied          = add(108)
	UserKickOut        = add(109)
	MobileNoVerfiy     = add(110)
	CsrfNotMatchErr    = add(111)
	ServiceUpdate      = add(112)

	NotModified         = add(304)
	TemporaryRedirect   = add(307)
	RequestErr          = add(400)
	Unauthorized        = add(401)
	AccessDenied        = add(403)
	NothingFound        = add(404)
	MethodNotAllowed    = add(405)
	Conflict            = add(409)
	RequestTooFast      = add(429)
	NotRepeatOperation  = add(439)
	ServerErr           = add(500)
	ServiceUnavailable  = add(503)
	Deadline            = add(504)
	LimitExceed         = add(509)
	FileTooLarge        = add(617)
	FailedTooManyTimes  = add(625)
	PasswordTooLeak     = add(628)
	PasswordErr         = add(629)
	TargetNumberLimit   = add(632)
	TargetBlocked       = add(643)
	UserDuplicate       = add(652)
	AccessTokenExpires  = add(658)
	PasswordHashExpires = add(662)
	TextTooLong         = add(700)

	Degrade     = add(1200)
	RPCNoClient = add(1201)
	RPCNoAuth   = add(1202)

	// SMS
	SMSTemplateIDNotExist    = add(20001)
	SMSTemplateInvalid       = add(20002)
	SMSTemplateExecuteFailed = add(20003)

	// Mail
	EmailSendGridFailure = add(21001)

	// Notification
	NeedUniqueID      = add(21101)
	UserNotNotifiable = add(21102)

	// OAuth
	ClientNotExist = add(22101)

	// Valcode
	ValcodeExpires = add(22201)
	ValcodeWrong   = add(22301)

	// User
	UserExist         = add(23001)
	UserPhoneNotExist = add(23002)
)

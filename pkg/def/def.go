package def

import "fmt"

func AccessTokenKey(token string) string {
	return fmt.Sprintf("ak_%s", token)
}

func RefreshTokenKey(token string) string {
	return fmt.Sprintf("rk_%s", token)
}

func MobileValcodeKey(vtype int32, mobile string) string {
	return fmt.Sprintf("rc_%d_%s", vtype, mobile)
}

func EmailValcodeKey(vtype int32, emailHex string) string {
	return fmt.Sprintf("rc_%d_%s", vtype, emailHex)
}

func AcuityDatesKey(doctorID int64, appointmentTypeID int32) string {
	return fmt.Sprintf("acuity_dates_%d_%d", doctorID, appointmentTypeID)
}
func AcuityDatesKeyByDoctorId(doctorID int64) string {
	return fmt.Sprintf("acuity_dates_%d", doctorID)
}
func AcuityTimesKey(doctorID int64, appointmentTypeID int32) string {
	return fmt.Sprintf("acuity_times_%d_%d", doctorID, appointmentTypeID)
}
func AcuityTimesKeyByDoctorId(doctorID int64) string {
	return fmt.Sprintf("acuity_times_%d", doctorID)
}

func AcuityBitmapKey(doctorID int64, dayStr string, typeId int32) string {
	return fmt.Sprintf("bitmap:%d:%v:%d", doctorID, dayStr, typeId)
}

func AcuityBitmapKeyByDoctorId(doctorID int64) string {
	return fmt.Sprintf("bitmap:%d", doctorID)
}

func AcuityTimeplanKey(doctorID int64, dayStr string, typeId int32) string {
	return fmt.Sprintf("timeplan:%d:%v:%d", doctorID, dayStr, typeId)
}

func AcuityTimeplanKeyByDoctorId(doctorID int64) string {
	return fmt.Sprintf("timeplan:%d", doctorID)
}

func AcuityBitmapLockKey(doctorID int64, dayStr string, typeId int32) string {
	return fmt.Sprintf("bitmap_lock:%d:%d:%d", doctorID, dayStr, typeId)
}

func AcuityDateKey(doctorID int64, typeId int32) string {
	return fmt.Sprintf("acuity_refactor_date:%d:%d", doctorID, typeId)
}

func AcuityDateKeyByDoctorId(doctorID int64) string {
	return fmt.Sprintf("acuity_refactor_date:%d", doctorID)
}

func TimePlanWeeklyArgsKey(doctorID int64) string {
	return fmt.Sprintf("time_plan_weekly_args:%d", doctorID)
}

func TimePlanProcessStatusKey(keySuffix int64) string {
	return fmt.Sprintf("time_plan_process_status:%d", keySuffix)
}

func FollowUpCalKeyByDoctorId(doctorID int64) string {
	return fmt.Sprintf("follow_up_cal:%d", doctorID)
}

func AppointmentUpcoming(appointment int64) string {
	return fmt.Sprintf("zoom_upcoming%d", appointment)
}

func AppointmentLateRecord(appointment int64, recordType, role string) string {
	return fmt.Sprintf("zoom_late%d_%s_%s", appointment, recordType, role)
}

func UserSignInKey(userId int64, month string) string {
	return fmt.Sprintf("user:sign:%d:%s", userId, month)
}

func TimePlanDaysKey(doctorID int64, date string) string {
	return fmt.Sprintf("time_plan_days:%d:%s", doctorID, date)
}

func AtZoomWaitingRoomTooLongKey(meetStartTime int64, meetingID, userName string) string {
	return fmt.Sprintf("zoom_waiting_too_long_reminder%s_%d_%s", meetingID, meetStartTime, userName)
}

func StripeOrderNew(orderID int64) string {
	return fmt.Sprintf("stripe_order_%d", orderID)
}

func AcuityDoctor(doctorUserID int64) string {
	return fmt.Sprintf("acuity_doctor_%d", doctorUserID)
}

func AcuityWeekRuleDoctor(doctorUserID int64) string {
	return fmt.Sprintf("acuity_week_rule_doctor_%d", doctorUserID)
}

func OrderPaymentFailPopUpKey(uid int64) string {
	return fmt.Sprintf("order_payment_fail_pop_up_%d", uid)
}
func AnnualPromotionKey(uid int64) string {
	return fmt.Sprintf("user:annual:pop:%d,%s", uid, "annual_promotion")
}

func AnnualFailKey(uid int64) string {
	return fmt.Sprintf("user:annual:fail:%d,%s", uid, "annual_fail_pop")
}

func ChartAuditOccupyKey(auditID int64) string {
	return fmt.Sprintf("privider:chart_audit:occupy:%d,%s", auditID, "quest_occupy")
}

func RefillRequestOccupyKey(requestId int64) string {
	return fmt.Sprintf("privider:refill:occupy:%d,%s", requestId, "quest_occupy")
}

func UrgentRefillPopFlag(uid int64) string {
	return fmt.Sprintf("privider:refill:popo:%d,%s", uid, "urgent_pop")
}

func PaymentMethodKey(number, expMonth, expYear, cvc string) string {
	return fmt.Sprintf("payment_method:%s,%s,%s,%s", number, expMonth, expYear, cvc)
}

func ReferralFriendEmailKey(uid int64, email string) string {
	return fmt.Sprintf("referral_friend_email_%d_%s", uid, email)
}

func GoogleCalendarCodeCheckKey(code string) string {
	return fmt.Sprintf("google_calendar_code_%s", code)
}

func PaypalOrderNew(orderID int64) string {
	return fmt.Sprintf("paypal_order_%d", orderID)
}

const CancelRxOrderCountKey = "cancel_rx_order_count_key"

const SystemTimeDeltaKey = "system_time_delta"

const FullEsTask = "full_es_task"

const VerifiableProviderKey = "verifiable_provider_key"

const ZoomWaitingRoomKey = "zoom_waiting_room_key"

const ZoomServerToServerToken = "zoom_server_to_server_token"

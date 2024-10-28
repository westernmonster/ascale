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
	UserExist                         = add(23001)
	UserPhoneNotExist                 = add(23002)
	NeedZipcode                       = add(23003)
	MembershipExpired                 = add(23004)
	AreaIsNotAvailable                = add(23005)
	EmailIsUsed                       = add(23006)
	InvalidState                      = add(23101)
	InvalidZipcode                    = add(23102)
	ZipcodeStateNotMatch              = add(23103)
	NoDoctorAssigned                  = add(23104)
	ReachedMaxConsultationNumber      = add(23105)
	UserEmailNotExist                 = add(23106)
	PhoneIsUsed                       = add(23107)
	Under18Years                      = add(23108)
	DeleteUserWhenIsMembership        = add(23109)
	PhoneNotChange                    = add(23110)
	EmailNotChange                    = add(23111)
	DoctorCannotUpdate                = add(23112)
	AreaNotOpen                       = add(23113)
	VerifiableProviderNotExist        = add(23114)
	VerifiableProviderLicenseNotExist = add(23115)
	VerifiableProviderDEANotExist     = add(23116)
	AreaIsNotSupportPullPDMP          = add(23117)
	EmailVerified                     = add(23118)
	DoctorNotVerifyIdentityPullPDMP   = add(23119)

	// Appointment
	InvalidTimezone                                  = add(23201)
	AlreadyHasAppointment                            = add(23202)
	TimeAlreadyBooked                                = add(23203)
	CancelAppointmentNotAllowed                      = add(23204)
	RescheduleAppointmentNotAllowed                  = add(23205)
	AppointmentAlreadySetDiagnosed                   = add(23206)
	AppointmentStatusNotMeetRequirements             = add(23207)
	NoInitialAppointment                             = add(23210)
	MustHaveDiagnosedAppointment                     = add(23211)
	MustBookWithSameDoctor                           = add(23212)
	MustHaveFollowUpAppointment                      = add(23213)
	MustHaveInitialNote                              = add(23214)
	NotResubmitDischargeNote                         = add(23215)
	AlreadyHasAppointmentRequest                     = add(23216)
	PatientHasAppointmentRequest                     = add(23217)
	NoSuchNoteType                                   = add(23218)
	AppointmentRequestHasBeenAccepted                = add(23219)
	AppointmentRequestPreferredTimesNotAvailable     = add(23220)
	AppointmentTimeNotWithinTherapySubscriptionCycle = add(23221)

	// Order
	EmptyPaymentMethod       = add(23301)
	ChargeFailed             = add(23302)
	PaymentNotReadyForRefund = add(23303) // Current payment status now allowed refund
	PaymentMethodSyncFailed  = add(23304)
	RefundAmountIncorrect    = add(23305)
	OnlyJoinOnceAnnual       = add(23306)
	AutoRefund               = add(23309)

	ChargeFailedNeedUpdate = add(23307)
	RepeatPaymentMethod    = add(23308)

	// Gcloud
	InvalidAccountKeyJSON  = add(23401)
	InvalidUploadTokenType = add(23402)

	// Intake Form
	InvalidAddressValue       = add(23501)
	InvalidIdentity           = add(23502)
	InvalidEmergencyContact   = add(23503)
	InvalidPrescriptionRecord = add(23504)
	UncompletedIntakeForm     = add(23505)
	InvalidProfile            = add(23506)

	// Acuity
	AcuityDefaultError                  = add(23600) // Acuity's error which are not in our list
	AcuityCancelTooClose                = add(23601) // Clients are not allowed to cancel this close to the time of the appointment.
	AcuityCancelNotAllowed              = add(23602) // Clients are not allowed to cancel.
	AcuityRescheduleNotAllowed          = add(23603) // Clients are not allowed to reschedule.
	AcuityRescheduleTooClose            = add(23605) // Clients are not allowed to reschedule this close to the time of the appointment.
	AcuityRescheduleSeries              = add(23606) // Class series may not be rescheduled.
	AcuityRescheduleCanceled            = add(23607) // Canceled appointments may not be rescheduled.
	AcuityInvalidCalendar               = add(23608) // The calendar does not exist.
	AcuityRequiredDatetime              = add(23609) // The parameter "datetime" is required.
	AcuityInvalidTimezone               = add(23610) // Invalid timezone.
	AcuityInvalidDatetime               = add(23611) // The datetime is invalid.
	AcuityNotAvailableMinHoursInAdvance = add(23612) // The time is not far enough in advance.
	AcuityNotAvailableMaxDaysInAdvance  = add(23613) // The time is too far in advance.
	AcuityNotAvailable                  = add(23614) // The time is not an available time slot.
	AcuityRequiredFirstName             = add(23615) // Attribute "firstName" is required.
	AcuityRequiredLastName              = add(23616) // Attribute "lastName" is required.
	AcuityRequiredEmail                 = add(23617) // Attribute "email" is required.
	AcuityInvalidEmail                  = add(23618) // Invalid "email" attribute value.
	AcuityInvalidFields                 = add(23619) // The field "1" does not exist on this appointment.
	AcuityRequiredField                 = add(23620) // The field "4" is required.
	AcuityRequiredAppointmentTypeID     = add(23621) // The parameter "appointmentTypeID" is required.
	AcuityInvalidAppointmentType        = add(23622) // The appointment type "987654321" does not exist.
	AcuityNoAvailableCalendar           = add(23627) // We could not find an available calendar.
	AcuityInvalidCertificate            = add(23631) // The certificate "NOCODENOPROBLEM" is invalid.
	AcuityExpiredCertificate            = add(23632) // The certificate "EXPIRED" is expired.
	AcuityCertificateUses               = add(23633) // The certificate "8013DA6F" has no remaining uses for appointment type "1".
	AcuityInvalidCertificateType        = add(23634) // The certificate "E5E5325C" is invalid for appointment type "5".
	AcuityDatetimeBlocked               = add(23635) // The time is blocked.

	// EHR Doctor
	PatientNotBelongsToYou            = add(23701)
	AppointmentNotConfirmedYet        = add(23702)
	AppointmentStatusNotAllowedRevert = add(23703)
	AppointmentAlreadyCanceled        = add(23704)
	AppointmentAlreadyNoShow          = add(23705)
	AppointmentAlreadyDiagnosis       = add(23706)

	// Treatment Request
	RefillTimeIsNotUpYet                              = add(23801)
	AlreadyHasRefillRequest                           = add(23802)
	NoTreatmentPlanYet                                = add(23803)
	ThreeMonthRenewalNotFinishedYet                   = add(23804)
	RefillRequestAlreadyNoteAdded                     = add(23805)
	RefillRequestAlreadyDeal                          = add(23806)
	RefillRequestBeboreNeededForceScheduleAppointment = add(23807)

	// Membership
	SystemBusy = add(23901) // The system is busy, Please try again later

	// Stripe
	StripeAccountAlreadyExists                   = add(24001) // The email address provided for the creation of a deferred account already has an account associated with it.
	StripeAccountCountryInvalidAddress           = add(24002) // The country of the business address provided does not match the country of the account.
	StripeAccountInvalid                         = add(24003) // The account ID provided as a value for the Stripe-Account header is invalid. Check that your requests are specifying a valid account ID.
	StripeAccountNumberInvalid                   = add(24004) // The bank account number provided is invalid.
	StripeAlipayUpgradeRequired                  = add(24005) // This method for creating Alipay payments is not supported anymore.
	StripeAmountTooLarge                         = add(24006) // The specified amount is greater than the maximum amount allowed.
	StripeAmountTooSmall                         = add(24007) // The specified amount is less than the minimum amount allowed.
	StripeApiKeyExpired                          = add(24008) // The API key provided has expired.
	StripeAuthenticationRequired                 = add(24009) // The payment requires authentication to proceed.
	StripeBalanceInsufficient                    = add(24010) // The transfer or payout could not be completed because the associated account does not have a sufficient balance available.
	StripeBankAccountDeclined                    = add(24011) // The bank account provided can not be used to charge, either because it is not verified yet or it is not supported.
	StripeBankAccountExists                      = add(24012) // The bank account provided already exists on the specified Customer object.
	StripeBankAccountUnusable                    = add(24013) // The bank account provided cannot be used for payouts.
	StripeBankAccountUnverified                  = add(24014) // Your Connect platform is attempting to share an unverified bank account with a connected account.
	StripeBankAccountVerificationFailed          = add(24015) // The bank account cannot be verified, either because the microdeposit amounts provided do not match the actual amounts, or because verification has failed too many times.
	StripeBitcoinUpgradeRequired                 = add(24016) // This method for creating Bitcoin payments is not supported anymore.
	StripeCardDeclineRateLimitExceeded           = add(24017) // This card has been declined too many times. You can try to charge this card again after 24 hours.
	StripeCardDeclined                           = add(24018) // The card has been declined.
	StripeChargeAlreadyCaptured                  = add(24019) // The charge you’re attempting to capture has already been captured.
	StripeChargeAlreadyRefunded                  = add(24020) // The charge you’re attempting to refund has already been refunded.
	StripeChargeDisputed                         = add(24021) // The charge you’re attempting to refund has been charged back.
	StripeChargeExceedsSourceLimit               = add(24022) // This charge would cause you to exceed your rolling-window processing limit for this source type.
	StripeChargeExpiredForCapture                = add(24023) // The charge cannot be captured as the authorization has expired.
	StripeChargeInvalidParameter                 = add(24024) // One or more provided parameters was not allowed for the given operation on the Charge.
	StripeClearingCodeUnsupported                = add(24025) // The clearing code provided is not supported.
	StripeCountryCodeInvalid                     = add(24026) // The country code provided was invalid.
	StripeCountryUnsupported                     = add(24027) // Your platform attempted to create a custom account in a country that is not yet supported.
	StripeCouponExpired                          = add(24028) // The coupon provided for a subscription or order has expired.
	StripeCustomerMaxPaymentMethods              = add(24029) // The maximum number of PaymentMethods for this Customer has been reached.
	StripeCustomerMaxSubscriptions               = add(24030) // The maximum number of subscriptions for a customer has been reached.
	StripeEmailInvalid                           = add(24031) // The email address is invalid (e.g., not properly formatted).
	StripeExpiredCard                            = add(24032) // The card has expired. Check the expiration date or use a different card.
	StripeIdempotencyKeyInUse                    = add(24033) // The idempotency key provided is currently being used in another request.
	StripeIncorrectAddress                       = add(24034) // The card’s address is incorrect.
	StripeIncorrectCvc                           = add(24035) // The card’s security code is incorrect.
	StripeIncorrectNumber                        = add(24036) // The card number is incorrect. Check the card’s number or use a different card.
	StripeIncorrectZip                           = add(24037) // The card’s ZIP code is incorrect. Check the card’s ZIP code or use a different card.
	StripeInstantPayoutsUnsupported              = add(24038) // This card is not eligible for Instant Payouts.  Try a debit card from a supported bank.
	StripeIntentInvalidState                     = add(24039) // Intent is not the state that is rquired to perform the operation.
	StripeIntentVerificationMethodMissing        = add(24040) // Intent does not have verification method specified in its PaymentMethodOptions object.
	StripeInvalidCardType                        = add(24041) // The card provided as an external account is not supported for payouts.
	StripeInvalidCharacters                      = add(24042) // This value provided to the field contains characters that are unsupported by the field.
	StripeInvalidChargeAmount                    = add(24043) // The specified amount is invalid.
	StripeInvalidCvc                             = add(24044) // The card’s security code is invalid. Check the card’s security code or use a different card.
	StripeInvalidExpiryMonth                     = add(24045) // The card’s expiration month is incorrect. Check the expiration date or use a different card.
	StripeInvalidExpiryYear                      = add(24046) // The card’s expiration year is incorrect. Check the expiration date or use a different card.
	StripeInvalidNumber                          = add(24047) // The card number is invalid. Check the card details or use a different card.
	StripeInvalidSourceUsage                     = add(24048) // The source cannot be used because it is not in the correct state
	StripeInvoiceNoCustomerLineItems             = add(24049) // An invoice cannot be generated for the specified customer as there are no pending invoice items.
	StripeInvoiceNoPaymentMethodTypes            = add(24050) // An invoice cannot be finalized because there are no payment method types available to process the payment.
	StripeInvoiceNoSubscriptionLineItems         = add(24051) // An invoice cannot be generated for the specified subscription as there are no pending invoice items.
	StripeInvoiceNotEditable                     = add(24052) // The specified invoice can no longer be edited. Instead, consider creating additional invoice items that will be applied to the next invoice.
	StripeInvoicePaymentIntentRequiresAction     = add(24053) // This payment requires additional user action before it can be completed successfully.
	StripeInvoiceUpcomingNone                    = add(24054) // There is no upcoming invoice on the specified customer to preview.
	StripeLivemodeMismatch                       = add(24055) // Test and live mode API keys, requests, and objects are only available within the mode they are in.
	StripeLockTimeout                            = add(24056) // This object cannot be accessed right now because another API request or Stripe process is currently accessing it. If you see this error intermittently, retry the request.
	StripeMissing                                = add(24057) // Both a customer and source ID have been provided, but the source has not been saved to the customer.
	StripeNotAllowedOnStandardAccount            = add(24058) // Transfers and payouts on behalf of a Standard connected account are not allowed.
	StripeOrderCreationFailed                    = add(24059) // The order could not be created. Check the order details and then try again.
	StripeOrderRequiredSettings                  = add(24060) // The order could not be processed as it is missing required information. Check the information provided and try again.
	StripeOrderStatusInvalid                     = add(24061) // The order cannot be updated because the status provided is either invalid or does not follow the order lifecycle.
	StripeOrderUpstreamTimeout                   = add(24062) // The request timed out. Try again later.
	StripeOutOfInventory                         = add(24063) // The SKU is out of stock. If more stock is available, update the SKU’s inventory quantity and try again.
	StripeParameterInvalidEmpty                  = add(24064) // One or more required values were not provided.
	StripeParameterInvalidInteger                = add(24065) // One or more of the parameters requires an integer, but the values provided were a different type.
	StripeParameterInvalidStringBlank            = add(24066) // One or more values provided only included whitespace.
	StripeParameterInvalidStringEmpty            = add(24067) // One or more required string values is empty.
	StripeParameterMissing                       = add(24068) // One or more required values are missing.
	StripeParameterUnknown                       = add(24069) // The request contains one or more unexpected parameters.
	StripeParametersExclusive                    = add(24070) // Two or more mutually exclusive parameters were provided.
	StripePaymentIntentActionRequired            = add(24071) // The provided payment method requires customer actions to complete, but error_on_requires_action was set. I
	StripePaymentIntentAuthenticationFailure     = add(24072) // The provided payment method has failed authentication. Provide a new payment method to attempt to fulfill this PaymentIntent again.
	StripePaymentIntentIncompatiblePaymentMethod = add(24073) // The PaymentIntent expected a payment method with different properties than what was provided.
	StripePaymentIntentInvalidParameter          = add(24074) // One or more provided parameters was not allowed for the given operation on the PaymentIntent.
	StripePaymentIntentPaymentAttemptFailed      = add(24075) // The latest payment attempt for the PaymentIntent has failed.
	StripePaymentIntentUnexpectedState           = add(24076) // The PaymentIntent’s state was incompatible with the operation you were trying to perform.
	StripePaymentMethodInvalidParameter          = add(24077) // Invalid parameter was provided in the payment method object.
	StripePaymentMethodProviderDecline           = add(24078) // The payment was declined by the issuer or customer.
	StripePaymentMethodProviderTimeout           = add(24079) // The payment method failed due to a timeout.
	StripePaymentMethodUnactivated               = add(24080) // The operation cannot be performed as the payment method used has not been activated.
	StripePaymentMethodUnexpectedState           = add(24081) // The provided payment method’s state was incompatible with the operation you were trying to perform.
	StripePayoutsNotAllowed                      = add(24082) // Payouts have been disabled on the connected account.
	StripePlatformApiKeyExpired                  = add(24083) // The API key provided by your Connect platform has expired.
	StripePostalCodeInvalid                      = add(24084) // The ZIP code provided was incorrect.
	StripeProcessingError                        = add(24085) // An error occurred while processing the card. Try again later or with a different payment method.
	StripeProductInactive                        = add(24086) // The product this SKU belongs to is no longer available for purchase.
	StripeRateLimit                              = add(24087) // Too many requests hit the API too quickly. We recommend an exponential backoff of your requests.
	StripeResourceAlreadyExists                  = add(24088) // A resource with a user-specified ID (e.g., plan or coupon) already exists.
	StripeResourceMissing                        = add(24089) // The ID provided is not valid. Either the resource does not exist, or an ID for a different resource has been provided.
	StripeRoutingNumberInvalid                   = add(24090) // The bank routing number provided is invalid.
	StripeSecretKeyRequired                      = add(24091) // The API key provided is a publishable key, but a secret key is required.
	StripeSepaUnsupportedAccount                 = add(24092) // Your account does not support SEPA payments.
	StripeSetupAttemptFailed                     = add(24093) // The latest setup attempt for the SetupIntent has failed.
	StripeSetupIntentAuthenticationFailure       = add(24094) // The provided payment method has failed authentication.
	StripeSetupIntentInvalidParameter            = add(24095) // One or more provided parameters was not allowed for the given operation on the SetupIntent.
	StripeSetupIntentUnexpectedState             = add(24096) // The SetupIntent’s state was incompatible with the operation you were trying to perform.
	StripeShippingCalculationFailed              = add(24097) // Shipping calculation failed as the information provided was either incorrect or could not be verified.
	StripeSkuInactive                            = add(24098) // The SKU is inactive and no longer available for purchase.
	StripeStateUnsupported                       = add(24099) // Occurs when providing the legal_entity information for a U.S. custom account, if the provided state is not supported.
	StripeTaxIdInvalid                           = add(24100) // The tax ID number provided is invalid (e.g., missing digits).
	StripeTaxesCalculationFailed                 = add(24101) // Tax calculation for the order failed.
	StripeTerminalLocationCountryUnsupported     = add(24102) // Terminal is currently only available in some countries.
	StripeTestmodeChargesOnly                    = add(24103) // Your account has not been activated and can only make test charges.
	StripeTlsVersionUnsupported                  = add(24104) // Your integration is using an older version of TLS that is unsupported. You must be using TLS 1.2 or above.
	StripeTokenAlreadyUsed                       = add(24105) // The token provided has already been used. You must create a new token before you can retry this request.
	StripeTokenInUse                             = add(24106) // The token provided is currently being used in another request. This occurs if your integration is making duplicate requests simultaneously.
	StripeTransfersNotAllowed                    = add(24107) // The requested transfer cannot be created.
	StripeUpstreamOrderCreationFailed            = add(24108) // The order could not be created. Check the order details and then try again.
	StripeUrlInvalid                             = add(24109) // The URL provided is invalid.

	StripeDeclineAuthenticationRequired         = add(24110) // The card was declined as the transaction requires authentication.
	StripeDeclineApproveWithId                  = add(24111) // The payment can’t be authorized.
	StripeDeclineCallIssuer                     = add(24112) // The card was declined for an unknown reason.
	StripeDeclineCardNotSupported               = add(24113) // The card does not support this type of purchase.
	StripeDeclineCardVelocityExceeded           = add(24114) // The customer has exceeded the balance or credit limit available on their card.
	StripeDeclineCurrencyNotSupported           = add(24115) // The card does not support the specified currency.
	StripeDeclineDoNotHonor                     = add(24116) // The card was declined for an unknown reason.
	StripeDeclineDoNotTryAgain                  = add(24117) // The card was declined for an unknown reason.
	StripeDeclineDuplicateTransaction           = add(24118) // A transaction with identical amount and credit card information was submitted very recently.
	StripeDeclineExpiredCard                    = add(24119) // The card has expired.
	StripeDeclineFraudulent                     = add(24120) // The payment was declined because Stripe suspects that it’s fraudulent.
	StripeDeclineGenericDecline                 = add(24121) // The card was declined for an unknown reason or possibly triggered by a blocked payment rule.
	StripeDeclineIncorrectNumber                = add(24122) // The card number is incorrect.
	StripeDeclineIncorrectCvc                   = add(24123) // The CVC number is incorrect.
	StripeDeclineIncorrectPin                   = add(24124) // The PIN entered is incorrect. This decline code only applies to payments made with a card reader.
	StripeDeclineIncorrectZip                   = add(24125) // The postal code is incorrect.
	StripeDeclineInsufficientFunds              = add(24126) // The card has insufficient funds to complete the purchase.
	StripeDeclineInvalidAccount                 = add(24127) // The card, or account the card is connected to, is invalid.
	StripeDeclineInvalidAmount                  = add(24128) // The payment amount is invalid, or exceeds the amount that’s allowed.
	StripeDeclineInvalidCvc                     = add(24129) // The CVC number is incorrect.
	StripeDeclineInvalidExpiryMonth             = add(24130) // The expiration month is invalid.
	StripeDeclineInvalidExpiryYear              = add(24131) // The expiration year is invalid.
	StripeDeclineInvalidNumber                  = add(24132) // The card number is incorrect.
	StripeDeclineInvalidPin                     = add(24133) // The PIN entered is incorrect. This decline code only applies to payments made with a card reader.
	StripeDeclineIssuerNotAvailable             = add(24134) // The card issuer couldn’t be reached, so the payment couldn’t be authorized.
	StripeDeclineLostCard                       = add(24135) // The payment was declined because the card is reported lost.
	StripeDeclineMerchantBlacklist              = add(24136) // The payment was declined because it matches a value on the Stripe user’s block list.
	StripeDeclineNewAccountInformationAvailable = add(24137) // The card, or account the card is connected to, is invalid.
	StripeDeclineNoActionTaken                  = add(24138) // The card was declined for an unknown reason.
	StripeDeclineNotPermitted                   = add(24139) // The payment isn’t permitted.
	StripeDeclineOfflinePinRequired             = add(24140) // The card was declined because it requires a PIN.
	StripeDeclineOnlineOrOfflinePinRequired     = add(24141) // The card was declined as it requires a PIN.
	StripeDeclinePickupCard                     = add(24142) // The customer can’t use this card to make this payment (it’s possible it was reported lost or stolen).
	StripeDeclinePinTryExceeded                 = add(24143) // The allowable number of PIN tries was exceeded.
	StripeDeclineProcessingError                = add(24144) // An error occurred while processing the card.
	StripeDeclineReenterTransaction             = add(24145) // The payment couldn’t be processed by the issuer for an unknown reason.
	StripeDeclineRestrictedCard                 = add(24146) // The customer can’t use this card to make this payment (it’s possible it was reported lost or stolen).
	StripeDeclineRevocationOfAllAuthorizations  = add(24147) // The card was declined for an unknown reason.
	StripeDeclineRevocationOfAuthorization      = add(24148) // The card was declined for an unknown reason.
	StripeDeclineSecurityViolation              = add(24149) // The card was declined for an unknown reason.
	StripeDeclineServiceNotAllowed              = add(24150) // The card was declined for an unknown reason.
	StripeDeclineStolenCard                     = add(24151) // The payment was declined because the card is reported stolen.
	StripeDeclineStopPaymentOrder               = add(24152) // The card was declined for an unknown reason.
	StripeDeclineTestmodeDecline                = add(24153) // A Stripe test card number was used.
	StripeDeclineTransactionNotAllowed          = add(24154) // The card was declined for an unknown reason.
	StripeDeclineTryAgainLater                  = add(24155) // The card was declined for an unknown reason.
	StripeDeclineWithdrawalCountLimitExceeded   = add(24156) // The customer has exceeded the balance or credit limit available on their card.
	StripeRequestTooClose                       = add(24440) // The interval between two deductions takes 5 minutes to prevent repeated deductions.
	NoStripAccount                              = add(24510)
	AmountMustLargeZero                         = add(24511)
	// Admin
	DoctorHasActivePatients                              = add(25001)
	TransferToSameDoctorNotAllowed                       = add(25002)
	TargetDoctorNotSupportThisState                      = add(25003)
	TransferAppointmentReschduleIsRequired               = add(25004)
	DoctorAppointmentTypeIDUsed                          = add(25005)
	DoctorFollowUpAppointmentTypeIDUsed                  = add(25006)
	DoctorExternalCalendarIDUsed                         = add(25007)
	MustCancelAllAppointment                             = add(25008)
	NeedFillInPersonInfoFirst                            = add(25009)
	DoctorHasActiveMembershipPatients                    = add(25010)
	EmptyTransferAction                                  = add(25011)
	NoHasCompletedInitialAppointmentCanNotAssignProvider = add(25012)

	// rx order cancel request
	RxOrderCannotCancel = add(25101)

	// Ayva
	NeedSync2AyvaFirst = add(26001)

	// Doctor
	DoctorAddressRequired                 = add(26100)
	DoctorNavigationInstructionRequired   = add(26101)
	DoctorInPersonalAppointmentIDRequired = add(26102)
	SupervisingDoctorIsNotCollabType      = add(26103)

	// UserMedicine
	UserMedicineExist         = add(27100)
	UserMedicinePillTimeError = add(27101)
	UserMedicinePillLeftError = add(27102)

	PharmacyBlacklistExist = add(28100)

	PharmacyInBackList = add(28110)

	// Membership
	MembershipDischarged                   = add(30001)
	ReactivateRequestAlreadySent           = add(30002)
	MembershipResumeFailed                 = add(30003)
	AlreadyHasReactivateRequest            = add(30004)
	NoSuitProvider                         = add(30005)
	HasActiveMembership                    = add(30006)
	AlreadySelectedTherapySubscriptionPlan = add(30007)
	HasPendingTherapySubscription          = add(30008)
	InactiveTherapySubscription            = add(30009)
	ExtendMembershipWeekToLarge            = add(30010)
	OnlyActiveMembershipCanExtend          = add(30011)
	TargetMembershipHasExtend              = add(30012)
	NoMembershipForNextMonth               = add(30013)
	AnnualMembershipNotAllowExtend         = add(30014)
	// Twilio
	VerificationDefaultSendError    = add(60000)
	VerificationDefaultReceiveError = add(60001)
	InvalidParameter                = add(60200)
	MaxCheckAttemptsReached         = add(60202)
	MaxSendAttemptsReached          = add(60203)
	LandlinePhoneNumberError        = add(60205)

	// Fax
	SendFaxError             = add(60301)
	FetchAyvaEncounterFailed = add(60302)
	EmptyPharmacyID          = add(60303)
	FindPharmacyFailed       = add(60304)
	ViewFaxError             = add(60305)
	ResendFaxError           = add(60306)

	// Rxnt
	NeedSync2RxntFirst              = add(60401)
	RxntAccountLoginError           = add(60402)
	RxntDoctorNotFound              = add(60403)
	RxntDoctorSecondAccountNotFound = add(60404)

	//freeMembership
	FreeMemberShipCondition = add(80100)

	//tags
	TagsSameName  = add(90100)
	TagsStillUsed = add(90200)

	//user clock-in gift
	GiftObtainConditionNotMatch = add(91100)

	//auto reschedule
	OutOfAutoRescheduleRange    = add(114111) //Reschedule Time Must be limited to 5 days
	ConflictTime                = add(114112) //additional time conflict with booked time
	RescheduleStatusMustConfirm = add(114113) // appointment status must be confirm
	NoneAvailableTimeSlot       = add(114114) // none available time slot for patient
	AlreadyAcceptReschedule     = add(114115) // your already accept

	//in person rules
	FailedToCalculateInPersonEndDate = add(115000) //Failed to calculate In Person End Date

	//Apero
	FailedToSyncPatientInfoToApero      = add(116000)
	FailedToGetPatientInfoFromApero     = add(116001)
	FailedToGetProviderInfoFromApero    = add(116002)
	FailedToSyncProviderInfoToApero     = add(116003)
	FailedToSyncFacilityInfoToApero     = add(116004)
	FailedToSyncAppointmentInfoToApero  = add(116005)
	FailedToGetAppointmentInfoFromApero = add(116006)
	FailedToSyncInsuranceInfoToApero    = add(116007)
	FailedToSyncInvoiceInfoToApero      = add(116008)

	AperoInvoiceAlreadyExists       = add(116009)
	PatientHomeBaseNotFound         = add(116010)
	DoctorLocationNotFound          = add(116011)
	AperoAppointmentIDNotFound      = add(116012)
	UserAddressNotFound             = add(116013)
	FailedToSyncDiagnosisToApero    = add(116014)
	FailedToSyncLineItemInfoToApero = add(116015)

	//zoho
	ZohoRefreshTokenFail = add(117000) //Failed to calculate In Person End Date

	//referral code
	ReferralNothingFound = add(120404)
	//referral expired
	ReferralExpired = add(120504)
	//code exist
	ReferralExist = add(120600)

	//Auto scheduling
	OldBookedNotMatch = add(120700)

	ExceedDailyLimit = add(120800)

	ExceedYearlyLimit = add(120801)

	ReferNotActiveMembership = add(120802)
	//thirdparty interface
	IncorrectTimeInterval = add(140000)
	// Candid
	CandidGetAuthTokenFailed  = add(130000)
	CandidPostEncounterFailed = add(130001)

	//paypal
	PayPalError             = add(118000) // paypal create authorization unexpected, please try again later
	PayPalAgreementInactive = add(118001) // paypal agreement inactive
	PayPalRequestTooClose   = add(118002) // The interval between two deductions takes 5 minutes to prevent repeated deductions.
	PayPalRefundFail        = add(118003) // Refund fail

)

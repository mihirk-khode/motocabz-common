package error

// ErrorCode represents a custom error code
type ErrorCode int

const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeValidation
	ErrorCodeNotFound
	ErrorCodeUnauthorized
	ErrorCodeForbidden
	ErrorCodeConflict
	ErrorCodeInternal
	ErrorCodeTimeout
	ErrorCodeRateLimit

	ErrorCodeServiceUnavailable
	ErrorCodeDatabaseError
	ErrorCodeNetworkError
	ErrorCodeConfigurationError

	ErrorCodeTripNotFound
	ErrorCodeTripAlreadyExists
	ErrorCodeInvalidTripStatus
	ErrorCodeTripCancelled
	ErrorCodeTripExpired

	ErrorCodeDriverNotFound
	ErrorCodeDriverOffline
	ErrorCodeDriverBusy
	ErrorCodeInvalidDriverStatus

	ErrorCodeRiderNotFound
	ErrorCodeRiderInactive
	ErrorCodeInvalidRiderStatus

	ErrorCodeBiddingSessionNotFound
	ErrorCodeBiddingSessionExpired
	ErrorCodeInvalidBidAmount
	ErrorCodeBiddingNotAllowed

	ErrorCodeInvalidLocation
	ErrorCodeLocationNotFound
	ErrorCodeLocationOutOfRange

	ErrorCodePaymentFailed
	ErrorCodePaymentNotFound
	ErrorCodeInvalidPaymentMethod
	ErrorCodeInsufficientFunds
)

var ErrorCodeNames = map[ErrorCode]string{
	ErrorCodeUnknown:                "UNKNOWN_ERROR",
	ErrorCodeValidation:             "VALIDATION_ERROR",
	ErrorCodeNotFound:               "NOT_FOUND",
	ErrorCodeUnauthorized:           "UNAUTHORIZED",
	ErrorCodeForbidden:              "FORBIDDEN",
	ErrorCodeConflict:               "CONFLICT",
	ErrorCodeInternal:               "INTERNAL_ERROR",
	ErrorCodeTimeout:                "TIMEOUT",
	ErrorCodeRateLimit:              "RATE_LIMIT_EXCEEDED",
	ErrorCodeServiceUnavailable:     "SERVICE_UNAVAILABLE",
	ErrorCodeDatabaseError:          "DATABASE_ERROR",
	ErrorCodeNetworkError:           "NETWORK_ERROR",
	ErrorCodeConfigurationError:     "CONFIGURATION_ERROR",
	ErrorCodeTripNotFound:           "TRIP_NOT_FOUND",
	ErrorCodeTripAlreadyExists:      "TRIP_ALREADY_EXISTS",
	ErrorCodeInvalidTripStatus:      "INVALID_TRIP_STATUS",
	ErrorCodeTripCancelled:          "TRIP_CANCELLED",
	ErrorCodeTripExpired:            "TRIP_EXPIRED",
	ErrorCodeDriverNotFound:         "DRIVER_NOT_FOUND",
	ErrorCodeDriverOffline:          "DRIVER_OFFLINE",
	ErrorCodeDriverBusy:             "DRIVER_BUSY",
	ErrorCodeInvalidDriverStatus:    "INVALID_DRIVER_STATUS",
	ErrorCodeRiderNotFound:          "RIDER_NOT_FOUND",
	ErrorCodeRiderInactive:          "RIDER_INACTIVE",
	ErrorCodeInvalidRiderStatus:     "INVALID_RIDER_STATUS",
	ErrorCodeBiddingSessionNotFound: "BIDDING_SESSION_NOT_FOUND",
	ErrorCodeBiddingSessionExpired:  "BIDDING_SESSION_EXPIRED",
	ErrorCodeInvalidBidAmount:       "INVALID_BID_AMOUNT",
	ErrorCodeBiddingNotAllowed:      "BIDDING_NOT_ALLOWED",
	ErrorCodeInvalidLocation:        "INVALID_LOCATION",
	ErrorCodeLocationNotFound:       "LOCATION_NOT_FOUND",
	ErrorCodeLocationOutOfRange:     "LOCATION_OUT_OF_RANGE",
	ErrorCodePaymentFailed:          "PAYMENT_FAILED",
	ErrorCodePaymentNotFound:        "PAYMENT_NOT_FOUND",
	ErrorCodeInvalidPaymentMethod:   "INVALID_PAYMENT_METHOD",
	ErrorCodeInsufficientFunds:      "INSUFFICIENT_FUNDS",
}

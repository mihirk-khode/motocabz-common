package validation

import "github.com/motocabz/common"

// IsValidTripStatus validates trip status using standardized constants
func IsValidTripStatus(status string) bool {
	validStatuses := []string{
		common.TripStatusDraft,
		common.TripStatusRequested,
		common.TripStatusSearchingForDriver,
		common.TripStatusOfferPhase,
		common.TripStatusAssigned,
		common.TripStatusDriverEnRoute,
		common.TripStatusDriverArrived,
		common.TripStatusInProgress,
		common.TripStatusCompleted,
		common.TripStatusCancelled,
		common.TripStatusFailed,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidPaymentStatus validates payment status using standardized constants
func IsValidPaymentStatus(status string) bool {
	validStatuses := []string{
		common.PaymentStatusPending,
		common.PaymentStatusProcessing,
		common.PaymentStatusCompleted,
		common.PaymentStatusFailed,
		common.PaymentStatusCancelled,
		common.PaymentStatusRefunded,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidFareModel validates fare model using standardized constants
func IsValidFareModel(model string) bool {
	validModels := []string{
		common.FareModelInstantMatch,
		common.FareModelFlexFare,
		common.FareModelAutomaticFare,
	}
	for _, m := range validModels {
		if m == model {
			return true
		}
	}
	return false
}

// IsValidBiddingStatus validates bidding status using standardized constants
func IsValidBiddingStatus(status string) bool {
	validStatuses := []string{
		common.BiddingStatusActive,
		common.BiddingStatusExpired,
		common.BiddingStatusAssigned,
		common.BiddingStatusCancelled,
		common.BiddingStatusCompleted,
		common.BiddingStatusFailed,
		common.BiddingStatusPending,
		common.BiddingStatusStarted,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidDriverStatus validates driver status using standardized constants
func IsValidDriverStatus(status string) bool {
	validStatuses := []string{
		common.DriverStatusOnline,
		common.DriverStatusOffline,
		common.DriverStatusAvailable,
		common.DriverStatusBusy,
		common.DriverStatusInactive,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidDriverApprovalStatus validates driver approval status using standardized constants
func IsValidDriverApprovalStatus(status string) bool {
	validStatuses := []string{
		common.DriverApprovalStatusPendingReview,
		common.DriverApprovalStatusApproved,
		common.DriverApprovalStatusRejected,
		common.DriverApprovalStatusSuspended,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidVehicleClass validates vehicle class using standardized constants
func IsValidVehicleClass(class string) bool {
	validClasses := []string{
		common.VehicleClassClassic,
		common.VehicleClassComfort,
		common.VehicleClassLuxury,
		common.VehicleClassMinivan,
	}
	for _, c := range validClasses {
		if c == class {
			return true
		}
	}
	return false
}

// IsValidPaymentMethod validates payment method using standardized constants
func IsValidPaymentMethod(method string) bool {
	validMethods := []string{
		common.PaymentMethodCash,
		common.PaymentMethodWallet,
		common.PaymentMethodTelebirr,
		common.PaymentMethodCard,
		common.PaymentMethodStripe,
		common.PaymentMethodPayPal,
	}
	for _, m := range validMethods {
		if m == method {
			return true
		}
	}
	return false
}

// IsValidComplaintStatus validates complaint status using standardized constants
func IsValidComplaintStatus(status string) bool {
	validStatuses := []string{
		common.ComplaintStatusOpen,
		common.ComplaintStatusAssigned,
		common.ComplaintStatusInProgress,
		common.ComplaintStatusResolved,
		common.ComplaintStatusClosed,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidComplaintPriority validates complaint priority using standardized constants
func IsValidComplaintPriority(priority string) bool {
	validPriorities := []string{
		common.ComplaintPriorityLow,
		common.ComplaintPriorityMedium,
		common.ComplaintPriorityHigh,
		common.ComplaintPriorityUrgent,
	}
	for _, p := range validPriorities {
		if p == priority {
			return true
		}
	}
	return false
}

// IsValidAdminRole validates admin role using standardized constants
func IsValidAdminRole(role string) bool {
	validRoles := []string{
		common.AdminRoleSuperAdmin,
		common.AdminRoleAdmin,
		common.AdminRoleOps,
		common.AdminRoleSupport,
	}
	for _, r := range validRoles {
		if r == role {
			return true
		}
	}
	return false
}

// IsValidUserRole validates user role using standardized constants
func IsValidUserRole(role string) bool {
	validRoles := []string{
		common.UserRoleAdmin,
		common.UserRoleDriver,
		common.UserRoleRider,
	}
	for _, r := range validRoles {
		if r == role {
			return true
		}
	}
	return false
}

// IsValidDocumentType validates document type using standardized constants
func IsValidDocumentType(docType string) bool {
	validTypes := []string{
		common.DocumentTypeLicense,
		common.DocumentTypeTIN,
		common.DocumentTypeID,
		common.DocumentTypePassport,
		common.DocumentTypeInsurance,
		common.DocumentTypeRegistration,
	}
	for _, t := range validTypes {
		if t == docType {
			return true
		}
	}
	return false
}

// IsValidAddressType validates address type using standardized constants
func IsValidAddressType(addrType string) bool {
	validTypes := []string{
		common.AddressTypeHome,
		common.AddressTypeWork,
		common.AddressTypeOther,
	}
	for _, t := range validTypes {
		if t == addrType {
			return true
		}
	}
	return false
}

// IsValidProvider validates provider using standardized constants
func IsValidProvider(provider string) bool {
	validProviders := []string{
		common.ProviderGoogle,
		common.ProviderManual,
		common.ProviderPhone,
	}
	for _, p := range validProviders {
		if p == provider {
			return true
		}
	}
	return false
}

// IsValidDeviceType validates device type using standardized constants
func IsValidDeviceType(deviceType string) bool {
	validTypes := []string{
		common.DeviceTypeAndroid,
		common.DeviceTypeIOS,
		common.DeviceTypeWeb,
	}
	for _, t := range validTypes {
		if t == deviceType {
			return true
		}
	}
	return false
}

// IsValidSubscriptionType validates subscription type using standardized constants
func IsValidSubscriptionType(subType string) bool {
	validTypes := []string{
		common.SubscriptionTypeNone,
		common.SubscriptionTypeBasic,
		common.SubscriptionTypePremium,
		common.SubscriptionTypeEnterprise,
	}
	for _, t := range validTypes {
		if t == subType {
			return true
		}
	}
	return false
}

// IsValidNegotiationStatus validates negotiation status using standardized constants
func IsValidNegotiationStatus(status string) bool {
	validStatuses := []string{
		common.NegotiationStatusOffered,
		common.NegotiationStatusCounter,
		common.NegotiationStatusPending,
		common.NegotiationStatusAccepted,
		common.NegotiationStatusRejected,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidWalletStatus validates wallet status using standardized constants
func IsValidWalletStatus(status string) bool {
	validStatuses := []string{
		common.WalletStatusActive,
		common.WalletStatusSuspended,
		common.WalletStatusClosed,
		common.WalletStatusFrozen,
	}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidWalletTransactionType validates wallet transaction type using standardized constants
func IsValidWalletTransactionType(txType string) bool {
	validTypes := []string{
		common.WalletTransactionTypeCredit,
		common.WalletTransactionTypeDebit,
		common.WalletTransactionTypeFreeze,
		common.WalletTransactionTypeUnfreeze,
	}
	for _, t := range validTypes {
		if t == txType {
			return true
		}
	}
	return false
}

// IsValidWalletTransactionReason validates wallet transaction reason using standardized constants
func IsValidWalletTransactionReason(reason string) bool {
	validReasons := []string{
		common.WalletTransactionReasonRidePayment,
		common.WalletTransactionReasonRefund,
		common.WalletTransactionReasonBonus,
		common.WalletTransactionReasonTopUp,
		common.WalletTransactionReasonWithdrawal,
		common.WalletTransactionReasonPromo,
	}
	for _, r := range validReasons {
		if r == reason {
			return true
		}
	}
	return false
}

// ValidateTripStatus validates trip status and returns ValidationError if invalid
func ValidateTripStatus(status string) *ValidationError {
	if !IsValidTripStatus(status) {
		return &ValidationError{
			Field:   "tripStatus",
			Message: "invalid trip status",
			Value:   status,
		}
	}
	return nil
}

// ValidatePaymentStatus validates payment status and returns ValidationError if invalid
func ValidatePaymentStatus(status string) *ValidationError {
	if !IsValidPaymentStatus(status) {
		return &ValidationError{
			Field:   "paymentStatus",
			Message: "invalid payment status",
			Value:   status,
		}
	}
	return nil
}

// ValidateFareModel validates fare model and returns ValidationError if invalid
func ValidateFareModel(model string) *ValidationError {
	if !IsValidFareModel(model) {
		return &ValidationError{
			Field:   "fareModel",
			Message: "invalid fare model",
			Value:   model,
		}
	}
	return nil
}

// ValidateBiddingStatus validates bidding status and returns ValidationError if invalid
func ValidateBiddingStatus(status string) *ValidationError {
	if !IsValidBiddingStatus(status) {
		return &ValidationError{
			Field:   "biddingStatus",
			Message: "invalid bidding status",
			Value:   status,
		}
	}
	return nil
}

// ValidateDriverStatus validates driver status and returns ValidationError if invalid
func ValidateDriverStatus(status string) *ValidationError {
	if !IsValidDriverStatus(status) {
		return &ValidationError{
			Field:   "driverStatus",
			Message: "invalid driver status",
			Value:   status,
		}
	}
	return nil
}

// ValidateDriverApprovalStatus validates driver approval status and returns ValidationError if invalid
func ValidateDriverApprovalStatus(status string) *ValidationError {
	if !IsValidDriverApprovalStatus(status) {
		return &ValidationError{
			Field:   "driverApprovalStatus",
			Message: "invalid driver approval status",
			Value:   status,
		}
	}
	return nil
}

// ValidateVehicleClass validates vehicle class and returns ValidationError if invalid
func ValidateVehicleClass(class string) *ValidationError {
	if !IsValidVehicleClass(class) {
		return &ValidationError{
			Field:   "vehicleClass",
			Message: "invalid vehicle class",
			Value:   class,
		}
	}
	return nil
}

// ValidatePaymentMethod validates payment method and returns ValidationError if invalid
func ValidatePaymentMethod(method string) *ValidationError {
	if !IsValidPaymentMethod(method) {
		return &ValidationError{
			Field:   "paymentMethod",
			Message: "invalid payment method",
			Value:   method,
		}
	}
	return nil
}

package common

// Service Names
const (
	ServiceTrip     = "trip-service"
	ServiceIdentity = "identity-service"
	ServiceDriver   = "driver-service"
	ServiceRider    = "rider-service"
	ServicePayment  = "payment-service"
	ServiceAdmin    = "admin-service"
)

// HTTP Methods
const (
	HTTPMethodGET    = "GET"
	HTTPMethodPOST   = "POST"
	HTTPMethodPUT    = "PUT"
	HTTPMethodDELETE = "DELETE"
	HTTPMethodPATCH  = "PATCH"
)

// API Routes
const (
	APIVersionV1 = "/v1"
	HealthCheck  = "/health"
	Healthz      = "/healthz"
	API          = "/api"
	Auth         = "/auth"
	OAuth        = "/oauth"
	Rider        = "/rider"
	Driver       = "/driver"
	Trip         = "/trip"
	Payment      = "/payment"
	Admin        = "/admin"
	WebSocket    = "/ws"
)

// Database Field Names
const (
	FieldID        = "id_uuid"
	FieldCreatedAt = "createdat"
	FieldUpdatedAt = "updatedat"
	FieldDeletedAt = "deletedat"
)

// Domain Names
const (
	DomainTrip           = "trip"
	DomainDriver         = "driver"
	DomainRider          = "rider"
	DomainIdentity       = "identity"
	DomainBiddingSession = "biddingSession"
	DomainBooking        = "booking"
	DomainBidding        = "bidding"
	DomainPayment        = "payment"
	DomainAdmin          = "admin"
)

// Environment Variables
const (
	EnvDBPort     = "DB_PORT"
	EnvDBUsername = "DB_USERNAME"
	EnvDBPassword = "DB_PASSWORD"
	EnvDBHost     = "DB_HOST"
	EnvDBName     = "DB_NAME"
	EnvServerPort = "PORT"
	EnvGRPCPort   = "GRPC_PORT"
	EnvDBSSLMODE  = "DB_SSLMODE"
	EnvGinMode    = "GIN_MODE"
	EnvJWTSecret  = "JWT_SECRET"

	// Dapr Configuration
	EnvDaprHTTPPort   = "DAPR_HTTP_PORT"
	EnvDaprGRPCPort   = "DAPR_GRPC_PORT"
	EnvDaprAppID      = "DAPR_APP_ID"
	EnvDaprAppPort    = "DAPR_APP_PORT"
	EnvDaprPubsubName = "DAPR_PUBSUB_NAME"

	// Redis Configuration
	EnvRedisHost     = "REDIS_HOST"
	EnvRedisPort     = "REDIS_PORT"
	EnvRedisPassword = "REDIS_PASSWORD"
	EnvRedisDB       = "REDIS_DB"
)

// WebSocket Message Types
const (
	WSMessageTypeBiddingStarted        = "bidding_started"
	WSMessageTypeBidReceived           = "bid_received"
	WSMessageTypeBiddingEnded          = "bidding_ended"
	WSMessageTypeDriverAssigned        = "driver_assigned"
	WSMessageTypeTimerUpdate           = "timer_update"
	WSMessageTypeError                 = "error"
	WSMessageTypePing                  = "ping"
	WSMessageTypePong                  = "pong"
	WSMessageTypeConnectionEstablished = "connection_established"
	WSMessageTypeSystemMessage         = "system_message"
)

const (
	NotificationTypeNewRideRequest   = "NEW_RIDE_REQUEST"
	NotificationTypeBiddingUpdate    = "BIDDING_UPDATE"
	NotificationTypeDriverAssigned   = "DRIVER_ASSIGNED"
	NotificationTypeTripCancelled    = "TRIP_CANCELLED"
	NotificationTypeTripCompleted    = "TRIP_COMPLETED"
	NotificationTypeTripNotification = "TRIP_NOTIFICATION"
	NotificationTypeTripStatusUpdate = "TRIP_STATUS_UPDATE"
	NotificationTypeDriverLocation   = "DRIVER_LOCATION_UPDATE"
	NotificationTypeNoDriverFound    = "NO_DRIVER_FOUND"
	NotificationTypeBiddingStarted   = "BIDDING_STARTED"
	NotificationTypeBidReceived      = "BID_RECEIVED"
	NotificationTypeBidAccepted      = "BID_ACCEPTED"
	NotificationTypeBidRejected      = "BID_REJECTED"
)

// Aggregate Types
const (
	AggregateTypeTrip           = "Trip"
	AggregateTypeBiddingSession = "BiddingSession"
	AggregateTypeBooking        = "Booking"
	AggregateTypeBidding        = "Bidding"
)

// User Types / Roles
const (
	UserTypeDriver = "driver"
	UserTypeRider  = "rider"
	UserTypeAdmin  = "admin"

	// User Roles (alternative naming)
	UserRoleAdmin  = "admin"
	UserRoleDriver = "driver"
	UserRoleRider  = "rider"
)

// Trip Status - Standardized values matching database enums
const (
	TripStatusDraft              = "DRAFT"
	TripStatusRequested          = "REQUESTED"
	TripStatusSearchingForDriver = "SEARCHING_FOR_DRIVER"
	TripStatusOfferPhase         = "OFFER_PHASE"
	TripStatusAssigned           = "ASSIGNED"
	TripStatusDriverEnRoute      = "DRIVER_EN_ROUTE"
	TripStatusDriverArrived      = "DRIVER_ARRIVED"
	TripStatusInProgress         = "IN_PROGRESS"
	TripStatusCompleted          = "COMPLETED"
	TripStatusCancelled          = "CANCELLED"
	TripStatusFailed             = "FAILED"
)

// Payment Status - Standardized values
const (
	PaymentStatusPending    = "PENDING"
	PaymentStatusProcessing = "PROCESSING"
	PaymentStatusCompleted  = "COMPLETED"
	PaymentStatusFailed     = "FAILED"
	PaymentStatusCancelled  = "CANCELLED"
	PaymentStatusRefunded   = "REFUNDED"
)

// Fare Models / Price Models - Standardized values
const (
	FareModelInstantMatch  = "INSTANT_MATCH"
	FareModelFlexFare      = "FLEX_FARE"
	FareModelAutomaticFare = "AUTOMATIC_FARE"
)

// Bidding Session Status - Standardized values
const (
	BiddingStatusActive    = "active"
	BiddingStatusExpired   = "expired"
	BiddingStatusAssigned  = "assigned"
	BiddingStatusCancelled = "cancelled"
	BiddingStatusCompleted = "completed"
	BiddingStatusFailed    = "failed"
	BiddingStatusPending   = "pending"
	BiddingStatusStarted   = "started"
)

// Location Validation
const (
	MinLatitude  = -90.0
	MaxLatitude  = 90.0
	MinLongitude = -180.0
	MaxLongitude = 180.0
)

// Timeouts and Intervals
const (
	BiddingTimerDuration  = 20 // seconds
	WebSocketPingInterval = 30 // seconds
	WebSocketWriteTimeout = 10 // seconds
	WebSocketReadTimeout  = 10 // seconds
	DefaultRequestTimeout = 30 // seconds
)

// Pub/Sub Topics
const (
	TopicTripEvents          = "trip.events"
	TopicBookingEvents       = "booking.events"
	TopicDriverEvents        = "driver.events"
	TopicIdentityEvents      = "identity.events"
	TopicPaymentEvents       = "payment.events"
	TopicBiddingEvents       = "bidding.events"
	TopicDriverNotifications = "driver.notifications"
	TopicRiderNotifications  = "rider.notifications"
)

// Dapr Components
const (
	DaprPubsubName  = "pubsub"
	DaprStateStore  = "statestore"
	DaprSecretStore = "secretstore"
	DaprConfigStore = "configstore"
)

// Redis Channels
const (
	RedisChannelDriverNotify  = "driver:notify:ch:"
	RedisChannelDriverBidding = "driver:bidding:ch:"
	RedisChannelRiderNotify   = "rider:notify:ch:"
)

// Error Messages
const (
	ErrMsgTripIDRequired      = "trip ID is required"
	ErrMsgUserIDRequired      = "user ID is required"
	ErrMsgDriverIDRequired    = "driver ID is required"
	ErrMsgRiderIDRequired     = "rider ID is required"
	ErrMsgInvalidLatitude     = "invalid latitude: must be between -90 and 90"
	ErrMsgInvalidLongitude    = "invalid longitude: must be between -180 and 180"
	ErrMsgServiceNotAvailable = "service not available"
	ErrMsgInvalidRequest      = "invalid request"
	ErrMsgUnauthorized        = "unauthorized access"
	ErrMsgForbidden           = "access forbidden"
	ErrMsgNotFound            = "not found"
	ErrMsgInternalError       = "internal server error"
)

// Success Messages
const (
	MsgTripCreated           = "Trip created successfully"
	MsgTripAccepted          = "Trip accepted successfully"
	MsgTripCancelled         = "Trip cancelled successfully"
	MsgTripCompleted         = "Trip completed successfully"
	MsgBidSubmitted          = "Bid submitted successfully"
	MsgBiddingStarted        = "Bidding started successfully"
	MsgDriverAssigned        = "Driver assigned successfully"
	MsgConnectionEstablished = "Connection established successfully"
	MsgDriverApproved        = "Driver approved successfully"
	MsgDriverRejected        = "Driver rejected successfully"
	MsgPaymentProcessed      = "Payment processed successfully"
)

// ==================== DRIVER STATUS ====================
// Driver operational status
const (
	DriverStatusOnline    = "online"
	DriverStatusOffline   = "offline"
	DriverStatusAvailable = "available"
	DriverStatusBusy      = "busy"
	DriverStatusInactive  = "inactive"
)

// ==================== DRIVER APPROVAL STATUS ====================
// Driver approval/verification status (for admin management)
const (
	DriverApprovalStatusPendingReview = "pending_review"
	DriverApprovalStatusApproved      = "approved"
	DriverApprovalStatusRejected      = "rejected"
	DriverApprovalStatusSuspended     = "suspended"
)

// ==================== VEHICLE CLASS / CATEGORY ====================
const (
	VehicleClassClassic = "CLASSIC"
	VehicleClassComfort = "COMFORT"
	VehicleClassLuxury  = "LUXURY"
	VehicleClassMinivan = "MINIVAN"
)

// ==================== PAYMENT METHODS ====================
const (
	PaymentMethodCash     = "CASH"
	PaymentMethodWallet   = "WALLET"
	PaymentMethodTelebirr = "TELEBIRR"
	PaymentMethodCard     = "CARD"
	PaymentMethodStripe   = "STRIPE"
	PaymentMethodPayPal   = "PAYPAL"
)

// ==================== COMPLAINT STATUS ====================
// Support ticket/complaint status
const (
	ComplaintStatusOpen       = "open"
	ComplaintStatusAssigned   = "assigned"
	ComplaintStatusInProgress = "in_progress"
	ComplaintStatusResolved   = "resolved"
	ComplaintStatusClosed     = "closed"
)

// ==================== COMPLAINT PRIORITY ====================
const (
	ComplaintPriorityLow    = "low"
	ComplaintPriorityMedium = "medium"
	ComplaintPriorityHigh   = "high"
	ComplaintPriorityUrgent = "urgent"
)

// ==================== ADMIN ROLES ====================
const (
	AdminRoleSuperAdmin = "super_admin"
	AdminRoleAdmin      = "admin"
	AdminRoleOps        = "ops"
	AdminRoleSupport    = "support"
)

// ==================== DOCUMENT TYPES ====================
const (
	DocumentTypeLicense      = "LICENSE"
	DocumentTypeTIN          = "TIN"
	DocumentTypeID           = "ID"
	DocumentTypePassport     = "PASSPORT"
	DocumentTypeInsurance    = "INSURANCE"
	DocumentTypeRegistration = "REGISTRATION"
)

// ==================== ADDRESS TYPES ====================
const (
	AddressTypeHome  = "HOME"
	AddressTypeWork  = "WORK"
	AddressTypeOther = "OTHER"
)

// ==================== PROVIDERS ====================
const (
	ProviderGoogle = "GOOGLE"
	ProviderManual = "MANUAL"
	ProviderPhone  = "PHONE"
)

// ==================== DEVICE TYPES ====================
const (
	DeviceTypeAndroid = "ANDROID"
	DeviceTypeIOS     = "IOS"
	DeviceTypeWeb     = "WEB"
)

// ==================== SUBSCRIPTION TYPES ====================
const (
	SubscriptionTypeNone       = "NONE"
	SubscriptionTypeBasic      = "BASIC"
	SubscriptionTypePremium    = "PREMIUM"
	SubscriptionTypeEnterprise = "ENTERPRISE"
)

// ==================== NEGOTIATION STATUS ====================
// Flex Fare negotiation status
const (
	NegotiationStatusOffered  = "OFFERED"
	NegotiationStatusCounter  = "COUNTERED"
	NegotiationStatusPending  = "PENDING"
	NegotiationStatusAccepted = "ACCEPTED"
	NegotiationStatusRejected = "REJECTED"
)

// ==================== WALLET STATUS ====================
const (
	WalletStatusActive    = "ACTIVE"
	WalletStatusSuspended = "SUSPENDED"
	WalletStatusClosed    = "CLOSED"
	WalletStatusFrozen    = "FROZEN"
)

// ==================== WALLET TRANSACTION TYPES ====================
const (
	WalletTransactionTypeCredit   = "CREDIT"
	WalletTransactionTypeDebit    = "DEBIT"
	WalletTransactionTypeFreeze   = "FREEZE"
	WalletTransactionTypeUnfreeze = "UNFREEZE"
)

// ==================== WALLET TRANSACTION REASONS ====================
const (
	WalletTransactionReasonRidePayment = "ride_payment"
	WalletTransactionReasonRefund      = "refund"
	WalletTransactionReasonBonus       = "bonus"
	WalletTransactionReasonTopUp       = "top_up"
	WalletTransactionReasonWithdrawal  = "withdrawal"
	WalletTransactionReasonPromo       = "promo"
)

// ==================== COMMISSION RATES ====================
const (
	DefaultCommissionPercent      = 15.0 // 15% default commission
	SubscriptionCommissionPercent = 0.0  // 0% for subscription drivers
)

// ==================== CURRENCY ====================
const (
	DefaultCurrency = "ETB"
)

// ==================== EVENT TYPES ====================
// Event types matching Common/events/registry.go (lowercase with dots)
const (
	EventTypeTripCreated    = "trip.created"
	EventTypeTripAccepted   = "trip.accepted"
	EventTypeTripCancelled  = "trip.cancelled"
	EventTypeTripCompleted  = "trip.completed"
	EventTypeTripInProgress = "trip.in_progress"
	EventTypeTripUpdated    = "trip.updated"

	EventTypeBookingCreated   = "booking.created"
	EventTypeBookingConfirmed = "booking.confirmed"
	EventTypeBookingCancelled = "booking.cancelled"

	EventTypeDriverOnline   = "driver.online"
	EventTypeDriverOffline  = "driver.offline"
	EventTypeDriverLocation = "driver.location"
	EventTypeDriverStatus   = "driver.status"

	EventTypeRiderCreated = "rider.created"
	EventTypeRiderUpdated = "rider.updated"

	EventTypePaymentInitiated = "payment.initiated"
	EventTypePaymentCompleted = "payment.completed"
	EventTypePaymentFailed    = "payment.failed"
	EventTypePaymentRefunded  = "payment.refunded"

	EventTypeBiddingStarted = "bidding.started"
	EventTypeBiddingEnded   = "bidding.ended"
	EventTypeBidSubmitted   = "bid.submitted"
	EventTypeBidAccepted    = "bid.accepted"
	EventTypeBidRejected    = "bid.rejected"

	EventTypeUserCreated = "user.created"
	EventTypeUserUpdated = "user.updated"
	EventTypeUserDeleted = "user.deleted"

	EventTypeDriverNotification = "driver.notification"
	EventTypeRiderNotification  = "rider.notification"
)

// ==================== OPEN TELEMETRY ENVIRONMENT VARIABLES ====================
const (
	EnvOTELExporterEndpoint = "OTEL_EXPORTER_OTLP_ENDPOINT"
	EnvOTELSamplingRate     = "OTEL_SAMPLING_RATE"
	EnvEnvironment          = "ENVIRONMENT"
)

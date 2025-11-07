package common

// Service Names
const (
	ServiceTrip     = "trip-service"
	ServiceIdentity = "identity-service"
	ServiceDriver   = "driver-service"
	ServiceRider    = "rider-service"
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

// Event Types
const (
	EventTypeTripCreated           = "TripCreated"
	EventTypeTripUpdated           = "TripUpdated"
	EventTypeTripCancelled         = "TripCancelled"
	EventTypeTripCompleted         = "TripCompleted"
	EventTypeTripOptionsUpdated    = "TripOptionsUpdated"
	EventTypeBiddingSessionStarted = "BiddingSessionStarted"
	EventTypeBidReceived           = "BidReceived"
	EventTypeBiddingSessionEnded   = "BiddingSessionEnded"
	EventTypeDriverAssigned        = "DriverAssigned"
	EventTypeBidSubmitted          = "BidSubmitted"
	EventTypeBidAccepted           = "BidAccepted"
	EventTypeBidRejected           = "BidRejected"
	EventTypeBidCountered          = "BidCountered"
	EventTypeInstantMatched        = "InstantMatched"
)

// Aggregate Types
const (
	AggregateTypeTrip           = "Trip"
	AggregateTypeBiddingSession = "BiddingSession"
	AggregateTypeBooking        = "Booking"
	AggregateTypeBidding        = "Bidding"
)

// User Types
const (
	UserTypeDriver = "driver"
	UserTypeRider  = "rider"
	UserTypeAdmin  = "admin"
)

// Trip Status
const (
	TripStatusUnspecified = "TRIP_STATUS_UNSPECIFIED"
	TripStatusPending     = "TRIP_STATUS_PENDING"
	TripStatusAccepted    = "TRIP_STATUS_ACCEPTED"
	TripStatusInProgress  = "TRIP_STATUS_IN_PROGRESS"
	TripStatusCompleted   = "TRIP_STATUS_COMPLETED"
	TripStatusCancelled   = "TRIP_STATUS_CANCELLED"
)

// Payment Status
const (
	PaymentStatusUnspecified = "PAYMENT_STATUS_UNSPECIFIED"
	PaymentStatusPending     = "PAYMENT_STATUS_PENDING"
	PaymentStatusCompleted   = "PAYMENT_STATUS_COMPLETED"
	PaymentStatusFailed      = "PAYMENT_STATUS_FAILED"
	PaymentStatusRefunded    = "PAYMENT_STATUS_REFUNDED"
)

// Price Models
const (
	PriceModelAutomaticFare = "AUTOMATIC_FARE"
	PriceModelFlexFare      = "FLEX_FARE"
	PriceModelInstantMatch  = "INSTANT_MATCH"
)

// Bidding Session Status
const (
	BiddingStatusActive    = "active"
	BiddingStatusExpired   = "expired"
	BiddingStatusAssigned  = "assigned"
	BiddingStatusCancelled = "cancelled"
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
	TopicDriverNotifications = "driver.notifications"
	TopicRiderNotifications  = "rider.notifications"
	TopicBiddingEvents       = "bidding.events"
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
)

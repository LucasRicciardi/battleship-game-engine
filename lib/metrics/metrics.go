package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Game metrics
	GameStarts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "battleship_game_starts_total",
		Help: "Total number of games started",
	})

	GamesCompleted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "battleship_games_completed_total",
		Help: "Total number of games completed",
	})

	GamesActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "battleship_games_active",
		Help: "Number of currently active games",
	})

	// Shot metrics
	ShotsFired = promauto.NewCounter(prometheus.CounterOpts{
		Name: "battleship_shots_fired_total",
		Help: "Total number of shots fired",
	})

	ShotsHit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "battleship_shots_hit_total",
		Help: "Total number of shots that hit",
	})

	ShotsMiss = promauto.NewCounter(prometheus.CounterOpts{
		Name: "battleship_shots_miss_total",
		Help: "Total number of shots that missed",
	})

	// Ship metrics
	ShipsSunk = promauto.NewCounter(prometheus.CounterOpts{
		Name: "battleship_ships_sunk_total",
		Help: "Total number of ships sunk",
	})

	// Latency metrics
	ShotLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "battleship_shot_latency_seconds",
		Help:    "Latency of shoot operations",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
	})

	StartGameLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "battleship_start_game_latency_seconds",
		Help:    "Latency of start game operations",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
	})

	GetGameStateLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "battleship_get_game_state_latency_seconds",
		Help:    "Latency of get game state operations",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
	})

	// Error metrics
	ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "battleship_errors_total",
		Help: "Total number of errors by error type",
	}, []string{"error_type"})

	// Request metrics
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "battleship_requests_total",
		Help: "Total number of requests by endpoint and method",
	}, []string{"endpoint", "method"})

	RequestsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "battleship_requests_duration_seconds",
		Help:    "Duration of requests by endpoint",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
	}, []string{"endpoint"})
)

// RecordShot records shot metrics
func RecordShot(hit bool) {
	ShotsFired.Inc()
	if hit {
		ShotsHit.Inc()
	} else {
		ShotsMiss.Inc()
	}
}

// RecordShipSunk records ship sunk metrics
func RecordShipSunk() {
	ShipsSunk.Inc()
}

// RecordGameStart records game start metrics
func RecordGameStart() {
	GameStarts.Inc()
	GamesActive.Inc()
}

// RecordGameComplete records game completion metrics
func RecordGameComplete() {
	GamesCompleted.Inc()
	GamesActive.Dec()
}

// RecordError records an error by type
func RecordError(errorType string) {
	ErrorsTotal.WithLabelValues(errorType).Inc()
}

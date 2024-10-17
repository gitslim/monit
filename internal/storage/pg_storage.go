package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGStorage struct {
	db *pgxpool.Pool
}

func (s *PGStorage) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for name, metric := range s.GetAllMetrics() {
		tmp[name] = map[string]interface{}{
			"name":  metric.GetName(),
			"value": metric.GetValue(),
			"type":  metric.GetType(),
		}
	}
	return json.Marshal(tmp)
}

func (s *PGStorage) UnmarshalJSON(data []byte) error {
	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	for name, raw := range temp {
		var metricType struct {
			Type entities.MetricType `json:"type"`
		}
		if err := json.Unmarshal(raw, &metricType); err != nil {
			return err
		}

		switch metricType.Type {
		case entities.Gauge:
			var gauge entities.GaugeMetric
			if err := json.Unmarshal(raw, &gauge); err == nil {
				if err := s.UpdateOrCreateMetric(name, metricType.Type, gauge); err != nil {
					return errs.ErrInternal
				}
			}
		case entities.Counter:
			var counter entities.CounterMetric
			if err := json.Unmarshal(raw, &counter); err == nil {
				if err := s.UpdateOrCreateMetric(name, metricType.Type, counter); err != nil {
					return errs.ErrInternal
				}
			}
		default:
			return errs.ErrInvalidMetricType
		}
	}
	return nil
}

func NewPGStorage(pool *pgxpool.Pool) *PGStorage {
	return &PGStorage{
		db: pool,
	}
}

// UpdateOrCreateMetric обновляет значение метрики, если метрика отстутствует то создает ее (Upsert)
func (s *PGStorage) UpdateOrCreateMetric(name string, metricType entities.MetricType, value interface{}) error {
	ctx := context.Background()

	var query string
	switch metricType {
	case entities.Gauge:
		v, ok := value.(float64)
		if !ok {
			return errs.ErrInvalidMetricValue
		}
		query = `
INSERT INTO metrics (name, type, value)
VALUES ($1, $2, $3)
ON CONFLICT (name, type)
DO UPDATE SET value = EXCLUDED.value`
		_, err := s.db.Exec(ctx, query, name, metricType.String(), v)
		if err != nil {
			return errs.ErrInternal
		}

	case entities.Counter:
		v, ok := value.(int64)
		if !ok {
			return errs.ErrInvalidMetricValue
		}
		query = `
INSERT INTO metrics (name, type, counter)
VALUES ($1, $2, $3)
ON CONFLICT (name, type)
DO UPDATE SET counter = metrics.counter + EXCLUDED.counter`
		_, err := s.db.Exec(ctx, query, name, metricType.String(), v)
		if err != nil {
			//			fmt.Printf("db error: %v", err)
			return errs.ErrInternal
		}

	default:
		return errs.ErrInvalidMetricType
	}

	return nil
}

// GetMetric получает метрику по имени
func (s *PGStorage) GetMetric(mName string, mType string) (entities.Metric, error) {
	ctx := context.Background()

	metricType, err := entities.GetMetricType(mType)
	if err != nil {
		return nil, errs.ErrInvalidMetricType
	}

	switch metricType {
	case entities.Gauge:
		var value float64
		query := `SELECT value FROM metrics WHERE name=$1 AND type=$2`
		err := s.db.QueryRow(ctx, query, mName, mType).Scan(&value)
		if err != nil {
			fmt.Printf("Ошибка выполнения запроса: %v", err)
			return nil, errs.ErrMetricNotFound
		}
		return &entities.GaugeMetric{
			Name:  mName,
			Value: value,
		}, nil

	case entities.Counter:
		var counter int64
		query := `SELECT counter FROM metrics WHERE name=$1 AND type=$2`
		err := s.db.QueryRow(ctx, query, mName, mType).Scan(&counter)
		if err != nil {
			fmt.Printf("Ошибка выполнения запроса: %v", err)
			return nil, errs.ErrMetricNotFound
		}
		return &entities.CounterMetric{
			Name:  mName,
			Value: counter,
		}, nil

	default:
		return nil, errs.ErrInvalidMetricType
	}
}

// GetAllMetrics получает все метрики
func (s *PGStorage) GetAllMetrics() map[string]entities.Metric {
	ctx := context.Background()
	query := `SELECT name, type, value, counter FROM metrics`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		fmt.Printf("Ошибка выполнения запроса: %v", err)
		return nil
	}
	defer rows.Close()

	metrics := make(map[string]entities.Metric)

	for rows.Next() {
		var name, metricTypeStr string
		var value sql.NullFloat64
		var counter sql.NullInt64

		err = rows.Scan(&name, &metricTypeStr, &value, &counter)
		if err != nil {
			fmt.Printf("Ошибка при чтении строки: %v", err)
			continue
		}

		metricType, err := entities.GetMetricType(metricTypeStr)
		if err != nil {
			fmt.Printf("Неизвестный тип метрики: %s\n", metricTypeStr)
			continue
		}

		switch metricType {
		case entities.Gauge:
			if value.Valid {
				metrics[name] = &entities.GaugeMetric{
					Name:  name,
					Value: value.Float64,
				}
			}

		case entities.Counter:
			if counter.Valid {
				metrics[name] = &entities.CounterMetric{
					Name:  name,
					Value: counter.Int64,
				}
			}

		default:
			fmt.Printf("Неизвестный тип метрики: %s\n", metricTypeStr)
		}
	}

	fmt.Printf("METRICS: %v\n", metrics)
	return metrics
}

// Ping проверяет соединение с бд
func (s *PGStorage) Ping() error {
	if err := s.db.Ping(context.TODO()); err != nil {
		return fmt.Errorf("postgres connection error: %w", err)
	}
	return nil
}

func CreateConnPoll(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres database dsn: %w", err)
	}

	// настройка пула
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute

	// Подключение к бд
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed create postgres connection pool: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("postgres connection error: %w", err)
	}
	return pool, nil
}

// CreateMetricsTable создает таблицу для хранения метрик, если она не существует
func CreatePGSchema(ctx context.Context, db *pgxpool.Pool) error {
	query := `
    CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    value DOUBLE PRECISION,
    counter BIGINT,
    CONSTRAINT name_type_unique UNIQUE (name, type)
	)`
	_, err := db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы metrics: %w", err)
	}
	return nil
}

func (s *PGStorage) BatchUpdateOrCreateMetrics(metrics []*entities.MetricDTO) error {
	ctx := context.TODO()

	// Начинаем транзакцию
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// В случае ошибки откатываем транзакцию
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	for _, dto := range metrics {
		var query string
		mType, err := entities.GetMetricType(dto.MType)
		if err != nil {
			fmt.Printf("Bad metric type: %v\n", err)
		}
		switch mType {
		case entities.Gauge:
			query = `
INSERT INTO metrics (name, type, value)
VALUES ($1, $2, $3)
ON CONFLICT (name, type)
DO UPDATE SET value = EXCLUDED.value`
			_, err := tx.Exec(ctx, query, dto.ID, dto.MType, dto.Value)
			if err != nil {
				return errs.ErrInternal
			}

		case entities.Counter:
			query = `
INSERT INTO metrics (name, type, counter)
VALUES ($1, $2, $3)
ON CONFLICT (name, type)
DO UPDATE SET counter = metrics.counter + EXCLUDED.counter`
			_, err := tx.Exec(ctx, query, dto.ID, dto.MType, dto.Delta)
			if err != nil {
				return errs.ErrInternal
			}

		default:
			return errs.ErrInvalidMetricType
		}

	}
	// Коммитим транзакцию при успешном завершении
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

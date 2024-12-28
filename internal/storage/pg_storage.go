package storage

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gitslim/monit/internal/entities"
	"github.com/gitslim/monit/internal/errs"
	"github.com/gitslim/monit/internal/retry"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/*.sql
var sqlFS embed.FS

// Создаем переменные для хранения SQL-запросов
var (
	UpsertGaugeQuery   string
	UpsertCounterQuery string
	GetGaugeQuery      string
	GetCounterQuery    string
	GetAllMetricsQuery string
)

// PGStorage хранилище для PostgreSQL
type PGStorage struct {
	db *pgxpool.Pool
}

// loadQueries загружает SQL-запросы из файлов и присваивает их переменным
func loadQueries() {
	queries := map[string]*string{
		"upsert_gauge.sql":    &UpsertGaugeQuery,
		"upsert_counter.sql":  &UpsertCounterQuery,
		"get_gauge.sql":       &GetGaugeQuery,
		"get_counter.sql":     &GetCounterQuery,
		"get_all_metrics.sql": &GetAllMetricsQuery,
	}

	for file, qPtr := range queries {
		data, err := sqlFS.ReadFile(filepath.Join("sql", file))
		if err != nil {
			log.Fatalf("Ошибка загрузки SQL-запроса из файла %s: %v", file, err)
		}
		*qPtr = string(data)
	}
}

func init() {
	// загружаем SQL-запросы
	loadQueries()
}

// MarshalJSON возвращает JSON-объект с метриками из хранилища
func (s *PGStorage) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})

	metrics, err := s.GetAllMetrics()
	if err != nil {
		return nil, err
	}

	for name, metric := range metrics {
		tmp[name] = map[string]interface{}{
			"name":  metric.GetName(),
			"value": metric.GetValue(),
			"type":  metric.GetType(),
		}
	}

	return json.Marshal(tmp)
}

// UnmarshalJSON обновляет значения метрик из JSON-объекта и сохраняет их в хранилище
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

// NewPGStorage возвращает экземпляр хранилища с подключением к БД
func NewPGStorage(pool *pgxpool.Pool) *PGStorage {
	return &PGStorage{
		db: pool,
	}
}

// UpdateOrCreateMetric обновляет значение метрики, если метрика отстутствует то создает ее (Upsert)
func (s *PGStorage) UpdateOrCreateMetric(name string, metricType entities.MetricType, value interface{}) error {
	return retry.Retry(func() error {
		ctx := context.Background()

		switch metricType {
		case entities.Gauge:
			v, ok := value.(float64)
			if !ok {
				return errs.ErrInvalidMetricValue
			}
			_, err := s.db.Exec(ctx, UpsertGaugeQuery, name, metricType.String(), v)
			if err != nil {
				return errs.ErrInternal
			}

		case entities.Counter:
			v, ok := value.(int64)
			if !ok {
				return errs.ErrInvalidMetricValue
			}
			_, err := s.db.Exec(ctx, UpsertCounterQuery, name, metricType.String(), v)
			if err != nil {
				//			fmt.Printf("db error: %v", err)
				return errs.ErrInternal
			}

		default:
			return errs.ErrInvalidMetricType
		}

		return nil
	}, 3)
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
		err := s.db.QueryRow(ctx, GetGaugeQuery, mName, mType).Scan(&value)
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
		err := s.db.QueryRow(ctx, GetCounterQuery, mName, mType).Scan(&counter)
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
func (s *PGStorage) GetAllMetrics() (map[string]entities.Metric, error) {
	ctx := context.Background()

	rows, err := s.db.Query(ctx, GetAllMetricsQuery)
	if err != nil {
		return nil, err
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

	return metrics, nil
}

// Ping проверяет соединение с базой данных
func (s *PGStorage) Ping() error {
	if err := s.db.Ping(context.TODO()); err != nil {
		return fmt.Errorf("postgres connection error: %w", err)
	}
	return nil
}

// CreateConnPool создает пул соединений с базой данных
func CreateConnPool(dsn string) (*pgxpool.Pool, error) {
	// По дефолту запросы подготавливаются и кэшируются: default_query_exec_mode=cache_statement
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

// CreatePGSchema создает таблицу metrics в базе данных
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

// BatchUpdateOrCreateMetrics обновляет метрики в базе данных или создает их, если они не существуют
func (s *PGStorage) BatchUpdateOrCreateMetrics(metrics []*entities.MetricDTO) error {
	return retry.Retry(func() error {
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
			mType, err := entities.GetMetricType(dto.MType)
			if err != nil {
				fmt.Printf("Bad metric type: %v\n", err)
			}
			switch mType {
			case entities.Gauge:
				_, err := tx.Exec(ctx, UpsertGaugeQuery, dto.ID, dto.MType, dto.Value)
				if err != nil {
					return errs.ErrInternal
				}

			case entities.Counter:
				_, err := tx.Exec(ctx, UpsertCounterQuery, dto.ID, dto.MType, dto.Delta)
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
	}, 3)
}

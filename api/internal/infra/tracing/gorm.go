package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	gormSpanKey        = "otel:span"
	callBackBeforeName = "otel:before"
	callBackAfterName  = "otel:after"
)

// GormPlugin is a GORM plugin for OpenTelemetry tracing
type GormPlugin struct {
	tracer trace.Tracer
}

// NewGormPlugin creates a new GORM tracing plugin
func NewGormPlugin(serviceName string) *GormPlugin {
	return &GormPlugin{
		tracer: otel.Tracer(serviceName + "/gorm"),
	}
}

// Name returns the plugin name
func (p *GormPlugin) Name() string {
	return "otel-gorm"
}

// Initialize initializes the plugin
func (p *GormPlugin) Initialize(db *gorm.DB) error {
	// Register callbacks for all operations
	cb := db.Callback()

	// Create
	if err := cb.Create().Before("gorm:create").Register(callBackBeforeName, p.before("gorm.create")); err != nil {
		return err
	}
	if err := cb.Create().After("gorm:create").Register(callBackAfterName, p.after()); err != nil {
		return err
	}

	// Query
	if err := cb.Query().Before("gorm:query").Register(callBackBeforeName, p.before("gorm.query")); err != nil {
		return err
	}
	if err := cb.Query().After("gorm:query").Register(callBackAfterName, p.after()); err != nil {
		return err
	}

	// Update
	if err := cb.Update().Before("gorm:update").Register(callBackBeforeName, p.before("gorm.update")); err != nil {
		return err
	}
	if err := cb.Update().After("gorm:update").Register(callBackAfterName, p.after()); err != nil {
		return err
	}

	// Delete
	if err := cb.Delete().Before("gorm:delete").Register(callBackBeforeName, p.before("gorm.delete")); err != nil {
		return err
	}
	if err := cb.Delete().After("gorm:delete").Register(callBackAfterName, p.after()); err != nil {
		return err
	}

	// Row
	if err := cb.Row().Before("gorm:row").Register(callBackBeforeName, p.before("gorm.row")); err != nil {
		return err
	}
	if err := cb.Row().After("gorm:row").Register(callBackAfterName, p.after()); err != nil {
		return err
	}

	// Raw
	if err := cb.Raw().Before("gorm:raw").Register(callBackBeforeName, p.before("gorm.raw")); err != nil {
		return err
	}
	if err := cb.Raw().After("gorm:raw").Register(callBackAfterName, p.after()); err != nil {
		return err
	}

	return nil
}

func (p *GormPlugin) before(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.Context == nil {
			return
		}

		ctx, span := p.tracer.Start(db.Statement.Context, operation,
			trace.WithSpanKind(trace.SpanKindClient),
		)

		// Store span in context
		db.Statement.Context = ctx
		db.InstanceSet(gormSpanKey, span)
	}
}

func (p *GormPlugin) after() func(*gorm.DB) {
	return func(db *gorm.DB) {
		v, ok := db.InstanceGet(gormSpanKey)
		if !ok {
			return
		}

		span, ok := v.(trace.Span)
		if !ok {
			return
		}
		defer span.End()

		// Add attributes
		span.SetAttributes(
			attribute.String("db.system", "gorm"),
			attribute.String("db.statement", db.Statement.SQL.String()),
			attribute.Int64("db.rows_affected", db.Statement.RowsAffected),
		)

		if db.Statement.Table != "" {
			span.SetAttributes(attribute.String("db.table", db.Statement.Table))
		}

		// Record error if any
		if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
			span.RecordError(db.Error)
			span.SetStatus(codes.Error, db.Error.Error())
		}
	}
}

// WithTracing adds tracing to a GORM DB instance
func WithTracing(db *gorm.DB, serviceName string) error {
	return db.Use(NewGormPlugin(serviceName))
}

package postgresql

import (
	"github.com/lib/pq"

	"github.com/benchttp/worker/benchttp"
)

// nolint:gocognit
func (s InsertionService) Insert(stats benchttp.Stats) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	insertStatsDescriptor, err := tx.Prepare(`
	INSERT INTO public.stats_descriptor(
		id,
		user_id,
		finished_at)
		VALUES($1, $2, $3)`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	defer insertStatsDescriptor.Close()

	if _, err = insertStatsDescriptor.Exec(
		stats.Descriptor.ID,
		stats.Descriptor.UserID,
		stats.Descriptor.FinishedAt,
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	insertTimestats, err := tx.Prepare(`
	INSERT INTO public.timestats(
		stats_descriptor_id, 
		min, 
		max, 
		mean, 
		median, 
		standard_deviation, 
		deciles) 
	VALUES($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	defer insertTimestats.Close()

	if _, err = insertTimestats.Exec(
		stats.Descriptor.ID,
		stats.Time.Min,
		stats.Time.Max,
		stats.Time.Mean,
		stats.Time.Median,
		stats.Time.StdDev,
		pq.Array(stats.Time.Deciles),
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	insertCodestats, err := tx.Prepare(`
	INSERT INTO codestats(
		stats_descriptor_id,
		code_1xx,
		code_2xx,
		code_3xx,
		code_4xx,
		code_5xx
	) VALUES($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	defer insertCodestats.Close()

	if _, err = insertCodestats.Exec(
		stats.Descriptor.ID,
		stats.Code.Status1xx,
		stats.Code.Status2xx,
		stats.Code.Status3xx,
		stats.Code.Status4xx,
		stats.Code.Status5xx,
	); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return nil
}

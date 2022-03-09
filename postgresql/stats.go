package postgresql

import (
	"github.com/lib/pq"

	"github.com/benchttp/worker/stats"
)

// nolint:gocognit
func (s StatsService) Create(statsToStore stats.Stats, statsID, userID, tag string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	insertStatsDescriptor, err := tx.Prepare(`
	INSERT INTO public.stats_descriptor(
		id,
		user_id,
		tag,
		finished_at)
		VALUES($1, $2, $3, '2022-03-03 17:36:38-02')`)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return ErrExecutingRollback
		}
		return ErrPreparingStmt
	}
	defer insertStatsDescriptor.Close()

	if _, err = insertStatsDescriptor.Exec(statsID, userID, tag); err != nil {
		err = tx.Rollback()
		if err != nil {
			return ErrExecutingRollback
		}
		return ErrExecutingPreparedStmt
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
		err = tx.Rollback()
		if err != nil {
			return ErrExecutingRollback
		}
		return ErrPreparingStmt
	}
	defer insertTimestats.Close()

	if _, err = insertTimestats.Exec(
		statsID,
		statsToStore.Min,
		statsToStore.Max,
		statsToStore.Mean,
		statsToStore.Median,
		statsToStore.StdDev,
		pq.Array(statsToStore.Deciles)); err != nil {
		err = tx.Rollback()
		if err != nil {
			return ErrExecutingRollback
		}
		return ErrExecutingPreparedStmt
	}

	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return ErrExecutingRollback
		}
		return err
	}
	return nil
}

package postgresql

import (
	"github.com/lib/pq"

	"github.com/benchttp/worker/stats"
)

// nolint:gocognit
func (s InsertionService) Insert(statsToStore stats.Stats, statsID, userID, tag string) error {
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
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	defer insertStatsDescriptor.Close()

	if _, err = insertStatsDescriptor.Exec(
		statsID,
		userID,
		tag,
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
		statsID,
		statsToStore.Min,
		statsToStore.Max,
		statsToStore.Mean,
		statsToStore.Median,
		statsToStore.StdDev,
		pq.Array(statsToStore.Deciles),
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

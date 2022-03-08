package postgresql

import (
	"github.com/benchttp/worker/stats"
)

func (s StatsService) Create(stats stats.Stats, statsId, userId, tag string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	insertIntoStatsDescriptor, err := tx.Prepare(`INSERT INTO public.stats_descriptor(id, user_id, tag, finished_at) VALUES($1, $2, $3, '2022-03-03 17:36:38-02')`)
	if err != nil {
		tx.Rollback()
		return ErrPreparingStmt
	}
	defer insertIntoStatsDescriptor.Close()

	if _, err = insertIntoStatsDescriptor.Exec(statsId, userId, tag); err != nil {
		tx.Rollback()
		return ErrExecutingPreparedStmt
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

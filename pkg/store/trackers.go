package store

func (s *Store) GetTrackers() ([]Tracker, error) {
	rows, err := s.DB.Query("SELECT id, address FROM trackers")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var trackers []Tracker
	for rows.Next() {
		var tracker Tracker
		err = rows.Scan(&tracker.ID, &tracker.Address)
		if err != nil {
			return nil, err
		}
		trackers = append(trackers, tracker)
	}

	return trackers, nil
}

func (s *Store) IsTracking(address string) (bool, error) {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM trackers WHERE address = $1", address).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) SaveTracker(tracker Tracker) error {
	_, err := s.DB.Exec("INSERT INTO trackers (address) VALUES ($1)", tracker.Address)
	if err != nil {
		return err
	}
	return nil
}

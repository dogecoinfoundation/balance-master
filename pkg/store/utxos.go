package store

import "fmt"

func (s *Store) GetUtxos(address string) ([]UTXO, error) {
	rows, err := s.DB.Query("SELECT tx_id, vout, address, amount FROM utxos WHERE address = $1", address)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var utxos []UTXO
	for rows.Next() {
		var utxo UTXO
		err = rows.Scan(&utxo.TxID, &utxo.VOut, &utxo.Address, &utxo.Amount)
		if err != nil {
			return nil, err
		}
		utxos = append(utxos, utxo)
	}

	return utxos, nil
}

func (s *Store) UpdateUtxos(inputs []InputRef, outputs []UTXO) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove spent inputs
	for _, in := range inputs {
		_, err := tx.Exec(
			`DELETE FROM utxos WHERE tx_id = $1 AND vout = $2`,
			in.TxID, in.VOut,
		)
		if err != nil {
			return fmt.Errorf("failed to delete input (%s:%d): %w", in.TxID, in.VOut, err)
		}
	}

	// Insert new outputs as UTXOs
	for _, out := range outputs {
		_, err := tx.Exec(
			`INSERT INTO utxos (tx_id, vout, address, amount) VALUES ($1, $2, $3, $4)`,
			out.TxID, out.VOut, out.Address, out.Amount,
		)
		if err != nil {
			return fmt.Errorf("failed to insert output (%s:%d): %w", out.TxID, out.VOut, err)
		}
	}

	return tx.Commit()
}

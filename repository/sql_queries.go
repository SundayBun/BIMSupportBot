package repository

const (
	findById = `SELECT * 
				FROM account 
				WHERE id = $1`
)

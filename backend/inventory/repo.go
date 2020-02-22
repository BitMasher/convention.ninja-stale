package inventory

import "database/sql"

type Repo struct {
	*sql.DB
}


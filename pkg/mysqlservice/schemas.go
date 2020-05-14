package mysqlservice

var schemas map[string]string

func init() {
	schemas = map[string]string{
		"projects": `
			id VARCHAR(128) UNIQUE NOT NULL,
			title VARCHAR(100) UNIQUE NOT NULL,
			description TEXT NOT NULL,
			created DATETIME NOT NULL,
			updated DATETIME NOT NULL
		`,
	}
}

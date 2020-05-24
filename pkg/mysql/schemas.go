package mysql

var schemas map[string]string

func init() {
	schemas = map[string]string{
		"projects": `
			(nu MEDIUMINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			id VARCHAR(50) UNIQUE NOT NULL,
			title VARCHAR(100) UNIQUE NOT NULL,
			description TEXT NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated TIMESTAMP  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)
			AUTO_INCREMENT=1;
		`,
	}
}

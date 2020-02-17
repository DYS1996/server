package pgsql

// PGConfig contains necessary of a PGSQL
type PGConfig struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DBName string
}

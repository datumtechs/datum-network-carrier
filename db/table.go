package db

type table struct {
	db     Database
	prefix string
}

// NewTable returns a Database object that prefixes all keys with a given string.
func NewTable(db Database, prefix string) Database {
	return &table{
		db:     db,
		prefix: prefix,
	}
}

func (t *table) Put(key []byte, value []byte) error {
	return t.db.Put(append([]byte(t.prefix), key...), value)
}

func (t *table) Has(key []byte) (bool, error) {
	return t.db.Has(append([]byte(t.prefix), key...))
}

func (t *table) Get(key []byte) ([]byte, error) {
	return t.db.Get(append([]byte(t.prefix), key...))
}

func (t *table) Delete(key []byte) error {
	return t.db.Delete(append([]byte(t.prefix), key...))
}

func (t *table) Close() {
	//
}

func (t *table) NewIteratorWithPrefixAndStart(prefix []byte, start []byte) Iterator {
	return nil
}

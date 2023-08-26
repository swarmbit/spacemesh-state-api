package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/cosmos/btcutil/bech32"
	sqlite "github.com/go-llsqlite/llsqlite"
	"github.com/go-llsqlite/llsqlite/sqlitex"
)

var (
	// ErrNoConnection is returned if pooled connection is not available.
	ErrNoConnection = errors.New("database: no free connection")
	// ErrNotFound is returned if requested record is not found.
	ErrNotFound = errors.New("database: not found")
	// ErrObjectExists is returned if database constraints didn't allow to insert an object.
	ErrObjectExists = errors.New("database: object exists")
)

const (
	beginDefault   = "BEGIN;"
	beginImmediate = "BEGIN IMMEDIATE;"
)

// Executor is an interface for executing raw statement.
type Executor interface {
	Exec(string, Encoder, Decoder) (int, error)
}

// Statement is an sqlite statement.
type Statement = sqlite.Stmt

// Encoder for parameters.
// Both positional parameters:
// select block from blocks where id = ?1;
//
// and named parameters are supported:
// select blocks from blocks where id = @id;
//
// For complete information see https://www.sqlite.org/c3ref/bind_blob.html.
type Encoder func(*Statement)

// Decoder for sqlite rows.
type Decoder func(*Statement) bool

func defaultConf() *conf {
	return &conf{
		connections: 16,
	}
}

type conf struct {
	flags       sqlite.OpenFlags
	connections int
}

// WithConnections overwrites number of pooled connections.
func WithConnections(n int) Opt {
	return func(c *conf) {
		c.connections = n
	}
}

// Opt for configuring database.
type Opt func(c *conf)

// InMemory database for testing.
func InMemory(opts ...Opt) *Database {
	opts = append(opts, WithConnections(1))
	db, err := Open("file::memory:?mode=memory", opts...)
	if err != nil {
		panic(err)
	}
	return db
}

// Open database with options.
//
// Database is opened in WAL mode and pragma synchronous=normal.
// https://sqlite.org/wal.html
// https://www.sqlite.org/pragma.html#pragma_synchronous
func Open(uri string, opts ...Opt) (*Database, error) {
	config := defaultConf()
	for _, opt := range opts {
		opt(config)
	}
	pool, err := sqlitex.Open(uri, config.flags, config.connections)
	if err != nil {
		return nil, fmt.Errorf("open db %s: %w", uri, err)
	}
	db := &Database{pool: pool}
	for i := 0; i < config.connections; i++ {
		conn := pool.Get(context.Background())
		if err := registerFunctions(conn); err != nil {
			return nil, err
		}
		defer pool.Put(conn)
	}
	return db, nil
}

// Database is an instance of sqlite database.
type Database struct {
	pool *sqlitex.Pool

	closed   bool
	closeMux sync.Mutex
}

func (db *Database) getTx(ctx context.Context, initstmt string) (*Tx, error) {
	conn := db.pool.Get(ctx)
	if conn == nil {
		return nil, ErrNoConnection
	}
	tx := &Tx{db: db, conn: conn}
	if err := tx.begin(initstmt); err != nil {
		return nil, err
	}
	return tx, nil
}

func (db *Database) withTx(ctx context.Context, initstmt string, exec func(*Tx) error) error {
	tx, err := db.getTx(ctx, initstmt)
	if err != nil {
		return err
	}
	defer tx.Release()
	if err := exec(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// Tx creates deferred sqlite transaction.
//
// Deferred transactions are not started until the first statement.
// Transaction may be started in read mode and automatically upgraded to write mode
// after one of the write statements.
//
// https://www.sqlite.org/lang_transaction.html
func (db *Database) Tx(ctx context.Context) (*Tx, error) {
	return db.getTx(ctx, beginDefault)
}

// WithTx will pass initialized deferred transaction to exec callback.
// Will commit only if error is nil.
func (db *Database) WithTx(ctx context.Context, exec func(*Tx) error) error {
	return db.withTx(ctx, beginImmediate, exec)
}

// TxImmediate creates immediate transaction.
//
// IMMEDIATE cause the database connection to start a new write immediately, without waiting
// for a write statement. The BEGIN IMMEDIATE might fail with SQLITE_BUSY if another write
// transaction is already active on another database connection.
func (db *Database) TxImmediate(ctx context.Context) (*Tx, error) {
	return db.getTx(ctx, beginImmediate)
}

// WithTxImmediate will pass initialized immediate transaction to exec callback.
// Will commit only if error is nil.
func (db *Database) WithTxImmediate(ctx context.Context, exec func(*Tx) error) error {
	return db.withTx(ctx, beginImmediate, exec)
}

// Exec statement using one of the connection from the pool.
//
// If you care about atomicity of the operation (for example writing rewards to multiple accounts)
// Tx should be used. Otherwise sqlite will not guarantee that all side-effects of operations are
// applied to the database if machine crashes.
//
// Note that Exec will block until database is closed or statement has finished.
// If application needs to control statement execution lifetime use one of the transaction.
func (db *Database) Exec(query string, encoder Encoder, decoder Decoder) (int, error) {
	conn := db.pool.Get(context.Background())
	if conn == nil {
		return 0, ErrNoConnection
	}
	defer db.pool.Put(conn)
	return exec(conn, query, encoder, decoder)
}

// Close closes all pooled connections.
func (db *Database) Close() error {
	db.closeMux.Lock()
	defer db.closeMux.Unlock()
	if db.closed {
		return nil
	}
	if err := db.pool.Close(); err != nil {
		return fmt.Errorf("close pool %w", err)
	}
	db.closed = true
	return nil
}

func exec(conn *sqlite.Conn, query string, encoder Encoder, decoder Decoder) (int, error) {
	stmt, err := conn.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("prepare %s: %w", query, err)
	}
	if encoder != nil {
		encoder(stmt)
	}
	defer stmt.ClearBindings()

	rows := 0
	for {
		row, err := stmt.Step()
		if err != nil {
			code := sqlite.ErrCode(err)
			if code == sqlite.SQLITE_CONSTRAINT_PRIMARYKEY {
				return 0, ErrObjectExists
			}
			return 0, fmt.Errorf("step %d: %w", rows, err)
		}
		if !row {
			return rows, nil
		}
		rows++
		// exhaust iterator
		if decoder == nil {
			continue
		}
		if !decoder(stmt) {
			if err := stmt.Reset(); err != nil {
				return rows, fmt.Errorf("statement reset %w", err)
			}
			return rows, nil
		}
	}
}

// Tx is wrapper for database transaction.
type Tx struct {
	db        *Database
	conn      *sqlite.Conn
	committed bool
	err       error
}

func (tx *Tx) begin(initstmt string) error {
	stmt := tx.conn.Prep(initstmt)
	_, err := stmt.Step()
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	return nil
}

// Commit transaction.
func (tx *Tx) Commit() error {
	stmt := tx.conn.Prep("COMMIT;")
	_, tx.err = stmt.Step()
	if tx.err != nil {
		return tx.err
	}
	tx.committed = true
	return nil
}

// Release transaction. Every transaction that was created must be released.
func (tx *Tx) Release() error {
	defer tx.db.pool.Put(tx.conn)
	if tx.committed {
		return nil
	}
	stmt := tx.conn.Prep("ROLLBACK")
	_, tx.err = stmt.Step()
	return tx.err
}

// Exec query.
func (tx *Tx) Exec(query string, encoder Encoder, decoder Decoder) (int, error) {
	return exec(tx.conn, query, encoder, decoder)
}

func registerFunctions(conn *sqlite.Conn) error {
	// sqlite doesn't provide native support for uint64,
	// it is a problem if we want to sort items using actual uint64 value
	// or do arithmetic operations on uint64 in database
	// for that we have to add custom functions, another example https://stackoverflow.com/a/8503318
	if err := conn.CreateFunction("add_uint64", true, 2, func(ctx sqlite.Context, values ...sqlite.Value) {
		ctx.ResultInt64(int64(uint64(values[0].Int64()) + uint64(values[1].Int64())))
	}, nil, nil); err != nil {
		return fmt.Errorf("registering add_uint64: %w", err)
	}
	return nil
}

func load(db Executor, address Address, query string, enc Encoder) (Account, error) {
	var account Account
	_, err := db.Exec(query, enc, func(stmt *Statement) bool {
		account.Balance = uint64(stmt.ColumnInt64(0))
		account.NextNonce = uint64(stmt.ColumnInt64(1))
		account.Layer = LayerID(uint32(stmt.ColumnInt64(2)))
		if stmt.ColumnLen(3) > 0 {
			account.TemplateAddress = &Address{}
			stmt.ColumnBytes(3, account.TemplateAddress[:])
			account.State = make([]byte, stmt.ColumnLen(4))
			stmt.ColumnBytes(4, account.State)
		}
		return false
	})
	if err != nil {
		return Account{}, err
	}
	account.Address = address
	return account, nil
}

type Account struct {
	Layer           LayerID
	Address         Address
	NextNonce       uint64
	Balance         uint64
	TemplateAddress *Address
	State           []byte `scale:"max=10000"`
}

type LayerID uint32
type Address [AddressLength]byte

const (
	// AddressLength is the expected length of the address.
	AddressLength = 24
	// AddressReservedSpace define how much bytes from top is reserved in address for future.
	AddressReservedSpace = 4
)

func Latest(db Executor, address Address) (Account, error) {
	account, err := load(db, address, "select balance, next_nonce, layer_updated, template, state from accounts where address = ?1;", func(stmt *Statement) {
		stmt.BindBytes(1, address.Bytes())
	})
	if err != nil {
		return Account{}, fmt.Errorf("failed to load %v: %w", address, err)
	}
	return account, nil
}

func (a Address) Bytes() []byte { return a[:] }

var (
	// ErrWrongAddressLength is returned when the length of the address is not correct.
	ErrWrongAddressLength = errors.New("wrong address length")
	// ErrUnsupportedNetwork is returned when a network is not supported.
	ErrUnsupportedNetwork = errors.New("unsupported network")
	// ErrDecodeBech32 is returned when an error occurs during decoding bech32.
	ErrDecodeBech32 = errors.New("error decoding bech32")
	// ErrMissingReservedSpace is returned if top bytes of address is not 0.
	ErrMissingReservedSpace = errors.New("missing reserved space")
)

// Config is the configuration of the address package.
var networkHrp = "sm"

func SetNetworkHRP(update string) {
	networkHrp = update
}

func NetworkHRP() string {
	return networkHrp
}

// StringToAddress returns a new Address from a given string like `sm1abc...`.
func StringToAddress(src string) (Address, error) {
	var addr Address
	hrp, data, err := bech32.DecodeNoLimit(src)
	if err != nil {
		return addr, fmt.Errorf("%w: %w", ErrDecodeBech32, err)
	}

	// for encoding bech32 uses slice of 5-bit unsigned integers. convert it back it 8-bit uints.
	dataConverted, err := bech32.ConvertBits(data, 5, 8, true)
	if err != nil {
		return addr, fmt.Errorf("error converting bech32 bits: %w", err)
	}

	// AddressLength+1 cause ConvertBits append empty byte to the end of the slice.
	if len(dataConverted) != AddressLength+1 {
		return addr, fmt.Errorf("expected %d bytes, got %d: %w", AddressLength, len(data), ErrWrongAddressLength)
	}
	if networkHrp != hrp {
		return addr, fmt.Errorf("wrong network id: expected `%s`, got `%s`: %w", NetworkHRP(), hrp, ErrUnsupportedNetwork)
	}
	// check that first 4 bytes are 0.
	for i := 0; i < AddressReservedSpace; i++ {
		if dataConverted[i] != 0 {
			return addr, fmt.Errorf("expected first %d bytes to be 0, got %d: %w", AddressReservedSpace, dataConverted[i], ErrMissingReservedSpace)
		}
	}

	copy(addr[:], dataConverted[:])
	return addr, nil
}

func main() {
	db, err := Open("file:/Users/brunovale/Dev/git/spacemesh/spacemesh-configs/docker/node/node-data/state.sql",
		WithConnections(5),
	)

	if err != nil {
		fmt.Print("Failed to open db")
	}

	address, err := StringToAddress("sm1qqqqqqqvsm9kslrgs4j6kf88hzr9qeqdk63f4gczaqnu4")
	if err != nil {
		fmt.Print("Failed to convert string to address")
	}

	account, err := Latest(db, address)
	if err != nil {
		fmt.Print("Failed to get account")
	}
	fmt.Printf("Account balance: %d\n", account.Balance)

}

package telnetbbs

import (
	"fmt"
	"net"
	"time"

	"github.com/iammegalith/tvbbs/pkg/database"
	"github.com/iammegalith/tvbbs/pkg/users"
)

type TelnetBBS struct {
	address string
	port    string
	ln      net.Listener
	db      *database.Database
}

func New(address, port, dbPath string) (*TelnetBBS, error) {
	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the database: %v", err)
	}

	return &TelnetBBS{
		address: address,
		port:    port,
		db:      db,
	}, nil
}

func (t *TelnetBBS) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", t.address, t.port))
	if err != nil {
		return err
	}
	t.ln = ln

	fmt.Printf("Telnet BBS started on %s:%s\n", t.address, t.port)

	for {
		conn, err := t.ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go t.HandleConnection(conn)
	}
}

func (t *TelnetBBS) HandleConnection(conn net.Conn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	_, err := users.Login(conn, t.db.conn)
	if err != nil {
		fmt.Printf("User login or registration failed: %v\n", err)
		return
	}

	// Continue with your main application logic here
}

func (t *TelnetBBS) Stop() {
	if t.ln != nil {
		t.ln.Close()
	}

	if t.db != nil {
		t.db.Close()
	}
}

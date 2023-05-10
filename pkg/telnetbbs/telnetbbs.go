package telnetbbs

import (
	"net"
	"github.com/iammegalith/tvbbs/pkg/commands"
	"github.com/iammegalith/tvbbs/pkg/users"
	"github.com/iammegalith/tvbbs/pkg/utils"
)

type TelnetBBS struct {
	addr string
}

func NewTelnetBBS(addr string) *TelnetBBS {
	return &TelnetBBS{addr: addr}
}

func (bbs *TelnetBBS) Run() error {
	listener, err := net.Listen("tcp", bbs.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go bbs.handleConnection(conn)
	}
}

func (bbs *TelnetBBS) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Show the pre-login text
	utils.ShowText(conn, "ascii/prelog.txt")

	// Perform the user login or registration process
	username, loggedIn := users.Login(conn)
	if !loggedIn {
		utils.ShowText(conn, "ansi/later.txt")
		return
	}

	// Show the main menu and handle user commands
	for {
		ansiEnabled := false // Check whether the user has ANSI enabled

		// Show the main menu
		if ansiEnabled {
			utils.ShowText(conn, "ansi/mainmenu.ans")
		} else {
			utils.ShowText(conn, "ascii/mainmenu.asc")
		}

		// Read a single character for the menu option
		menuOption := utils.ReadMenuOption(conn)

		// Execute the corresponding command or menu based on the menu option
		switch menuOption {
		case 'A', 'a':
			utils.AnsiToggle(username)
		case 'B', 'b':
			utils.Menu(conn, "bulletins")
		case 'C', 'c':
			commands.MultiUserChat(conn)
		case 'D', 'd':
			utils.Menu(conn, "doors")
		case 'G', 'g':
			if utils.GoodBye(conn) {
				return
			}
		case 'I', 'i':
			if ansiEnabled {
				utils.ShowText(conn, "ansi/systeminfo.ans")
			} else {
				utils.ShowText(conn, "ascii/systeminfo.asc")
			}
		default:
			conn.Write([]byte("Invalid option. Please try again.\n"))
		}
	}
}

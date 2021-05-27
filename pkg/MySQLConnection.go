package pkg

import (
	"MySQLHelper/internal"
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

type viaSSHDialer struct {
	client *ssh.Client
}

func (sshDialer *viaSSHDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return sshDialer.client.Dial("tcp", addr)
}

type MySQLConfig struct {
	DbHost string
	DbPass string
	DbUser string
	DbPort string
	DbName string
	UseSSH bool
	SshKeyPath string
	SshHost string
	SshUser string
	SshPort string
}

func (m *MySQLConfig) Connect() (*sql.DB, error) {
	var dialContext string = "tcp"

	if m.UseSSH {
		sshtun, err := internal.SSHClient(&m.SshHost, &m.SshUser, &m.SshKeyPath, &m.SshPort)
		if err != nil {
			return nil, err
		}

		dialContext = "mysql+tcp"
		mysql.RegisterDialContext(dialContext, (&viaSSHDialer{sshtun}).Dial)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true&columnsWithAlias=true", m.DbUser, m.DbPass, dialContext, m.DbHost, m.DbPort, m.DbName))

	return db, err
}

package executor

import (
	// "bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/pingcap/errors"
)

type Executor struct {
	ExecutionPath string

	Addr     string
	User     string
	Password string
	Database string
	Sql      string

	ErrOut io.Writer
}

func NewExecutor(executionPath string, addr string, user string, password string, database string, sql string) (*Executor, error) {
	if len(executionPath) == 0 {
		return nil, nil
	}

	path, err := exec.LookPath(executionPath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	d := new(Executor)
	d.ExecutionPath = path
	d.Addr = addr
	d.User = user
	d.Password = password
	d.Database = database
	d.Sql = sql

	d.ErrOut = os.Stderr

	return d, nil
}

func (d *Executor) SetErrOut(o io.Writer) {
	d.ErrOut = o
}

func (d *Executor) MongoExec(w io.Writer) error {
	args := make([]string, 0, 6)
	args = append(args, d.Addr+"/"+d.Database)
	args = append(args, fmt.Sprintf("--username=%s", d.User))
	args = append(args, fmt.Sprintf("--password=%s", d.Password))
	args = append(args, fmt.Sprintf("--authenticationDatabase=admin"))
	args = append(args, fmt.Sprintf("--quiet"))

	var tmp string = "/tmp/" + strconv.FormatInt(time.Now().UnixNano(), 16)
	var f, err = os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open file error %v\n", errors.ErrorStack(err))
		//os.Exit(1)
	}
	//_,err=f.Write([]byte(d.Sql))
	f.WriteString(d.Sql)

	args = append(args, tmp)

	cmd := exec.Command(d.ExecutionPath, args...)

	cmd.Stderr = d.ErrOut
	cmd.Stdout = w

	return cmd.Run()
}

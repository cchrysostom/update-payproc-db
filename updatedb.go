package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	sourceDbPath := os.Args[1]
	targetDbPath := os.Args[2]

	fmt.Println(sourceDbPath, targetDbPath)

	sourceDb, err := sql.Open("sqlite3", sourceDbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer sourceDb.Close()

	targetDb, err := sql.Open("sqlite3", targetDbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer targetDb.Close()

	selectPayAddressNotForwarded := `SELECT paymentAddress FROM PaymentAddress WHERE forwarded = 1 AND paymentAddress != ''`
	updateTargetForwarded := `UPDATE PaymentAddress SET forwarded = 1 WHERE paymentAddress = ?`

	sourceCountRows, srcCountErr := sourceDb.Query(selectPayAddressNotForwarded)
	if srcCountErr != nil {
		Error.Fatal(srcCountErr)
	}
	defer sourceCountRows.Close()

	for sourceCountRows.Next() {
		var pymtAddr string
		sourceCountRows.Scan(&pymtAddr)
		fmt.Println(updateTargetForwarded, pymtAddr)

		_, tgtErr := targetDb.Exec(updateTargetForwarded, pymtAddr)
		if tgtErr != nil {
			Error.Fatal(tgtErr)
		}
	}
}

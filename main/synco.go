package main

import "github.com/sqs/go-synco/synco"
import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var query = flag.String("query", "after:2012/09/12", "query to limit fetch")
var mbox = flag.String("mbox", "inbox", "name of mail box/folder from which you want to get mail")
var destBox = flag.String("dbox", "", "name of mail box/folder where you want to move mail")
var jobSize = flag.Int("jobsize", 2, "Number of Emails to be processed at a time")

func usage() {
	fmt.Fprintf(os.Stderr, "usage: synco [server] [port] [username] [password]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 4 {
		usage()
	}
	portI, _ := strconv.Atoi(args[1])
	port := uint16(portI)
	server := &synco.IMAPServer{args[0], port}
	acct := &synco.IMAPAccount{args[2], args[3], server}

	mails, err := synco.GetEMails(acct, *query, *mbox, *jobSize)
	if err != nil {
		fmt.Println("Error while Getting mails ", err)
		return
	}
	var uids []uint32
	fmt.Println("Fetched Mails ", len(mails))
	fmt.Println("UID		|	From		|	To		|		Subject		|Body	| HTMLBODY	|GPGBody")
	for _, msg := range mails {
		//PRocess Emails here
		errP := processEmail(msg)
		if errP != nil {
			continue
		}
		//If successfull then append them to be moved to processed
		uids = append(uids, msg.Imap_uid)
	}

	if *destBox != "" {
		err = synco.MoveEmails(acct, *mbox, *destBox, uids, *jobSize)
		if err != nil {
			fmt.Println("Eror while moving ", err)
		}
	}

}

func processEmail(msg synco.MsgData) (err error) {
	fmt.Println("[" + strconv.Itoa(int(msg.Imap_uid)) + "]  |  " + msg.From + "  |  " + msg.To + "  |  " + msg.Subject + "  |  " +
		msg.Body + " | " + msg.HtmlBody + " | " + msg.GpgBody)
	return
}

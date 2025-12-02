package imap

import (
	"fmt"
	"mail-cleaner/internal/config"
	"mail-cleaner/internal/rules"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Client struct {
	config *config.Config
	client *client.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		client: nil,
	}
}

func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.config.IMAPServer, c.config.IMAPPort)
	fmt.Printf("Connecting to IMAP server at %s with user %s\n", addr, c.config.Email)
	client, err := client.DialTLS(addr, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %v", err)
	}

	c.client = client

	if err := c.client.Login(c.config.Email, c.config.Password); err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}

	fmt.Println("Connected and logged in successfully")

	return nil
}

func (c *Client) Disconnect() error {
	if c.client != nil {
		fmt.Println("Disconnecting from IMAP server")
		err := c.client.Logout()
		if err != nil {
			fmt.Printf("Error during logout: %v\n", err)
			return fmt.Errorf("failed to logout: %v", err)
		}
		fmt.Println("Disconnected successfully")
	} else {
		fmt.Println("No active IMAP client to disconnect")
	}

	return nil
}

func (c *Client) ProcessEmails(handler func(*imap.Message) error) error {
	mbox, err := c.client.Select("INBOX", false)
	if err != nil {
		return err
	}

	if mbox.Messages == 0 {
		fmt.Println("No messages in INBOX")
		return nil
	}

	fmt.Printf("Total messages in INBOX: %d\n", mbox.Messages)

	// For UidFetch use range "1:*" (all UIDs)
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, 0)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	// run fetch in a goroutine
	go func() {
		done <- c.client.UidFetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}, messages)
	}()

	for msg := range messages {
		if err := handler(msg); err != nil {
			return err
		}
	}

	return <-done
}

func (c *Client) MarkForDeletion(uid uint32) error {
	// set up flag for deletion using UID
	seqset := new(imap.SeqSet)
	seqset.AddNum(uid)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}

	return c.client.UidStore(seqset, item, flags, nil)
}

func (c *Client) ExpungeMarked() error {
	return c.client.Expunge(nil)
}

func (c *Client) CleanEmails(rules *rules.Rules) error {
	var toDelete []uint32
	processed := 0

	err := c.ProcessEmails(func(msg *imap.Message) error {
		processed++
		if processed%100 == 0 {
			fmt.Printf("Processed %d emails...\n", processed)
		}

		if rules.ShouldDelete(msg) {
			toDelete = append(toDelete, msg.Uid)
			if msg.Envelope != nil && len(msg.Envelope.From) > 0 {
				fmt.Printf("Marking for deletion: %s - %s\n",
					msg.Envelope.From[0].MailboxName+"@"+msg.Envelope.From[0].HostName,
					msg.Envelope.Subject)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("\nTotal emails to delete: %d\n", len(toDelete))

	//mark emails for deletion
	for _, uid := range toDelete {
		if err := c.MarkForDeletion(uid); err != nil {
			return fmt.Errorf("failed to mark UID %d: %v", uid, err)
		}
	}

	fmt.Println("Expunging marked emails...")
	return c.ExpungeMarked()
}

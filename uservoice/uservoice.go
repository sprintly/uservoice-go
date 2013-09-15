package uservoice

type UservoiceApi interface {
	GetTicketByNumber(number int) (UservoiceTicket, error)
	PostNote(ticket_id int, text string) error
	UrlForTicket(number int) string
}

type UservoiceUser struct {
	Id    int    `xml:"id"`
	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type Message struct {
	Id            int           `xml:"id"`
	PlaintextBody string        `xml:"plaintext_body"`
	Sender        UservoiceUser `xml:"sender"`
	Attachments   []Attachment  `xml:"attachments>attachment"`
}

type Attachment struct {
	Url         string `xml:"url"`
	Name        string `xml:"name"`
	SizeInBytes int    `xml:"size_in_bytes"`
	// CreatedAt time.Time
}

type UservoiceTicket struct {
	Id       int       `xml:"id"`
	Number   int       `xml:"ticket_number"`
	Subject  string    `xml:"subject"`
	Messages []Message `xml:"messages>message"`
}

type Metadata struct {
	Query string `xml:"query"`
}

type SearchResponse struct {
	Tickets  []UservoiceTicket `xml:"tickets>ticket"`
	MetaData Metadata          `xml:"response_data"`
}

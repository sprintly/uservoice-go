package uservoice

type MockUservoiceApi struct{}

func (u MockUservoiceApi) GetTicketByNumber(number int) (UservoiceTicket, error) {
	return UservoiceTicket{Subject: "test title"}, nil
}
func (u MockUservoiceApi) PostNote(ticket_id int, text string) error {
	return nil
}
func (u MockUservoiceApi) UrlForTicket(number int) string {
	return ""
}

func NewMockUservoiceApi() UservoiceApi {
	return MockUservoiceApi{}
}

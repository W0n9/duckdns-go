package duckdns

//Domain structure
type Domain struct {
	Name []string
}

//DomainsService structure
type DomainsService struct {
	client  *Client
	domains Domain
}

//SetDomains function
func (s *DomainsService) SetDomains(domains Domain) {
	s.domains = domains
}

//GetDomains function
func (s *DomainsService) GetDomains() Domain {
	return s.domains
}

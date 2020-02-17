package ad

import ldap "gopkg.in/ldap.v3"

func addOU(OUname string, dnOfOU string, adConn *ldap.Conn) error {

	addRequest := ldap.NewAddRequest(dnOfOU, nil)
	addRequest.Attribute("objectClass", []string{"organizationalUNit"})
	addRequest.Attribute("sAMAccountName", []string{OUname})               //login name
	
	err := adConn.Add(addRequest)
	if err != nil {
		return err
	}
	return nil
}

func deleteOU(dnOfOU string, adConn *ldap.Conn) error {
	delRequest := ldap.NewDelRequest(dnOfOU, nil)                        //creates a delete request for the given DN
	err := adConn.Del(delRequest)
	if err != nil {
		return err
	}
	return nil
}

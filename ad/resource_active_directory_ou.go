package ad

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	ldap "gopkg.in/ldap.v3"

	"github.com/hashicorp/terraform/helper/schema"                       //ensures consistency between providers
)

func resourceOU() *schema.Resource {
	return &schema.Resource{
		Create: resourceADouCreate,
		Read:   resourceADouRead,
		Delete: resourceADouDelete,
		Schema: map[string]*schema.Schema{                      //mapping attribute to what it is
			"ou_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,                          //API doesn’t allow you to change the name of an item after it is created			
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
func resourceADouCreate(d *schema.ResourceData, m interface{}) error {      //called at terraform apply
	client := m.(*ldap.Conn)                                            // m is our client to talk to server
	ouName := d.Get("ou_name").(string)
	domain := d.Get("domain").(string)
	var dnOfOU string
	dnOfOU += "OU=" + ouName                                            //object's entire path to the root
	domainArr := strings.Split(domain, ".")
	for _, item := range domainArr {
		dnOfOU += ",dc=" + item                                      //dc =domain-component
	}
	log.Printf("[DEBUG] dnOfOU: %s ", dnOfOU)
	log.Printf("[DEBUG] Adding OU : %s ", ouName)
	err := addOU(ouName, dnOfOU, client)
	if err != nil {
		log.Printf("[ERROR] Error while adding OU: %s ", err)
		return fmt.Errorf("Error while adding OU %s", err)
	}
	log.Printf("[DEBUG] OU Added successfully: %s", ouName)
	d.SetId(domain + "/" + ouName)
	return nil
}

func resourceADouRead(d *schema.ResourceData, m interface{}) error {  
	client := m.(*ldap.Conn)
	ouName := d.Get("ou_name").(string)
	domain := d.Get("domain").(string)
	var dnOfOU string
	dnOfOU += "OU=" + ouName
	domainArr := strings.Split(domain, ".")
	dnOfOU += "dc=" + domainArr[0]                              
	for index, i := range domainArr {
		if index == 0 {
			continue
		}
		dnOfOU += ",dc=" + i
	}
	log.Printf("[DEBUG] dnOfOU: %s ", dnOfOU)
	log.Printf("[DEBUG] Deleting OU : %s ", ouName)

	NewReq := ldap.NewSearchRequest( //represents the search request send to the server
		dnOfOU, // base dnOfOU.
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0,
		false,
		"(&(objectClass=organizationalUNit)(cn="+ouName+"))", //applied filter
		[]string{"dnOfOU", "cn"},
		nil,
	)

	sr, err := client.Search(NewReq)
	if err != nil {
		log.Printf("[ERROR] while seaching OU : %s", err)
		return fmt.Errorf("Error while searching  OU : %s", err)
	}

	fmt.Println("[ERROR] Found " + strconv.Itoa(len(sr.Entries)) + " Entries")
	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))

	}

	if len(sr.Entries) == 0 {
		log.Println("[ERROR] OU not found")
		d.SetId("")
	}
	return nil
}

func resourceADouDelete(d *schema.ResourceData, m interface{}) error {              // m is interface, we have to assert the type as the client type
	log.Println("[ERROR] Finding OU")
	resourceADouRead(d, m)
	if d.Id() == "" {
		log.Println("[ERROR] Cannot find OU in the specified AD")
		return fmt.Errorf("[ERROR] Cannot find OU in the specified AD")
	}
	client := m.(*ldap.Conn)

	ouName := d.Get("ou_name").(string)
	domain := d.Get("domain").(string)

	var dnOfOU string
	dnOfOU += "OU=" + ouName
	domainArr := strings.Split(domain, ".")
	for _, item := range domainArr {
		dnOfOU += ",dc=" + item
	}

	log.Printf("[DEBUG] Name of the DN is : %s ", dnOfOU)
	log.Printf("[DEBUG] Deleting the OU from the AD : %s ", ouName)

	err := deleteOU(dnOfOU, client)
	if err != nil {
		log.Printf("[ERROR] Error while Deleting OU from AD : %s ", err)
		return fmt.Errorf("Error while Deleting OU from AD %s", err)
	}
	log.Printf("[DEBUG] OU deleted from AD successfully: %s", ouName)
	return nil
}

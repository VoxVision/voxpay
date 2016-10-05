package robokassa

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const delim = ":"
const robokasaURL = "https://auth.robokassa.ru/Merchant/Index.aspx"

type Kassa struct {
	MrchLogin string
	MrchPass1 string
	MrchPass2 string
	IsTest    bool
}

func Init(MrchLogin, MrchPass1, MrchPass2 string) *Kassa {
	var k = Kassa{}
	k.MrchLogin = MrchLogin
	k.MrchPass1 = MrchPass1
	k.MrchPass2 = MrchPass2
	return &k
}

func (k *Kassa) PaymentURL(amount, orderID int) string {
	s := "?MerchantLogin=" + k.MrchLogin + "&OutSum=" + strconv.Itoa(amount) + ".00" + "&InvoiceID=" + strconv.Itoa(orderID) + "&SignatureValue=" + k.SignatureValue(amount, orderID)

	if k.IsTest {
		s += "&IsTest=1"
	}

	return robokasaURL + s
}

func (k *Kassa) SignatureValue(amount, orderID int) string {
	//$mrh_login:$out_summ:$inv_id:$mrh_pass1
	s := []string{k.MrchLogin, strconv.Itoa(amount) + ".00", strconv.Itoa(orderID), k.MrchPass1}
	h := md5.New()
	io.WriteString(h, strings.Join(s, delim))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (k *Kassa) ResultCRC(outSum, invID string) string {
	crc := outSum + delim + invID + delim + k.MrchPass2
	h := md5.New()
	io.WriteString(h, crc)
	result := fmt.Sprintf("%x", h.Sum(nil))
	return strings.ToUpper(result)
}

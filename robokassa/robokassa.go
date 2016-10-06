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
	mrchLogin     string
	mrchPass1     string
	mrchPass2     string
	IsTest        bool
	testMrchPass1 string
	testMrchPass2 string
}

func Init(MrchLogin, MrchPass1, MrchPass2, TestMrchPass1, TestMrchPass2 string) *Kassa {
	var k = Kassa{}
	k.mrchLogin = MrchLogin
	k.mrchPass1 = MrchPass1
	k.mrchPass2 = MrchPass2
	k.testMrchPass1 = TestMrchPass1
	k.testMrchPass2 = TestMrchPass2
	return &k
}

func (k *Kassa) PaymentURL(amount, orderID int) string {
	s := "?MerchantLogin=" + k.mrchLogin + "&OutSum=" + strconv.Itoa(amount) + ".00" + "&InvoiceID=" + strconv.Itoa(orderID) + "&SignatureValue=" + k.SignatureValue(amount, orderID)

	if k.IsTest {
		s += "&IsTest=1"
	}

	return robokasaURL + s
}

func (k *Kassa) SignatureValue(amount, orderID int) string {
	//$mrh_login:$out_summ:$inv_id:$mrh_pass1
	s := []string{k.MrchLogin, strconv.Itoa(amount) + ".00", strconv.Itoa(orderID), k.MrchPass1()}
	h := md5.New()
	io.WriteString(h, strings.Join(s, delim))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (k *Kassa) ResultCRC(outSum, invID string) string {
	crc := outSum + delim + invID + delim + k.MrchPass2()
	h := md5.New()
	io.WriteString(h, crc)
	result := fmt.Sprintf("%x", h.Sum(nil))
	return strings.ToUpper(result)
}

func (k *Kassa) MrchPass1() string {
	if k.IsTest {
		return k.testMrchPass1
	}
	return k.mrchPass1
}

func (k *Kassa) MrchPass2() string {
	if k.IsTest {
		return k.testMrchPass2
	}
	return k.mrchPass2
}

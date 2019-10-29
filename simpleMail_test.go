package simpleMail

import (
	"bytes"
	"testing"
//	"io/ioutil"
)


func TestSendMail(t *testing.T){
  // Please rewrite and test the following.
  auth      := LoginAuth("account", "password")
  from      := "from@domain.com"
  recipents := Recipients{  To: []string{"to@domain.com"}}
  mx        := "mx.hostname.com:587"
  
  attachment := MailAttachement{}
  attachment.Filename = "test.txt"

  // test for file to attachment
  //filedata := bytes.NewBufferString("test text file\n") 
  //out  := new(bytes.Buffer)
  //buf ,_ := ioutil.ReadAll(filedata)
  //out.Write(buf)
  
  // text to attachemnt
  out := bytes.NewBuffer([]byte("test text to file\n"))
  attachment.Data = *out
  
  msg := MakeMessage(from, recipents, "subject", "body string", "==delemiter", &attachment)

  if err := SendMail(mx, auth, from, recipents.To, msg); err != nil{
    t.Error(err)
  }

}


# simpleMail
TLS未対応のSMTPサーバ用のメールクライアント    
golang の net/smtp の SendMailはTLS認証が必須のため   


### Usage
詳しい使い方は　、testコードを参照してください。

```go
import (
        "bytes"
	"github.com/ishiokakeiichi/simpleMail"
)
// パスワード設定 (パスワード認証不要の場合は nil)
auth      := simpleMail.LoginAuth("account", "password")
from      := "from@domain.com"
// 送信先設定 (To, Cc, Bcc) それぞれ指定できます
recipents := simpleMail.Recipients{  To: []string{"to@domain.com"}}
mx        := "mx.hostname.com:587"

// 添付ファイル
attachment := simpleMail.MailAttachement{}
attachment.Filename = "test.txt"
out := bytes.NewBuffer([]byte("test text to file\n"))
attachment.Data = *out

// メッセージ作成
msg := simpleMail.MakeMessage(from, recipents, "subject", "body string", "==delemiter", &attachment)

// メール送信
if err := simpleMail.SendMail(mx, auth, from, recipents.To, msg); err != nil{
  t.Error(err)
}

```


### Author
keiichi ishioka


package user

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func SendMail() error {
	form := url.Values{}
	form.Add("subject", "My letter subject")
	form.Add("name", "My Name")
	form.Add("html", "<html><head></head><body><p>My text</p></body></html>")
	form.Add("from", "emgushovs.ru")
	form.Add("to", "emgushovs@mail.ru")
	form.Add("to_name", "Name TO")
	form.Add("headers", "[{ 'x-tag': 'my_newsletter_ids' }]")
	form.Add("text", "Text version message")

	req, err := http.NewRequest("POST", "https://api.smtp.bz/v1/smtp/send", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("authorization", "vaJoDLwCpxfBobPlLupVcm0o9Ha8Ji15iziM")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

package apiutils
 
import (
  "net/http"
  "net/url"
  "strings"
  "os"
)

var (
    accountSid = os.Getenv("TWILIO_ACCOUNT_SID")
    authToken = os.Getenv("TWILIO_AUTH_TOKEN")
    fromPhone = os.Getenv("TWILIO_PHONE")
)
 
func SendMMSMessage(toPhone, body, mediaURL string) error {
    // Set initial variables
    urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
    
    // Build out the data for our message
    v := url.Values{}
    v.Set("To", toPhone)
    v.Set("From", fromPhone)
    v.Set("Body", body)
    v.Set("MediaUrl",mediaURL)
    rb := *strings.NewReader(v.Encode())
    
    // Create client
    client := &http.Client{}
    
    req, err := http.NewRequest("POST", urlStr, &rb)
    if err != nil {
        return err
    }
    
    req.SetBasicAuth(accountSid, authToken)
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    
    // Make request
    _, err2 := client.Do(req)
    if err2 != nil {
        return err2
    }
    
    return nil
}

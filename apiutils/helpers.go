package apiutils

import (
    "github.com/tmaiaroto/aegis/lambda"
    
    "encoding/base64"
)


func GetBodyFromEvent(evt *lambda.Event) ([]byte, error) { 

    body := evt.Body
    var bodyByte []byte
    
    if evt.IsBase64Encoded {
        if data, err := base64.StdEncoding.DecodeString(body.(string)); err != nil {
            
            return bodyByte, err
            
        } else {
            bodyByte = data
        }
        
    } else {
    
        bodyByte = []byte(body.(string))
        
    }
    
    return bodyByte, nil
    
}
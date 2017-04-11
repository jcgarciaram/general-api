package apiutils

import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws"
    
    "bytes"
    "time"
)

func SaveToS3(bucket, key string, body []byte) error {
    sess := session.Must(session.NewSession())

    svc := s3.New(sess)

    params := &s3.PutObjectInput{
        Bucket:             aws.String(bucket), // Required
        Key:                aws.String(key),  // Required
        // ACL:                aws.String("ObjectCannedACL"),
        Body:               bytes.NewReader(body),
        // CacheControl:       aws.String("CacheControl"),
        // ContentDisposition: aws.String("ContentDisposition"),
        // ContentEncoding:    aws.String("ContentEncoding"),
        // ContentLanguage:    aws.String("ContentLanguage"),
        // ContentLength:      aws.Int64(1),
        ContentType:        aws.String("image/png"),
        Expires:            aws.Time(time.Now().Add(time.Hour)),
        // GrantFullControl:   aws.String("GrantFullControl"),
        // GrantRead:          aws.String("GrantRead"),
        // GrantReadACP:       aws.String("GrantReadACP"),
        // GrantWriteACP:      aws.String("GrantWriteACP"),
        // Metadata: map[string]*string{
            // "Key": aws.String("MetadataValue"), // Required
            // More values...
        // },
        // RequestPayer:            aws.String("RequestPayer"),
        // SSECustomerAlgorithm:    aws.String("SSECustomerAlgorithm"),
        // SSECustomerKey:          aws.String("SSECustomerKey"),
        // SSECustomerKeyMD5:       aws.String("SSECustomerKeyMD5"),
        // SSEKMSKeyId:             aws.String("SSEKMSKeyId"),
        // ServerSideEncryption:    aws.String("ServerSideEncryption"),
        // StorageClass:            aws.String("StorageClass"),
        // Tagging:                 aws.String("TaggingHeader"),
        // WebsiteRedirectLocation: aws.String("WebsiteRedirectLocation"),
    }
    
    _, err := svc.PutObject(params)

    return err
}

// GetDownloadURL generates a pre-signed URL that enables you to temporarily share a file without making it public. Anyone with access to the URL can view the file. 
func GetDownloadURL(bucket, key string) (string, error) {
    svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
    req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    
    return req.Presign(15 * time.Minute)
}


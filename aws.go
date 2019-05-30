package minioanalytics

import "time"

const (
	S3ObjectAccessedGet  = "s3:ObjectAccessed:Get"
	S3ObjectAccessedHead = "s3:ObjectAccessed:Head"
)

type AWSRecord struct {
	EventVersion string    `json:"eventVersion"`
	EventSource  string    `json:"eventSource"`
	AwsRegion    string    `json:"awsRegion"`
	EventTime    time.Time `json:"eventTime"`
	EventName    string    `json:"eventName"`
	UserIdentity struct {
		PrincipalID string `json:"principalId"`
	} `json:"userIdentity"`
	RequestParameters struct {
		AccessKey       string `json:"accessKey"`
		Region          string `json:"region"`
		SourceIPAddress string `json:"sourceIPAddress"`
	} `json:"requestParameters"`
	ResponseElements struct {
		XAmzRequestID        string `json:"x-amz-request-id"`
		XAmzID2              string `json:"x-amz-id-2"`
		ContentLength        string `json:"content-length"`
		XMinioDeploymentID   string `json:"x-minio-deployment-id"`
		XMinioOriginEndpoint string `json:"x-minio-origin-endpoint"`
	} `json:"responseElements"`
	S3               S3     `json:"s3"`
	Source           Source `json:"source"`
	GlacierEventData struct {
		RestoreEventData struct {
			LifecycleRestorationExpiryTime time.Time `json:"lifecycleRestorationExpiryTime"`
			LifecycleRestoreStorageClass   string    `json:"lifecycleRestoreStorageClass"`
		} `json:"restoreEventData"`
	} `json:"glacierEventData"`
}

type Source struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	UserAgent string `json:"userAgent"`
}

type S3 struct {
	S3SchemaVersion string   `json:"s3SchemaVersion"`
	ConfigurationID string   `json:"configurationId"`
	Bucket          S3Bucket `json:"bucket"`
	Object          S3Object `json:"object"`
}

type S3Bucket struct {
	Name          string `json:"name"`
	OwnerIdentity struct {
		PrincipalID string `json:"principalId"`
	} `json:"ownerIdentity"`
	Arn string `json:"arn"`
}

type S3Object struct {
	Key          string `json:"key"`
	Size         int    `json:"size"`
	ETag         string `json:"eTag"`
	ContentType  string `json:"contentType"`
	UserMetadata struct {
		ContentType string `json:"content-type"`
	} `json:"userMetadata"`
	VersionID string `json:"versionId"`
	Sequencer string `json:"sequencer"`
}

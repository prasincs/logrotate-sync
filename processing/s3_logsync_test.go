package processing

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

func Test_splitS3URL(t *testing.T) {
	type args struct {
		rs string
	}
	tests := []struct {
		name       string
		args       args
		wantBucket string
		wantKey    string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBucket, gotKey := splitS3URL(tt.args.rs)
			if gotBucket != tt.wantBucket {
				t.Errorf("splitS3URL() gotBucket = %v, want %v", gotBucket, tt.wantBucket)
			}
			if gotKey != tt.wantKey {
				t.Errorf("splitS3URL() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestNewS3Logsync(t *testing.T) {
	type args struct {
		config Config
		path   string
	}
	tests := []struct {
		name    string
		args    args
		want    *S3Logsync
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewS3Logsync(tt.args.config, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewS3Logsync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewS3Logsync() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3Logsync_remoteFileKey(t *testing.T) {
	type fields struct {
		config   Config
		bucket   string
		prefix   string
		s3Client s3iface.S3API
	}
	type args struct {
		fm *FileMatch
	}
	filePath := "foo/bar/server.log.2017-10-30-04"
	matches := matchFileNames(filePath, defaultLogPattern)
	fm, _ := NewFileMatch(filePath, matches)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "returns the remote file key with the date/time properly set",
			fields: fields{
				config: Config{
					UploadDatePattern: "2006/01/02/15",
					Hostname:          "myhost",
				},
				prefix: "serverlog",
			},
			args: args{fm: fm},
			want: "serverlog/2017/10/30/04/myhost-server.log.2017-10-30-04",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &S3Logsync{
				config:   tt.fields.config,
				bucket:   tt.fields.bucket,
				prefix:   tt.fields.prefix,
				s3Client: tt.fields.s3Client,
			}
			got, err := s.remoteFileKey(tt.args.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("S3Logsync.remoteFileKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("S3Logsync.remoteFileKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3Logsync_Process(t *testing.T) {
	type fields struct {
		config   Config
		bucket   string
		prefix   string
		s3Client s3iface.S3API
	}
	type args struct {
		deleteLocalAfterUpload bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &S3Logsync{
				config:   tt.fields.config,
				bucket:   tt.fields.bucket,
				prefix:   tt.fields.prefix,
				s3Client: tt.fields.s3Client,
			}
			if err := s.Process(tt.args.deleteLocalAfterUpload); (err != nil) != tt.wantErr {
				t.Errorf("S3Logsync.Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

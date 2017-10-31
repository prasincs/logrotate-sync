package processing

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_matchFileNames(t *testing.T) {
	type args struct {
		name         string
		matchPattern *regexp.Regexp
	}
	logPattern := regexp.MustCompile(`(?P<log_type>\S+).log(-|.)(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<extra>\d+).?(?P<compression>gz)?`)
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Matches for files logrotate < 3.9.0 dateformat",
			args: args{
				name:         "test.log-2017-10-31-1509432925.gz",
				matchPattern: logPattern,
			},
			want: map[string]string{"log_type": "test", "year": "2017", "month": "10", "day": "31", "extra": "1509432925", "compression": "gz", "": "-"},
		},
		{
			name: "Matches for files logrotate < 3.9.0 dateformat without compression",
			args: args{
				name:         "test.log-2017-10-31-1509432925",
				matchPattern: logPattern,
			},
			want: map[string]string{"log_type": "test", "year": "2017", "month": "10", "day": "31", "extra": "1509432925", "compression": "", "": "-"},
		},
		{
			name: "Matches kafka hourly server logs",
			args: args{
				name:         "server.log.2017-10-09-09",
				matchPattern: logPattern,
			},
			want: map[string]string{"log_type": "server", "": ".", "year": "2017", "month": "10", "day": "09", "extra": "09", "compression": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchFileNames(tt.args.name, tt.args.matchPattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matchFileNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

var samples = []string{
	"log-cleaner.log.2017-10-11-18", "server.log.2017-10-09-18", "server.log.2017-10-19-23", "server.log.2017-10-30-04", "log-cleaner.log.2017-10-15-16", "server.log.2017-10-11-01", "server.log.2017-10-21-06", "state-change.log.2017-10-07-12",
}

// Uses the defaultLogPattern to avoid extra allocation
func BenchmarkMatchFileNames(b *testing.B) {
	for n := 0; n < b.N; n++ {
		s := samples[n%len(samples)]
		matchFileNames(s, defaultLogPattern)
	}
}

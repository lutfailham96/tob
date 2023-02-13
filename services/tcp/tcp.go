package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/telkomdev/tob/util"
)

// TCP service
type TCP struct {
	url           string
	recovered     bool
	serviceName   string
	lastDownTime  string
	enabled       bool
	verbose       bool
	logger        *log.Logger
	conn          net.Conn
	errConn       error
	checkInterval int
	stopChan      chan bool
}

// TCP's constructor
func NewTCP(verbose bool, logger *log.Logger) *TCP {
	stopChan := make(chan bool, 1)
	return &TCP{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (t *TCP) Name() string {
	return "tcp"
}

// Ping will try to ping the service
func (t *TCP) Ping() []byte {
	// Try to reconnect
	if t.errConn != nil {
		if t.conn != nil {
			t.conn.Close()
		}
		t.Connect()
	}

	if t.conn == nil {
		return []byte("NOT_OK")
	}

	// Try to send data to server
	if _, err := fmt.Fprintf(t.conn, "PING\n"); err != nil {
		t.errConn = errors.New("cannot write ping packet")
		return []byte("NOT_OK")
	}

	scanner := bufio.NewScanner(t.conn)
	for scanner.Scan() {
		if scanner.Text() == "" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (t *TCP) SetURL(url string) {
	t.url = url
}

// Connect to service if needed
func (t *TCP) Connect() error {
	if t.verbose {
		t.logger.Println("connecting to TCP server")
	}

	conn, err := net.Dial("tcp", t.url)
	if err != nil {
		t.errConn = errors.New("cannot connect to TCP server")
		return nil
	}

	if t.verbose {
		t.logger.Println("connecting to TCP server succeed")
	}

	// set connected conn
	t.conn = conn
	t.errConn = nil

	return nil
}

// Close will close the service resources if needed
func (t *TCP) Close() error {
	if t.verbose {
		t.logger.Println("close Web")
	}

	return nil
}

// SetRecover will set recovered status
func (t *TCP) SetRecover(recovered bool) {
	t.recovered = recovered
}

// IsRecover will return recovered status
func (t *TCP) IsRecover() bool {
	return t.recovered
}

// LastDownTime will set last down time of service to current time
func (t *TCP) SetLastDownTimeNow() {
	if t.recovered {
		t.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (t *TCP) GetDownTimeDiff() string {
	return util.TimeDifference(t.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (t *TCP) SetCheckInterval(interval int) {
	t.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (t *TCP) GetCheckInterval() int {
	return t.checkInterval
}

// Enable will set enabled status to service
func (t *TCP) Enable(enabled bool) {
	t.enabled = enabled
}

// IsEnabled will return enable status
func (t *TCP) IsEnabled() bool {
	return t.enabled
}

// Stop will receive stop channel
func (t *TCP) Stop() chan bool {
	return t.stopChan
}

package structs

import (
	"bytes"
	"errors"
	"net"
	"testing"
)

type mockConn struct {
	net.Conn
	writeErr error
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return len(b), nil
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return bytes.NewReader(b).Read(b)
}

func (m *mockConn) Close() error {
	return nil
}

// Your tests will go here.
func TestNewConnection(t *testing.T) {

	testCases := []struct {
		name   string
		conn   net.Conn
		expect *Connection
	}{
		{
			name:   "ValidConnection",
			conn:   &mockConn{},
			expect: &Connection{conn: &mockConn{}, Status: ConnNew, readBuffer: ReadBuffer},
		},
		// Add more test cases here for all possible corner cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := NewConnection(tc.conn)
			if actual != tc.expect {
				t.Errorf("NewConnection(%v) = %v; expected %v",
					tc.conn, actual, tc.expect)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		writeErr    error
		expectedErr error
	}{
		{
			name:        "HappyPath",
			status:      "Connected",
			writeErr:    nil,
			expectedErr: nil,
		},
		{
			name:        "ConnectionDeadBeforeWrite",
			status:      ConnLinkDead,
			writeErr:    nil,
			expectedErr: nil,
		},
		{
			name:        "ConnectionWriteError",
			status:      "Connected",
			writeErr:    errors.New("write error"),
			expectedErr: errors.New("write error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &mockConn{
				writeErr: tt.writeErr,
			}
			c := &Connection{
				conn:   mockConn,
				Status: tt.status,
			}

			err := c.Write("Test data")

			if err != nil {
				if tt.expectedErr == nil {
					t.Errorf("unexpected error: %v", err)
				} else if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				}
			}

			if tt.writeErr != nil && c.Status != ConnLinkDead {
				t.Errorf("expected status: %v, got: %v", ConnLinkDead, c.Status)
			}
		})
	}
}

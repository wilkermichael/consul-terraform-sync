package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/consul-terraform-sync/logging"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDefaultLogLevel = "INFO"
)

func TestNewPrinter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		expectError bool
		config      *PrinterConfig
	}{
		{
			"error nil config",
			true,
			nil,
		},
		{
			"happy path debug log level",
			false,
			&PrinterConfig{
				WorkingDir: "path/to/wd",
				Workspace:  "ws",
			},
		},
		{
			"happy path",
			false,
			&PrinterConfig{
				WorkingDir: "path/to/wd",
				Workspace:  "ws",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NewPrinter(tc.config)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, actual)
			assert.Equal(t, testDefaultLogLevel, actual.logLevel)
			assert.Equal(t, tc.config.WorkingDir, actual.workingDir)
			assert.Equal(t, tc.config.Workspace, actual.workspace)
		})
	}
}

func DefaultTestPrinter(buf *bytes.Buffer) (*Printer, error) {
	printer, err := NewTestPrinter(&PrinterConfig{
		WorkingDir: "path/to/wd",
		Workspace:  "ws",
		Writer:     buf,
	})
	if err != nil {
		return nil, err
	}

	return printer, nil
}

func NewTestPrinter(config *PrinterConfig) (*Printer, error) {
	printer, err := NewPrinter(config)
	if err != nil {
		return nil, err
	}

	return printer, nil
}

func TestPrinterSetEnv(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	p, err := DefaultTestPrinter(&buf)
	assert.NoError(t, err)

	p.SetEnv(make(map[string]string))
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "client.printer")
	assert.Contains(t, buf.String(), "set")
	assert.Contains(t, buf.String(), "env")
	fmt.Println(buf.String())
}

func TestPrinterSetStdout(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	p, err := DefaultTestPrinter(&buf)
	assert.NoError(t, err)

	p.SetStdout(&buf)
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "client.printer")
	assert.Contains(t, buf.String(), "set")
	assert.Contains(t, buf.String(), "standard out")
}

func TestPrinterInit(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	p, err := DefaultTestPrinter(&buf)
	assert.NoError(t, err)

	ctx := context.Background()
	p.Init(ctx)
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "client.printer")
	assert.Contains(t, buf.String(), "init")
}

func TestPrinterLogLevel(t *testing.T) {
	cases := []struct {
		name        string
		expectWrite bool
		logLevel    string
		config      *PrinterConfig
	}{
		{
			"no write info log level",
			false,
			"INFO",
			&PrinterConfig{
				WorkingDir: "path/to/wd",
				Workspace:  "ws",
			},
		},
		{
			"write trace log level",
			true,
			"TRACE",
			&PrinterConfig{
				WorkingDir: "path/to/wd",
				Workspace:  "ws",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logging.LogLevel = tc.logLevel
			var buf bytes.Buffer
			tc.config.Writer = &buf
			p, err := NewTestPrinter(tc.config)
			assert.NoError(t, err)

			p.logger.Printf("[TRACE] Test Message")
			if tc.expectWrite {
				assert.NotEmpty(t, buf.String())
			} else {
				assert.Empty(t, buf.String())
			}
		})
	}
}

func TestPrinterApply(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	p, err := DefaultTestPrinter(&buf)
	assert.NoError(t, err)

	ctx := context.Background()
	p.Apply(ctx)
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "client.printer")
	assert.Contains(t, buf.String(), "apply")
}

func TestPrinterPlan(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	p, err := DefaultTestPrinter(&buf)
	assert.NoError(t, err)

	ctx := context.Background()
	p.Plan(ctx)
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "client.printer")
	assert.Contains(t, buf.String(), "plan")
}

func TestPrinterValidate(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	p, err := DefaultTestPrinter(&buf)
	assert.NoError(t, err)

	ctx := context.Background()
	p.Validate(ctx)
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "client.printer")
	assert.Contains(t, buf.String(), "validating")
}

func TestPrinterGoString(t *testing.T) {
	cases := []struct {
		name    string
		printer *Printer
	}{
		{
			"nil printer",
			nil,
		},
		{
			"happy path",
			&Printer{
				workingDir: "path/to/wd",
				workspace:  "ws",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.printer == nil {
				assert.Contains(t, tc.printer.GoString(), "nil")
				return
			}

			assert.Contains(t, tc.printer.GoString(), "&Printer")
			assert.Contains(t, tc.printer.GoString(), tc.printer.workingDir)
			assert.Contains(t, tc.printer.GoString(), tc.printer.workspace)
		})
	}
}

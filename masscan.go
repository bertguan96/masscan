package masscan

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ScanRunner represents something that can run a scan.
type ScanRunner interface {
	Run() (result interface{}, err error)
}

type Scanner struct {
	args        []string        // 执行参数
	binaryPath  string          // 二进制文件位置
	ctx         context.Context // 上下文
	rootRuntime bool            // root模式运行
	pid         int             // os.getpid()
}

type Option func(*Scanner)

// NewScanner creates a new Scanner, and can take options to apply to the scanner.
func NewScanner(options ...Option) (*Scanner, error) {
	scanner := &Scanner{}

	if len(options) == 0 {
		return nil, OptionsIsNull
	}

	for _, option := range options {
		option(scanner)
	}

	if scanner.binaryPath == "" {
		var err error
		scanner.binaryPath, err = exec.LookPath("bin/masscan")
		if err != nil {
			return nil, ErrMasscanNotInstalled
		}
	}

	if scanner.ctx == nil {
		scanner.ctx = context.Background()
	}

	return scanner, nil
}

func (scanner *Scanner) Run() (result *MasscanResult, warnings []string, err error) {
	var (
		cmd    *exec.Cmd
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	if scanner.rootRuntime {
		execPath := []string{scanner.binaryPath}
		execPath = append(execPath, scanner.args...)
		cmd = exec.Command("sudo", execPath...)
	} else {
		cmd = exec.Command(scanner.binaryPath, scanner.args...)
	}
	path, _ := os.Getwd()
	cmd.Dir = path // 绑定当前路径
	log.Printf("exec cmd: %s\n", cmd.String())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, warnings, err
	}

	scanner.pid = cmd.Process.Pid

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-scanner.ctx.Done():

		_ = cmd.Process.Kill()

		return nil, warnings, ErrScanTimeout
	case <-done:

		if stderr.Len() > 0 {
			warnings = strings.Split(strings.Trim(stderr.String(), "\n"), "\n")
		}

		if stdout.Len() > 0 {
			result, err := ParseScanResult(stdout.Bytes())
			if err != nil {
				warnings = append(warnings, err.Error())
				return nil, warnings, ErrParseOutput
			}
			return result, warnings, err
		}

	}
	return nil, nil, nil
}

// WithBinaryPath 设置二进制文件地址
func WithBinaryPath(path string) Option {
	return func(s *Scanner) {
		binary, err := filepath.Abs(path)
		if err != nil {
			log.Fatalf("set binary path error, %s", err)
		}
		s.binaryPath = binary
	}
}

// WithRate 设置访问频率
func WithRate(rate int) Option {
	return func(s *Scanner) {
		rate := fmt.Sprintf("--rate=%d", rate)
		s.args = append(s.args, rate)
	}
}

// WithPort 设置访问频率
func WithPort(ports string) Option {
	return func(s *Scanner) {
		rate := fmt.Sprintf("-p %s", ports)
		s.args = append(s.args, rate)
	}
}

// WithRandomizeHosts 随机host
func WithRandomizeHosts() Option {
	return func(s *Scanner) {
		s.args = append(s.args, "--randomize-hosts")
	}
}

// WithOutputJson 输出为JSON
func WithOutputJson() Option {
	return func(s *Scanner) {
		s.args = append(s.args, "-oJ")
		s.args = append(s.args, "-")
	}
}

// WithRetryTime 重试次数
func WithRetryTime(retryTime int) Option {
	return func(s *Scanner) {
		s.args = append(s.args, fmt.Sprintf("--retries=%d", retryTime))
	}
}

// WithTarget 设置待扫描的目标
func WithTarget(targets []string) Option {
	return func(s *Scanner) {
		s.args = append(s.args, targets...)
	}
}

func WithRoot() Option {
	return func(s *Scanner) {
		s.rootRuntime = true
	}
}

// AddOptions sets more scan options after the scan is created.
func (scanner *Scanner) AddOptions(options ...Option) *Scanner {
	for _, option := range options {
		option(scanner)
	}
	return scanner
}

// Args return the list of nmap args.
func (scanner *Scanner) Args() []string {
	return scanner.args
}

// AddArgs return the list of nmap args.
func (scanner *Scanner) AddArgs(val string) {
	scanner.args = append(scanner.args, val)
}

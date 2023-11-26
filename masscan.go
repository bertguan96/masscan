package masscan

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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
}

type Option func(*Scanner)

// NewScanner creates a new Scanner, and can take options to apply to the scanner.
func NewScanner(ctx context.Context, options ...Option) (*Scanner, error) {
	scanner := &Scanner{
		ctx: ctx,
	}

	if len(options) == 0 {
		return nil, ErrMasscanNotInstalled
	}

	for _, option := range options {
		option(scanner)
	}

	if scanner.binaryPath == "" {
		var err error
		scanner.binaryPath, err = exec.LookPath("bin/masscan")
		if err != nil {
			return nil, OptionsIsNull
		}
	}

	return scanner, nil
}

func (scanner *Scanner) Run() (err error) {
	start := time.Now()
	var cmd *exec.Cmd
	if scanner.rootRuntime {
		execPath := []string{scanner.binaryPath}
		execPath = append(execPath, scanner.args...)
		cmd = exec.Command("sudo", execPath...)
	} else {
		cmd = exec.Command(scanner.binaryPath, scanner.args...)
	}

	log.Printf("exec cmd: %s\n", cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	if err := cmd.Wait(); err != nil {
		return err
	}

	// 关闭输出流
	if err := stdout.Close(); err != nil {
		return err
	}
	if err := stderr.Close(); err != nil {
		return err
	}
	log.Printf("scan finish, cost: %ss\n", time.Since(start))
	return
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
		rate := fmt.Sprintf("--rate %d", rate)
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
func WithOutputJson(fileName string) Option {
	return func(s *Scanner) {
		s.args = append(s.args, fmt.Sprintf("-oJ %s", fileName))
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

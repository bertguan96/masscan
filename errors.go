package masscan

import "errors"

var (
	ErrMasscanNotInstalled = errors.New("masscan binary was not found")
	OptionsIsNull          = errors.New("options is nul") // 参数为空
)

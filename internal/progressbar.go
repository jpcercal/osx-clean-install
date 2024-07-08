package internal

import (
	"fmt"

	"github.com/pterm/pterm"
)

const (
	prefixText = "┗❯"

	infoText      = "info"
	infoTextColor = pterm.FgYellow

	successText      = "ok"
	successTextColor = pterm.FgGreen

	warningText      = "warning"
	warningTextColor = pterm.FgLightRed

	errorText      = "error"
	errorTextColor = pterm.FgRed
)

func NewProgressBar(total int) (*pterm.ProgressbarPrinter, error) {
	pterm.Info = pterm.PrefixPrinter{
		MessageStyle: &pterm.Style{pterm.FgGray, pterm.Italic},
		Prefix: pterm.Prefix{
			Style: &pterm.Style{pterm.FgDarkGray, pterm.Bold},
			Text:  prefixText,
		},
		Scope: pterm.Scope{
			Style: &pterm.Style{infoTextColor},
			Text:  infoText,
		},
	}

	pterm.Success = pterm.PrefixPrinter{
		MessageStyle: &pterm.Style{pterm.FgGray, pterm.Italic},
		Prefix: pterm.Prefix{
			Style: &pterm.Style{pterm.FgDarkGray, pterm.Bold},
			Text:  prefixText,
		},
		Scope: pterm.Scope{
			Style: &pterm.Style{successTextColor},
			Text:  successText,
		},
	}

	pterm.Warning = pterm.PrefixPrinter{
		MessageStyle: &pterm.Style{pterm.FgGray, pterm.Italic},
		Prefix: pterm.Prefix{
			Style: &pterm.Style{pterm.FgDarkGray, pterm.Bold},
			Text:  prefixText,
		},
		Scope: pterm.Scope{
			Style: &pterm.Style{warningTextColor},
			Text:  warningText,
		},
	}

	pterm.Error = pterm.PrefixPrinter{
		MessageStyle: &pterm.Style{pterm.FgGray, pterm.Italic},
		Prefix: pterm.Prefix{
			Style: &pterm.Style{pterm.FgDarkGray, pterm.Bold},
			Text:  prefixText,
		},
		Scope: pterm.Scope{
			Style: &pterm.Style{errorTextColor},
			Text:  errorText,
		},
	}

	p, err := pterm.DefaultProgressbar.
		WithMaxWidth(0).
		WithTotal(total).
		WithTitleStyle(&pterm.Style{pterm.FgCyan, pterm.Bold}).
		Start()
	if err != nil {
		return nil, fmt.Errorf("could not create progress bar %v", err)
	}

	return p, nil
}

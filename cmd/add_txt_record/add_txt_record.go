package main

import (
	"flag"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	motmedelLog "github.com/Motmedel/utils_go/pkg/log"
	"github.com/Motmedel/utils_go/pkg/net/domain_breakdown"
	"github.com/altshiftab/loopia_utils/pkg/loopia_utils"
	loopiaUtilsTypes "github.com/altshiftab/loopia_utils/pkg/types"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	logger := slog.Default()

	var username string
	flag.StringVar(
		&username,
		"username",
		"",
		"The Loopia API username.",
	)

	var password string
	flag.StringVar(
		&password,
		"password",
		"",
		"The Loopia API password.",
	)

	var domain string
	flag.StringVar(
		&domain,
		"domain",
		"",
		"The domain for which to add a record.",
	)

	var record string
	flag.StringVar(
		&record,
		"record",
		"",
		"The record to add.",
	)

	flag.Parse()

	if username == "" {
		motmedelLog.LogFatalWithExitingMessage("The username is empty.", nil, logger)
	}

	if password == "" {
		motmedelLog.LogFatalWithExitingMessage("The password is empty.", nil, logger)
	}

	if domain == "" {
		motmedelLog.LogFatalWithExitingMessage("The domain is empty.", nil, logger)
	}

	if record == "" {
		motmedelLog.LogFatalWithExitingMessage("The record is empty.", nil, logger)
	}

	if domainBreakdown := domain_breakdown.GetDomainBreakdown(domain); domainBreakdown == nil {
		motmedelLog.LogFatalWithExitingMessage("Invalid domain.", nil, logger)
	}

	client := loopia_utils.Client{
		ApiUser:     username,
		ApiPassword: password,
		HttpClient:  &http.Client{Timeout: 30 * time.Second},
	}

	_, err := client.AddRecord(
		&loopiaUtilsTypes.Record{Type: "TXT", Ttl: loopia_utils.DefaultTtlValue, Rdata: record},
		domain,
	)
	if err != nil {
		msg := "An error occurred when adding a TXT record."
		motmedelLog.LogFatalWithExitingMessage(
			msg,
			&motmedelErrors.InputError{
				Message: msg,
				Cause:   err,
				Input:   []any{record, domain},
			},
			logger,
		)
	}
}

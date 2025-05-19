package main

import (
	"context"
	"flag"
	"fmt"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	motmedelErrorLogger "github.com/Motmedel/utils_go/pkg/log/error_logger"
	"github.com/Motmedel/utils_go/pkg/net/domain_breakdown"
	"github.com/altshiftab/loopia_utils/pkg/loopia_utils"
	loopiaUtilsTypes "github.com/altshiftab/loopia_utils/pkg/types"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := motmedelErrorLogger.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger.Logger)

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
		logger.FatalWithExitingMessage("The username is empty.", nil)
	}

	if password == "" {
		logger.FatalWithExitingMessage("The password is empty.", nil)
	}

	if domain == "" {
		logger.FatalWithExitingMessage("The domain is empty.", nil)
	}

	if record == "" {
		logger.FatalWithExitingMessage("The record is empty.", nil)
	}

	if domainBreakdown := domain_breakdown.GetDomainBreakdown(domain); domainBreakdown == nil {
		logger.FatalWithExitingMessage("Invalid domain.", nil)
	}

	client := loopia_utils.Client{
		Client:      &http.Client{Timeout: 30 * time.Second},
		ApiUser:     username,
		ApiPassword: password,
	}

	err := client.AddRecord(
		context.Background(),
		&loopiaUtilsTypes.Record{Type: "TXT", Ttl: loopia_utils.DefaultTtlValue, Rdata: record},
		domain,
	)
	if err != nil {
		logger.FatalWithExitingMessage(
			"An error occurred when adding a TXT record.",
			motmedelErrors.NewWithTrace(
				fmt.Errorf("client add record: %w", err),
				domain,
				record,
			),
		)
	}
}

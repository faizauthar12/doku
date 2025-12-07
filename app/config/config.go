package config

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/faizauthar12/doku/app/utils/helper"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Doku struct {
		ClientID   string
		SecretKey  string
		PrivateKey string
	}
	TransactionFee struct {
		Cards struct {
			PercentageRate float64 // e.g. 2.8%
			FlatFee        int64   // e.g. IDR 2000
		}
		VirtualAccount struct {
			FlatFee int64 // e.g. IDR 4000
		}
		ConvenienceStore struct {
			Alfamart struct {
				FlatFee int64 //e.g. IDR 5000
			}
			Indomaret struct {
				FlatFee int64 //e.g. IDR 6500
			}
		}
		QR struct {
			FlatFee int64 // e.g. IDR 700
		}
		EWallet struct {
			ShopeePay struct {
				PercentageRate float64 // e.g. 2%
			}
			OVO struct {
				PercentageRate float64 // e.g. 2%
			}
			LinkAja struct {
				PercentageRate float64 // e.g. 2%
			}
			Doku struct {
				PercentageRate float64 // e.g. 1.5%
			}
			Dana struct {
				PercentageRate float64 // e.g. 1.5%
			}
		}
		DirectDebit struct {
			BRI struct {
				PercentageRate float64 // e.g. 2%
			}
			AlloBank struct {
				PercentageRate float64 // e.g. 2%
			}
			OctoCash struct {
				PercentageRate float64 // e.g. 2%
			}
		}
		DigitalBanking struct {
			JeniusPay struct {
				PercentageRate float64 // e.g. 1.5%
			}
		}
		PayLater struct {
			Akulaku struct {
				PercentageRate float64 // e.g. 1.5%
			}
			Kredivo struct {
				PercentageRate float64 // e.g. 2.3%
			}
			Ceria struct {
				PercentageRate float64 // e.g. 1.5%
			}
			Indodana struct {
				PercentageRate float64 // e.g. 2.3%
			}
		}
		Tax int64 // e.g. 11%
	}
}

var cfg Configuration

func GetEnvString(key string, dflt string) string {
	value := os.Getenv(key)
	if value == "" {
		return dflt
	}
	return value
}

func GetEnvSliceString(key string, dflt string) []string {
	values := os.Getenv(key)
	if values == "" {
		return strings.Split(dflt, ",")
	}
	return strings.Split(values, ",")
}

func GetEnvInt(key string, dflt int) int {
	value := os.Getenv(key)
	i, err := strconv.ParseInt(value, 10, 64)
	if value == "" && err != nil {
		return dflt
	}
	return int(i)
}

func GetEnvBool(key string, dflt bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return dflt
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return dflt
	}
	return b
}

func GetEnvFloat64(key string, dflt float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return dflt
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return dflt
	}
	return f
}

func Get() Configuration {
	return cfg
}

// loadConfigFromEnv populates cfg from environment variables (internal helper)
func loadConfigFromEnv() {
	// Doku Config
	cfg.Doku.ClientID = GetEnvString("DOKU_API_CLIENT_ID", "")
	cfg.Doku.SecretKey = GetEnvString("DOKU_API_SECRET_KEY", "")
	cfg.Doku.PrivateKey = GetEnvString("DOKU_API_PRIVATE_KEY", "")

	// Transaction Fee Config
	// Cards
	cfg.TransactionFee.Cards.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_CARDS_PERCENTAGE_RATE", 2.8)
	cfg.TransactionFee.Cards.FlatFee = int64(GetEnvInt("TRANSACTION_FEE_CARDS_FLAT_FEE", 2000))

	// Virtual Account (Transfer Bank)
	cfg.TransactionFee.VirtualAccount.FlatFee = int64(GetEnvInt("TRANSACTION_FEE_VIRTUAL_ACCOUNT_FLAT_FEE", 4000))

	// Convenience Store
	cfg.TransactionFee.ConvenienceStore.Alfamart.FlatFee = int64(GetEnvInt("TRANSACTION_FEE_ALFAMART_FLAT_FEE", 5000))
	cfg.TransactionFee.ConvenienceStore.Indomaret.FlatFee = int64(GetEnvInt("TRANSACTION_FEE_INDOMARET_FLAT_FEE", 6500))

	// QR
	cfg.TransactionFee.QR.FlatFee = int64(GetEnvInt("TRANSACTION_FEE_QR_FLAT_FEE", 700))

	// E-Wallet
	cfg.TransactionFee.EWallet.ShopeePay.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_SHOPEEPAY_PERCENTAGE_RATE", 2.0)
	cfg.TransactionFee.EWallet.OVO.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_OVO_PERCENTAGE_RATE", 2.0)
	cfg.TransactionFee.EWallet.LinkAja.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_LINKAJA_PERCENTAGE_RATE", 2.0)
	cfg.TransactionFee.EWallet.Doku.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_DOKU_WALLET_PERCENTAGE_RATE", 1.5)
	cfg.TransactionFee.EWallet.Dana.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_DANA_PERCENTAGE_RATE", 1.5)

	// Direct Debit
	cfg.TransactionFee.DirectDebit.BRI.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_BRI_DIRECT_DEBIT_PERCENTAGE_RATE", 2.0)
	cfg.TransactionFee.DirectDebit.AlloBank.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_ALLOBANK_DIRECT_DEBIT_PERCENTAGE_RATE", 2.0)
	cfg.TransactionFee.DirectDebit.OctoCash.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_OCTOCASH_DIRECT_DEBIT_PERCENTAGE_RATE", 2.0)

	// Digital Banking
	cfg.TransactionFee.DigitalBanking.JeniusPay.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_JENIUSPAY_PERCENTAGE_RATE", 1.5)

	// PayLater
	cfg.TransactionFee.PayLater.Akulaku.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_AKULAKU_PERCENTAGE_RATE", 1.5)
	cfg.TransactionFee.PayLater.Kredivo.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_KREDIVO_PERCENTAGE_RATE", 2.3)
	cfg.TransactionFee.PayLater.Ceria.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_CERIA_PERCENTAGE_RATE", 1.5)
	cfg.TransactionFee.PayLater.Indodana.PercentageRate = GetEnvFloat64("TRANSACTION_FEE_INDODANA_PERCENTAGE_RATE", 2.3)

	// Tax
	cfg.TransactionFee.Tax = int64(GetEnvInt("TRANSACTION_FEE_TAX", 11))
}

// InitConfig loads .env file(s) and then populates config.
// Use this when running doku as a standalone application.
func InitConfig(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, helper.DefaultStatusText[http.StatusInternalServerError])
		log.Fatalf("Error loading .env file")
		log.Fatalf("Error: %v", logData)
	}

	loadConfigFromEnv()
}

// InitConfigFromEnv populates config from existing environment variables.
// Use this when doku is used as a module and the parent project has already
// loaded the .env file (e.g., backend project calls godotenv.Load() first).
func InitConfigFromEnv() {
	loadConfigFromEnv()
}

// InitConfigWithStruct allows setting the configuration directly with a struct.
// Use this when the parent project wants full control over configuration values.
func InitConfigWithStruct(configuration Configuration) {
	cfg = configuration
}

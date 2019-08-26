package i18n

//go:generate gotext -srclang=en update -out=catalog/catalog.go -lang=en,zh
import (
	"github.com/hyperledger/fabric/common/flogging"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	logger     = flogging.MustGetLogger("i18n")
	printerMap = make(map[LocaleType]*message.Printer)
)

// Error this error will adapt to i18n, and errors implementing this interface must be output using message.Printer
type Error interface {
	error
	Code() string
}

// LocaleType locale type
type LocaleType int

// now is only can use chinese and english
const (
	Chinese LocaleType = iota
	English
)

// GetI18nPrinter get an I18n printer to fit the browser of each country
func GetI18nPrinter(locale LocaleType) *message.Printer {
	if _, ok := printerMap[locale]; !ok {
		switch locale {
		case Chinese:
			printerMap[locale] = message.NewPrinter(language.Chinese)
		case English:
			printerMap[locale] = message.NewPrinter(language.English)
		default:
			return message.NewPrinter(language.Chinese)
		}
	}

	return printerMap[locale]
}

// GetDefaultPrinter get a default printer, now is Chinese
func GetDefaultPrinter() *message.Printer {
	return GetI18nPrinter(Chinese)
}

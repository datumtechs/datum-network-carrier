// Package logutil creates a Multi writer instance that
// write all logs that are written to stdout.
package logutil

import (
	"github.com/RosettaFlow/Carrier-Go/common/params"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func addLogWriter(w io.Writer) {
	mw := io.MultiWriter(logrus.StandardLogger().Out, w)
	logrus.SetOutput(mw)
}

// ConfigurePersistentLogging adds a log-to-file writer. File content is identical to stdout.
func ConfigurePersistentLogging(logFileName string) error {
	logrus.WithField("logFileName", logFileName).Info("Logs will be made persistent")
	f, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, params.CarrierIoConfig().ReadWritePermissions)
	if err != nil {
		return err
	}

	addLogWriter(f)

	logrus.Info("File logging initialized")
	return nil
}

// Masks the url credentials before logging for security purpose
// [scheme:][//[userinfo@]host][/]path[?query][#fragment] -->  [scheme:][//[***]host][/***][#***]
// if the format is not matched nothing is done, string is returned as is.
/*func MaskCredentialsLogging(currUrl string) string {
	// error if the input is not a URL
	MaskedUrl := currUrl
	u, err := url.Parse(currUrl)
	if err != nil {
		return currUrl // Not a URL, nothing to do
	}
	// Mask the userinfo and the URI (path?query or opaque?query ) and fragment, leave the scheme and host(host/port)  untouched
	if u.GetUser != nil {
		MaskedUrl = strings.Replace(MaskedUrl, u.GetUser.String(), "***", 1)
	}
	if len(u.RequestURI()) > 1 { // Ignore the '/'
		MaskedUrl = strings.Replace(MaskedUrl, u.RequestURI(), "/***", 1)
	}
	if len(u.Fragment) > 0 {
		MaskedUrl = strings.Replace(MaskedUrl, u.RawFragment, "***", 1)
	}
	return MaskedUrl
}*/

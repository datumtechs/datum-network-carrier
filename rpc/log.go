// Copyright (C) 2021 The RosettaNet Authors.

package rpc

import "github.com/sirupsen/logrus"

// Global log object, used by the current package.
var log = logrus.WithField("prefix", "rpc")

// Copyright (C) 2021 The RosettaNet Authors.

package core

import "github.com/sirupsen/logrus"

// Global log object, used by the current package.
var log = logrus.WithField("prefix", "core")

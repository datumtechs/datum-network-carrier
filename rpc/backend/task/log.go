// Copyright (C) 2021 The datum-network Authors.

package task

import "github.com/sirupsen/logrus"

// Global log object, used by the current package.
var log = logrus.WithField("prefix", "rpc.backend.task")

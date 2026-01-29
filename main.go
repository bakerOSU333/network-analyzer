package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/bakerOSU333/network-analyzer/cronjob"
	"github.com/bakerOSU333/network-analyzer/network"
	"github.com/bakerOSU333/network-analyzer/ping"
	"github.com/bakerOSU333/network-analyzer/speedtest"
)
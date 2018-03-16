package utilities

import (
    "testing"
)

func TestUtilities(t *testing.T) {
    ips := MyIpAddress()
    t.Logf("Ips:%v\n", ips)
}
package tun

import (
	"fmt"
	"github.com/2pizzzza/tunShiffer/internal/config"
	"github.com/2pizzzza/tunShiffer/internal/logger"
	"github.com/2pizzzza/tunShiffer/pkg/packet"

	"github.com/songgao/water"
	"net"
	"os/exec"
)

type TunHandler struct {
	iface  *water.Interface
	config *config.Config
	logger *logger.Logger
}

func NewTunHandler(cfg *config.Config, logger *logger.Logger) (*TunHandler, error) {
	cmd := exec.Command("ip", "link", "delete", cfg.TunName)
	if err := cmd.Run(); err != nil {
		logger.Log(fmt.Sprintf("Warning: failed to delete existing interface %s: %v", cfg.TunName, err))
	}

	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = cfg.TunName

	iface, err := water.New(config)
	if err != nil {
		return nil, err
	}

	err = setupInterface(cfg.TunName, cfg.TunIP, cfg.TunMask, logger)
	if err != nil {
		iface.Close()
		return nil, fmt.Errorf("failed to setup interface: %v", err)
	}

	logger.Log(fmt.Sprintf("TUN interface %s created and configured with IP %s/%s", cfg.TunName, cfg.TunIP, cfg.TunMask))

	return &TunHandler{
		iface:  iface,
		config: cfg,
		logger: logger,
	}, nil
}

func setupInterface(tunName, ip, mask string, logger *logger.Logger) error {
	cmd := exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%s", ip, mask), "dev", tunName)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Log(fmt.Sprintf("Failed to add IP address: %v, output: %s", err, output))
		return err
	}

	cmd = exec.Command("ip", "link", "set", tunName, "up")
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Log(fmt.Sprintf("Failed to set interface up: %v, output: %s", err, output))
		return err
	}

	ipNet := calculateNetwork(ip, mask)

	cmd = exec.Command("ip", "route", "del", ipNet, "dev", tunName)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Log(fmt.Sprintf("Warning: failed to delete route: %v, output: %s", err, output))
	}

	cmd = exec.Command("ip", "route", "add", ipNet, "dev", tunName)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Log(fmt.Sprintf("Failed to add route: %v, output: %s", err, output))
		return err
	}

	return nil
}

func calculateNetwork(ip, mask string) string {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return "0.0.0.0/0"
	}
	maskInt := net.ParseIP(mask)
	if maskInt == nil {
		_, ipNet, _ := net.ParseCIDR(fmt.Sprintf("%s/%s", ip, mask))
		if ipNet != nil {
			return ipNet.String()
		}
	}
	return "10.0.0.0/24"
}

func (t *TunHandler) Start() error {
	buffer := make([]byte, 65535)
	for {
		n, err := t.iface.Read(buffer)
		if err != nil {
			return err
		}

		packetData := buffer[:n]
		packetInfo, err := packet.ParsePacket(packetData)
		if err != nil {
			t.logger.Log(fmt.Sprintf("Failed to parse packet: %v", err))
			continue
		}

		t.logger.LogPacket("Packet received", packetInfo, packetData)
	}
}

func (t *TunHandler) Close() {
	t.iface.Close()
}

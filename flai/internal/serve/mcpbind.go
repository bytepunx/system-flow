package serve

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// DockerDesktop is what docker info says of the operating system of Docker
// Desktop's engine, which forwards the host gateway to the host's loopback.
const DockerDesktop = "Docker Desktop"

// mcpBindAuto is the bind that leaves the HTTP MCP server's addresses to
// flai, config.MCPBindAuto's value, kept here so that serve reads the word
// without depending on the config's own spelling of it.
const mcpBindAuto = "auto"

// mcpLoopback is the address an automatic bind always listens on.
const mcpLoopback = "127.0.0.1"

// gatewayFormat is the docker network inspect format that prints a network's
// gateways on one line, space-separated.
const gatewayFormat = "{{range .IPAM.Config}}{{.Gateway}} {{end}}"

// DockerEngine is the operating system docker info names for the engine
// program, the docker client, reaches from dir.
func DockerEngine(docker execx.Runner, dir, program string) (string, error) {
	out, err := docker.Run(dir, program, "info", "--format", "{{.OperatingSystem}}")
	engine := strings.TrimSpace(out)
	if err != nil || engine == "" {
		return "", fmt.Errorf("%s info could not name the Docker engine (%s): start Docker, or check that %s reaches it", program, dockerSaid(out, err), program)
	}
	return engine, nil
}

// NetworkGateway is the first gateway docker network inspect names for
// network, where "" is the bridge network the container harness runs on.
func NetworkGateway(docker execx.Runner, dir, program, network string) (string, error) {
	if network == "" {
		network = "bridge"
	}
	out, err := docker.Run(dir, program, "network", "inspect", network, "--format", gatewayFormat)
	fields := strings.Fields(out)
	if err != nil {
		return "", fmt.Errorf("%s network inspect %s could not name the network's gateway (%s): check that the network exists with %s network ls", program, network, dockerSaid(out, err), program)
	}
	if len(fields) == 0 {
		return "", fmt.Errorf("%s network inspect %s names no gateway: the network has no IPAM gateway, so run the container on a network that has one", program, network)
	}
	return fields[0], nil
}

// ContainerGateways is every gateway of every Docker network program reaches,
// each once, in the order docker names them.
func ContainerGateways(docker execx.Runner, dir, program string) ([]string, error) {
	out, err := docker.Run(dir, program, "network", "ls", "--format", "{{.Name}}")
	if err != nil {
		return nil, fmt.Errorf("%s network ls could not list Docker's networks (%s): start Docker, or check that %s reaches it", program, dockerSaid(out, err), program)
	}
	names := strings.Fields(out)
	if len(names) == 0 {
		return nil, nil
	}
	args := append([]string{"network", "inspect"}, names...)
	out, err = docker.Run(dir, program, append(args, "--format", gatewayFormat)...)
	if err != nil {
		return nil, fmt.Errorf("%s network inspect could not name the gateways of the networks %s (%s): check them with %s network ls", program, strings.Join(names, ", "), dockerSaid(out, err), program)
	}
	var gateways []string
	for _, gateway := range strings.Fields(out) {
		if !slices.Contains(gateways, gateway) {
			gateways = append(gateways, gateway)
		}
	}
	return gateways, nil
}

// ResolveMCPBinds is the addresses the HTTP MCP server listens on for binds,
// as the host config writes them, with each auto expanded for the container
// harness's network, each address once and in order; when docker cannot say
// what auto needs, the addresses are still returned, auto as loopback alone,
// with an error saying what could not be read.
func ResolveMCPBinds(docker execx.Runner, dir, program, network string, binds []string) ([]string, error) {
	var (
		resolved []string
		auto     []string
		autoErr  error
		asked    bool
	)
	add := func(addr string) {
		if !slices.Contains(resolved, addr) {
			resolved = append(resolved, addr)
		}
	}
	for _, bind := range binds {
		if bind != mcpBindAuto {
			add(bind)
			continue
		}
		if !asked {
			auto, autoErr = autoBinds(docker, dir, program, network)
			asked = true
		}
		for _, addr := range auto {
			add(addr)
		}
	}
	return resolved, autoErr
}

// autoBinds is what auto listens on: loopback, and, on an engine other than
// Docker Desktop and a network other than the host's, the network's gateway,
// which such an engine routes host.docker.internal to; with an error, and
// loopback alone, when docker cannot say.
func autoBinds(docker execx.Runner, dir, program, network string) ([]string, error) {
	binds := []string{mcpLoopback}
	if network == "host" {
		return binds, nil
	}
	engine, err := DockerEngine(docker, dir, program)
	if err != nil {
		return binds, fmt.Errorf("auto listens on %s alone: %w", mcpLoopback, err)
	}
	if engine == DockerDesktop {
		return binds, nil
	}
	gateway, err := NetworkGateway(docker, dir, program, network)
	if err != nil {
		return binds, fmt.Errorf("auto listens on %s alone, though this Docker engine (%s) routes %s to a network's gateway: %w", mcpLoopback, engine, containerHost, err)
	}
	return append(binds, gateway), nil
}

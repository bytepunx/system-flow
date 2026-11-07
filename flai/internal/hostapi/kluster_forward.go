package hostapi

// Port-forwards (S-0586, the operator's answer 3 on TH-0861): flai serve runs
// kubectl port-forward for a service port of a namespace of the project's
// kluster, bound to 0.0.0.0 so that a hosted flai can reach it, keeps each
// one until the dashboard stops it, its namespace is torn down through the
// dashboard, its process exits, or flai serve stops, and kluster.deployments
// lists them beside each namespace's services.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
)

// KlusterForwardAddress is the address every forward listens on: every
// interface of the host, so that the dashboard's browser, or a hosted flai,
// reaches it on the host's own address. Anything that reaches the host's
// network reaches the forward, which has no sign-in of its own.
const KlusterForwardAddress = "0.0.0.0"

// klusterForwardReady bounds how long a start waits for kubectl to say it
// forwards.
const klusterForwardReady = 15 * time.Second

// klusterForwardStop bounds how long a stop waits for kubectl to exit.
const klusterForwardStop = 5 * time.Second

// klusterForwardOutput is the most of a forward's output kept, from its end.
const klusterForwardOutput = 2048

// KlusterForward is one port-forward flai serve runs, or ran, to a service
// port of a namespace.
type KlusterForward struct {
	Namespace string `json:"namespace"`
	Service   string `json:"service"`
	// Port is the service's port in the cluster.
	Port int `json:"port"`
	// LocalPort is the host's port it listens on, on Address.
	LocalPort int    `json:"local_port"`
	Address   string `json:"address"`
	StartedAt string `json:"started_at"`
	// Running is whether kubectl still runs; a forward whose kubectl exited
	// is kept, not running, with the end of its output, until it is started
	// or stopped again.
	Running bool   `json:"running"`
	EndedAt string `json:"ended_at,omitempty"`
	// Output is the end of kubectl's output once it exited.
	Output string `json:"output,omitempty"`
}

type forwardKey struct {
	root, namespace, service string
	port                     int
}

// klusterForward is a forward and its kubectl.
type klusterForward struct {
	KlusterForward
	cancel context.CancelFunc
	done   chan struct{}
	// ready is closed once kubectl says it forwards.
	ready <-chan struct{}
}

// klusterForwards are the forwards of every project, kept by KlusterJobs.
// The zero value is ready to use.
type klusterForwards struct {
	mu sync.Mutex
	m  map[forwardKey]*klusterForward
}

// klusterForwardCommand is the kubectl a forward runs, in the project's
// main checkout; tests replace it.
var klusterForwardCommand = func(ctx context.Context, root string, args []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Dir = root
	return cmd
}

// klusterFreePort is a port of the host nothing listens on, as the system
// gives one; tests replace it.
var klusterFreePort = func() (int, error) {
	l, err := net.Listen("tcp", KlusterForwardAddress+":0")
	if err != nil {
		return 0, err
	}
	defer func() { _ = l.Close() }()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// forwardWriter keeps the end of kubectl's output and says once it forwards.
type forwardWriter struct {
	mu    sync.Mutex
	buf   []byte
	ready chan struct{}
	once  sync.Once
}

func (w *forwardWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf = append(w.buf, p...)
	if len(w.buf) > 4*klusterForwardOutput {
		w.buf = append([]byte{}, w.buf[len(w.buf)-klusterForwardOutput:]...)
	}
	if bytes.Contains(w.buf, []byte("Forwarding from")) {
		w.once.Do(func() { close(w.ready) })
	}
	return len(p), nil
}

func (w *forwardWriter) tail() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	b := w.buf
	if len(b) > klusterForwardOutput {
		b = b[len(b)-klusterForwardOutput:]
		if i := bytes.IndexByte(b, '\n'); i >= 0 && i < len(b)-1 {
			b = b[i+1:]
		}
	}
	return string(b)
}

// StartForward starts kubectl port-forward of service's port in namespace
// of the cluster whose kube context is kubeContext, for the project at root,
// on a free port of the host bound to KlusterForwardAddress, and answers it
// once kubectl says it forwards. A forward of the same port that runs is
// answered as it is. It fails, with kubectl's output, when kubectl cannot
// start, exits first, or says nothing within klusterForwardReady.
func (k *KlusterJobs) StartForward(root, kubeContext, namespace, service string, port int) (KlusterForward, error) {
	f := &k.forwards
	key := forwardKey{root, namespace, service, port}
	f.mu.Lock()
	if f.m == nil {
		f.m = map[forwardKey]*klusterForward{}
	}
	if fw, ok := f.m[key]; ok && fw.Running {
		f.mu.Unlock()
		// a start of the same port may still be waiting for kubectl: answer
		// once it forwards, or say it ended, never a port not listening yet
		timer := time.NewTimer(klusterForwardReady)
		defer timer.Stop()
		select {
		case <-fw.ready:
		case <-fw.done:
		case <-timer.C:
		}
		f.mu.Lock()
		got := fw.KlusterForward
		f.mu.Unlock()
		if !got.Running {
			return KlusterForward{}, fmt.Errorf("kubectl port-forward of service/%s port %d in the namespace %s ended before it forwarded: %s; start it again", service, port, namespace, strings.TrimSpace(got.Output))
		}
		return got, nil
	}
	local, err := klusterFreePort()
	if err != nil {
		f.mu.Unlock()
		return KlusterForward{}, fmt.Errorf("find a free port on the host to forward service/%s port %d of the namespace %s to: %w", service, port, namespace, err)
	}
	args := []string{"--kubeconfig", kluster.KubeconfigPath(root), "--context", kubeContext, "--namespace", namespace,
		"port-forward", "--address", KlusterForwardAddress, "service/" + service, fmt.Sprintf("%d:%d", local, port)}
	ctx, cancel := context.WithCancel(k.ctx)
	cmd := klusterForwardCommand(ctx, root, args)
	out := &forwardWriter{ready: make(chan struct{})}
	cmd.Stdout, cmd.Stderr = out, out
	if err := cmd.Start(); err != nil {
		cancel()
		f.mu.Unlock()
		return KlusterForward{}, fmt.Errorf("start kubectl port-forward of service/%s port %d in the namespace %s: %w; install kubectl on the host, with kluster setup, then start it again", service, port, namespace, err)
	}
	fw := &klusterForward{KlusterForward: KlusterForward{Namespace: namespace, Service: service, Port: port, LocalPort: local, Address: KlusterForwardAddress,
		StartedAt: k.now().UTC().Format(time.RFC3339), Running: true}, cancel: cancel, done: make(chan struct{}), ready: out.ready}
	f.m[key] = fw
	f.mu.Unlock()
	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		_ = cmd.Wait()
		cancel()
		f.mu.Lock()
		fw.Running, fw.EndedAt, fw.Output = false, k.now().UTC().Format(time.RFC3339), out.tail()
		f.mu.Unlock()
		close(fw.done)
		if k.logger != nil {
			k.logger.Info("kluster port-forward ended", "component", "hostapi", "namespace", namespace, "service", service, "port", port, "local_port", local)
		}
	}()
	timer := time.NewTimer(klusterForwardReady)
	defer timer.Stop()
	select {
	case <-out.ready:
		if k.logger != nil {
			k.logger.Info("kluster port-forward started", "component", "hostapi", "namespace", namespace, "service", service, "port", port, "local_port", local, "address", KlusterForwardAddress)
		}
		f.mu.Lock()
		defer f.mu.Unlock()
		return fw.KlusterForward, nil
	case <-fw.done:
		return KlusterForward{}, fmt.Errorf("kubectl port-forward of service/%s port %d in the namespace %s exited before it forwarded: %s", service, port, namespace, strings.TrimSpace(out.tail()))
	case <-timer.C:
		cancel()
		<-fw.done
		return KlusterForward{}, fmt.Errorf("kubectl port-forward of service/%s port %d in the namespace %s said nothing of forwarding within %s, so flai stopped it: %s", service, port, namespace, klusterForwardReady, strings.TrimSpace(out.tail()))
	}
}

// StopForward stops the forward of service's port in namespace for the
// project at root, waits for kubectl to exit, forgets it, and answers it as
// it ended; false when none is kept.
func (k *KlusterJobs) StopForward(root, namespace, service string, port int) (KlusterForward, bool) {
	f := &k.forwards
	key := forwardKey{root, namespace, service, port}
	f.mu.Lock()
	fw, ok := f.m[key]
	if ok {
		delete(f.m, key)
	}
	f.mu.Unlock()
	if !ok {
		return KlusterForward{}, false
	}
	fw.cancel()
	select {
	case <-fw.done:
	case <-time.After(klusterForwardStop):
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return fw.KlusterForward, true
}

// StopForwards stops every forward into namespace for the project at root,
// as its teardown starts, and answers how many it stopped.
func (k *KlusterJobs) StopForwards(root, namespace string) int {
	f := &k.forwards
	f.mu.Lock()
	var keys []forwardKey
	for key := range f.m {
		if key.root == root && key.namespace == namespace {
			keys = append(keys, key)
		}
	}
	f.mu.Unlock()
	for _, key := range keys {
		k.StopForward(key.root, key.namespace, key.service, key.port)
	}
	return len(keys)
}

// Forwards are the forwards kept for namespace of the project at root,
// running or ended, by service and port.
func (k *KlusterJobs) Forwards(root, namespace string) []KlusterForward {
	f := &k.forwards
	f.mu.Lock()
	defer f.mu.Unlock()
	out := []KlusterForward{}
	for key, fw := range f.m {
		if key.root == root && key.namespace == namespace {
			out = append(out, fw.KlusterForward)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Service != out[j].Service {
			return out[i].Service < out[j].Service
		}
		return out[i].Port < out[j].Port
	})
	return out
}

// klusterForwardParams are kluster.forward.start's and kluster.forward.stop's
// params.
type klusterForwardParams struct {
	Namespace string `json:"namespace"`
	Service   string `json:"service"`
	Port      int    `json:"port"`
}

// dnsLabel is a Kubernetes name: a DNS label, as namespaces and services are.
var dnsLabel = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?$`)

// klusterForwardIn reads and checks a forward's params.
func klusterForwardIn(raw json.RawMessage) (klusterForwardParams, *channel.Error) {
	in, e := decode[klusterForwardParams](raw)
	if e != nil {
		return in, e
	}
	switch {
	case !dnsLabel.MatchString(in.Namespace):
		return in, bad("%q is not a namespace; give one kluster.deployments lists, such as s-0586", in.Namespace)
	case !dnsLabel.MatchString(in.Service):
		return in, bad("%q is not a service's name; give one kluster.deployments lists under its namespace", in.Service)
	case in.Port < 1 || in.Port > 65535:
		return in, bad("%d is not a port; give one of the service's ports, as kluster.deployments lists them", in.Port)
	}
	return in, nil
}

// noForwards refuses a forward asked of a flai that runs no kluster jobs.
func noForwards() *channel.Error {
	return &channel.Error{Code: Rule, Message: "this flai runs no port-forwards: flai serve runs them, so ask through the dashboard flai serve answers"}
}

// klusterForwardStartDo answers kluster.forward.start: the forward started,
// or the one that runs, refused when the project has no block, its cluster
// has no kubeconfig yet, or the cluster lists no such service port.
func klusterForwardStartDo(ctx context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
	in, e := klusterForwardIn(raw)
	if e != nil {
		return nil, e
	}
	if host.Kluster == nil {
		return nil, noForwards()
	}
	m, e := projectManifest(p)
	if e != nil {
		return nil, e
	}
	if m.Kluster == nil {
		return nil, &channel.Error{Code: Rule, Message: kluster.ErrNoBlock.Error()}
	}
	c, err := openKluster(p.Root, m)
	if err != nil {
		return nil, &channel.Error{Code: Rule, Message: klusterRelative(p.Root, err.Error())}
	}
	named := c.Named()
	if _, err := os.Stat(kluster.KubeconfigPath(p.Root)); err != nil {
		return nil, &channel.Error{Code: Rule, Message: fmt.Sprintf("the cluster %s has no kubeconfig yet, so there is nothing to forward; bring it up on the Klusters page", named.Name)}
	}
	sctx, cancel := context.WithTimeout(ctx, klusterServicesTimeout)
	services, err := c.Services(sctx)
	cancel()
	if err != nil {
		return nil, failed(errors.New(klusterRelative(p.Root, err.Error())))
	}
	if !servesPort(services, in) {
		return nil, &channel.Error{Code: Rule, Message: fmt.Sprintf("the cluster %s lists no service %s with the port %d in the namespace %s; pick one kluster.deployments lists", named.Name, in.Service, in.Port, in.Namespace)}
	}
	fw, err := host.Kluster.StartForward(p.Root, named.Context, in.Namespace, in.Service, in.Port)
	if err != nil {
		return nil, failed(errors.New(klusterRelative(p.Root, err.Error())))
	}
	return written(fw)
}

// servesPort says whether services hold in's service with in's port.
func servesPort(services []KlusterService, in klusterForwardParams) bool {
	for _, s := range services {
		if s.Namespace != in.Namespace || s.Name != in.Service {
			continue
		}
		for _, port := range s.Ports {
			if port.Port == in.Port {
				return true
			}
		}
	}
	return false
}

// klusterForwardStopDo answers kluster.forward.stop: the forward stopped,
// refused when none is kept.
func klusterForwardStopDo(_ context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
	in, e := klusterForwardIn(raw)
	if e != nil {
		return nil, e
	}
	if host.Kluster == nil {
		return nil, noForwards()
	}
	fw, ok := host.Kluster.StopForward(p.Root, in.Namespace, in.Service, in.Port)
	if !ok {
		return nil, &channel.Error{Code: Rule, Message: fmt.Sprintf("no port-forward of service/%s port %d in the namespace %s is running; nothing to stop", in.Service, in.Port, in.Namespace)}
	}
	return written(fw)
}

// forwardTarget names a forward for the journal.
func forwardTarget(fw KlusterForward) string {
	return "service/" + fw.Service + " port " + strconv.Itoa(fw.Port) + " in the namespace " + fw.Namespace
}

// describeKlusterForward says for the journal which forward was started or
// stopped, and where it listens.
func describeKlusterForward(started bool) func(any, *channel.Error) (string, string) {
	return func(res any, err *channel.Error) (string, string) {
		if err != nil {
			return "failed", err.Message
		}
		w, _ := res.(Written)
		var fw KlusterForward
		_ = json.Unmarshal(w.Data, &fw)
		if started {
			return "done", fmt.Sprintf("forwarding %s on %s:%d", forwardTarget(fw), fw.Address, fw.LocalPort)
		}
		return "done", "stopped forwarding " + forwardTarget(fw)
	}
}

// klusterForwardSpecs are kluster.forward.start and kluster.forward.stop,
// writes under the kluster host action answered in flai serve.
func klusterForwardSpecs() map[string]spec {
	params := func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		_, e := klusterForwardIn(raw)
		return nil, "", e
	}
	return map[string]spec{
		// kluster.forward.start: kubectl port-forward of a service port of a
		// namespace, on a free port of the host bound to 0.0.0.0 (S-0586).
		"kluster.forward.start": {action: ActionKluster, describe: describeKlusterForward(true), build: params, do: klusterForwardStartDo},
		// kluster.forward.stop: that forward stopped.
		"kluster.forward.stop": {action: ActionKluster, describe: describeKlusterForward(false), build: params, do: klusterForwardStopDo},
	}
}

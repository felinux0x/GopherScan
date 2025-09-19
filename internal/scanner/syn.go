package scanner

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	time "time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/user/pentscan/internal/types"
)

// SYNScanner implementa a varredura usando pacotes SYN raw.
type SYNScanner struct {
	iface       *net.Interface
	srcIP, dstIP net.IP
	srcPort     layers.TCPPort
	handle      *pcap.Handle
	results     sync.Map // Mapeia "ip:porta" para chan types.ScanResult
	mu          sync.Mutex
}

// NewSYNScanner cria um novo scanner SYN.
func NewSYNScanner() (*SYNScanner, error) {
	if os.Geteuid() != 0 {
		return nil, errors.New("a varredura SYN requer privilégios de root/administrador")
	}

	iface, err := getOutputInterface()
	if err != nil {
		return nil, fmt.Errorf("falha ao encontrar a interface de saída: %w", err)
	}

	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir o handle do pcap: %w", err)
	}

	ss := &SYNScanner{
		iface:   iface,
		handle:  handle,
		srcPort: 54321, // Porta de origem fixa por enquanto
	}

	go ss.listenForResponses()

	return ss, nil
}

// Scan envia um pacote SYN e espera por uma resposta ou timeout.
func (ss *SYNScanner) Scan(target types.Target, timeout time.Duration) types.ScanResult {
	resultChan := make(chan types.ScanResult, 1)
	key := fmt.Sprintf("%s:%d", target.Host, target.Port)
	ss.results.Store(key, resultChan)
	defer ss.results.Delete(key)

	ss.mu.Lock() // Bloqueia para obter o IP de origem correto
	srcIP, err := getSourceIP(ss.iface)
	if err != nil {
		ss.mu.Unlock()
		return types.ScanResult{Target: target, Status: types.StatusUnknown, Error: err.Error()}
	}
	ss.mu.Unlock()

	dstIP := net.ParseIP(target.Host)
	if dstIP == nil {
		return types.ScanResult{Target: target, Status: types.StatusUnknown, Error: "host inválido"}
	}

	if err := ss.sendSYN(srcIP, dstIP.To4(), layers.TCPPort(target.Port)); err != nil {
		return types.ScanResult{Target: target, Status: types.StatusUnknown, Error: err.Error()}
	}

	select {
	case result := <-resultChan:
		return result
	case <-time.After(timeout):
		return types.ScanResult{Target: target, Status: types.StatusFiltered}
	}
}

// sendSYN constrói e envia um único pacote SYN.
func (ss *SYNScanner) sendSYN(srcIP, dstIP net.IP, dstPort layers.TCPPort) error {
	ipLayer := &layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolTCP,
	}
	tcpLayer := &layers.TCP{
		SrcPort: ss.srcPort,
		DstPort: dstPort,
		SYN:     true,
		Window:  1024,
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}
	if err := gopacket.SerializeLayers(buf, opts, tcpLayer); err != nil {
		return err
	}

	return ss.handle.WritePacketData(buf.Bytes())
}

// listenForResponses é uma goroutine que captura e processa respostas.
func (ss *SYNScanner) listenForResponses() {
	packetSource := gopacket.NewPacketSource(ss.handle, ss.handle.LinkType())
	for packet := range packetSource.Packets() {
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp, _ := tcpLayer.(*layers.TCP)

		if tcp.DstPort != ss.srcPort {
			continue
		}

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}
		ipv4, _ := ipLayer.(*layers.IPv4)

		key := fmt.Sprintf("%s:%d", ipv4.SrcIP, tcp.SrcPort)
		if resultChan, ok := ss.results.Load(key); ok {
			result := types.ScanResult{Target: types.Target{Host: ipv4.SrcIP.String(), Port: int(tcp.SrcPort)}}

			if tcp.SYN && tcp.ACK { // SYN/ACK -> Porta Aberta
				result.Status = types.StatusOpen
			} else if tcp.RST { // RST -> Porta Fechada
				result.Status = types.StatusClosed
			}

			resultChan.(chan types.ScanResult) <- result
		}
	}
}

func (ss *SYNScanner) Close() {
	ss.handle.Close()
}

// getOutputInterface encontra a interface de rede usada para tráfego externo.
func getOutputInterface() (*net.Interface, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localIP := conn.LocalAddr().(*net.UDPAddr).IP

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.Equal(localIP) {
					return &iface, nil
				}
			}
		}
	}

	return nil, errors.New("não foi possível encontrar uma interface correspondente para o IP local")
}

// getSourceIP encontra o endereço IPv4 de uma interface.
func getSourceIP(iface *net.Interface) (net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.To4(), nil
			}
		}
	}
	return nil, errors.New("nenhum endereço IPv4 encontrado para a interface")
}
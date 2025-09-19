# GopherScan

![Go Version](https://img.shields.io/badge/go-1.21%2B-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**GopherScan** é uma ferramenta de descoberta e varredura de rede moderna, escrita em Go, projetada para profissionais de segurança (pentesters, red team) e administradores de sistema. O objetivo é combinar a velocidade de varreduras em larga escala com a precisão necessária para a pós-análise.

---

### Índice
- [Recursos Principais](#recursos-principais)
- [⚠️ Aviso de Uso Ético](#%EF%B8%8F-aviso-de-uso-ético)
- [Instalação](#instalação)
- [Uso](#uso)
  - [Opções de Linha de Comando](#opções-de-linha-de-comando)
  - [Exemplos Práticos](#exemplos-práticos)
- [Formatos de Saída](#formatos-de-saída)
- [Arquitetura](#arquitetura)
- [Como Contribuir](#como-contribuir)
- [Licença](#licença)

---

### Recursos Principais

- ⚡ **Varredura Rápida e Concorrente**: Utiliza um pool de workers para escanear milhares de portas em segundos.
- 🎯 **Múltiplos Tipos de Alvo**: Suporte para IPs, nomes de host, ranges CIDR e leitura de alvos a partir de arquivos.
- 📜 **Parsing Flexível de Portas**: Especifique portas individuais, listas separadas por vírgula e ranges (ex: `80,443,8080-8090`).
- 🕵️ **Múltiplos Modos de Scan**: Inclui `Connect Scan` (padrão) e `SYN Scan` (furtivo, requer privilégios de administrador).
- 🔍 **Banner Grabbing Inteligente**: Identifica serviços comuns (HTTP, SSH) e captura banners para análise.
- 📄 **Formatos de Saída Versáteis**: Salve os resultados em `txt` (legível por humanos), `json` (para processamento) ou `csv` (para planilhas).
- 📊 **Métricas em Tempo Real**: Exponha métricas no formato Prometheus para monitoramento do progresso da varredura.
- ⚙️ **Controle Total**: Ajuste o número de workers, a taxa de pacotes e o timeout das conexões.
- 🛑 **Cancelamento Gracioso**: Interrompa a varredura a qualquer momento com `Ctrl+C` sem perder os resultados já processados.

### ⚠️ Aviso de Uso Ético

> Esta ferramenta foi desenvolvida para fins educacionais e para uso em ambientes controlados e autorizados. A varredura de redes sem permissão explícita do proprietário é **ilegal** e antiética.
> - **NÃO USE ESTA FERRAMENTA EM REDES PÚBLICAS OU CORPORATIVAS SEM AUTORIZAÇÃO.**
> - Os desenvolvedores não se responsabilizam por qualquer dano, perda de dados ou consequências legais resultantes do mau uso desta ferramenta.
> - Ao usar o GopherScan, você concorda em utilizá-lo de forma responsável e legal.

### Instalação

Você precisa ter o **Go (versão 1.21 ou superior)** instalado.

```bash
# 1. Clone o repositório (ou use o código-fonte que você já possui)
git clone https://github.com/user/gopherscan.git
cd gopherscan

# 2. Compile o binário
go build -o gopherscan.exe ./cmd/pentscan/
```

### Uso

#### Opções de Linha de Comando

| Flag Curto | Flag Longo        | Descrição                                         | Exemplo                |
| :--------- | :---------------- | :------------------------------------------------ | :--------------------- |
| `-H`       | `--hosts`         | Lista de hosts, IPs ou CIDRs a escanear.          | `example.com,192.168.1.1` |
| `-t`       | `--targets`       | Caminho para um arquivo com alvos (um por linha). | `targets.txt`          |
| `-p`       | `--ports`         | Portas a escanear.                                | `80,443,1000-2000`     |
| `-o`       | `--output`        | Arquivo para salvar os resultados.                | `results.json`         |
|            | `--output-format` | Formato da saída (`txt`, `json`, `csv`).          | `csv`                  |
|            | `--scan-type`     | Tipo de varredura (`connect`, `syn`).             | `syn`                  |
| `-w`       | `--workers`       | Número de workers concorrentes.                   | `100`                  |
| `-r`       | `--rate`          | Limite de pacotes/conexões por segundo.           | `500`                  |
| `-T`       | `--timeout`       | Timeout para cada tentativa de conexão.           | `3s`                   |
|            | `--metrics-addr`  | Endereço para expor métricas Prometheus.          | `:9090`                |
| `-h`       | `--help`          | Exibe a mensagem de ajuda.                        |                        |

#### Exemplos Práticos

**1. Varredura rápida de portas comuns em sua máquina local:**
```bash
./gopherscan.exe --hosts 127.0.0.1 --ports 80,443,8080
```

**2. Varredura de um range de portas em múltiplos hosts, salvando em CSV:**
```bash
./gopherscan.exe --hosts 192.168.1.1,example.com --ports 1-1024 -o results.csv --output-format csv
```

**3. Ler alvos de um arquivo (incluindo CIDR) e aumentar a performance:**
Crie um arquivo `targets.txt`:
```
192.168.1.10
192.168.1.0/30
scanme.nmap.org
```
Execute o comando com mais workers e controle de taxa:
```bash
./gopherscan.exe --targets targets.txt --ports 80,8080-8090 --workers 100 --rate 500
```

**4. Varredura SYN (requer privilégios de root/administrador):**
```bash
# No Linux/macOS
sudo ./gopherscan --hosts 192.168.1.1 --ports 22,80,443 --scan-type syn

# No Windows, execute o terminal como Administrador
.\gopherscan.exe --hosts 192.168.1.1 --ports 22,80,443 --scan-type syn
```

**5. Expondo métricas Prometheus durante a varredura:**
```bash
./gopherscan.exe --hosts 127.0.0.1 --ports 1-1024 --metrics-addr :9090
# Em outro terminal ou navegador, acesse http://localhost:9090/metrics
```

### Formatos de Saída

**`txt` (Padrão)**: Formato colorido e otimizado para leitura humana, focado em exibir as informações mais relevantes.
```
[+] 127.0.0.1:80         open       http       HTTP/1.1 200 OK
[?] 127.0.0.1:443        filtered   Nenhuma conexão pôde ser feita...
```

**`json`**: Ideal para integração com outras ferramentas e scripts. Cada linha é um objeto JSON válido.
```json
{"schema_version":"1.1","target":{"Host":"127.0.0.1","Port":80},"status":1,"service_name":"http","banner":"HTTP/1.1 200 OK"}
```

**`csv`**: Perfeito para importar em planilhas ou bancos de dados.
```csv
host,port,status,service_name,banner,error
127.0.0.1,80,open,http,"HTTP/1.1 200 OK",
127.0.0.1,443,filtered,,,Nenhuma conexão pôde ser feita...
```

### Arquitetura

- **`cmd/pentscan`**: Ponto de entrada da CLI. Responsável por parsear flags e orquestrar a execução.
- **`internal/engine`**: Gerencia a concorrência (worker pool) e distribui os alvos.
- **`internal/scanner`**: Contém a interface `Scanner` e implementações para `ConnectScan` e `SYNScan`.
- **`internal/probes`**: Define a interface `Probe` e implementações para `HTTPProbe`, `SSHProbe`.
- **`internal/writer`**: Lida com a formatação e escrita dos resultados (JSON, CSV, TXT).
- **`internal/types`**: Define as estruturas de dados compartilhadas (`Target`, `ScanResult`).
- **`internal/metrics`**: Define e gerencia as métricas Prometheus.

### Como Contribuir

Contribuições são bem-vindas! Se você encontrar um bug ou tiver uma sugestão de melhoria, por favor, abra uma *issue* no repositório do projeto.

### Licença

Este projeto está sob a licença MIT.

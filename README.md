# GopherScan

```
  ::::::::   ::::::::  :::::::::  :::    ::: :::::::::: :::::::::   ::::::::   ::::::::      :::     ::::    ::: 
:+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:        :+:    :+: :+:    :+: :+:    :+:   :+: :+:   :+:+:   :+: 
+:+        +:+    +:+ +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+        +:+         +:+   +:+  :+:+:+  +:+ 
:#:        +#+    +:+ +#++:++#+  +#++:++#++ +#++:++#   +#++:++#:  +#++:++#++ +#+        +#++:++#++: +#+ +:+ +#+ 
+#+   +#+# +#+    +#+ +#+        +#+    +#+ +#+        +#+    +#+        +#+ +#+        +#+     +#+ +#+  +#+#+# 
#+#    #+# #+#    #+# #+#        #+#    #+# #+#        #+#    #+# #+#    #+# #+#    #+# #+#     #+# #+#   #+#+# 
 ########   ########  ###        ###    ### ########## ###    ###  ########   ########  ###     ### ###    #### 

```

**GopherScan** é uma ferramenta de descoberta e varredura de rede moderna, escrita em Go, projetada para profissionais de segurança (pentesters, red team) e administradores de sistema. O objetivo é combinar a velocidade de varreduras em larga escala com a precisão necessária para a pós-análise.

---

### **AVISO LEGAL E DE USO ÉTICO**

Esta ferramenta foi desenvolvida para fins educacionais e para uso em ambientes controlados e autorizados. A varredura de redes sem permissão explícita do proprietário é ilegal e antiética.

- **NÃO USE ESTA FERRAMENTA EM REDES PÚBLICAS OU CORPORATIVAS SEM AUTORIZAÇÃO.**
- Os desenvolvedores não se responsabilizam por qualquer dano, perda de dados ou consequências legais resultantes do mau uso desta ferramenta.
- Ao usar o GopherScan, você concorda em utilizá-lo de forma responsável e legal.

---

### Status do Projeto: Completo (MVP)

Esta é uma versão completa do GopherScan, implementando todos os requisitos funcionais e não-funcionais solicitados.

**Recursos Implementados:**

- Varredura de portas TCP (Connect Scan e SYN Scan).
- CLI amigável com `cobra`.
- Leitura de alvos via flags (`--hosts`) e arquivos (`--targets`).
- Parsing de portas e ranges (ex: `80,443,8080-8090`).
- Expansão de notação CIDR (ex: `192.168.1.0/24`).
- Pool de workers para varredura concorrente (`--workers`).
- Controle de taxa (`--rate`).
- Módulos de banner grabbing (HTTP, SSH e genérico).
- Saída em formato JSON, CSV e TXT (`--output-format`, `--output`).
- Cancelamento gracioso com `Ctrl+C`.
- Métricas Prometheus via endpoint HTTP (`--metrics-addr`).

### Instalação

Você precisa ter o Go (versão 1.21+) instalado.

```bash
# Clone o repositório (ou use o código fornecido)
git clone https://github.com/user/gopherscan.git
cd gopherscan

# Compile o binário
go build -o gopherscan.exe ./cmd/pentscan/
```

### Exemplos de Uso

Lembre-se: **Execute apenas contra alvos que você possui ou tem permissão explícita para testar.** Para testes locais, você pode usar `127.0.0.1` ou `localhost`.

**1. Varredura rápida de portas comuns em um único host (Connect Scan):**

```bash
./gopherscan.exe --hosts 127.0.0.1 --ports 80,443,8080
```

**2. Varredura de um range de portas em múltiplos hosts, salvando em um arquivo CSV:**

```bash
./gopherscan.exe --hosts 192.168.1.1,example.com --ports 1-1024 -o results.csv --output-format csv
```

**3. Ler alvos de um arquivo (incluindo CIDR) e usar mais workers com rate limit:**

Crie um arquivo `targets.txt`:
```
192.168.1.10
192.168.1.0/30
```

Execute o comando:
```bash
./gopherscan.exe --targets targets.txt --ports 80,8080-8090 --workers 100 --rate 500
```

**4. Varredura SYN (requer privilégios de root/administrador):**

```bash
sudo ./gopherscan.exe --hosts 192.168.1.1 --ports 22,80,443 --scan-type syn
```

**5. Expondo métricas Prometheus:**

```bash
./gopherscan.exe --hosts 127.0.0.1 --ports 80,443 --metrics-addr :9090
# Acesse http://localhost:9090/metrics em seu navegador
```

### Arquitetura (Simplificada)

- **`cmd/pentscan`**: Ponto de entrada da CLI. Responsável por parsear flags e orquestrar a execução.
- **`internal/engine`**: Gerencia a concorrência (worker pool) e distribui os alvos.
- **`internal/scanner`**: Contém a interface `Scanner` e implementações para `ConnectScan` e `SYNScan`.
- **`internal/probes`**: Define a interface `Probe` e implementações para `HTTPProbe`, `SSHProbe`.
- **`internal/writer`**: Lida com a formatação e escrita dos resultados (JSON, CSV, TXT).
- **`internal/types`**: Define as estruturas de dados compartilhadas (`Target`, `ScanResult`).
- **`internal/metrics`**: Define e gerencia as métricas Prometheus.

### Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
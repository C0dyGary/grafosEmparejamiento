# Emparejamiento en Grafos (Matching) — App de Terminal

Aplicación de terminal en Go que pide **vértices y aristas por teclado**, grafica el grafo en ASCII y calcula el **emparejamiento maximal** (greedy) y el **emparejamiento máximo** (algoritmo Blossom de Edmonds).

`MaximalMatching` — recorre las aristas una vez y toma las que no colisionan. Rápido, pero no garantiza tamaño óptimo.

`MaximumMatching` — implementa Blossom para obtener el emparejamiento de cardinalidad máxima.

## Requisitos

- Go 1.26.1 o superior (según `go.mod`)

## 1. Instalar Go

### Linux (Debian/Ubuntu)

```bash
sudo apt update
sudo apt install -y wget tar
wget https://go.dev/dl/go1.26.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Linux (Arch/Fedora)

```bash
# Arch
sudo pacman -S go

# Fedora
sudo dnf install golang
```

### Windows

1. Descargar el instalador desde https://go.dev/dl/
2. Ejecutar `go1.26.1.windows-amd64.msi`
3. Abrir una nueva terminal y verificar:

```powershell
go version
```

### macOS

```bash
# Con Homebrew
brew install go

# O descargar el .pkg desde https://go.dev/dl/
```

### Verificar instalación

```bash
go version
```

Debe mostrar algo como `go version go1.26.1 linux/amd64`.

## Estructura (código separado)

```text
emparejamiento/
├── main.go              # solo llama a terminal.Run(), sin lógica mezclada
├── go.mod
├── graph/
│   └── graph.go         # Edge, Graph, MaximalMatching, Blossom/MaximumMatching
└── terminal/
    └── terminal.go      # app de terminal: input, gráfica ASCII y resultados
```

## 2. Ejecutar este código

Desde la carpeta del proyecto (`emparejamiento/`):

```bash
# 1. Dejar el módulo consistente (no hay dependencias externas)
go mod tidy

# 2. Ejecutar la app de terminal (usa todos los paquetes)
go run .
```

## 3. Generar los ejecutables

Compílalos tú directamente en esta carpeta:

```text
emparejamiento/
├── emparejamiento      # se genera en Linux (ELF 64-bit)
└── emparejamiento.exe  # se genera para Windows (PE32+ 64-bit)
```

Uso directo sin compilar:

```bash
# Linux
./emparejamiento

# Windows (en PowerShell o cmd)
.\emparejamiento.exe
```

Para regenerarlos:

```bash
# Linux
go build -o emparejamiento .

# Windows desde Linux/macOS (compilación cruzada)
GOOS=windows GOARCH=amd64 go build -o emparejamiento.exe .
```

En Windows (compilación nativa):

```powershell
go build -o emparejamiento.exe .
.\emparejamiento.exe
```

## 4. Cómo usar la app

Al ejecutar verás este menú:

```text
=== Emparejamiento en Grafos ===
1) Ingresar grafo por teclado
2) Usar grafo de ejemplo (5 vertices)
3) Salir
```

Si eliges `1`, la app te pide:

```text
Numero de vertices (>=1): 5
Numero de aristas (>=0): 6
Ingresa cada arista como: u v   (vertices de 0 a 4)
Arista 1/6: 0 1
Arista 2/6: 0 2
Arista 3/6: 1 2
Arista 4/6: 1 3
Arista 5/6: 2 4
Arista 6/6: 3 4
```

Reglas de entrada:

- Vértices numerados de `0` a `n-1`.
- Cada arista en una línea como `u v` (ejemplo: `0 1`).
- No se admiten bucles (`1 1`) ni aristas duplicadas.

## 5. Salida esperada

La app muestra lista de adyacencia, el grafo dibujado en ASCII y ambos emparejamientos:

```text
--- Grafo: 5 vertices, 6 aristas ---
Lista de adyacencia:
  0: 1 2
  1: 0 2 3
  2: 0 1 4
  3: 1 4
  4: 2 3

Grafo (distribucion circular):
                                 0
  ...dibujo ASCII del grafo...

Leyenda: numero = vertice, . = arista, # = arista del emparejamiento maximo

Maximal matching (greedy):
  (0, 1)
  (2, 4)
  Tamano: 2

Maximum matching (Blossom):
  (0, 1)
  (2, 4)
  Tamano: 2

Vertices expuestos (sin emparejar en el maximo):
  3
  Total: 1 de 5 vertices
```

## Comandos útiles

```bash
go vet ./...      # analizar el código
gofmt -l .        # verificar formato
go run .          # ejecutar sin compilar
```

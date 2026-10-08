package terminal

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"matching/graph"
)

// Run inicia el menu interactivo de terminal.
func Run() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Emparejamiento en Grafos ===")
	fmt.Println("Maximal (greedy) vs Maximum (Blossom de Edmonds)")

	for {
		fmt.Println()
		fmt.Println("1) Ingresar grafo por teclado")
		fmt.Println("2) Usar grafo de ejemplo (5 vertices)")
		fmt.Println("3) Salir")
		fmt.Print("Elige una opcion: ")

		op := leerLinea(reader)
		switch op {
		case "1":
			g := LeerGrafo(reader)
			MostrarResultados(g)
		case "2":
			g := graph.Graph{
				Vertices: 5,
				Edges: []graph.Edge{
					{U: 0, V: 1},
					{U: 0, V: 2},
					{U: 1, V: 2},
					{U: 1, V: 3},
					{U: 2, V: 4},
					{U: 3, V: 4},
				},
			}
			MostrarResultados(g)
		case "3":
			fmt.Println("Adios.")
			return
		default:
			fmt.Println("Opcion invalida. Usa 1, 2 o 3.")
		}
	}
}

// ---------- Entrada por terminal ----------

func leerLinea(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func leerEntero(reader *bufio.Reader, prompt string, min int) int {
	for {
		fmt.Print(prompt)
		txt := leerLinea(reader)
		n, err := strconv.Atoi(txt)
		if err != nil || n < min {
			fmt.Printf("Ingresa un entero >= %d.\n", min)
			continue
		}
		return n
	}
}

// LeerGrafo pide vertices y aristas por teclado y valida el formato "u v".
func LeerGrafo(reader *bufio.Reader) graph.Graph {
	n := leerEntero(reader, "Numero de vertices (>=1): ", 1)
	m := leerEntero(reader, "Numero de aristas (>=0): ", 0)

	edges := make([]graph.Edge, 0, m)
	vistas := make(map[graph.Edge]bool)

	fmt.Println("Ingresa cada arista como: u v   (vertices de 0 a", n-1, ")")
	for len(edges) < m {
		fmt.Printf("Arista %d/%d: ", len(edges)+1, m)
		txt := leerLinea(reader)
		parts := strings.Fields(txt)
		if len(parts) != 2 {
			fmt.Println("Formato invalido. Ejemplo: 0 1")
			continue
		}
		u, err1 := strconv.Atoi(parts[0])
		v, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Ambos valores deben ser enteros.")
			continue
		}
		if u < 0 || u >= n || v < 0 || v >= n {
			fmt.Printf("Vertices fuera de rango [0, %d].\n", n-1)
			continue
		}
		if u == v {
			fmt.Println("No se permiten bucles (u == v).")
			continue
		}
		e := graph.Normalize(u, v)
		if vistas[e] {
			fmt.Println("Arista duplicada, ingresa otra.")
			continue
		}
		vistas[e] = true
		edges = append(edges, e)
	}

	sort.Slice(edges, func(i, j int) bool {
		if edges[i].U == edges[j].U {
			return edges[i].V < edges[j].V
		}
		return edges[i].U < edges[j].U
	})

	return graph.Graph{Vertices: n, Edges: edges}
}

// ---------- Salida / grafica ASCII ----------

func MostrarResultados(g graph.Graph) {
	fmt.Printf("\n--- Grafo: %d vertices, %d aristas ---\n", g.Vertices, len(g.Edges))
	imprimirListaAdyacencia(g)

	maximal := graph.MaximalMatching(g)
	maximum := graph.MaximumMatching(g)

	fmt.Println("\nGrafo (distribucion circular):")
	fmt.Println(DibujarGrafo(g, aSet(maximum)))

	fmt.Println("Leyenda: numero = vertice, . = arista, # = arista del emparejamiento maximo")
	fmt.Println("\nMaximal matching (greedy):")
	imprimirMatching(maximal)
	fmt.Println("\nMaximum matching (Blossom):")
	imprimirMatching(maximum)
	imprimirExpuestos(g, maximum)
}

func imprimirExpuestos(g graph.Graph, maximum []graph.Edge) {
	emparejado := make([]bool, g.Vertices)
	for _, e := range maximum {
		if e.U >= 0 && e.U < g.Vertices {
			emparejado[e.U] = true
		}
		if e.V >= 0 && e.V < g.Vertices {
			emparejado[e.V] = true
		}
	}
	var expuestos []int
	for v := 0; v < g.Vertices; v++ {
		if !emparejado[v] {
			expuestos = append(expuestos, v)
		}
	}
	fmt.Println("\nVertices expuestos (sin emparejar en el maximo):")
	if len(expuestos) == 0 {
		fmt.Println("  Ninguno (emparejamiento perfecto).")
		return
	}
	fmt.Print(" ")
	for _, v := range expuestos {
		fmt.Printf(" %d", v)
	}
	fmt.Printf("\n  Total: %d de %d vertices\n", len(expuestos), g.Vertices)
}

func imprimirMatching(m []graph.Edge) {
	if len(m) == 0 {
		fmt.Println("  (vacio)")
	}
	for _, e := range m {
		fmt.Printf("  (%d, %d)\n", e.U, e.V)
	}
	fmt.Printf("  Tamaño: %d\n", len(m))
}

func imprimirListaAdyacencia(g graph.Graph) {
	ady := make([][]int, g.Vertices)
	for _, e := range g.Edges {
		ady[e.U] = append(ady[e.U], e.V)
		ady[e.V] = append(ady[e.V], e.U)
	}
	fmt.Println("Lista de adyacencia:")
	for v := 0; v < g.Vertices; v++ {
		fmt.Printf("  %d:", v)
		for _, w := range ady[v] {
			fmt.Printf(" %d", w)
		}
		fmt.Println()
	}
}

func aSet(m []graph.Edge) map[graph.Edge]bool {
	s := make(map[graph.Edge]bool, len(m))
	for _, e := range m {
		s[graph.Normalize(e.U, e.V)] = true
	}
	return s
}

// DibujarGrafo dibuja vertices en circulo y aristas con Bresenham sobre una grilla de texto.
func DibujarGrafo(g graph.Graph, destacados map[graph.Edge]bool) string {
	const W = 65
	const H = 21
	grid := make([][]rune, H)
	for y := range grid {
		grid[y] = make([]rune, W)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	if g.Vertices == 0 {
		return "(grafo vacio)"
	}

	cx := float64(W) / 2.0
	cy := float64(H) / 2.0
	rx := float64(W)/2.0 - 6
	ry := float64(H)/2.0 - 2

	px := make([]int, g.Vertices)
	py := make([]int, g.Vertices)
	for v := 0; v < g.Vertices; v++ {
		var ang float64
		if g.Vertices == 1 {
			ang = -math.Pi / 2
		} else {
			ang = 2*math.Pi*float64(v)/float64(g.Vertices) - math.Pi/2
		}
		px[v] = int(math.Round(cx + rx*math.Cos(ang)))
		py[v] = int(math.Round(cy + ry*math.Sin(ang)))
	}

	// Aristas primero, vertices despues para no tapar etiquetas.
	for _, e := range g.Edges {
		ch := '.'
		if destacados[graph.Normalize(e.U, e.V)] {
			ch = '#'
		}
		trazarLinea(grid, px[e.U], py[e.U], px[e.V], py[e.V], ch)
	}

	for v := 0; v < g.Vertices; v++ {
		dibujarEtiqueta(grid, px[v], py[v], strconv.Itoa(v))
	}

	var sb strings.Builder
	for y := 0; y < H; y++ {
		sb.WriteString(strings.TrimRight(string(grid[y]), " "))
		if y < H-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func trazarLinea(grid [][]rune, x0, y0, x1, y1 int, ch rune) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	for {
		// Solo celdas vacias; las etiquetas se dibujan despues.
		if y0 >= 0 && y0 < len(grid) && x0 >= 0 && x0 < len(grid[0]) && grid[y0][x0] == ' ' {
			grid[y0][x0] = ch
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func dibujarEtiqueta(grid [][]rune, x, y int, s string) {
	start := x - len(s)/2
	for i, r := range s {
		xx := start + i
		if y >= 0 && y < len(grid) && xx >= 0 && xx < len(grid[0]) {
			grid[y][xx] = r
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

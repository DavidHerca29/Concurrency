/*
VERSIÓN RECORTADA POR ITZ PARA OMITIR DETALLES ALGORÍTMICOS
*/

/*
Notas:
1. Para una visualización correcta del gráfico de barras, por favor ejecutar el programa
en una terminal a pantalla completa.
https://www.geeksforgeeks.org/iterative-quick-sort/
https://tecadmin.net/get-current-date-time-golang/
*/

package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/lxn/win"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

type IndexValue struct {
	index int
	value int
}
type stats struct {
	intercambios  int
	comparaciones int
	evaluaciones  int
	tiempo        time.Duration
}

const (
	BAR_WIDTH       = 1
	FONT_WIDTH      = 8
	FONT_HEIGHT     = 15
	MAX_NUMBER_SIZE = 32
)

var (
	width           int = int(win.GetSystemMetrics(win.SM_CXSCREEN) / FONT_WIDTH)
	height          int = int(win.GetSystemMetrics(win.SM_CYSCREEN) / (FONT_HEIGHT))
	bubblesChart    widgets.BarChart
	quicksChart     widgets.BarChart
	heapsChart      widgets.BarChart
	selectionChart  widgets.BarChart
	insertionsChart widgets.BarChart
	m               sync.Mutex

	bubbleStats        stats // estadísticas de bubble
	quickSortStats     stats // estadísticas de quicksort
	HeapSortStats      stats // estadísticas de heapsort
	insertionSortStats stats // estadísticas de insertion
	selectionSortStats stats // estadísticas de selection

)

func main() {
	barNumber := width/(BAR_WIDTH*2) - 1
	fmt.Print("Indique la cantidad de numeros(Se recomienda " + strconv.Itoa(barNumber) + " maximo para una visualizacion correcta): ")
	var size int
	fmt.Scanln(&size)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	/*
		Generates a 3 digit number between 0-599 from the system hour
	*/
	baseSlice := make([]float64, 100)
	baseSlice = generarEnteroCLM(size, time.Now().UnixMilli()%600) // enviamos un entero entre 0-599
	//baseSlice = []float64{9, 6, 23,3,4, 4, 7, 12, 15, 22, 17}
	fmt.Println(baseSlice)
	initBubblesChart(baseSlice)
	initHeapsChart(baseSlice)
	initInsertionsChart(baseSlice)
	initSelectionsChart(baseSlice)
	initQuicksChart(baseSlice)
	ui.Render(&bubblesChart)
	ui.Render(&selectionChart)
	ui.Render(&heapsChart)
	ui.Render(&insertionsChart)
	ui.Render(&quicksChart)
	go bsChartDrawer(baseSlice)
	go quicksChartDrawer(baseSlice)
	go heapsChartDrawer(baseSlice)
	go selectionsChartDrawer(baseSlice)
	go insertionsChartDrawer(baseSlice)
	fmt.Scanln() //end until any key is pressed
	ui.Close()
}

/*
Aux function to swap two numbers
*/
func swap(a *float64, b *float64) {
	temp := *a
	*a = *b
	*b = temp
}

/*
Aux function to remove an element from a slice
*/
func remove(slice [][]int, index int) [][]int {
	return append(slice[:index], slice[index+1:]...)
}

/*func partition(slice []float64, start int, end int, pair chan []int) int {
	// CÓDIGO OMITIDO
	for i := start; i < end; i++{
		// CÓDIGO OMITIDO
		if slice[i] <= pivot{
			pair <- []int{i,index}
			// CÓDIGO OMITIDO
		}
	}
	pair <- []int{index,end}
	// CÓDIGO OMITIDO
	return index
}*/

/*
Iterative quicksort
*/ /*
func quickSort(slice []float64, size int, pair chan []int) {
	// CÓDIGO OMITIDO
	for len(stack) > 0{
		// CÓDIGO OMITIDO
		pivot := partition(slice, start, end, pair)
		// CÓDIGO OMITIDO
		if pivot-1 > start {
			stack = append(stack, []int{start,pivot-1})
		}
		// CÓDIGO OMITIDO
		if pivot+1 < end {
			stack = append(stack, []int{pivot+1,end})
		}
		// CÓDIGO OMITIDO
	}
	close(pair)
	// CÓDIGO OMITIDO
}
*/
/*
Bubblesort graphic drawer
*/
func bsChartDrawer(slice []float64) {
	bubblesChart.Data = make([]float64, len(slice))
	copy(bubblesChart.Data, slice)
	pairsChannel := make(chan []int, 2000)
	go callBubble(&slice, &bubbleStats, pairsChannel)
	for pair := range pairsChannel {
		swap(&bubblesChart.Data[pair[0]], &bubblesChart.Data[pair[1]])
		m.Lock()
		ui.Render(&bubblesChart)
		m.Unlock()
	}
	bubblesChart.Title = "BubbleSort-Finalizado-" /*+
	"Tiempo:"+strconv.FormatInt(bsTime.Milliseconds(),10)+"ms-" +
	"Swaps:"+strconv.Itoa(bsSwaps)+"-" +
	"Comparaciones:"+strconv.Itoa(bsComparisons)+"-"+
	"Iteraciones:"+strconv.Itoa(bsIterations)*/
	m.Lock()
	ui.Render(&bubblesChart)
	m.Unlock()
}

/*
Quicksort graphic drawer
*/
func quicksChartDrawer(slice []float64) {
	quicksChart.Data = make([]float64, len(slice))
	copy(quicksChart.Data, slice)
	copySlice := make([]float64, len(slice))
	copy(copySlice, slice)
	pairsChannel := make(chan []int)
	go quickSortIterative(&copySlice, 0, len(copySlice)-1, pairsChannel)
	for pair := range pairsChannel {
		swap(&quicksChart.Data[pair[0]], &quicksChart.Data[pair[1]])
		m.Lock()
		ui.Render(&quicksChart)
		m.Unlock()
	}
	quicksChart.Title = "QuickSort-Finalizado-" /*+
	"Tiempo:"+strconv.FormatInt(qsTime.Milliseconds(),10)+"ms-" +
	"Swaps:"+strconv.Itoa(qsSwaps)+"-" +
	"Comparaciones:"+strconv.Itoa(qsComparisons)+"-"+
	"Iteraciones:"+strconv.Itoa(qsIterations)*/
	m.Lock()
	ui.Render(&quicksChart)
	m.Unlock()
}
func insertionsChartDrawer(slice []float64) {
	insertionsChart.Data = make([]float64, len(slice))
	copy(insertionsChart.Data, slice)
	copySlice := make([]float64, len(slice))
	copy(copySlice, slice)
	pairsChannel := make(chan IndexValue, 2000)
	go insertion(copySlice, pairsChannel)
	for pair := range pairsChannel {
		insertionsChart.Data[pair.index] = float64(pair.value)
		//swap(&, &insertionsChart.Data[pair.index])
		m.Lock()
		ui.Render(&insertionsChart)
		m.Unlock()
	}
	insertionsChart.Title = "InsertionSort-Finalizado-" /* +
	"Tiempo:"+strconv.FormatInt(qsTime.Milliseconds(),10)+"ms-" +
	"Swaps:"+strconv.Itoa(qsSwaps)+"-" +
	"Comparaciones:"+strconv.Itoa(qsComparisons)+"-"+
	"Iteraciones:"+strconv.Itoa(qsIterations)*/
	m.Lock()
	ui.Render(&insertionsChart)
	m.Unlock()
}
func selectionsChartDrawer(slice []float64) {
	selectionChart.Data = make([]float64, len(slice))
	copy(selectionChart.Data, slice)
	copySlice := make([]float64, len(slice))
	copy(copySlice, slice)
	selChannel := make(chan []int)
	go selection(copySlice, selChannel)
	for pair := range selChannel {
		swap(&selectionChart.Data[pair[0]], &selectionChart.Data[pair[1]])
		m.Lock()
		ui.Render(&selectionChart)
		m.Unlock()
	}
	selectionChart.Title = "SelectionSort-Finalizado-" /* +
	"Tiempo:"+strconv.FormatInt(qsTime.Milliseconds(),10)+"ms-" +
	"Swaps:"+strconv.Itoa(qsSwaps)+"-" +
	"Comparaciones:"+strconv.Itoa(qsComparisons)+"-"+
	"Iteraciones:"+strconv.Itoa(qsIterations)*/
	m.Lock()
	ui.Render(&selectionChart)
	m.Unlock()
}
func heapsChartDrawer(slice []float64) {
	heapsChart.Data = make([]float64, len(slice))
	copy(heapsChart.Data, slice)
	copySlice := make([]float64, len(slice))
	copy(copySlice, slice)
	pairsChannel := make(chan []int, 2000)
	go heapsort(&copySlice, &HeapSortStats, pairsChannel)
	for pair := range pairsChannel {
		swap(&heapsChart.Data[pair[0]], &heapsChart.Data[pair[1]])
		m.Lock()
		ui.Render(&heapsChart)
		m.Unlock()
	}
	heapsChart.Title = "HeapSort-Finalizado-" /* +
	"Tiempo:"+strconv.FormatInt(qsTime.Milliseconds(),10)+"ms-" +
	"Swaps:"+strconv.Itoa(qsSwaps)+"-" +
	"Comparaciones:"+strconv.Itoa(qsComparisons)+"-"+
	"Iteraciones:"+strconv.Itoa(qsIterations)*/
	m.Lock()
	ui.Render(&heapsChart)
	m.Unlock()
}
func generateLabels(slice []float64) []string {
	var labels = make([]string, len(slice))
	for i := range slice {
		labels[i] = strconv.Itoa(i)
	}
	return labels
}

func initBubblesChart(slice []float64) {
	bubblesChart = *widgets.NewBarChart()
	bubblesChart.Data = slice
	fmt.Println(bubblesChart.Data)
	var s int
	fmt.Scanln(&s)
	bubblesChart.Title = "BubbleSort"
	bubblesChart.SetRect(0, 0, width/2-1, height/3-3)
	bubblesChart.BarWidth = BAR_WIDTH
	bubblesChart.BarGap = 0
	//bubblesChart.Labels = generateLabels(slice)
	bubblesChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	bubblesChart.BorderBottom = true
	bubblesChart.BarColors = []ui.Color{ui.ColorRed}
	bubblesChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}
func initInsertionsChart(slice []float64) {
	insertionsChart = *widgets.NewBarChart()
	insertionsChart.Data = slice
	insertionsChart.Title = "Insertion Sort"
	insertionsChart.SetRect(0, height/3-3, width/2-1, height/3*2-4)
	insertionsChart.BarWidth = BAR_WIDTH
	insertionsChart.BarGap = 0
	//insertionsChart.Labels = generateLabels(slice)
	insertionsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	insertionsChart.BorderBottom = true
	insertionsChart.BarColors = []ui.Color{ui.ColorCyan}
	insertionsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}

func initQuicksChart(slice []float64) {
	quicksChart = *widgets.NewBarChart()
	quicksChart.Data = slice
	quicksChart.Title = "QuickSort"
	quicksChart.SetRect(width/2+1, 0, width-3, height/3-3)
	quicksChart.BarWidth = BAR_WIDTH
	quicksChart.BarGap = 0
	//quicksChart.Labels = generateLabels(slice)
	quicksChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	quicksChart.BarColors = []ui.Color{ui.ColorBlue}
	quicksChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}
func initHeapsChart(slice []float64) {
	heapsChart = *widgets.NewBarChart()
	heapsChart.Data = slice
	heapsChart.Title = "Heap Sort"
	heapsChart.SetRect(width/2+1, height/3-3, width-3, height/3*2-4)
	heapsChart.BarWidth = BAR_WIDTH
	heapsChart.BarGap = 0
	//heapsChart.Labels = generateLabels(slice)
	heapsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	heapsChart.BorderBottom = true
	heapsChart.BarColors = []ui.Color{ui.ColorGreen}
	heapsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}
func initSelectionsChart(slice []float64) {
	selectionChart = *widgets.NewBarChart()
	selectionChart.Data = slice
	selectionChart.Title = "Selection Sort"
	selectionChart.SetRect(0, (height/3*2)-4, width-3, height-9)
	selectionChart.BarWidth = BAR_WIDTH
	selectionChart.BarGap = 0
	//selectionChart.Labels = generateLabels(slice)
	selectionChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	selectionChart.BorderBottom = true
	selectionChart.BarColors = []ui.Color{ui.ColorYellow}
	selectionChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}

// METODOS DE ORDENAMIENTO
func insertion(lista []float64, channel chan IndexValue) {
	largo := len(lista)
	for i := 1; i < largo; i++ { // Recorrer todos los elementos de la lista
		llave := lista[i] //Encontrar el elemento mínimo en la lista no ordenada
		j := i - 1
		insertionSortStats.evaluaciones++
		insertionSortStats.comparaciones++
		for j >= 0 && llave <= lista[j] {
			insertionSortStats.evaluaciones += 2
			insertionSortStats.comparaciones += 2
			insertionSortStats.intercambios++
			channel <- IndexValue{
				index: j + 1,
				value: int(lista[j]),
			}
			lista[j+1] = lista[j]
			j--
		}
		channel <- IndexValue{
			index: j + 1,
			value: int(llave),
		}
		insertionSortStats.intercambios++
		lista[j+1] = llave // Intercambiar el elemento mínimo encontrado con el primer elemento
	}
	close(channel)
}
func heap(lista []float64, n int, i int, s *stats, channel chan []int) {

	largest := i
	l := 2*i + 1
	r := 2*i + 2

	s.comparaciones++
	if l < n && lista[largest] < lista[l] {
		s.comparaciones++
		largest = l
	}
	s.comparaciones++
	if r < n && lista[largest] < lista[r] {
		s.comparaciones++
		largest = r
	}
	s.comparaciones++
	if largest != i {
		s.intercambios += 2
		channel <- []int{largest, i}
		/*channel <- IndexValue{
			index: i,
			value: int(lista[largest]),
		}
		channel <- IndexValue{
			index: largest,
			value: int(lista[i]),
		}*/
		lista[i], lista[largest] = lista[largest], lista[i]
		heap(lista, n, largest, s, channel)
	}

}

func heapsort(list *[]float64, s *stats, channel chan []int) {
	lista := *list
	largo := len(lista)
	for i := largo/2 - 1; i >= 0; i-- {
		s.evaluaciones++
		s.comparaciones++
		heap(lista, largo, i, s, channel)
	}
	for i := largo - 1; i > 0; i-- {
		s.evaluaciones++
		s.comparaciones++
		s.intercambios += 2
		channel <- []int{0, i}
		/*channel <- IndexValue{
			index: i,
			value: int(lista[0]),
		}
		channel <- IndexValue{
			index: 0,
			value: int(lista[i]),
		}*/
		lista[i], lista[0] = lista[0], lista[i]
		heap(lista, i, 0, s, channel)
	}
	fmt.Println("heap", lista)
	close(channel)
}
func bubble(list *[]float64, s *stats, channel chan []int) {
	lista := *list
	largo := len(lista)
	fmt.Println("largo ", largo)
	for i := 0; i < largo-1; i++ { // Recorrer todos los elementos de la lista
		// Encuentrar el elemento mínimo en la lista no ordenada
		s.evaluaciones++
		s.comparaciones++
		for j := 0; j < largo-i-1; j++ {
			s.evaluaciones++
			s.comparaciones += 2
			if lista[j] >= lista[j+1] {
				s.intercambios += 2
				channel <- []int{j, j + 1}
				/*channel <- IndexValue{
					index: j+1,
					value: int(lista[j]),
				}*/
				lista[j], lista[j+1] = lista[j+1], lista[j]
			}
		} // Intercambiar el elemento mínimo encontrado con el primer elemento
	}
	//fmt.Println(lista)
}
func selection(lista []float64, selChannel chan []int) {

	largo := len(lista)
	for i := 0; i < largo; i++ { // Recorrer todos los elementos de la lista
		selectionSortStats.evaluaciones++
		selectionSortStats.comparaciones++
		min := i //Encontrar el elemento mínimo en la lista no ordenada
		for j := i; j < largo; j++ {
			selectionSortStats.evaluaciones++
			selectionSortStats.comparaciones += 2
			if lista[j] <= lista[min] {
				min = j
			}
		}
		selChannel <- []int{i, min}
		//fmt.Println("envia ", i, min)
		/*channel <-IndexValue{
			index: i,
			value: int(lista[min]),
		}
		channel <- IndexValue{
			index: min,
			value: int(lista[i]),
		}*/
		temp := lista[min]
		lista[min] = lista[i]
		lista[i] = temp // Intercambiar el elemento mínimo encontrado con el primer elemento
		//fmt.Println("Selection",lista)
		selectionSortStats.intercambios += 1
	}
	//fmt.Println("Selection",lista)
	close(selChannel)
}

func callBubble(list *[]float64, s *stats, channel chan []int) {
	bubble(list, s, channel)
	close(channel)
}

/*
creates a N size slice with random numbers based on the linear congruential method
output: slice with N random integers
*/
func generarEnteroCLM(tamanoArreglo int, time int64) []float64 {
	respuesta := make([]float64, tamanoArreglo, tamanoArreglo)
	// Declaramos los  probables valores para la semilla que se va a utilizar
	probablesSemillas := []int{11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101}

	// Elegimos la semilla que será un valor primo entre el 11 y el 101 que se encuentra en
	// el arreglo previamente definido.
	index := int(time) % len(probablesSemillas)
	semilla := probablesSemillas[index]

	// Declaramos constantes del algoritmo
	PERIODO := int(math.Pow(2, 32))
	const (
		INCREMENTO    int = 11
		MULTIPLICADOR int = 8121
	)

	for num := 0; num < tamanoArreglo; num++ {
		semilla = (semilla*MULTIPLICADOR + INCREMENTO) % PERIODO
		respuesta[num] = float64(semilla % 30) // hacemos este modulo para tener el rango de valores de 0-29
	}

	return respuesta
}

func partition(array *[]float64, l int, h int, channel chan []int) int {
	arr := *array
	i := l - 1
	x := arr[h]

	for j := l; j <= h-1; j++ {
		if arr[j] <= x {
			// increment index of smaller element
			i = i + 1
			channel <- []int{i, j}
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	channel <- []int{i + 1, h}
	arr[i+1], arr[h] = arr[h], arr[i+1]
	return i + 1
}

// Function to do Quick sort
// arr[] --> Array to be sorted,
// l  --> Starting index,
// h  --> Ending index

func quickSortIterative(array *[]float64, l int, h int, channel chan []int) {
	// Create an auxiliary stack
	stack := make([]int, h-l+1)
	// initialize top of stack
	top := -1
	// push	initial	values of l and	h to stack
	top = top + 1
	stack[top] = l
	top = top + 1
	stack[top] = h
	// Keep popping from stack while is not empty
	for top >= 0 {
		// Pop h and l
		h = stack[top]
		top = top - 1
		l = stack[top]
		top = top - 1
		// Set pivot element at its correct position in
		// sorted array
		p := partition(array, l, h, channel)
		// If there are elements on left side of pivot,
		// then push left side to stack
		if p-1 > l {
			top = top + 1
			stack[top] = l
			top = top + 1
			stack[top] = p - 1
		}
		// If there are elements on right side of pivot,
		// then push right side to stack
		if p+1 < h {
			top = top + 1
			stack[top] = p + 1
			top = top + 1
			stack[top] = h
		}
	}
	close(channel)
}

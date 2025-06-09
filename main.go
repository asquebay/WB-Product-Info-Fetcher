package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Структуры данных для парсинга ответа
type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Brand      string  `json:"brand"`
	PriceU     int     `json:"priceU"`
	SalePriceU int     `json:"salePriceU"`
	Rating     float64 `json:"rating"`
	Feedbacks  int     `json:"feedbacks"`
	TotalQty   int     `json:"totalQuantity"`
}

type Response struct {
	Data struct {
		Products []Product `json:"products"`
	} `json:"data"`
}

// структура для передачи результата
type Result struct {
	Product Product
	RawJSON string
	Error   error
}

func fetchProductInfo(ctx context.Context, article string) Result {
	url := fmt.Sprintf("https://card.wb.ru/cards/detail?&dest=-1257786&nm=%s", article)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{Error: err}
	}

	// установка заголовков для сокрытия запросов (обход блокировки на бота)
	// это данные с моего браузера
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:137.0) Gecko/20100101 Firefox/137.0")
	req.Header.Set("Host", "card.wb.ru")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{Error: err}
	}
	defer resp.Body.Close()

	// распаковка gzip (т.к. я обращаюсь к бинарнику)
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return Result{Error: err}
		}
		defer reader.Close()
	} else {
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return Result{Error: err}
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Result{Error: err}
	}

	if len(response.Data.Products) == 0 {
		return Result{Error: fmt.Errorf("нет данных о товаре")}
	}

	return Result{
		Product: response.Data.Products[0],
		RawJSON: string(body),
	}
}

func main() {
	// проверка аргументов командной строки
	switch len(os.Args) {
	case 1:
		fmt.Fprintf(os.Stderr, "Usage: %s <article_number> [field]\n", os.Args[0])
		os.Exit(1)
	case 2:
		// только артикул - ок
	case 3:
		// артикул и поле - проверка
		allowedFields := map[string]bool{
			"body":      true,
			"name":      true,
			"price":     true,
			"salePrice": true,
			"rating":    true,
		}
		if !allowedFields[os.Args[2]] {
			fmt.Fprintln(os.Stderr, "Error: Field not found. Expected: body, name, price, salePrice, rating\nОшибка: Поле не найдено. Ожидалось: body, name, price, salePrice, rating")
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Error: Too Many Arguments. Expected: 1 or 2 arguments\nОшибка: Слишком много аргументов. Ожидалось: 1 или 2 аргумента")
		os.Exit(1)
	}

	article := os.Args[1]
	var requestedField string
	if len(os.Args) == 3 {
		requestedField = os.Args[2]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// запрос
	result := fetchProductInfo(ctx, article)
	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", result.Error)
		os.Exit(1)
	}

	// обработка результатов
	price := float64(result.Product.PriceU) / 100.0
	salePrice := float64(result.Product.SalePriceU) / 100.0

	if requestedField == "" {
		fmt.Printf("Название: %s\nЦена: %.2f ₽\nАкционная цена: %.2f ₽\nРейтинг: %.1f ★\n",
			result.Product.Name, price, salePrice, result.Product.Rating)
		return
	}

	switch requestedField {
	case "body":
		fmt.Println(result.RawJSON)
	case "name":
		fmt.Println(result.Product.Name)
	case "price":
		fmt.Printf("%.2f\n", price)
	case "salePrice":
		fmt.Printf("%.2f\n", salePrice)
	case "rating":
		fmt.Printf("%.1f\n", result.Product.Rating)
	}
}

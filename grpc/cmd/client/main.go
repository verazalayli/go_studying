package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// gRPC-клиентская библиотека.
	// Она умеет устанавливать HTTP/2 соединение, кодировать/декодировать
	// сообщения, отправлять RPC-запросы и получать ответы.
	"google.golang.org/grpc"

	// Это ПАКЕТ СО СГЕНЕРИРОВАННЫМИ ТИПАМИ из твоего proto-файла.
	// protoc с плагинами создал в нём:
	//  - интерфейс NoteServiceClient (клиентский "стаб"),
	//  - структуры сообщений (CreateNoteRequest/Response и т.д.),
	//  - код сериализации/десериализации protobuf.
	"github.com/verazalayli/go_studying/grpc/proto/pb"
)

func main() {
	// 1) Устанавливаем КЛИЕНТСКОЕ соединение с gRPC-сервером.
	//
	//   grpc.Dial(...) не "звонит" напрямую в функцию сервера. Он:
	//   - создаёт клиентский канал (ClientConn),
	//   - настраивает HTTP/2 транспорт,
	//   - готовит механизм вызова RPC (отправки/получения сообщений).
	//
	// Параметры:
	//  - WithInsecure(): используем небезопасное соединение (без TLS).
	//      В проде лучше использовать TLS: grpc.WithTransportCredentials(credentials.NewTLS(...)).
	//      WithInsecure() помечен как deprecated, но для локального примера это ок.
	//  - WithBlock(): Dial будет БЛОКИРОВАТЬСЯ, пока не установит соединение (или не выйдет по таймауту).
	//  - WithTimeout(3s): общий таймаут на установление соединения (вместе с WithBlock).
	//
	// ФАКТИЧЕСКИ: здесь создаётся и настраивается HTTP/2 клиент, открывается TCP-сокет
	// к localhost:50051, договаривается протокол gRPC поверх HTTP/2.
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(3*time.Second),
	)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	// Важно закрыть соединение, когда клиент завершает работу:
	// закроются все TCP/HTTP2 ресурсы, воркеры, пула коннекшенов и т.д.
	defer conn.Close()

	// 2) Создаём КЛИЕНТСКИЙ СТАБ (proxy-объект).
	//
	//   pb.NewNoteServiceClient(conn) возвращает объект, у которого есть методы
	//   CreateNote / GetNote / ListNotes и т.д.
	//
	// Что делает стаб при вызове метода?
	//   - берёт вашу Go-структуру запроса (например, *pb.CreateNoteRequest),
	//   - кодирует её в бинарный protobuf,
	//   - формирует gRPC-запрос (HTTP/2 фреймы), шлёт по открытому соединению на сервер,
	//   - ждёт ответ, читает фреймы, декодирует protobuf-ответ в Go-структуру,
	//   - возвращает её вам (или ошибку с gRPC-кодом).
	client := pb.NewNoteServiceClient(conn)

	// 3) Создаём контекст с ДЕДЛАЙНОМ для всех RPC ниже.
	//
	//   Контекст передаётся в каждый RPC-метод, и gRPC:
	//   - добавляет из него deadline в заголовки,
	//   - отменяет RPC на стороне клиента, если время вышло,
	//   - сервер тоже видит deadline и может завершить обработку раньше.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 4) ВЫЗОВ CreateNote дважды.
	//
	//   mustCreate ниже (см. функцию) вызывает client.CreateNote(ctx, req).
	//   Под капотом:
	//     - стаб сериализует req в protobuf-байты,
	//     - отправляет unary RPC на сервер (HTTP/2),
	//     - читает unary response,
	//     - десериализует в *pb.CreateNoteResponse.
	created1 := mustCreate(ctx, client, "First", "Hello world")
	created2 := mustCreate(ctx, client, "Second", "Another note")

	// 5) ВЫЗОВ GetNote — получение по id.
	//
	//   Здесь мы создаём Go-структуру запроса (*pb.GetNoteRequest) и передаём её
	//   в сгенерированный стаб. Всё остальное (сериализация, HTTP/2, десериализация)
	//   делает библиотека.
	got, err := client.GetNote(ctx, &pb.GetNoteRequest{Id: created1.GetNote().GetId()})
	if err != nil {
		// Если сервер вернул ошибку (например, codes.NotFound), стаб вернёт error,
		// в которой зашит gRPC status. Его можно разобрать через status.FromError(err).
		log.Fatalf("GetNote failed: %v", err)
	}

	// На этом этапе got — это уже РАСПАРСЕННЫЙ ответ (*pb.GetNoteResponse),
	// внутри которого protobuf-сообщение Note, но представленное как обычная Go-структура.
	fmt.Printf("GetNote: %+v\n", got.GetNote())

	// 6) ВЫЗОВ ListNotes — запрос списка.
	list, err := client.ListNotes(ctx, &pb.ListNotesRequest{})
	if err != nil {
		log.Fatalf("ListNotes failed: %v", err)
	}
	fmt.Println("ListNotes:")
	for _, n := range list.GetNotes() {
		// Каждая n — это *pb.Note (Go-структура, полученная после десериализации protobuf).
		fmt.Printf("- %s | %s\n", n.GetId(), n.GetTitle())
	}

	// 7) Повторно получаем вторую заметку — ещё один unary RPC.
	got2, _ := client.GetNote(ctx, &pb.GetNoteRequest{Id: created2.GetNote().GetId()})
	fmt.Printf("GetNote(2): %s: %s\n", got2.GetNote().GetId(), got2.GetNote().GetTitle())
}

// mustCreate — небольшая обёртка, чтобы не дублировать однотипный код.
// Здесь хорошо видно, ЧТО именно мы "отправляем" и ЧТО "получаем".
func mustCreate(ctx context.Context, c pb.NoteServiceClient, title, content string) *pb.CreateNoteResponse {
	// 1) СФОРМИРОВАТЬ ЗАПРОС:
	//    Мы создаём обычную Go-структуру *pb.CreateNoteRequest.
	//    Её поля соответствуют полям message CreateNoteRequest в proto.
	req := &pb.CreateNoteRequest{
		Title:   title,
		Content: content,
	}

	// 2) ВЫЗВАТЬ RPC:
	//    c.CreateNote(ctx, req) — это вызов МЕТОДА КЛИЕНТСКОГО СТАБА, а не "функции на сервере" напрямую.
	//    Стаб:
	//      - сериализует req → protobuf-байты,
	//      - отправляет по HTTP/2 в уже установленном соединении (conn),
	//      - ждёт ответ,
	//      - декодирует protobuf-ответ в *pb.CreateNoteResponse.
	resp, err := c.CreateNote(ctx, req)
	if err != nil {
		// Если сервер вернул ошибку (например, InvalidArgument при пустом title),
		// здесь будет error с gRPC статусом.
		log.Fatalf("CreateNote failed: %v", err)
	}

	// 3) ИСПОЛЬЗОВАТЬ ОТВЕТ:
	//    resp — это уже готовая Go-структура (десериализованный protobuf),
	//    можно спокойно читать поля.
	fmt.Printf("Created: %s => %s\n", resp.GetNote().GetId(), resp.GetNote().GetTitle())
	return resp
}

package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	// gRPC серверная библиотека (HTTP/2 транспорт, маршрутизация RPC, кодеки и т.д.)
	"google.golang.org/grpc"

	// Наш входной адаптер транспорта: gRPC-обработчик сервиса заметок.
	grpch "github.com/verazalayli/go_studying/grpc/pkg/handler/grpc"
	// Репозиторий в памяти — реализация интерфейса хранилища (порт прикладного слоя).
	"github.com/verazalayli/go_studying/grpc/pkg/repository/memory"
	// Прикладной слой (use cases): бизнес-логика и интерфейс порта NoteRepository.
	"github.com/verazalayli/go_studying/grpc/pkg/service"
)

func main() {
	// 1) Выбираем порт, на котором будет слушать gRPC-сервер.
	//    Можно задать через переменную окружения PORT.
	port := 50051
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}

	// 2) КОМПОЗИЦИЯ ЗАВИСИМОСТЕЙ (Composition Root).
	//    Склеиваем слои строго «снаружи вовнутрь»:
	//    transport(gRPC handler) -> service(use cases) -> repository(хранилище).
	//
	//    ↓ Хранилище: in-memory реализация. В проде легко заменить на Postgres/Mongo
	//      просто подставив другой пакет, реализующий тот же интерфейс service.NoteRepository.
	repo := memory.NewNoteRepo()

	//    ↓ Прикладной слой (use cases): инкапсулирует бизнес-правила.
	//      Он знает ТОЛЬКО про абстрактный NoteRepository (порт), а не про конкретную БД.
	svc := service.NewNoteService(repo)

	//    ↓ Транспортный адаптер: gRPC-хендлер, который:
	//      - получает protobuf-запросы,
	//      - маппит в доменные типы и вызывает svc,
	//      - маппит доменные результаты обратно в protobuf-ответы.
	handler := grpch.NewNoteHandler(svc)

	// 3) Создаём gRPC-сервер.
	//    Вызов grpc.NewServer() настраивает серверный рантайм:
	//     - HTTP/2 обработку фреймов,
	//     - регистрацию сервисов (ниже),
	//     - опционально интерсепторы, кредитный контроль, лимиты и т.д. (можно передавать опции).
	grpcServer := grpc.NewServer()

	// 4) Регистрируем наш gRPC-сервис в сервере.
	//    Внутри grpch.Register(...) вызывается сгенерённая функця pb.RegisterNoteServiceServer,
	//    которая "учит" gRPC-рунтайм: если придёт RPC NoteService.XYZ —
	//    дернуть соответствующий метод у нашего handler (NoteHandler).
	grpch.Register(grpcServer, handler)

	// 5) Открываем TCP-слушатель.
	//    net.Listen создаёт сокет и начинает слушать порт :<port>.
	//    Дальше grpcServer.Serve(lis) примет этот listener и будет:
	//      - принимать TCP соединения,
	//      - апгрейдить до HTTP/2 (h2c, если без TLS),
	//      - читать gRPC-фреймы,
	//      - диспатчить их в нужные зарегистрированные методы обработчика.
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	// 6) Грейсфул-шатдаун по сигналам ОС (Ctrl+C, docker stop, Kubernetes SIGTERM и т.п.).
	//    signal.NotifyContext вернёт контекст, который закроется при получении SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 7) Фоновая горутина, которая ждёт отмены контекста (сигнала) и мягко останавливает сервер.
	//    GracefulStop:
	//      - перестаёт принимать новые соединения,
	//      - ждёт завершения активных RPC,
	//      - закрывает слушатели и соединения корректно.
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	// 8) Запускаем главный цикл gRPC-сервера.
	//    Serve блокируется и:
	//      - принимает входящие соединения/стримы,
	//      - для каждого unary RPC:
	//          * читает HTTP/2 DATA фреймы, получает бинарный protobuf,
	//          * десериализует в *pb.<Request>,
	//          * вызывает соответствующий метод нашего handler'а,
	//          * получает от него *pb.<Response> или error,
	//          * сериализует resp в protobuf, отправляет в HTTP/2 ответ,
	//          * проставляет gRPC status (OK/ошибка) и метаданные.
	//
	//    Если вернулась ошибка — логируем фатально (обычно это проблемы на уровне listener'а).
	log.Printf("gRPC server starting on :%d\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
	log.Println("gRPC server stopped")
}

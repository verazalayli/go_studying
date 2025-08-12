package grpc

import (
	"context"
	"time"

	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/verazalayli/go_studying/grpc/pkg/service"
	"github.com/verazalayli/go_studying/grpc/proto/pb"
)

// NoteHandler — gRPC-обработчик, ничего не знает о репозитории, общается с сервисом.
type NoteHandler struct {
	pb.UnimplementedNoteServiceServer
	svc service.NoteService
}

func NewNoteHandler(svc service.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func Register(grpcServer *gogrpc.Server, h *NoteHandler) {
	pb.RegisterNoteServiceServer(grpcServer, h)
}

func (h *NoteHandler) CreateNote(ctx context.Context, req *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	n, err := h.svc.Create(ctx, req.GetTitle(), req.GetContent())
	if err != nil {
		switch err {
		case service.ErrBadRequest:
			return nil, status.Error(codes.InvalidArgument, "title is required")
		default:
			return nil, status.Errorf(codes.Internal, "create failed: %v", err)
		}
	}
	return &pb.CreateNoteResponse{Note: toPB(n)}, nil
}

func (h *NoteHandler) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	n, err := h.svc.Get(ctx, req.GetId())
	if err != nil {
		if err == service.ErrNotFound {
			return nil, status.Error(codes.NotFound, "note not found")
		}
		return nil, status.Errorf(codes.Internal, "get failed: %v", err)
	}
	return &pb.GetNoteResponse{Note: toPB(n)}, nil
}

func (h *NoteHandler) ListNotes(ctx context.Context, _ *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	list, err := h.svc.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	out := make([]*pb.Note, 0, len(list))
	for _, n := range list {
		out = append(out, toPB(n))
	}
	return &pb.ListNotesResponse{Notes: out}, nil
}

func toPB(n service.Note) *pb.Note {
	return &pb.Note{
		Id:        n.ID,
		Title:     n.Title,
		Content:   n.Content,
		CreatedAt: n.CreatedAt.Unix(),
	}
}

func fromPB(m *pb.Note) service.Note {
	return service.Note{
		ID:        m.GetId(),
		Title:     m.GetTitle(),
		Content:   m.GetContent(),
		CreatedAt: time.Unix(m.GetCreatedAt(), 0),
	}
}

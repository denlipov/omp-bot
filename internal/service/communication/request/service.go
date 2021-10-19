package request

import (
        comm "github.com/denlipov/omp-bot/internal/model/communication"
        "errors"
)

type RequestService interface {
        Describe(requestID uint64) (*comm.Request, error)
        List(cursor uint64, limit uint64) ([]comm.Request, error)
        Create(req comm.Request) (uint64, error)
        Update(requestID uint64, request comm.Request) error
        Remove(requestID uint64) (bool, error)
}


type DummyRequestService struct {
        requests map[uint64]comm.Request
        maxIdx uint64
}


func NewDummyRequestService() *DummyRequestService {
        return &DummyRequestService{
                requests: map[uint64]comm.Request{
                        0: comm.Request{
                                0, "request", "Вася", "no text",
                        },
                        1: comm.Request{
                                1, "request", "Петя", "no text",
                        },
                        2: comm.Request{
                                2, "request", "Миша", "no text",
                        },
                        3: comm.Request{
                                3, "request", "Вова", "no text",
                        },
                        4: comm.Request{
                                4, "request", "Дуся", "no text",
                        },
                        5: comm.Request{
                                5, "request", "Игорь", "no text",
                        },
                },
                maxIdx: 6,
        }
}


func (s *DummyRequestService) Describe(requestID uint64) (*comm.Request, error) {
        if r, exists := s.requests[requestID]; exists {
                return &r, nil
        } else {
                return nil, errors.New("Req id out of range")
        }
}


func (s *DummyRequestService) List(cursor uint64, limit uint64) ([]comm.Request, error) {
        result := make([]comm.Request, 0, limit)
        var totalProcessed uint64 = 0
        for i := cursor; i < s.maxIdx && totalProcessed < limit; i++ {
                if req, exists := s.requests[i]; exists {
                        result = append(result, req)
                        totalProcessed++
                }
        }
        return result, nil
}


func (s *DummyRequestService) Create(req comm.Request) (uint64, error) {
        req.Id = s.maxIdx
        s.requests[s.maxIdx] = req
        resIdx := s.maxIdx
        s.maxIdx++
        return resIdx, nil
}


func (s *DummyRequestService) Update(requestID uint64, request comm.Request) error {
        if _, exists := s.requests[requestID]; exists {
                request.Id = requestID
                s.requests[requestID] = request
                return nil
        } else {
                return errors.New("Req id out of range")
        }
}


func (s *DummyRequestService) Remove(requestID uint64) (bool, error) {
        if _, exists := s.requests[requestID]; exists {
                delete(s.requests, requestID)
                return true, nil
        } else {
                return false, errors.New("Req id out of range")
        }
}

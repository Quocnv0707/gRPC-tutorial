package main

import (
	"context"
	"fmt"
	"gRPC-tutorial/caculator/caculatorpb"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{
	caculatorpb.CaculatorServiceServer
}

const (
	port = ":9090"
)

func (*server)Sum(ctx context.Context, req *caculatorpb.SumRequest) (*caculatorpb.SumResponse, error){
	log.Println("Sum called")
	resp := &caculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}
	return resp, nil
}

func (*server) PrimeNumberDecomposition(req *caculatorpb.PNDRequest, stream caculatorpb.CaculatorService_PrimeNumberDecompositionServer) error {
	log.Println("PrimeNumberDecomposition called")
	k := int32(2)
	N := req.GetNumber()
	log.Printf("first value of k is: %d", k)
	for N > 1 {
		if N % k == 0 {
			N /= k
			stream.Send(&caculatorpb.PNDResponse{
				Result: k,
			})
			time.Sleep(time.Millisecond * 1000)
		} else {
			k++
			log.Printf("k increate to: %d", k)
		}
	}
	
	return nil
}

func (*server) Average(stream caculatorpb.CaculatorService_AverageServer) error{
	log.Println("Average called")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF{
			//tính trung bình trong đây
			resp := caculatorpb.AverageResponse{
				Result: total / float32(count),
			}
			return stream.SendAndClose(&resp)
		}
		if err != nil {
			log.Fatalf("error while Recv average: %v", err)
			return err
		}
		log.Printf("receive req %v", req)
		total += float32(req.GetNumber())
		count++
	}
}

func (*server) FindMax(stream caculatorpb.CaculatorService_FindMaxServer) error{
	log.Println("Average called")
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Request ended..")
			return nil
		}
		if err != nil {
			log.Fatalf("err while Recv request: %v", err)
		}
		num := req.GetNum()
		log.Printf("recv num: %d", num)
		if num > max {
			max = num;
		}
		err = stream.Send(&caculatorpb.FindMaxResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("send max err: %v", err)
			return err
		}
	}
}

func (*server) SumWithDeadline(ctx context.Context, req *caculatorpb.SumRequest) (*caculatorpb.SumResponse, error){
	return &caculatorpb.SumResponse{Result: 0}, nil
}

func (*server) Square(ctx context.Context, req *caculatorpb.SquareRequest) (*caculatorpb.SquareResponse, error){
	log.Println("Square API called")
	num := req.GetNum()
	if num < 0 {
		log.Printf("%d < 0, return InvalidArgument", num)
		return nil, status.Errorf(codes.InvalidArgument, "Expect num > 0, req num was %d", num)
	}
	return &caculatorpb.SquareResponse{
		SquareRoot: math.Sqrt(float64(num)),
		}, nil
}


func main() {
	lis, err := net.Listen("tcp", "0.0.0.0"+port)
	if err != nil {
		log.Fatalf("err while create listen: %v", err)
	}

	s := grpc.NewServer()

	caculatorpb.RegisterCaculatorServiceServer(s, &server{})

	fmt.Println("Server is running at %v", lis.Addr())

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("err while serve: %v ", err)
	}
}
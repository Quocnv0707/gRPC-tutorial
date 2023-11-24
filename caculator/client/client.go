package main

import (
	"context"
	"gRPC-tutorial/caculator/caculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const port = ":9090"

func main() {
	// conn, err := grpc.Dial("localhost:" + port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost" + port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("err while dial %v ", err)
	}
	defer conn.Close()

	client := caculatorpb.NewCaculatorServiceClient(conn)

	// log.Printf("service client: %v", client)
	// callSum(client)
	// callPrimeNumberDecomposition(client)
	// callAverage(client)
	// callFindMax(client)
	callSquareRoot(client, -9)

}

func callSum(c caculatorpb.CaculatorServiceClient){
	log.Println("calling sum api")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	resp, err := c.Sum(ctx, &caculatorpb.SumRequest{
		Num1: 26,
		Num2: 5,
	})
	if err != nil {
		log.Fatalf("err while call Sum function: %v", err)
	}

	log.Printf("sum api response: %v", resp.GetResult())
}

func callPrimeNumberDecomposition(c caculatorpb.CaculatorServiceClient){
	log.Println("calling PrimeNumberDecomposition api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &caculatorpb.PNDRequest{
		Number: 120,
	})
	if err != nil {
		log.Fatalf("err while call PrimeNumberDecomposition api: %v", err)
	}
	for{
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("finished api")
			return
		}

		if err != nil {
			log.Fatalf("err while call PrimeNumberDecomposition api: %v", err)
		}
		log.Printf("prime number: %d", resp.GetResult())
	}
}

func callAverage(c caculatorpb.CaculatorServiceClient){
	log.Println("calling Average api")
	stream, err := c.Average(context.Background())
	if err != nil {
		log. Fatalf("call average err: %v", err)
	}
	listReq := []caculatorpb.AverageRequest{
		caculatorpb.AverageRequest{Number: 5,},
		caculatorpb.AverageRequest{Number: 9,},
		caculatorpb.AverageRequest{Number: 19,},
		caculatorpb.AverageRequest{Number: 69,},
		caculatorpb.AverageRequest{Number: 99,},
	}

	for _, req := range listReq{
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("average request error: %v", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("average response error: %v", err)
	}
	log.Printf("average response: %v", resp)
}

func callFindMax(c caculatorpb.CaculatorServiceClient){
	log.Println("calling FindMax api")
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("call Findmax err: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		listReq := []caculatorpb.FindMaxRequest{
			caculatorpb.FindMaxRequest{Num: 5,},
			caculatorpb.FindMaxRequest{Num: 9,},
			caculatorpb.FindMaxRequest{Num: 19,},
			caculatorpb.FindMaxRequest{Num: 69,},
			caculatorpb.FindMaxRequest{Num: 9,},
		}
		for _, req := range listReq{
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("findMax request error: %v", err)
			}
			time.Sleep(time.Millisecond * 500)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("Ending call findmax api")
				break
			}
			if err != nil {
				log.Fatalf("err while recv: %v", err)
				break
			}
			log.Printf("Max = %d", resp.GetMax())			
		}
		close(waitc)
	}()
	<- waitc
}

func callSquareRoot(c caculatorpb.CaculatorServiceClient, num int32){
	log.Println("calling SquareRoot api")
	resp, err := c.Square(context.Background(), &caculatorpb.SquareRequest{
		Num: num,
	})
	if err != nil {
		log.Println("err while call Square: %v", err)
		if errStatus, ok := status.FromError(err); ok {
			log.Printf("err Code: %d\nerr Message: %v", errStatus.Code(), errStatus.Message())
			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("Invalid num: %d", num)
				return
			}
		}
	}
	log.Printf("square root response: %f", resp.GetSquareRoot())
}
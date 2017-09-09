package esp4g_test

import (
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"io"
)

type CalcService struct {
	allRequests [][]*calc.Operand
	allResponses []calc.Sum
}

func (c *CalcService)AddAll(stream calc.CalcService_AddAllServer) error {
	all := []*calc.Operand{}
	sum := calc.Sum{Z: 0}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		all = append(all, req)
		sum.Z = sum.Z + req.GetX() + req.GetY()
	}

	if err := stream.SendAndClose(&sum); err != nil {
		return err
	}

	c.allRequests = append(c.allRequests, all)
	c.allResponses = append(c.allResponses, sum)

	return nil
}

func (c *CalcService)AddDeffered(req *calc.OperandList, stream calc.CalcService_AddDefferedServer) error {
	return nil
}

func (c *CalcService)AddAsync(stream calc.CalcService_AddAsyncServer) error {
	return nil
}

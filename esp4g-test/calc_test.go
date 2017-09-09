package esp4g_test

import (
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"io"
)

type CalcService struct {
	allRequests [][]*calc.Operand
	allResponses []calc.Sum
	defferedRequests []*calc.OperandList
	defferedResponses [][]calc.Sum
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
	res := []calc.Sum{}
	for _, operand := range req.GetOperand() {
		sum := calc.Sum{Z: operand.GetX() + operand.GetY()}
		if err := stream.Send(&sum); err != nil {
			return err
		}
		res = append(res, sum)
	}

	c.defferedRequests = append(c.defferedRequests, req)
	c.defferedResponses = append(c.defferedResponses, res)

	return nil
}

func (c *CalcService)AddAsync(stream calc.CalcService_AddAsyncServer) error {
	return nil
}

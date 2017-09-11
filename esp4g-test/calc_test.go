package esp4g_test

import (
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"io"
)

type CalcService struct {
	lastAllRequests       []*calc.Operand
	lastAllResponse       *calc.Sum
	lastDefferedRequest   *calc.OperandList
	lastDefferedResponses []*calc.Sum
	lastAsyncRequests     []*calc.Operand
	lastAsyncResponses    []*calc.Sum
}

func (c *CalcService)AddAll(stream calc.CalcService_AddAllServer) error {
	c.lastAllRequests = []*calc.Operand{}
	c.lastAllResponse = &calc.Sum{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		c.lastAllRequests = append(c.lastAllRequests, req)
		c.lastAllResponse.Z = c.lastAllResponse.Z + req.GetX() + req.GetY()
	}

	if err := stream.SendAndClose(c.lastAllResponse); err != nil {
		return err
	}

	return nil
}

func (c *CalcService)AddDeffered(req *calc.OperandList, stream calc.CalcService_AddDefferedServer) error {
	c.lastDefferedResponses = []*calc.Sum{}

	for _, operand := range req.GetOperand() {
		sum := &calc.Sum{Z: operand.GetX() + operand.GetY()}
		if err := stream.Send(sum); err != nil {
			return err
		}
		c.lastDefferedResponses = append(c.lastDefferedResponses, sum)
	}

	c.lastDefferedRequest = req

	return nil
}

func (c *CalcService)AddAsync(stream calc.CalcService_AddAsyncServer) error {
	c.lastAsyncRequests = []*calc.Operand{}
	c.lastAsyncResponses = []*calc.Sum{}

	for {
		operand, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		sum := &calc.Sum{Z: operand.GetX() + operand.GetY()}
		if err := stream.Send(sum); err != nil {
			return err
		}

		c.lastAsyncRequests = append(c.lastAsyncRequests, operand)
		c.lastAsyncResponses = append(c.lastAsyncResponses, sum)
	}

	return nil
}

func makeTestCase() ([]*calc.Operand, *calc.OperandList, *calc.Sum, []*calc.Sum) {
	i := int64(0)

	as := []*calc.Operand{}
	a := &calc.OperandList{}
	b := &calc.Sum{}
	bs := []*calc.Sum{}

	for i < 5 {
		x, y := i * 2, (i * 2) + 1

		req := &calc.Operand{X: x, Y: y}

		as = append(as, req)
		a.Operand = append(a.Operand, req)

		b.Z = b.Z + x + y
		bs = append(bs, &calc.Sum{Z: x + y})

		i = i + 1
	}

	return as, a, b, bs
}

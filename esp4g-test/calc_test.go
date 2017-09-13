package esp4g_test

import (
	calc "github.com/nokamoto/esp4g/examples/calc/protobuf"
	"io"
	"sync"
)

type CalcService struct {
	_lastAllRequests       []*calc.Operand
	_lastAllResponse       *calc.Sum
	_lastDefferedRequest   *calc.OperandList
	_lastDefferedResponses []*calc.Sum
	_lastAsyncRequests     []*calc.Operand
	_lastAsyncResponses    []*calc.Sum

	mu *sync.Mutex
}

func (c *CalcService)lastAllRequests() []*calc.Operand {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._lastAllRequests
}

func (c *CalcService)lastAllResponse() *calc.Sum {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._lastAllResponse
}

func (c *CalcService)lastDefferedRequest() *calc.OperandList {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._lastDefferedRequest
}

func (c *CalcService)lastDefferedResponses() []*calc.Sum {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._lastDefferedResponses
}

func (c *CalcService)lastAsyncRequests() []*calc.Operand {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._lastAsyncRequests
}

func (c *CalcService)lastAsyncResponses() []*calc.Sum {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c._lastAsyncResponses
}

func (c *CalcService)AddAll(stream calc.CalcService_AddAllServer) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c._lastAllRequests = []*calc.Operand{}
	c._lastAllResponse = &calc.Sum{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		c._lastAllRequests = append(c._lastAllRequests, req)
		c._lastAllResponse.Z = c._lastAllResponse.Z + req.GetX() + req.GetY()
	}

	if err := stream.SendAndClose(c._lastAllResponse); err != nil {
		return err
	}

	return nil
}

func (c *CalcService)AddDeffered(req *calc.OperandList, stream calc.CalcService_AddDefferedServer) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c._lastDefferedResponses = []*calc.Sum{}

	for _, operand := range req.GetOperand() {
		sum := &calc.Sum{Z: operand.GetX() + operand.GetY()}
		if err := stream.Send(sum); err != nil {
			return err
		}
		c._lastDefferedResponses = append(c._lastDefferedResponses, sum)
	}

	c._lastDefferedRequest = req

	return nil
}

func (c *CalcService)AddAsync(stream calc.CalcService_AddAsyncServer) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c._lastAsyncRequests = []*calc.Operand{}
	c._lastAsyncResponses = []*calc.Sum{}

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

		c._lastAsyncRequests = append(c._lastAsyncRequests, operand)
		c._lastAsyncResponses = append(c._lastAsyncResponses, sum)
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

package SN

import (
	"fmt"
	"log"
	"sync"

	"github.com/nickng/scribble-go/runtime/session"
)

const A = "A"
const B = "B"

const LTimes = 1
const LEnd = 2

func check(e error) {
	if e != nil {
		log.Panic(e.Error())
	}
}

/*****************************************************************************/
/************ A API **********************************************************/
/*****************************************************************************/

type (
	A_Init struct {
		id  int
		ept *session.Endpoint
	}
	A_1 struct {
		id  int
		ept *session.Endpoint
	}
	A_2 struct {
		id  int
		ept *session.Endpoint
	}
	A_3 struct {
		id  int
		ept *session.Endpoint
	}
	A_4 struct {
		id  int
		ept *session.Endpoint
	}
	A_5 struct {
		id  int
		ept *session.Endpoint
	}
	A_6 struct {
		id  int
		ept *session.Endpoint
	}
	A_7 struct {
		id  int
		ept *session.Endpoint
	}
	A_8 struct {
		id  int
		ept *session.Endpoint
	}
	A_End struct {
	}
	mask struct {
		state uint16
		lock  sync.Mutex
	}
)

var a_Init []*A_Init
var a_1 []*A_1
var a_2 []*A_2
var a_3 []*A_3
var a_4 []*A_4
var a_5 []*A_5
var a_6 []*A_6
var a_7 []*A_7
var a_8 []*A_8
var a_End []*A_End

var stateA []*mask

func initializeA(id int) {
	for id >= len(stateA) {
		stateA = append(stateA, nil)
		a_Init = append(a_Init, nil)
		a_1 = append(a_1, nil)
		a_2 = append(a_2, nil)
		a_3 = append(a_3, nil)
		a_4 = append(a_4, nil)
		a_5 = append(a_5, nil)
		a_6 = append(a_6, nil)
		a_7 = append(a_7, nil)
		a_8 = append(a_8, nil)
		a_End = append(a_End, nil)
	}
	stateA[id] = &mask{0, sync.Mutex{}}
}

func toggleA(id int, n uint) {
	stateA[id].lock.Lock()
	defer stateA[id].lock.Unlock()
	stateA[id].state ^= 1 << n
}

func testA(id int, n uint) {
	stateA[id].lock.Lock()
	defer stateA[id].lock.Unlock()
	if stateA[id].state&(1<<n) == 0 {
		panic("Linear resource already used")
	}
	stateA[id].state = 0
}

func (ini *A_Init) test()         { testA(ini.id, 0) }
func (ini *A_Init) toggle(n uint) { toggleA(ini.id, n) }
func (ini *A_1) test()            { testA(ini.id, 1) }
func (ini *A_1) toggle(n uint)    { toggleA(ini.id, n) }
func (ini *A_2) test()            { testA(ini.id, 2) }
func (ini *A_2) toggle(n uint)    { toggleA(ini.id, 3) }
func (ini *A_3) test()            { testA(ini.id, 3) }
func (ini *A_3) toggle(n uint)    { toggleA(ini.id, 4) }
func (ini *A_4) test()            { testA(ini.id, 4) }
func (ini *A_4) toggle(n uint)    { toggleA(ini.id, 5) }
func (ini *A_5) test()            { testA(ini.id, 5) }
func (ini *A_5) toggle(n uint)    { toggleA(ini.id, 6) }
func (ini *A_6) test()            { testA(ini.id, 6) }
func (ini *A_6) toggle(n uint)    { toggleA(ini.id, 7) }
func (ini *A_7) test()            { testA(ini.id, 7) }
func (ini *A_7) toggle(n uint)    { toggleA(ini.id, 8) }
func (ini *A_8) test()            { testA(ini.id, 8) }
func (ini *A_8) toggle(n uint)    { toggleA(ini.id, 1) }

func NewA(id, numA, numB int) (*A_Init, error) {
	session.RoleRange(id, numA)
	conn, err := session.NewConn([]session.ParamRole{{B, numB}})
	if err != nil {
		return nil, err
	}

	initializeA(id)

	eptA := session.NewEndpoint(id, numA, conn)
	a_Init[id] = &A_Init{id, eptA}
	a_1[id] = &A_1{id, eptA}
	a_2[id] = &A_2{id, eptA}
	a_3[id] = &A_3{id, eptA}
	a_4[id] = &A_4{id, eptA}
	a_5[id] = &A_5{id, eptA}
	a_6[id] = &A_6{id, eptA}
	a_7[id] = &A_7{id, eptA}
	a_8[id] = &A_8{id, eptA}

	return a_Init[id], nil
}

func (ini *A_Init) Ept() *session.Endpoint {
	return ini.ept
}

func (ini *B_Init) Ept() *session.Endpoint {
	return ini.ept
}

func (ini *A_Init) Init() (*A_1, error) {
	ini.test()
	ini.ept.ConnMu.Lock()
	for n, _ := range ini.ept.Conn {
		for j, _ := range ini.ept.Conn[n] {
			for ini.ept.Conn[n][j] == nil {
			}
		}
	}
	ini.ept.ConnMu.Unlock()
	ini.toggle(1)
	return a_1[ini.id], nil
}

func (ini *A_1) SendTimes(pl []int) *A_2 {
	ini.test()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(LTimes))
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(2)
	return a_2[ini.id]
}

func (ini *A_1) SendEnd(pl []int) *A_End {
	ini.test()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(LEnd))
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(9)
	return a_End[ini.id]

}

func (ini *A_2) RecvDone() ([]int, *A_3) {
	ini.test()
	var tmp int
	pl := make([]int, len(ini.ept.Conn[B]))

	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		pl[i] = tmp
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(3)
	return pl, a_3[ini.id]

}

func (ini *A_3) SendNext(pl []int) *A_4 {
	ini.test()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(4)
	return a_4[ini.id]

}

func (ini *A_4) RecvDone() ([]int, *A_5) {
	ini.test()
	var tmp int
	pl := make([]int, len(ini.ept.Conn[B]))

	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		pl[i] = tmp
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(5)
	return pl, a_5[ini.id]

}

func (ini *A_5) SendTimesTr(pl []int) *A_6 {
	ini.test()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(6)
	return a_6[ini.id]

}

func (ini *A_6) RecvDone() ([]int, *A_7) {
	ini.test()
	var tmp int
	pl := make([]int, len(ini.ept.Conn[B]))

	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		pl[i] = tmp
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(7)
	return pl, a_7[ini.id]

}

func (ini *A_7) SendNext(pl []int) *A_8 {
	ini.test()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(pl[i]))
	}
	ini.toggle(8)
	return a_8[ini.id]

}

func (ini *A_8) RecvDone() ([]int, *A_1) {
	ini.test()
	var tmp int
	pl := make([]int, len(ini.ept.Conn[B]))

	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		pl[i] = tmp
	}
	ini.toggle(1)
	return pl, a_1[ini.id]
}

func (ini *A_Init) Run(f func(*A_1) *A_End) {
	ini.ept.CheckConnection()
	st1, err := ini.Init()
	check(err)
	f(st1)
}

/************ A API **********************************************************/

/*****************************************************************************/
/************ B API **********************************************************/
/*****************************************************************************/

type (
	B_Init struct {
		id  int
		ept *session.Endpoint
	}
	B_1 struct {
		id  int
		ept *session.Endpoint
	}
	B_2 struct {
		id  int
		ept *session.Endpoint
		Val int
	}
	B_3 struct {
		id  int
		ept *session.Endpoint
	}
	B_4 struct {
		id  int
		ept *session.Endpoint
		Val int
	}
	B_5 struct {
		id  int
		ept *session.Endpoint
	}
	B_6 struct {
		id  int
		ept *session.Endpoint
		Val int
	}
	B_7 struct {
		id  int
		ept *session.Endpoint
	}
	B_8 struct {
		id  int
		ept *session.Endpoint
		Val int
	}
	B_End struct {
		id  int
		Val int
	}
)

var b_Init []*B_Init
var b_1 []*B_1
var b_2 []*B_2
var b_3 []*B_3
var b_4 []*B_4
var b_5 []*B_5
var b_6 []*B_6
var b_7 []*B_7
var b_8 []*B_8
var b_End []*B_End

var stateB []*mask

func initializeB(id int) {
	for id >= len(stateB) {
		stateB = append(stateB, nil)
		b_Init = append(b_Init, nil)
		b_1 = append(b_1, nil)
		b_2 = append(b_2, nil)
		b_3 = append(b_3, nil)
		b_4 = append(b_4, nil)
		b_5 = append(b_5, nil)
		b_6 = append(b_6, nil)
		b_7 = append(b_7, nil)
		b_8 = append(b_8, nil)
		b_End = append(b_End, nil)
	}
	stateB[id] = &mask{0, sync.Mutex{}}
}

func toggleB(id int, n uint) {
	stateB[id].lock.Lock()
	defer stateB[id].lock.Unlock()
	stateB[id].state ^= 1 << n
}

func testB(id int, n uint) {
	stateB[id].lock.Lock()
	defer stateB[id].lock.Unlock()
	if stateB[id].state&(1<<n) == 0 {
		panic("Linear resource already used")
	}
	stateB[id].state = 0
}

func NewB(id, numB, numA int) (*B_Init, error) {
	session.RoleRange(id, numB)
	conn, err := session.NewConn([]session.ParamRole{{A, numA}})
	if err != nil {
		return nil, err
	}

	initializeB(id)

	eptB := session.NewEndpoint(id, numB, conn)

	b_Init[id] = &B_Init{id, eptB}
	b_1[id] = &B_1{id, eptB}
	b_2[id] = &B_2{id, eptB, 0}
	b_3[id] = &B_3{id, eptB}
	b_4[id] = &B_4{id, eptB, 0}
	b_5[id] = &B_5{id, eptB}
	b_6[id] = &B_6{id, eptB, 0}
	b_7[id] = &B_7{id, eptB}
	b_8[id] = &B_8{id, eptB, 0}
	b_End[id] = &B_End{id, 0}

	return b_Init[id], nil
}

func (ini *B_Init) test()         { testB(ini.id, 0) }
func (ini *B_Init) toggle(n uint) { toggleB(ini.id, n) }
func (ini *B_1) test()            { testB(ini.id, 1) }
func (ini *B_1) toggle(n uint)    { toggleB(ini.id, n) }
func (ini *B_2) test()            { testB(ini.id, 2) }
func (ini *B_2) toggle(n uint)    { toggleB(ini.id, 3) }
func (ini *B_3) test()            { testB(ini.id, 3) }
func (ini *B_3) toggle(n uint)    { toggleB(ini.id, 4) }
func (ini *B_4) test()            { testB(ini.id, 4) }
func (ini *B_4) toggle(n uint)    { toggleB(ini.id, 5) }
func (ini *B_5) test()            { testB(ini.id, 5) }
func (ini *B_5) toggle(n uint)    { toggleB(ini.id, 6) }
func (ini *B_6) test()            { testB(ini.id, 6) }
func (ini *B_6) toggle(n uint)    { toggleB(ini.id, 7) }
func (ini *B_7) test()            { testB(ini.id, 7) }
func (ini *B_7) toggle(n uint)    { toggleB(ini.id, 8) }
func (ini *B_8) test()            { testB(ini.id, 8) }
func (ini *B_8) toggle(n uint)    { toggleB(ini.id, 1) }

func (ini *B_Init) Init() (*B_1, error) {
	ini.test()
	for n, l := range ini.ept.Conn {
		for i, c := range l {
			if c == nil {
				return nil, fmt.Errorf("Invalid connection for worker %s[%d]", n, i)
			}
		}
	}
	return b_1[ini.id], nil
}

func (st1 *B_1) TimesOrEnd() interface{} {
	st1.test()
	var lbl int
	var res int

	conn := st1.ept.Conn[A][0]
	err := conn.Recv(&lbl)
	if err != nil {
		log.Fatalf("wrong label from server at %d: %s", st1.ept.Id, err)
	}

	if lbl == LTimes {
		err = conn.Recv(&res)
		if err != nil {
			log.Fatalf("wrong value(times) from server at %d: %s", st1.ept.Id, err)
		}
		b_2[st1.id].Val = res
		st1.toggle(2)
		return b_2[st1.id]
	}
	if lbl == LEnd {
		err = conn.Recv(&res)
		if err != nil {
			log.Fatalf("wrong value(end) from server at %d: %s", st1.ept.Id, err)
		}
		b_End[st1.id].Val = res
		st1.toggle(9)
		return b_End[st1.id]
	}

	log.Fatalf("wrong label from server at %d: %s", st1.ept.Id, err)
	return nil
}

func (ini *B_2) SendDone(pl int) *B_3 {
	ini.test()

	check(ini.ept.Conn[A][0].Send(pl))
	ini.toggle(3)
	return b_3[ini.id]
}

func (ini *B_3) RecvNext() *B_4 {
	ini.test()

	check(ini.ept.Conn[A][0].Recv(&b_4[ini.id].Val))
	ini.toggle(4)
	return b_4[ini.id]

}

func (ini *B_4) SendDone(pl int) *B_5 {
	ini.test()

	check(ini.ept.Conn[A][0].Send(pl))
	ini.toggle(5)
	return b_5[ini.id]

}

func (ini *B_5) RecvTimesTr() *B_6 {
	ini.test()

	check(ini.ept.Conn[A][0].Recv(&b_6[ini.id].Val))
	ini.toggle(6)
	return b_6[ini.id]

}

func (ini *B_6) SendDone(pl int) *B_7 {
	ini.test()

	check(ini.ept.Conn[A][0].Send(pl))
	ini.toggle(7)
	return b_7[ini.id]

}

func (ini *B_7) RecvNext() *B_8 {
	ini.test()

	check(ini.ept.Conn[A][0].Recv(&b_8[ini.id].Val))
	ini.toggle(8)
	return b_8[ini.id]

}

func (ini *B_8) SendDone(pl int) *B_1 {
	ini.test()

	check(ini.ept.Conn[A][0].Send(pl))
	ini.toggle(1)
	return b_1[ini.id]

}

func (ini *B_Init) Run(f func(*B_1) *B_End) {
	ini.ept.CheckConnection()
	st1, err := ini.Init()
	check(err)
	f(st1)
}

/************ B API **********************************************************/

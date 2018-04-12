package Regex

import (
	"fmt"
	"log"
	"sync"

	"github.com/nickng/scribble-go/runtime/session"
)

const A = "A"
const B = "B"
const C = "C"

func check(e error) {
	if e != nil {
		log.Panic(e.Error())
	}
}

type mask struct {
	state uint16
	lock  sync.Mutex
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
		Val []int
	}
	A_End struct {
		Val int
	}
)

var a_Init []*A_Init
var a_1 []*A_1
var a_2 []*A_2
var a_3 []*A_3
var a_4 []*A_4
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
		a_End = append(a_End, nil)
	}
	stateA[id] = &mask{1, sync.Mutex{}}
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

func (ini *A_Init) test()   { testA(ini.id, 0) }
func (ini *A_Init) toggle() { toggleA(ini.id, 0) }
func (ini *A_1) test()      { testA(ini.id, 1) }
func (ini *A_1) toggle()    { toggleA(ini.id, 1) }
func (ini *A_2) test()      { testA(ini.id, 2) }
func (ini *A_2) toggle()    { toggleA(ini.id, 2) }
func (ini *A_3) test()      { testA(ini.id, 3) }
func (ini *A_3) toggle()    { toggleA(ini.id, 3) }
func (ini *A_4) test()      { testA(ini.id, 4) }
func (ini *A_4) toggle()    { toggleA(ini.id, 4) }

func (ini *A_Init) Ept() *session.Endpoint {
	return ini.ept
}

func NewA(id, numA, numB, numC int) (*A_Init, error) {
	session.RoleRange(id, numA)
	conn, err := session.NewConn([]session.ParamRole{{B, numB}, {C, numC}})
	if err != nil {
		return nil, err
	}

	initializeA(id)

	eptA := session.NewEndpoint(id, numA, conn)
	a_Init[id] = &A_Init{id, eptA}
	a_1[id] = &A_1{id, eptA}
	a_2[id] = &A_2{id, eptA}
	a_3[id] = &A_3{id, eptA}
	a_4[id] = &A_4{id, eptA, nil}
	a_End[id] = &A_End{0}

	return a_Init[id], nil
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
	a_1[ini.id].toggle()
	return a_1[ini.id], nil
}

func (ini *A_1) Count(pl []string) *A_2 {
	ini.test()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' Count")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	a_2[ini.id].toggle()
	return a_2[ini.id]
}

func (ini *A_2) Measure(pl int) *A_3 {
	ini.test()
	if 1 != len(ini.ept.Conn[C]) {
		log.Panicf("Incorrect number of arguments to role 'C' Measure")
	}
	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[C][0].Send(pl))
	ini.ept.ConnMu.RUnlock()
	a_3[ini.id].toggle()
	return a_3[ini.id]
}

func (ini *A_3) Donec() *A_4 {
	ini.test()
	var tmp int
	pl := &a_4[ini.id].Val
	*pl = make([]int, len(ini.ept.Conn[B]))

	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		(*pl)[i] = tmp
	}
	ini.ept.ConnMu.RUnlock()
	a_4[ini.id].toggle()
	return a_4[ini.id]
}

func (ini *A_4) Len() *A_End {
	ini.test()
	var tmp int
	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[C][0].Recv(&tmp))
	ini.ept.ConnMu.RUnlock()
	a_End[ini.id].Val = tmp
	return a_End[ini.id]
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
		Val string
	}
	B_End struct {
	}
)

var b_Init []*B_Init
var b_1 []*B_1
var b_2 []*B_2
var b_End []*B_End

var stateB []*mask

func initializeB(id int) {
	for id >= len(stateB) {
		stateB = append(stateB, nil)
		b_Init = append(b_Init, nil)
		b_1 = append(b_1, nil)
		b_2 = append(b_2, nil)
		b_End = append(b_End, nil)
	}
	stateB[id] = &mask{1, sync.Mutex{}}
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

func (ini *B_Init) test()   { testB(ini.id, 0) }
func (ini *B_Init) toggle() { toggleB(ini.id, 0) }
func (ini *B_1) test()      { testB(ini.id, 1) }
func (ini *B_1) toggle()    { toggleB(ini.id, 1) }
func (ini *B_2) test()      { testB(ini.id, 2) }
func (ini *B_2) toggle()    { toggleB(ini.id, 2) }
func (ini *B_Init) Ept() *session.Endpoint {
	return ini.ept
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
	b_2[id] = &B_2{id, eptB, ""}
	b_End[id] = &B_End{}

	return b_Init[id], nil
}

func (ini *B_Init) Init() (*B_1, error) {
	ini.test()
	ini.ept.ConnMu.Lock()
	for n, l := range ini.ept.Conn {
		for i, c := range l {
			if c == nil {
				return nil, fmt.Errorf("nvalid connection for worker %s[%d] at B[%d]", n, i, ini.Ept().Id)
			}
		}
	}
	ini.ept.ConnMu.Unlock()
	b_1[ini.id].toggle()
	return b_1[ini.id], nil
}

func (ini *B_1) Count() *B_2 {
	ini.test()
	var tmp string

	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[A][0].Recv(&tmp))
	ini.ept.ConnMu.RUnlock()
	b_2[ini.id].Val = tmp
	b_2[ini.id].toggle()
	return b_2[ini.id]
}

func (ini *B_2) Donec(pl int) *B_End {
	ini.test()

	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[A][0].Send(pl))
	ini.ept.ConnMu.RUnlock()
	return b_End[ini.id]
}

func (ini *B_Init) Run(f func(*B_1) *B_End) {
	ini.ept.CheckConnection()
	st1, err := ini.Init()
	check(err)
	f(st1)
}

/************ B API **********************************************************/

/*****************************************************************************/
/************ C API **********************************************************/
/*****************************************************************************/

type (
	C_Init struct {
		id  int
		ept *session.Endpoint
	}
	C_1 struct {
		id  int
		ept *session.Endpoint
	}
	C_2 struct {
		id  int
		ept *session.Endpoint
		Val int
	}
	C_End struct {
	}
)

var c_Init []*C_Init
var c_1 []*C_1
var c_2 []*C_2
var c_End []*C_End

var stateC []*mask

func initializeC(id int) {
	for id >= len(stateC) {
		stateC = append(stateC, nil)
		c_Init = append(c_Init, nil)
		c_1 = append(c_1, nil)
		c_2 = append(c_2, nil)
		c_End = append(c_End, nil)
	}
	stateC[id] = &mask{1, sync.Mutex{}}
}

func toggleC(id int, n uint) {
	stateC[id].lock.Lock()
	defer stateC[id].lock.Unlock()
	stateC[id].state = 1 << n
}

func testC(id int, n uint) {
	stateC[id].lock.Lock()
	defer stateC[id].lock.Unlock()
	if stateC[id].state&(1<<n) == 0 {
		panic("Linear resource already used")
	}
	stateC[id].state = 0
}

func (ini *C_Init) test()   { testC(ini.id, 0) }
func (ini *C_Init) toggle() { toggleC(ini.id, 0) }
func (ini *C_1) test()      { testC(ini.id, 1) }
func (ini *C_1) toggle()    { toggleC(ini.id, 1) }
func (ini *C_2) test()      { testC(ini.id, 2) }
func (ini *C_2) toggle()    { toggleC(ini.id, 2) }
func (ini *C_Init) Ept() *session.Endpoint {
	return ini.ept
}

func NewC(id, numS, numA int) (*C_Init, error) {
	session.RoleRange(id, numS)
	conn, err := session.NewConn([]session.ParamRole{{A, numA}})
	if err != nil {
		return nil, err
	}

	initializeC(id)

	eptC := session.NewEndpoint(id, numS, conn)
	c_Init[id] = &C_Init{id, eptC}
	c_1[id] = &C_1{id, eptC}
	c_2[id] = &C_2{id, eptC, 0}
	c_End[id] = &C_End{}

	return c_Init[id], nil
}

func (ini *C_Init) Init() (*C_1, error) {
	ini.test()
	ini.ept.ConnMu.Lock()
	for n, l := range ini.ept.Conn {
		for i, l := range l {
			if l == nil {
				return nil, fmt.Errorf("Invalid connection for worker %s[%d]", n, i)
			}
		}
	}
	ini.ept.ConnMu.Unlock()
	c_1[ini.id].toggle()
	return c_1[ini.id], nil
}

func (ini *C_1) Measure() *C_2 {
	ini.test()
	var tmp int

	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[A][0].Recv(&tmp))
	ini.ept.ConnMu.RUnlock()
	c_2[ini.id].Val = tmp
	c_2[ini.id].toggle()
	return c_2[ini.id]
}

func (ini *C_2) Len(s int) *C_End {
	ini.test()
	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[A][0].Send(s))
	ini.ept.ConnMu.RUnlock()
	return c_End[ini.id]
}

func (ini *C_Init) Run(f func(*C_1) *C_End) {
	st1, err := ini.Init()
	check(err)
	f(st1)
}

/************ C API **********************************************************/

package KNuc

import (
	"fmt"
	"log"
	"sync"

	"github.com/nickng/scribble-go/runtime/session"
)

const A = "A"
const B = "B"
const S = "S"

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
		Val []string
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

func (ini *A_Init) Ept() *session.Endpoint {
	return ini.ept
}

func NewA(id, numA, numS, numB int) (*A_Init, error) {
	session.RoleRange(id, numA)
	conn, err := session.NewConn([]session.ParamRole{{S, numS}, {B, numB}})
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
	a_End[id] = &A_End{nil}

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
	ini.toggle(1)
	return a_1[ini.id], nil
}

func (ini *A_1) SendS(pl []int) *A_2 {
	ini.test()
	if len(pl) != len(ini.ept.Conn[S]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[S] {
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(2)
	return a_2[ini.id]
}

func (ini *A_2) SendB(pl []string) *A_3 {
	ini.test()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'B' SendB")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(3)
	return a_3[ini.id]
}

func (ini *A_3) RecvS() *A_4 {
	ini.test()
	var tmp int
	a_4[ini.id].Val = make([]int, len(ini.ept.Conn[S]))

	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[S] {
		check(c.Recv(&tmp))
		a_4[ini.id].Val[i] = tmp
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(4)
	return a_4[ini.id]
}

func (ini *A_4) RecvB() *A_End {
	var tmp string
	a_End[ini.id].Val = make([]string, len(ini.ept.Conn[B]))

	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		a_End[ini.id].Val[i] = tmp
	}
	ini.ept.ConnMu.RUnlock()
	ini.toggle(5)
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

func (ini *B_Init) test()         { testB(ini.id, 0) }
func (ini *B_Init) toggle(n uint) { toggleB(ini.id, n) }
func (ini *B_1) test()            { testB(ini.id, 1) }
func (ini *B_1) toggle(n uint)    { toggleB(ini.id, n) }
func (ini *B_2) test()            { testB(ini.id, 2) }
func (ini *B_2) toggle(n uint)    { toggleB(ini.id, 3) }

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
				return nil, fmt.Errorf("Invalid connection for worker %s[%d]", n, i)
			}
		}
	}
	ini.ept.ConnMu.Unlock()
	ini.toggle(1)
	return b_1[ini.id], nil
}

func (ini *B_1) Recv_BA() *B_2 {
	ini.test()

	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[A][0].Recv(&b_2[ini.id].Val))
	ini.ept.ConnMu.RUnlock()
	ini.toggle(2)
	return b_2[ini.id]
}

func (ini *B_2) Send_BA(pl string) *B_End {
	ini.test()

	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[A][0].Send(pl))
	ini.ept.ConnMu.RUnlock()
	ini.toggle(3)
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
/************ S API **********************************************************/
/*****************************************************************************/

type (
	S_Init struct {
		id  int
		ept *session.Endpoint
	}
	S_1 struct {
		id  int
		ept *session.Endpoint
	}
	S_2 struct {
		id  int
		ept *session.Endpoint
		Val int
	}
	S_End struct {
	}
)

var s_Init []*S_Init
var s_1 []*S_1
var s_2 []*S_2
var s_End []*S_End

var stateS []*mask

func initializeS(id int) {
	for id >= len(stateS) {
		stateS = append(stateS, nil)
		s_Init = append(s_Init, nil)
		s_1 = append(s_1, nil)
		s_2 = append(s_2, nil)
		s_End = append(s_End, nil)
	}
	stateS[id] = &mask{1, sync.Mutex{}}
}

func toggleS(id int, n uint) {
	stateS[id].lock.Lock()
	defer stateS[id].lock.Unlock()
	stateS[id].state ^= 1 << n
}

func testS(id int, n uint) {
	stateS[id].lock.Lock()
	defer stateS[id].lock.Unlock()
	if stateS[id].state&(1<<n) == 0 {
		panic("Linear resource already used")
	}
	stateS[id].state = 0
}

func (ini *S_Init) test()         { testS(ini.id, 0) }
func (ini *S_Init) toggle(n uint) { toggleS(ini.id, n) }
func (ini *S_1) test()            { testS(ini.id, 1) }
func (ini *S_1) toggle(n uint)    { toggleS(ini.id, n) }
func (ini *S_2) test()            { testS(ini.id, 2) }
func (ini *S_2) toggle(n uint)    { toggleS(ini.id, 3) }

func (ini *S_Init) Ept() *session.Endpoint {
	return ini.ept
}

func NewS(id, numS, numA int) (*S_Init, error) {
	session.RoleRange(id, numS)
	conn, err := session.NewConn([]session.ParamRole{{A, numA}})
	if err != nil {
		return nil, err
	}

	initializeS(id)

	eptS := session.NewEndpoint(id, numS, conn)
	s_Init[id] = &S_Init{id, eptS}
	s_1[id] = &S_1{id, eptS}
	s_2[id] = &S_2{id, eptS, 0}
	s_End[id] = &S_End{}

	return s_Init[id], nil
}

func (ini *S_Init) Init() (*S_1, error) {
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
	ini.toggle(1)
	return s_1[ini.id], nil
}

func (ini *S_1) Recv_SA() *S_2 {
	ini.test()

	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[A][0].Recv(&s_2[ini.id].Val))
	ini.ept.ConnMu.RUnlock()
	ini.toggle(2)
	return s_2[ini.id]
}

func (ini *S_2) Send_SA(s int) *S_End {
	ini.test()
	ini.ept.ConnMu.RLock()
	check(ini.ept.Conn[A][0].Send(s))
	ini.ept.ConnMu.RUnlock()
	ini.toggle(3)
	return s_End[ini.id]
}

func (ini *S_Init) Run(f func(*S_1) *S_End) {
	st1, err := ini.Init()
	check(err)
	f(st1)
}

/************ S API **********************************************************/

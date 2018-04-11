package SN

import (
	"fmt"
	"github.com/nickng/scribble-go/runtime/session"
	"log"
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
	A_Init struct{ ept *session.Endpoint }
	A_1    struct{ ept *session.Endpoint }
	A_2    struct{ ept *session.Endpoint }
	A_3    struct{ ept *session.Endpoint }
	A_4    struct{ ept *session.Endpoint }
	A_5    struct{ ept *session.Endpoint }
	A_6    struct{ ept *session.Endpoint }
	A_7    struct{ ept *session.Endpoint }
	A_8    struct{ ept *session.Endpoint }
	A_End  struct{}
	mask   = uint16
)

var eptA *session.Endpoint

var a_Init *A_Init
var a_1 *A_1
var a_2 *A_2
var a_3 *A_3
var a_4 *A_4
var a_5 *A_5
var a_6 *A_6
var a_7 *A_7
var a_8 *A_8
var a_End = A_End(struct{}{})

var stateA mask

func ma_Init() {
	stateA = 1
}

func toggleA(n uint) {
	stateA ^= 1 << n
}

func testA(n uint) {
	if stateA&(1<<n) == 0 {
		panic("Linear resource already used")
	}
	stateA = 0
}

func useIni() {
	testA(0)
	toggleA(1)
}

func use1_2() {
	testA(1)
	toggleA(2)
}

func use1_E() {
	testA(1)
}

func NewA(id, numA, numB int) (*A_Init, error) {
	session.RoleRange(id, numA)
	conn, err := session.NewConn([]session.ParamRole{{B, numB}})
	if err != nil {
		return nil, err
	}

	eptA = session.NewEndpoint(id, numA, conn)
	a_Init = &A_Init{eptA}
	a_1 = &A_1{eptA}
	a_2 = &A_2{eptA}
	a_3 = &A_3{eptA}
	a_4 = &A_4{eptA}
	a_5 = &A_5{eptA}
	a_6 = &A_6{eptA}
	a_7 = &A_7{eptA}
	a_8 = &A_8{eptA}

	ma_Init()

	return a_Init, nil
}

func (ini *A_Init) Ept() *session.Endpoint {
	return ini.ept
}

func (ini *B_Init) Ept() *session.Endpoint {
	return ini.ept
}

func (ini *A_Init) Init() (*A_1, error) {
	useIni()
	ini.ept.ConnMu.Lock()
	for n, _ := range ini.ept.Conn {
		for j, _ := range ini.ept.Conn[n] {
			for ini.ept.Conn[n][j] == nil {
			}
		}
	}
	ini.ept.ConnMu.Unlock()
	return a_1, nil
}

func (ini *A_1) SendTimes(pl []int) *A_2 {
	use1_2()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(LTimes))
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	return a_2
}

func (ini *A_1) SendEnd(pl []int) *A_End {
	use1_E()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(LEnd))
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	return &a_End
}

func use2() {
	testA(2)
	toggleA(3)
}

func (ini *A_2) RecvDone() ([]int, *A_3) {
	use2()
	var tmp int
	pl := make([]int, len(ini.ept.Conn[B]))

	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		pl[i] = tmp
	}
	ini.ept.ConnMu.RUnlock()
	return pl, a_3
}

func use3() {
	testA(3)
	toggleA(4)
}

func (ini *A_3) SendNext(pl []int) *A_4 {
	use3()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	return a_4
}

func use4() {
	testA(4)
	toggleA(5)
}

func (ini *A_4) RecvDone() ([]int, *A_5) {
	use4()
	var tmp int
	pl := make([]int, len(ini.ept.Conn[B]))

	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		pl[i] = tmp
	}
	ini.ept.ConnMu.RUnlock()
	return pl, a_5
}

func use5() {
	testA(5)
	toggleA(6)
}

func (ini *A_5) SendTimesTr(pl []int) *A_6 {
	use5()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(pl[i]))
	}
	ini.ept.ConnMu.RUnlock()
	return a_6
}

func use6() {
	testA(6)
	toggleA(7)
}

func (ini *A_6) RecvDone() ([]int, *A_7) {
	use6()
	var tmp int
	pl := make([]int, len(ini.ept.Conn[B]))

	ini.ept.ConnMu.RLock()
	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		pl[i] = tmp
	}
	ini.ept.ConnMu.RUnlock()
	return pl, a_7
}

func use7() {
	testA(7)
	toggleA(8)
}

func (ini *A_7) SendNext(pl []int) *A_8 {
	use7()
	if len(pl) != len(ini.ept.Conn[B]) {
		log.Panicf("Incorrect number of arguments to role 'A' SendS")
	}
	for i, c := range ini.ept.Conn[B] {
		check(c.Send(pl[i]))
	}
	return a_8
}

func use8() {
	testA(8)
	toggleA(1)
}

func (ini *A_8) RecvDone() ([]int, *A_1) {
	use8()
	var tmp int
	pl := make([]int, len(ini.ept.Conn[B]))

	for i, c := range ini.ept.Conn[B] {
		check(c.Recv(&tmp))
		pl[i] = tmp
	}
	return pl, a_1
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

var eptB []*session.Endpoint
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

var stateB []mask

func mb_Init(r int) {
	stateB[r] = 1
}

func toggleB(id int, n uint) {
	stateB[id] ^= 1 << n
	return
}

func testB(id int, n uint) {
	if (stateB[id] & (1 << n)) == 0 {
		panic("Linear resource already used")
	}
	stateB[id] = 0
	return
}

func NewB(id, numB, numA int) (*B_Init, error) {
	session.RoleRange(id, numB)
	conn, err := session.NewConn([]session.ParamRole{{A, numA}})
	if err != nil {
		return nil, err
	}

	for id >= len(stateB) {
		stateB = append(stateB, 0)
		eptB = append(eptB, nil)
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

	stateB[id] = 0

	eptB[id] = session.NewEndpoint(id, numB, conn)

	b_Init[id] = &B_Init{id, eptB[id]}
	b_1[id] = &B_1{id, eptB[id]}
	b_2[id] = &B_2{id, eptB[id], 0}
	b_3[id] = &B_3{id, eptB[id]}
	b_4[id] = &B_4{id, eptB[id], 0}
	b_5[id] = &B_5{id, eptB[id]}
	b_6[id] = &B_6{id, eptB[id], 0}
	b_7[id] = &B_7{id, eptB[id]}
	b_8[id] = &B_8{id, eptB[id], 0}
	b_End[id] = &B_End{id, 0}

	mb_Init(id)

	return b_Init[id], nil
}

func (ini *B_Init) use() {
	testB(ini.id, 0)
	toggleB(ini.id, 1)
}

func (ini *B_Init) Init() (*B_1, error) {
	ini.use()
	for n, l := range ini.ept.Conn {
		for i, c := range l {
			if c == nil {
				return nil, fmt.Errorf("Invalid connection for worker %s[%d]", n, i)
			}
		}
	}
	return b_1[ini.id], nil
}

func (ini *B_1) use() {
	testB(ini.id, 1)
}

func (ini *B_1) go2() {
	toggleB(ini.id, 2)
}

func (ini *B_1) goEnd() {
	toggleB(ini.id, 9)
}

func (st1 *B_1) TimesOrEnd() interface{} {
	st1.use()
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
		st1.go2()
		return b_2[st1.id]
	}
	if lbl == LEnd {
		err = conn.Recv(&res)
		if err != nil {
			log.Fatalf("wrong value(end) from server at %d: %s", st1.ept.Id, err)
		}
		b_End[st1.id].Val = res
		st1.goEnd()
		return b_End[st1.id]
	}

	log.Fatalf("wrong label from server at %d: %s", st1.ept.Id, err)
	return nil
}

func (ini *B_2) use() {
	testB(ini.id, 2)
	toggleB(ini.id, 3)
}

func (ini *B_3) use() {
	testB(ini.id, 3)
	toggleB(ini.id, 4)
}

func (ini *B_4) use() {
	testB(ini.id, 4)
	toggleB(ini.id, 5)
}

func (ini *B_5) use() {
	testB(ini.id, 5)
	toggleB(ini.id, 6)
}

func (ini *B_6) use() {
	testB(ini.id, 6)
	toggleB(ini.id, 7)
}

func (ini *B_7) use() {
	testB(ini.id, 7)
	toggleB(ini.id, 8)
}

func (ini *B_8) use() {
	testB(ini.id, 8)
	toggleB(ini.id, 1)
}

func (ini *B_2) SendDone(pl int) *B_3 {
	ini.use()

	check(ini.ept.Conn[A][0].Send(pl))
	return b_3[ini.id]
}

func (ini *B_3) RecvNext() *B_4 {
	ini.use()

	check(ini.ept.Conn[A][0].Recv(&b_4[ini.id].Val))
	return b_4[ini.id]

}

func (ini *B_4) SendDone(pl int) *B_5 {
	ini.use()

	check(ini.ept.Conn[A][0].Send(pl))
	return b_5[ini.id]

}

func (ini *B_5) RecvTimesTr() *B_6 {
	ini.use()

	check(ini.ept.Conn[A][0].Recv(&b_6[ini.id].Val))
	return b_6[ini.id]

}

func (ini *B_6) SendDone(pl int) *B_7 {
	ini.use()

	check(ini.ept.Conn[A][0].Send(pl))
	return b_7[ini.id]

}

func (ini *B_7) RecvNext() *B_8 {
	ini.use()

	check(ini.ept.Conn[A][0].Recv(&b_8[ini.id].Val))
	return b_8[ini.id]

}

func (ini *B_8) SendDone(pl int) *B_1 {
	ini.use()

	check(ini.ept.Conn[A][0].Send(pl))
	return b_1[ini.id]

}

func (ini *B_Init) Run(f func(*B_1) *B_End) {
	ini.ept.CheckConnection()
	st1, err := ini.Init()
	check(err)
	f(st1)
}

/************ B API **********************************************************/

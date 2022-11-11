package main

import (
	"errors"
	"log"
	"sync"
)

// 连接池简单实现

var ErrPoolClosed error = errors.New("Conn Pool Closed")

type Pool struct {
	m       sync.Mutex
	res     chan *MyConn            // 存储连接的channel
	factory func() (*MyConn, error) // 创建连接工厂
	closed  bool                    // 连接池关闭标志
}

func NewPool(fn func() (*MyConn, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size to small")
	}
	return &Pool{
		factory: fn,
		res:     make(chan *MyConn, size),
	}, nil
}

// 从连接池中获取一个连接
func (p *Pool) Acquire() (*MyConn, error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		log.Println("产生新连接")
		return p.factory()
	}
}

// 释放连接
// 释放连接首先得有个前提，就是连接池还没有关闭。
// 如果连接池已经关闭再往res里面送连接的话就好触发panic。
func (p *Pool) Release(r *MyConn) {
	p.m.Lock() // 保证closed字段的线程安全
	defer p.m.Unlock()

	// 连接池关闭了
	if p.closed {
		r.Conn.Close()
		return
	}

	select {
	case p.res <- r:
		log.Println("放入连接池")
	default:
		log.Println("连接池已满")
		r.Conn.Close()
	}
}

func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}
	p.closed = true

	close(p.res)

	// 关闭channel中的连接
	for r := range p.res {
		r.Conn.Close()
	}
}

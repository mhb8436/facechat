package connect

import "sync"

type Bucket struct {
	cLock         sync.RWMutex
	chs           map[string]*Channel
	bucketOptions BucketOptions
	broadcast     chan []byte
}

type BucketOptions struct {
	ChannelSize int
}

func NewBucket(bucketOptions BucketOptions) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, bucketOptions.ChannelSize)
	b.bucketOptions = bucketOptions
	return
}

func (b *Bucket) Put(userId string, ch *Channel) (err error) {
	b.cLock.Lock()
	ch.userId = userId
	b.chs[userId] = ch
	b.cLock.Unlock()
	return
}

func (b *Bucket) DeleteChannel(ch *Channel) {
	var ok bool
	b.cLock.RLock()
	if ch, ok = b.chs[ch.userId]; ok {
		delete(b.chs, ch.userId)
	}
	b.cLock.RUnlock()
}

func (b *Bucket) Channel(userId string) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[userId]
	b.cLock.RUnlock()
	return
}

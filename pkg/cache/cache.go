package cache

import (
	"bytes"
	"encoding/json"
	"time"

	"bingo/pkg/log"
	"bingo/pkg/utils"
)

type (
	// Cache struct
	Cache struct {
		// id   int64
		file string
		size int
		d    time.Duration
		key  bool

		head interface{}
		// hbytes     []byte
		headWriter func()
		// change  bool
		// values utils.Map
		// index  utils.Map
		Map
		c *cache
	}

	InstFn func() interface{}
	loadFn func(v interface{}) (interface{}, bool)
)

var (
	tag = []byte("<--CACHE-HEAD_SPLITER-->")
)

// New returns a new cache with a given default expiration duration and cleanup
// interval. If the expiration duration is less than one (or NoExpiration),
// the items in the cache never expire (by default), and must be deleted
// manually. If the cleanup interval is less than one, expired items are not
// deleted from the cache before calling c.DeleteExpired().
func New(fileName string, saveInterval time.Duration) *Cache {
	// func New(fileName string, saveInterval, defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	return &Cache{
		// id:    utils.GetID(),
		file: fileName,
		d:    saveInterval,
		Map:  Map{change: false, Map: utils.NewMap()},
		c:    newCacheWithJanitor(DefaultExpiration, 0, items),
		// c:    newCacheWithJanitor(defaultExpiration, cleanupInterval, items),
	}
}

// Init used to load list cache from file
// Be aware of that 'instancePtr' should return a instance pointor of struct
// And put in the pointor to 'loader' while it returns a key value
// func (c *Cache) Init(instancePtr func() interface{}, loader func(v interface{}) (interface{}, bool)) {
func (c *Cache) Init(instancePtr InstFn, loader loadFn) {
	c.InitWithHead(nil, nil, instancePtr, loader)
}

// InitWithHead used to load list cache from file
// Be aware of that 'instancePtr' should return a instance pointor of struct
// And put in the pointor to 'loader' while it returns a key value
func (c *Cache) InitWithHead(header interface{}, writer func(), instancePtr InstFn, loader loadFn) {
	// func (c *Cache) InitWithHead(header InstFn, headLoader func(), instancePtr InstFn, loader loadFn) {
	defer utils.Recover()

	if header != nil {
		c.head = header
	}
	if writer != nil {
		c.headWriter = writer
	}
	// log.Warning("1 c.head %#v", c.head)

	if d, err := utils.ReadFile(c.file); err == nil {
		lines := bytes.Split(d, []byte("\n"))
		if c.head == nil {
			c.init(lines, instancePtr, loader)
		} else {
			for i, l := range lines {
				if bytes.Equal(l, tag) {
					if i > 0 {
						hbytes := bytes.Join(lines[:i], []byte{})
						// log.Warning("set head %p: %s", c.head, c.hbytes)
						if e := json.Unmarshal(hbytes, &c.head); e != nil {
							log.Error("load head from cache file error %#v in %s line:%d %s", header, c.file, i+1, e)
							// } else {
							// 	headLoader()
						}
						// log.PrintJSON(c.head)
						// log.Info("loaded head: %+v", c.head)
					}
					if len(lines) > i {
						c.init(lines[i+1:], instancePtr, loader)
					}
					return
				}
			}
		}
	}
}
func (c *Cache) init(lines [][]byte, instancePtr InstFn, loader loadFn) {
	start := time.Now()
	// var lines [][]byte

	list := []interface{}{}
	oldFmt := json.Unmarshal(bytes.Join(lines, []byte("\n")), &list) == nil
	if oldFmt {
		for _, l := range list {
			if l != nil {
				d, _ := json.Marshal(&l)
				lines = append(lines, d)
			}
		}
		// } else {
		// new format
		// lines = bytes.Split(d, []byte("\n"))
	}

	// bads := []int{}
	for i, d := range lines {
		if len(bytes.TrimSpace(d)) == 0 {
			continue
		}

		inst := instancePtr()
		// log.Notice("raw:%s type:%#v", d, inst)
		if e := json.Unmarshal(d, &inst); e != nil {
			// bads = append(bads, i)
			log.Warning("parse cache file error in %s line:%d(%s) %s", c.file, i+1, d, e)
		} else if k, ok := loader(inst); ok {
			// log.Notice("key:%#v, value:%#v", k, inst)
			c.Store(k, inst)
		} else {
			log.Error("load from cache file error %#v in %s line:%d(%s) %s", inst, c.file, i+1, d, e)
		}
	}
	if oldFmt {
		c.Flush()
	}
	c.Map.reset()
	c.Map.change = false
	log.Info("finished to load cache. file:%s total:%d cost:%s", c.file, c.Len(), utils.Since(start))

	if c.d > 0 {
		utils.Looper(c.d, true, false, nil, c.Flush)
	}
}

// func (c *Cache) LoadHead(h interface{}) error { return json.Unmarshal(c.hbytes, &h) }
// func (c *Cache) SetHeadWriter(fn func()) { c.headWriter = fn }

// Flush saves new cache item to file
func (c *Cache) Flush() {
	if !c.Map.change {
		return
	}

	utils.WriteNewFile(c.file, []byte{})
	// log.Warning("c.head %#v", c.head)
	if c.head != nil {
		if c.headWriter != nil {
			c.headWriter()
		}
		if d, e := json.Marshal(c.head); len(d) > 0 && e == nil {
			// log.Warning("c.head %s", d)
			utils.AppendBytesLine(c.file, d)
			utils.AppendBytesLine(c.file, tag)
		}
	}
	start := time.Now()
	var tc, ec int
	c.Range(func(k, v interface{}) bool {
		if v != nil {
			if d, e := json.Marshal(v); e == nil {
				utils.AppendBytesLine(c.file, d)
				tc++
			} else {
				log.Error("save error: %s. k:%+v, v:%+v", e, k, v)
				ec++
			}
		}
		return true
	})
	log.Info("finished to save cache. file:%s total:%d/%d cost:%s", c.file, tc, ec, utils.Since(start))
}

// func (c *Cache) SetHeader(header InstFn) { c.head = header() }

// Index used to
func (c *Cache) Index(keyName, keyValue interface{}) (v interface{}, ok bool) {
	return
}

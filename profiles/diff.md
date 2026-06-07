# Diff


## Diff profile

```bash
go tool pprof -top -diff_base=./profiles/base.pprof ./profiles/result.pprof

File: shortener
Build ID: 4baed019e76ecd5629a54e83bd1551c6a4d6c97c
Type: cpu
Time: 2026-06-07 14:57:09 MSK
Duration: 120.03s, Total samples = 2.48s ( 2.07%)
Showing nodes accounting for 0.41s, 16.53% of 2.48s total
Dropped 5 nodes (cum <= 0.01s)
flat  flat%   sum%        cum   cum%
0.07s  2.82%  2.82%      0.07s  2.82%  runtime.memclrNoHeapPointers
-0.05s  2.02%  0.81%     -0.05s  2.02%  runtime.spanClass.sizeclass (inline)
0.05s  2.02%  2.82%     -0.02s  0.81%  runtime.tryDeferToSpanScan
0.04s  1.61%  4.44%      0.04s  1.61%  internal/runtime/syscall/linux.Syscall6
-0.04s  1.61%  2.82%     -0.04s  1.61%  runtime.typePointers.next
-0.03s  1.21%  1.61%     -0.02s  0.81%  runtime.findObject
0.03s  1.21%  2.82%      0.18s  7.26%  runtime.newobject
0.03s  1.21%  4.03%      0.03s  1.21%  runtime.nextFreeFast (inline)
-0.03s  1.21%  2.82%     -0.03s  1.21%  runtime.procyieldAsm
0.03s  1.21%  4.03%      0.02s  0.81%  runtime.spanSetScans
0.02s  0.81%  4.84%      0.05s  2.02%  compress/flate.(*compressor).deflate
0.02s  0.81%  5.65%      0.11s  4.44%  compress/flate.(*compressor).init
0.02s  0.81%  6.45%      0.01s   0.4%  crypto/internal/fips140/hmac.(*HMAC).Sum
-0.02s  0.81%  5.65%     -0.02s  0.81%  crypto/internal/fips140/sha256.blockSHANI
0.02s  0.81%  6.45%      0.01s   0.4%  encoding/json.indirect
0.02s  0.81%  7.26%      0.02s  0.81%  indexbytebody
-0.02s  0.81%  6.45%     -0.02s  0.81%  internal/runtime/atomic.(*Uint32).CompareAndSwap (inline)
0.02s  0.81%  7.26%      0.02s  0.81%  internal/runtime/atomic.(*Uint32).Load (inline)
-0.02s  0.81%  6.45%     -0.01s   0.4%  internal/runtime/maps.(*Map).getWithoutKeySmallFastStr
0.02s  0.81%  7.26%      0.02s  0.81%  internal/sync.(*Mutex).Lock (inline)
0.02s  0.81%  8.06%      0.02s  0.81%  internal/sync.(*Mutex).Unlock (inline)
-0.02s  0.81%  7.26%     -0.02s  0.81%  net/http.Header.writeSubset
0.02s  0.81%  8.06%      0.02s  0.81%  runtime.(*gcBitsArena).tryAlloc (inline)
0.02s  0.81%  8.87%      0.05s  2.02%  runtime.(*mcentral).cacheSpan
0.02s  0.81%  9.68%      0.01s   0.4%  runtime.(*mheap).initSpan
-0.02s  0.81%  8.87%     -0.02s  0.81%  runtime.(*spanScanOwnership).or (inline)
-0.02s  0.81%  8.06%     -0.01s   0.4%  runtime.(*spanSet).pop
-0.02s  0.81%  7.26%     -0.02s  0.81%  runtime.findfunc
0.02s  0.81%  8.06%      0.02s  0.81%  runtime.mPark (inline)
0.02s  0.81%  8.87%      0.12s  4.84%  runtime.makeslice
0.02s  0.81%  9.68%      0.08s  3.23%  runtime.mallocgcSmallScanNoHeader
-0.02s  0.81%  8.87%     -0.03s  1.21%  runtime.markrootSpans
-0.02s  0.81%  8.06%     -0.02s  0.81%  runtime.memmove
-0.02s  0.81%  7.26%     -0.02s  0.81%  runtime.nanotime (inline)
-0.02s  0.81%  6.45%      0.02s  0.81%  runtime.scanObjectsSmall
-0.02s  0.81%  5.65%     -0.02s  0.81%  runtime.step
0.02s  0.81%  6.45%      0.02s  0.81%  runtime.typePointers.nextFast (inline)
0.01s   0.4%  6.85%      0.01s   0.4%  aeshashbody
-0.01s   0.4%  6.45%      0.05s  2.02%  bufio.(*Reader).fill
0.01s   0.4%  6.85%      0.06s  2.42%  compress/flate.(*compressor).initDeflate (inline)
-0.01s   0.4%  6.45%     -0.01s   0.4%  compress/flate.(*deflateFast).encode
0.01s   0.4%  6.85%      0.01s   0.4%  compress/flate.(*huffmanEncoder).bitCounts
0.01s   0.4%  7.26%      0.01s   0.4%  compress/flate.(*huffmanEncoder).bitLength (inline)
-0.01s   0.4%  6.85%      0.02s  0.81%  context.AfterFunc.func1
-0.01s   0.4%  6.45%     -0.01s   0.4%  context.WithValue
-0.01s   0.4%  6.05%     -0.01s   0.4%  crypto/internal/fips140.RecordApproved
0.01s   0.4%  6.45%      0.01s   0.4%  crypto/internal/fips140/hmac.(*HMAC).Write (inline)
0.01s   0.4%  6.85%      0.03s  1.21%  crypto/internal/fips140/pbkdf2.Key[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }]
-0.01s   0.4%  6.45%     -0.03s  1.21%  crypto/internal/fips140/sha256.(*Digest).Write
0.01s   0.4%  6.85%     -0.01s   0.4%  crypto/internal/fips140/sha256.(*Digest).checkSum
0.01s   0.4%  7.26%     -0.01s   0.4%  crypto/internal/fips140/sha256.block
-0.01s   0.4%  6.85%     -0.01s   0.4%  crypto/sha256.New
0.01s   0.4%  7.26%      0.04s  1.61%  encoding/json.(*decodeState).value
-0.01s   0.4%  6.85%     -0.01s   0.4%  encoding/json.(*encodeState).marshal.func1
-0.01s   0.4%  6.45%      0.04s  1.61%  encoding/json.Unmarshal
-0.01s   0.4%  6.05%     -0.01s   0.4%  encoding/json.stateInString
0.01s   0.4%  6.45%      0.01s   0.4%  encoding/json.stringEncoder
0.01s   0.4%  6.85%      0.01s   0.4%  encoding/json.unquoteBytes
0.01s   0.4%  7.26%      0.03s  1.21%  fmt.(*pp).handleMethods
0.01s   0.4%  7.66%     -0.01s   0.4%  github.com/go-chi/chi/v5/middleware.(*defaultLogEntry).Write
0.01s   0.4%  8.06%      0.03s  1.21%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).Build
-0.01s   0.4%  7.66%     -0.06s  2.42%  github.com/jackc/pgx/v5/pgconn.(*PgConn).ExecStatement
-0.01s   0.4%  7.26%      0.03s  1.21%  github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage
0.01s   0.4%  7.66%      0.04s  1.61%  github.com/jackc/pgx/v5/pgconn/internal/bgreader.(*BGReader).Read
0.01s   0.4%  8.06%      0.03s  1.21%  github.com/jackc/pgx/v5/pgxpool.(*Tx).Commit
-0.01s   0.4%  7.66%      0.01s   0.4%  github.com/liebeSonne/shortlink/internal/repository/database.(*shortLinkRepository).StoreAll.func1
-0.01s   0.4%  7.26%     -0.01s   0.4%  github.com/liebeSonne/shortlink/internal/service.(*shortIDGenerator).GenerateID
-0.01s   0.4%  6.85%      0.03s  1.21%  github.com/liebeSonne/shortlink/internal/service.(*shortLinkService).nextID
0.01s   0.4%  7.26%      0.03s  1.21%  github.com/liebeSonne/shortlink/internal/service.(*shortLinkService).nextID.func1
0.01s   0.4%  7.66%      0.01s   0.4%  internal/abi.(*Type).Pointers (inline)
-0.01s   0.4%  7.26%     -0.01s   0.4%  internal/abi.Name.IsExported (inline)
-0.01s   0.4%  6.85%     -0.01s   0.4%  internal/abi.Name.ReadVarint (inline)
-0.01s   0.4%  6.45%     -0.01s   0.4%  internal/bytealg.IndexByte
0.01s   0.4%  6.85%      0.01s   0.4%  internal/bytealg.IndexByteString
-0.01s   0.4%  6.45%     -0.05s  2.02%  internal/poll.(*FD).Write
0.01s   0.4%  6.85%      0.01s   0.4%  internal/poll.(*fdMutex).rwlock
-0.01s   0.4%  6.45%     -0.01s   0.4%  internal/reflectlite.resolveNameOff
0.01s   0.4%  6.85%      0.01s   0.4%  internal/runtime/atomic.(*Uint64).CompareAndSwap (inline)
-0.01s   0.4%  6.45%     -0.01s   0.4%  internal/runtime/atomic.(*Uint8).Load (inline)
-0.01s   0.4%  6.05%     -0.01s   0.4%  internal/runtime/atomic.(*UnsafePointer).Load (inline)
0.01s   0.4%  6.45%      0.01s   0.4%  internal/runtime/gc/scan.FilterNilAVX512
-0.01s   0.4%  6.05%     -0.01s   0.4%  internal/runtime/gc/scan.ScanSpanPackedAVX512
-0.01s   0.4%  5.65%     -0.01s   0.4%  internal/runtime/gc/scan.scanSpanPackedAVX512
-0.01s   0.4%  5.24%     -0.01s   0.4%  internal/runtime/maps.(*Map).Used (inline)
-0.01s   0.4%  4.84%     -0.01s   0.4%  internal/strconv.formatBase10
0.01s   0.4%  5.24%      0.01s   0.4%  internal/strconv.formatDigits
0.01s   0.4%  5.65%      0.02s  0.81%  internal/stringslite.Cut
0.01s   0.4%  6.05%      0.05s  2.02%  io.ReadAtLeast
0.01s   0.4%  6.45%      0.01s   0.4%  net/http.(*Cookie).String
-0.01s   0.4%  6.05%     -0.01s   0.4%  net/http.(*Request).expectsContinue
-0.01s   0.4%  5.65%     -0.02s  0.81%  net/http.(*chunkWriter).writeHeader
0.01s   0.4%  6.05%      0.01s   0.4%  net/http.(*connReader).setInfiniteReadLimit (inline)
-0.01s   0.4%  5.65%      0.42s 16.94%  net/http.HandlerFunc.ServeHTTP
0.01s   0.4%  6.05%     -0.02s  0.81%  net/http.Redirect
-0.01s   0.4%  5.65%     -0.02s  0.81%  net/http.parseCookieValue
0.01s   0.4%  6.05%      0.01s   0.4%  net/http.readCookies
-0.01s   0.4%  5.65%     -0.01s   0.4%  net/http.validCookieValueByte (inline)
0.01s   0.4%  6.05%      0.01s   0.4%  net/textproto.(*Reader).upcomingHeaderKeys
0.01s   0.4%  6.45%      0.01s   0.4%  net/textproto.MIMEHeader.Add (inline)
-0.01s   0.4%  6.05%     -0.01s   0.4%  net/textproto.trim (inline)
-0.01s   0.4%  5.65%     -0.01s   0.4%  net/url.parse
-0.01s   0.4%  5.24%     -0.01s   0.4%  reflect.(*structType).Field
0.01s   0.4%  5.65%      0.01s   0.4%  reflect.TypeOf (inline)
-0.01s   0.4%  5.24%     -0.01s   0.4%  reflect.Value.Equal
0.01s   0.4%  5.65%      0.01s   0.4%  reflect.Value.SetMapIndex
0.01s   0.4%  6.05%      0.01s   0.4%  reflect.Value.SetString
0.01s   0.4%  6.45%      0.01s   0.4%  reflect.implements
0.01s   0.4%  6.85%      0.01s   0.4%  runtime.(*_panic).nextDefer
-0.01s   0.4%  6.45%     -0.02s  0.81%  runtime.(*activeSweep).begin (inline)
0.01s   0.4%  6.85%      0.01s   0.4%  runtime.(*bucket).mp
0.01s   0.4%  7.26%      0.01s   0.4%  runtime.(*fixalloc).alloc
0.01s   0.4%  7.66%      0.01s   0.4%  runtime.(*gcControllerState).heapGoalInternal
-0.01s   0.4%  7.26%     -0.02s  0.81%  runtime.(*gcWork).balance
0.01s   0.4%  7.66%      0.01s   0.4%  runtime.(*gcWork).dispose
0.01s   0.4%  8.06%      0.03s  1.21%  runtime.(*mcache).releaseAll
-0.01s   0.4%  7.66%     -0.04s  1.61%  runtime.(*mheap).nextSpanForSweep
0.01s   0.4%  8.06%      0.01s   0.4%  runtime.(*mspan).heapBitsSmallForAddr
0.01s   0.4%  8.47%      0.01s   0.4%  runtime.(*mspan).moveInlineMarks
0.01s   0.4%  8.87%      0.01s   0.4%  runtime.(*mspan).nextFreeIndex
0.01s   0.4%  9.27%      0.01s   0.4%  runtime.(*mspan).refillAllocCache
-0.01s   0.4%  8.87%     -0.01s   0.4%  runtime.(*scavengeIndex).alloc
-0.01s   0.4%  8.47%     -0.01s   0.4%  runtime.(*stackScanState).getPtr
-0.01s   0.4%  8.06%     -0.01s   0.4%  runtime.(*stkframe).getStackMap
0.01s   0.4%  8.47%      0.01s   0.4%  runtime.(*sweepLocker).tryAcquire
0.01s   0.4%  8.87%     -0.01s   0.4%  runtime.(*timers).check
-0.01s   0.4%  8.47%     -0.01s   0.4%  runtime.acquireSudog
-0.01s   0.4%  8.06%     -0.01s   0.4%  runtime.addb (inline)
0.01s   0.4%  8.47%      0.01s   0.4%  runtime.arenaIndex (inline)
0.01s   0.4%  8.87%      0.01s   0.4%  runtime.asyncPreempt
0.01s   0.4%  9.27%      0.01s   0.4%  runtime.binarySearchTree
0.01s   0.4%  9.68%      0.01s   0.4%  runtime.findBitRange64 (inline)
0.01s   0.4% 10.08%     -0.04s  1.61%  runtime.gcBgMarkWorker
0.01s   0.4% 10.48%      0.04s  1.61%  runtime.gcNextMarkRoot
0.01s   0.4% 10.89%      0.01s   0.4%  runtime.gcmarknewobject
-0.01s   0.4% 10.48%     -0.02s  0.81%  runtime.getGCMask (inline)
-0.01s   0.4% 10.08%     -0.01s   0.4%  runtime.getGCMaskOnDemand
-0.01s   0.4%  9.68%     -0.01s   0.4%  runtime.getempty
0.01s   0.4% 10.08%      0.01s   0.4%  runtime.gorecover
0.01s   0.4% 10.48%      0.01s   0.4%  runtime.greyobject
0.01s   0.4% 10.89%      0.01s   0.4%  runtime.headTailIndex.head (inline)
-0.01s   0.4% 10.48%     -0.01s   0.4%  runtime.ifaceeq
0.01s   0.4% 10.89%      0.01s   0.4%  runtime.limiterEventStamp.typ (inline)
0.01s   0.4% 11.29%      0.06s  2.42%  runtime.mallocgcSmallNoscan
-0.01s   0.4% 10.89%     -0.02s  0.81%  runtime.mallocgcSmallScanHeader
-0.01s   0.4% 10.48%     -0.01s   0.4%  runtime.mapaccess2_faststr
0.01s   0.4% 10.89%     -0.01s   0.4%  runtime.markroot
0.01s   0.4% 11.29%     -0.01s   0.4%  runtime.mcall
-0.01s   0.4% 10.89%     -0.01s   0.4%  runtime.newArenaMayUnlock
0.01s   0.4% 11.29%     -0.04s  1.61%  runtime.newstack
0.01s   0.4% 11.69%      0.01s   0.4%  runtime.objptr.objIndex (inline)
0.01s   0.4% 12.10%      0.01s   0.4%  runtime.offAddr.add (inline)
0.01s   0.4% 12.50%      0.01s   0.4%  runtime.offAddrToLevelIndex (inline)
-0.01s   0.4% 12.10%     -0.01s   0.4%  runtime.pageIndexOf (inline)
0.01s   0.4% 12.50%     -0.01s   0.4%  runtime.pcvalue
0.01s   0.4% 12.90%      0.01s   0.4%  runtime.save
-0.01s   0.4% 12.50%     -0.12s  4.84%  runtime.scanObject
0.01s   0.4% 12.90%      0.01s   0.4%  runtime.scanObjectSmall
0.01s   0.4% 13.31%      0.01s   0.4%  runtime.spanOf (inline)
0.01s   0.4% 13.71%      0.01s   0.4%  runtime.spanOfUnchecked (inline)
0.01s   0.4% 14.11%      0.03s  1.21%  runtime.stackcache_clear
0.01s   0.4% 14.52%      0.04s  1.61%  runtime.stealWork
-0.01s   0.4% 14.11%     -0.06s  2.42%  runtime.sweepone
0.01s   0.4% 14.52%      0.01s   0.4%  runtime.usleep
0.01s   0.4% 14.92%      0.01s   0.4%  sync.(*Pool).Put
0.01s   0.4% 15.32%      0.01s   0.4%  sync.(*poolChain).popTail
0.01s   0.4% 15.73%      0.01s   0.4%  sync.indexLocal (inline)
0.01s   0.4% 16.13%      0.01s   0.4%  sync.runtime_notifyListAdd
0.01s   0.4% 16.53%      0.01s   0.4%  sync/atomic.CompareAndSwapPointer
-0.01s   0.4% 16.13%     -0.01s   0.4%  time.Time.appendFormat
0.01s   0.4% 16.53%      0.01s   0.4%  time.runtimeNow
0     0% 16.53%      0.05s  2.02%  bufio.(*Reader).Peek
0     0% 16.53%     -0.03s  1.21%  bufio.(*Writer).Flush
0     0% 16.53%      0.01s   0.4%  bytes.(*Buffer).Write
0     0% 16.53%      0.01s   0.4%  bytes.(*Buffer).grow
0     0% 16.53%     -0.01s   0.4%  bytes.(*Buffer).readSlice
0     0% 16.53%     -0.01s   0.4%  bytes.IndexByte (inline)
0     0% 16.53%      0.01s   0.4%  bytes.growSlice
0     0% 16.53%      0.01s   0.4%  bytes.growSlice.func1
0     0% 16.53%      0.05s  2.02%  compress/flate.(*Writer).Close (inline)
0     0% 16.53%     -0.01s   0.4%  compress/flate.(*Writer).Write (inline)
0     0% 16.53%      0.05s  2.02%  compress/flate.(*compressor).close
0     0% 16.53%     -0.01s   0.4%  compress/flate.(*compressor).encSpeed
0     0% 16.53%     -0.01s   0.4%  compress/flate.(*compressor).write
0     0% 16.53%      0.02s  0.81%  compress/flate.(*compressor).writeBlock
0     0% 16.53%      0.01s   0.4%  compress/flate.(*huffmanBitWriter).fixedSize (inline)
0     0% 16.53%      0.01s   0.4%  compress/flate.(*huffmanBitWriter).indexTokens
0     0% 16.53%      0.02s  0.81%  compress/flate.(*huffmanBitWriter).writeBlock
0     0% 16.53%      0.01s   0.4%  compress/flate.(*huffmanEncoder).generate
0     0% 16.53%      0.19s  7.66%  compress/flate.NewWriter (inline)
0     0% 16.53%      0.03s  1.21%  compress/flate.newHuffmanBitWriter (inline)
0     0% 16.53%      0.02s  0.81%  compress/flate.newHuffmanEncoder (inline)
0     0% 16.53%      0.04s  1.61%  compress/gzip.(*Writer).Close
0     0% 16.53%      0.18s  7.26%  compress/gzip.(*Writer).Write
0     0% 16.53%      0.02s  0.81%  context.(*afterFuncCtx).cancel
0     0% 16.53%      0.01s   0.4%  context.(*cancelCtx).cancel
0     0% 16.53%     -0.01s   0.4%  context.(*cancelCtx).propagateCancel
0     0% 16.53%     -0.01s   0.4%  context.AfterFunc
0     0% 16.53%      0.01s   0.4%  context.removeChild
0     0% 16.53%     -0.01s   0.4%  crypto.Hash.New
0     0% 16.53%     -0.01s   0.4%  crypto/hmac.New.UnwrapNew[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }].func1
0     0% 16.53%      0.03s  1.21%  crypto/pbkdf2.Key[go.shape.interface { BlockSize int; Reset; Size int; Sum []uint8; Write  }]
0     0% 16.53%      0.01s   0.4%  encoding/hex.DecodeString
0     0% 16.53%     -0.01s   0.4%  encoding/json.(*Decoder).readValue
0     0% 16.53%      0.08s  3.23%  encoding/json.(*Encoder).Encode
0     0% 16.53%      0.01s   0.4%  encoding/json.(*decodeState).literalStore
0     0% 16.53%      0.04s  1.61%  encoding/json.(*decodeState).object
0     0% 16.53%      0.04s  1.61%  encoding/json.(*decodeState).unmarshal
0     0% 16.53%     -0.02s  0.81%  encoding/json.(*encodeState).marshal
0     0% 16.53%     -0.01s   0.4%  encoding/json.(*encodeState).reflectValue
0     0% 16.53%     -0.02s  0.81%  encoding/json.Marshal
0     0% 16.53%      0.01s   0.4%  encoding/json.NewDecoder (inline)
0     0% 16.53%     -0.01s   0.4%  encoding/json.appendString[go.shape.string]
0     0% 16.53%      0.01s   0.4%  encoding/json.interfaceEncoder
0     0% 16.53%     -0.01s   0.4%  encoding/json.newEncodeState
0     0% 16.53%     -0.01s   0.4%  encoding/json.structEncoder.encode
0     0% 16.53%      0.01s   0.4%  errors.(*joinError).Error
0     0% 16.53%     -0.02s  0.81%  errors.As
0     0% 16.53%      0.02s  0.81%  fmt.(*buffer).writeString (inline)
0     0% 16.53%      0.02s  0.81%  fmt.(*fmt).fmtS
0     0% 16.53%      0.02s  0.81%  fmt.(*fmt).padString
0     0% 16.53%      0.01s   0.4%  fmt.(*pp).badVerb
0     0% 16.53%      0.04s  1.61%  fmt.(*pp).doPrintf
0     0% 16.53%      0.02s  0.81%  fmt.(*pp).fmtString
0     0% 16.53%      0.04s  1.61%  fmt.(*pp).printArg
0     0% 16.53%      0.01s   0.4%  fmt.(*pp).printValue
0     0% 16.53%      0.02s  0.81%  fmt.Errorf (inline)
0     0% 16.53%      0.04s  1.61%  fmt.Fprintf
0     0% 16.53%     -0.01s   0.4%  fmt.Fprintln
0     0% 16.53%      0.02s  0.81%  fmt.Printf (inline)
0     0% 16.53%      0.01s   0.4%  fmt.Sprintf
0     0% 16.53%      0.02s  0.81%  fmt.errorf
0     0% 16.53%      0.01s   0.4%  fmt.newPrinter
0     0% 16.53%      0.04s  1.61%  github.com/avast/retry-go.Do
0     0% 16.53%     -0.01s   0.4%  github.com/go-chi/chi/v5.(*Mux).Mount.func1
0     0% 16.53%      0.28s 11.29%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
0     0% 16.53%      0.28s 11.29%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
0     0% 16.53%      0.02s  0.81%  github.com/go-chi/chi/v5/middleware.(*DefaultLogFormatter).NewLogEntry
0     0% 16.53%      0.20s  8.06%  github.com/go-chi/chi/v5/middleware.(*basicWriter).Write
0     0% 16.53%      0.01s   0.4%  github.com/go-chi/chi/v5/middleware.(*basicWriter).WriteHeader
0     0% 16.53%     -0.01s   0.4%  github.com/go-chi/chi/v5/middleware.NoCache.func1
0     0% 16.53%      0.01s   0.4%  github.com/go-chi/chi/v5/middleware.cW
0     0% 16.53%      0.28s 11.29%  github.com/go-chi/chi/v5/middleware.init.0.RequestLogger.func1.1
0     0% 16.53%     -0.01s   0.4%  github.com/go-chi/chi/v5/middleware.init.0.RequestLogger.func1.1.1
0     0% 16.53%      0.01s   0.4%  github.com/golang-jwt/jwt/v4.(*NumericDate).UnmarshalJSON
0     0% 16.53%      0.04s  1.61%  github.com/golang-jwt/jwt/v4.(*Parser).ParseUnverified
0     0% 16.53%      0.06s  2.42%  github.com/golang-jwt/jwt/v4.(*Parser).ParseWithClaims
0     0% 16.53%     -0.01s   0.4%  github.com/golang-jwt/jwt/v4.(*SigningMethodHMAC).Sign
0     0% 16.53%      0.01s   0.4%  github.com/golang-jwt/jwt/v4.(*SigningMethodHMAC).Verify
0     0% 16.53%     -0.04s  1.61%  github.com/golang-jwt/jwt/v4.(*Token).SignedString
0     0% 16.53%     -0.02s  0.81%  github.com/golang-jwt/jwt/v4.(*Token).SigningString
0     0% 16.53%      0.06s  2.42%  github.com/golang-jwt/jwt/v4.ParseWithClaims
0     0% 16.53%      0.01s   0.4%  github.com/google/uuid.New (inline)
0     0% 16.53%      0.01s   0.4%  github.com/google/uuid.NewRandom
0     0% 16.53%      0.01s   0.4%  github.com/google/uuid.NewRandomFromReader
0     0% 16.53%     -0.01s   0.4%  github.com/google/uuid.UUID.String (inline)
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5.(*Conn).Exec
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5.(*Conn).Ping (inline)
0     0% 16.53%     -0.01s   0.4%  github.com/jackc/pgx/v5.(*Conn).Prepare
0     0% 16.53%     -0.04s  1.61%  github.com/jackc/pgx/v5.(*Conn).Query
0     0% 16.53%     -0.04s  1.61%  github.com/jackc/pgx/v5.(*Conn).QueryRow (inline)
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5.(*Conn).exec
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5.(*Conn).execPrepared
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5.(*Conn).execSimpleProtocol
0     0% 16.53%     -0.01s   0.4%  github.com/jackc/pgx/v5.(*Conn).getStatementDescription
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).appendParam
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).encodeExtendedParamValue
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5.(*dbTx).Commit
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5.(*dbTx).Exec
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5.ConnectConfig
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5.ParseConfig (inline)
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5.ParseConfigWithOptions
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5.connect
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).Close (inline)
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).NextResult
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).receiveMessage
0     0% 16.53%     -0.01s   0.4%  github.com/jackc/pgx/v5/pgconn.(*PgConn).Exec
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgconn.(*PgConn).Ping
0     0% 16.53%     -0.01s   0.4%  github.com/jackc/pgx/v5/pgconn.(*PgConn).Prepare
0     0% 16.53%     -0.05s  2.02%  github.com/jackc/pgx/v5/pgconn.(*PgConn).execExtendedSuffix
0     0% 16.53%     -0.04s  1.61%  github.com/jackc/pgx/v5/pgconn.(*PgConn).flushWithPotentialWriteReadDeadlock
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5/pgconn.(*PgConn).scramAuth
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5/pgconn.(*ResultReader).Close
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5/pgconn.(*ResultReader).Read
0     0% 16.53%     -0.02s  0.81%  github.com/jackc/pgx/v5/pgconn.(*ResultReader).readUntilRowDescription
0     0% 16.53%     -0.01s   0.4%  github.com/jackc/pgx/v5/pgconn.(*ResultReader).receiveMessage
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5/pgconn.(*scramClient).clientFinalMessage
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5/pgconn.ConnectConfig
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgconn.ParseConfigWithOptions
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5/pgconn.connectOne
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5/pgconn.connectPreferred
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgconn.defaultSettings
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5/pgconn/ctxwatch.(*ContextWatcher).Unwatch
0     0% 16.53%     -0.01s   0.4%  github.com/jackc/pgx/v5/pgconn/ctxwatch.(*ContextWatcher).Watch
0     0% 16.53%     -0.01s   0.4%  github.com/jackc/pgx/v5/pgproto3.(*ErrorResponse).Decode
0     0% 16.53%     -0.04s  1.61%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).Flush
0     0% 16.53%      0.04s  1.61%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgproto3.(*ParameterStatus).Decode
0     0% 16.53%      0.04s  1.61%  github.com/jackc/pgx/v5/pgproto3.(*chunkReader).Next
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5/pgtype.(*Map).Encode
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgtype.(*Map).Scan
0     0% 16.53%      0.02s  0.81%  github.com/jackc/pgx/v5/pgtype.(*encodePlanDriverValuer).Encode
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgtype.(*pointerEmptyInterfaceScanPlan).Scan
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgtype.UUIDCodec.DecodeValue
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgtype.codecScan
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgtype.parseUUID
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgtype.scanPlanTextAnyToUUIDScanner.Scan
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgxpool.(*Conn).Ping
0     0% 16.53%     -0.04s  1.61%  github.com/jackc/pgx/v5/pgxpool.(*Conn).QueryRow
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgxpool.(*Pool).Ping
0     0% 16.53%     -0.04s  1.61%  github.com/jackc/pgx/v5/pgxpool.(*Pool).QueryRow
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgxpool.(*Tx).Exec
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgxpool.New
0     0% 16.53%      0.03s  1.21%  github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func3
0     0% 16.53%      0.01s   0.4%  github.com/jackc/pgx/v5/pgxpool.ParseConfig
0     0% 16.53%      0.03s  1.21%  github.com/jackc/puddle/v2.(*Pool[go.shape.*uint8]).initResourceValue.func1
0     0% 16.53%     -0.04s  1.61%  github.com/liebeSonne/shortlink/internal/auth.(*tokenServiceImpl).Create
0     0% 16.53%      0.06s  2.42%  github.com/liebeSonne/shortlink/internal/auth.(*tokenServiceImpl).Parse
0     0% 16.53%      0.02s  0.81%  github.com/liebeSonne/shortlink/internal/handler.(*databaseHandler).HandlePing
0     0% 16.53%      0.01s   0.4%  github.com/liebeSonne/shortlink/internal/handler.(*loggingResponseWriter).WriteHeader
0     0% 16.53%      0.15s  6.05%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).HandleCreate
0     0% 16.53%      0.07s  2.82%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).HandleCreateShorten
0     0% 16.53%      0.08s  3.23%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).HandleCreateShortenBatch
0     0% 16.53%     -0.07s  2.82%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).HandleGet
0     0% 16.53%      0.04s  1.61%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).HandleGetUserUrls
0     0% 16.53%      0.02s  0.81%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).createShortLink
0     0% 16.53%     -0.01s   0.4%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).createShortLinkURL (inline)
0     0% 16.53%      0.03s  1.21%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).responseError
0     0% 16.53%     -0.01s   0.4%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).sendAuditEvent
0     0% 16.53%      0.38s 15.32%  github.com/liebeSonne/shortlink/internal/handler/compress.NewCompressorMiddleware.NewEncodingMiddleware.NewGzipHandlerMiddleware.func1
0     0% 16.53%      0.04s  1.61%  github.com/liebeSonne/shortlink/internal/handler/compress.NewCompressorMiddleware.NewEncodingMiddleware.NewGzipHandlerMiddleware.func1.3
0     0% 16.53%      0.04s  1.61%  github.com/liebeSonne/shortlink/internal/handler/compress/gzip.(*gzipWriter).Close
0     0% 16.53%      0.20s  8.06%  github.com/liebeSonne/shortlink/internal/handler/compress/gzip.(*gzipWriter).Write
0     0% 16.53%      0.01s   0.4%  github.com/liebeSonne/shortlink/internal/handler/cookie.(*cookieServiceImpl).GetAuthToken
0     0% 16.53%      0.02s  0.81%  github.com/liebeSonne/shortlink/internal/handler/cookie.(*cookieServiceImpl).SetAuthToken
0     0% 16.53%      0.01s   0.4%  github.com/liebeSonne/shortlink/internal/logger.(*zapLoggerImpl).Errorf
0     0% 16.53%      0.02s  0.81%  github.com/liebeSonne/shortlink/internal/logger.(*zapLoggerImpl).Infow
0     0% 16.53%      0.01s   0.4%  github.com/liebeSonne/shortlink/internal/repository.(*ErrConflictURL).Error
0     0% 16.53%      0.01s   0.4%  github.com/liebeSonne/shortlink/internal/repository.NewErrConflictURL (inline)
0     0% 16.53%      0.02s  0.81%  github.com/liebeSonne/shortlink/internal/repository/database.(*database).Ping
0     0% 16.53%     -0.01s   0.4%  github.com/liebeSonne/shortlink/internal/repository/database.(*shortLinkRepository).Find
0     0% 16.53%     -0.02s  0.81%  github.com/liebeSonne/shortlink/internal/repository/database.(*shortLinkRepository).FindByURL
0     0% 16.53%      0.05s  2.02%  github.com/liebeSonne/shortlink/internal/repository/database.(*shortLinkRepository).Store
0     0% 16.53%      0.06s  2.42%  github.com/liebeSonne/shortlink/internal/repository/database.(*shortLinkRepository).StoreAll
0     0% 16.53%      0.04s  1.61%  github.com/liebeSonne/shortlink/internal/service.(*shortLinkService).Create
0     0% 16.53%      0.05s  2.02%  github.com/liebeSonne/shortlink/internal/service.(*shortLinkService).CreateBatch
0     0% 16.53%      0.01s   0.4%  github.com/liebeSonne/shortlink/internal/service.(*userServiceImpl).NextID
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap.(*Logger).Check (inline)
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap.(*Logger).check
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap.(*SugaredLogger).Errorf (inline)
0     0% 16.53%      0.02s  0.81%  go.uber.org/zap.(*SugaredLogger).Infow (inline)
0     0% 16.53%      0.03s  1.21%  go.uber.org/zap.(*SugaredLogger).log
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap.(*SugaredLogger).sweetenFields
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap.getMessage
0     0% 16.53%     -0.01s   0.4%  go.uber.org/zap/buffer.(*Buffer).AppendTime (inline)
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap/internal/pool.(*Pool[go.shape.*uint8]).Get (inline)
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap/zapcore.(*CheckedEntry).AddCore (inline)
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap/zapcore.(*ioCore).Check
0     0% 16.53%     -0.01s   0.4%  go.uber.org/zap/zapcore.(*jsonEncoder).AddTime
0     0% 16.53%     -0.01s   0.4%  go.uber.org/zap/zapcore.(*jsonEncoder).AppendTime
0     0% 16.53%     -0.01s   0.4%  go.uber.org/zap/zapcore.(*jsonEncoder).AppendTimeLayout
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap/zapcore.Field.AddTo
0     0% 16.53%     -0.01s   0.4%  go.uber.org/zap/zapcore.ISO8601TimeEncoder
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap/zapcore.addFields (inline)
0     0% 16.53%     -0.01s   0.4%  go.uber.org/zap/zapcore.encodeTimeLayout
0     0% 16.53%      0.01s   0.4%  go.uber.org/zap/zapcore.getCheckedEntry
0     0% 16.53%     -0.01s   0.4%  internal/abi.Name.Name
0     0% 16.53%      0.02s  0.81%  internal/poll.(*FD).Accept
0     0% 16.53%      0.08s  3.23%  internal/poll.(*FD).Read
0     0% 16.53%      0.01s   0.4%  internal/poll.(*FD).writeLock (inline)
0     0% 16.53%      0.02s  0.81%  internal/poll.accept
0     0% 16.53%      0.03s  1.21%  internal/poll.ignoringEINTRIO (inline)
0     0% 16.53%     -0.02s  0.81%  internal/reflectlite.implements
0     0% 16.53%     -0.02s  0.81%  internal/reflectlite.rtype.Implements
0     0% 16.53%     -0.01s   0.4%  internal/reflectlite.rtype.nameOff (inline)
0     0% 16.53%     -0.01s   0.4%  internal/runtime/gc/scan.ScanSpanPacked (inline)
0     0% 16.53%      0.01s   0.4%  internal/runtime/maps.(*Map).Delete
0     0% 16.53%      0.01s   0.4%  internal/runtime/maps.(*Map).deleteSmall
0     0% 16.53%      0.01s   0.4%  internal/runtime/maps.(*Map).growToSmall
0     0% 16.53%      0.01s   0.4%  internal/runtime/maps.newGroups (inline)
0     0% 16.53%      0.01s   0.4%  internal/runtime/maps.newarray
0     0% 16.53%      0.01s   0.4%  internal/runtime/maps.typedmemclr
0     0% 16.53%      0.01s   0.4%  internal/strconv.FormatFloat (inline)
0     0% 16.53%     -0.01s   0.4%  internal/strconv.FormatInt
0     0% 16.53%      0.01s   0.4%  internal/strconv.genericFtoa
0     0% 16.53%      0.01s   0.4%  internal/stringslite.Index
0     0% 16.53%      0.01s   0.4%  internal/stringslite.IndexByte (inline)
0     0% 16.53%      0.01s   0.4%  io.ReadAll
0     0% 16.53%      0.01s   0.4%  io.ReadFull (inline)
0     0% 16.53%     -0.01s   0.4%  log.(*Logger).Print
0     0% 16.53%     -0.01s   0.4%  log.(*Logger).output
0     0% 16.53%      0.42s 16.94%  main.initRouter.LoggingMiddleware.func5
0     0% 16.53%      0.34s 13.71%  main.initRouter.NewAuthCookieMiddleware.func4
0     0% 16.53%      0.31s 12.50%  main.initRouter.NewAuthMiddleware.func3
0     0% 16.53%      0.03s  1.21%  main.runApp.func1
0     0% 16.53%     -0.01s   0.4%  net.(*Dialer).DialContext
0     0% 16.53%      0.01s   0.4%  net.(*TCPConn).SetKeepAliveConfig
0     0% 16.53%      0.03s  1.21%  net.(*TCPListener).Accept
0     0% 16.53%      0.03s  1.21%  net.(*TCPListener).accept
0     0% 16.53%      0.10s  4.03%  net.(*conn).Read
0     0% 16.53%     -0.05s  2.02%  net.(*conn).Write
0     0% 16.53%      0.08s  3.23%  net.(*netFD).Read
0     0% 16.53%     -0.05s  2.02%  net.(*netFD).Write
0     0% 16.53%      0.02s  0.81%  net.(*netFD).accept
0     0% 16.53%     -0.01s   0.4%  net.(*sysDialer).dialParallel
0     0% 16.53%     -0.01s   0.4%  net.(*sysDialer).dialSerial
0     0% 16.53%     -0.01s   0.4%  net.(*sysDialer).dialSingle
0     0% 16.53%     -0.01s   0.4%  net.(*sysDialer).dialTCP
0     0% 16.53%     -0.01s   0.4%  net.(*sysDialer).doDialTCP (inline)
0     0% 16.53%     -0.01s   0.4%  net.(*sysDialer).doDialTCPProto
0     0% 16.53%      0.01s   0.4%  net.setKeepAliveIdle
0     0% 16.53%     -0.01s   0.4%  net.setNoDelay
0     0% 16.53%      0.01s   0.4%  net/http.(*Request).Cookie (inline)
0     0% 16.53%     -0.01s   0.4%  net/http.(*Request).wantsClose
0     0% 16.53%      0.03s  1.21%  net/http.(*Server).ListenAndServe
0     0% 16.53%      0.03s  1.21%  net/http.(*Server).Serve
0     0% 16.53%     -0.01s   0.4%  net/http.(*body).Read
0     0% 16.53%     -0.01s   0.4%  net/http.(*body).readLocked
0     0% 16.53%     -0.02s  0.81%  net/http.(*chunkWriter).Write
0     0% 16.53%      0.02s  0.81%  net/http.(*conn).readRequest
0     0% 16.53%      0.46s 18.55%  net/http.(*conn).serve
0     0% 16.53%      0.01s   0.4%  net/http.(*conn).serve.func1
0     0% 16.53%      0.01s   0.4%  net/http.(*conn).setState
0     0% 16.53%      0.06s  2.42%  net/http.(*connReader).Read
0     0% 16.53%      0.04s  1.61%  net/http.(*connReader).backgroundRead
0     0% 16.53%      0.01s   0.4%  net/http.(*connReader).lock
0     0% 16.53%      0.01s   0.4%  net/http.(*response).WriteHeader
0     0% 16.53%     -0.03s  1.21%  net/http.(*response).finishRequest
0     0% 16.53%      0.02s  0.81%  net/http.Error
0     0% 16.53%      0.01s   0.4%  net/http.Header.Add (inline)
0     0% 16.53%      0.01s   0.4%  net/http.Header.Clone (inline)
0     0% 16.53%      0.01s   0.4%  net/http.Header.Set (inline)
0     0% 16.53%     -0.02s  0.81%  net/http.Header.WriteSubset (inline)
0     0% 16.53%      0.02s  0.81%  net/http.SetCookie
0     0% 16.53%     -0.01s   0.4%  net/http.checkConnErrorWriter.Write
0     0% 16.53%      0.01s   0.4%  net/http.parseRequestLine
0     0% 16.53%      0.02s  0.81%  net/http.readRequest
0     0% 16.53%      0.42s 16.94%  net/http.serverHandler.ServeHTTP
0     0% 16.53%     -0.01s   0.4%  net/http/pprof.collectProfile
0     0% 16.53%     -0.01s   0.4%  net/http/pprof.handler.ServeHTTP
0     0% 16.53%     -0.01s   0.4%  net/http/pprof.handler.serveDeltaProfile
0     0% 16.53%      0.01s   0.4%  net/textproto.(*Reader).ReadMIMEHeader (inline)
0     0% 16.53%     -0.01s   0.4%  net/textproto.(*Reader).readContinuedLineSlice
0     0% 16.53%      0.01s   0.4%  net/textproto.MIMEHeader.Set (inline)
0     0% 16.53%      0.01s   0.4%  net/textproto.readMIMEHeader
0     0% 16.53%     -0.01s   0.4%  net/url.ParseRequestURI
0     0% 16.53%      0.01s   0.4%  os.Stat
0     0% 16.53%      0.01s   0.4%  os.ignoringEINTR (inline)
0     0% 16.53%      0.01s   0.4%  os.statNolog
0     0% 16.53%      0.01s   0.4%  os.statNolog.func1 (inline)
0     0% 16.53%      0.01s   0.4%  reflect.(*rtype).AssignableTo
0     0% 16.53%     -0.01s   0.4%  reflect.(*rtype).Field
0     0% 16.53%     -0.01s   0.4%  reflect.Value.Field
0     0% 16.53%      0.01s   0.4%  reflect.Value.Interface (inline)
0     0% 16.53%      0.01s   0.4%  reflect.packEface (inline)
0     0% 16.53%      0.01s   0.4%  reflect.packEfaceData
0     0% 16.53%      0.01s   0.4%  reflect.unsafe_New
0     0% 16.53%      0.01s   0.4%  reflect.valueInterface
0     0% 16.53%     -0.01s   0.4%  runtime.(*activeSweep).end
0     0% 16.53%      0.02s  0.81%  runtime.(*gcCPULimiterState).update
0     0% 16.53%      0.02s  0.81%  runtime.(*gcCPULimiterState).updateLocked
0     0% 16.53%     -0.01s   0.4%  runtime.(*gcControllerState).findRunnableGCWorker
0     0% 16.53%      0.01s   0.4%  runtime.(*gcControllerState).trigger
0     0% 16.53%     -0.02s  0.81%  runtime.(*gcWork).tryStealSpan
0     0% 16.53%      0.02s  0.81%  runtime.(*limiterEvent).consume
0     0% 16.53%      0.06s  2.42%  runtime.(*mcache).nextFree
0     0% 16.53%      0.05s  2.02%  runtime.(*mcache).prepareForSweep
0     0% 16.53%      0.05s  2.02%  runtime.(*mcache).refill
0     0% 16.53%      0.02s  0.81%  runtime.(*mcentral).uncacheSpan
0     0% 16.53%      0.02s  0.81%  runtime.(*mheap).alloc
0     0% 16.53%      0.02s  0.81%  runtime.(*mheap).alloc.func1
0     0% 16.53%      0.01s   0.4%  runtime.(*mheap).allocManual
0     0% 16.53%      0.03s  1.21%  runtime.(*mheap).allocSpan
0     0% 16.53%      0.02s  0.81%  runtime.(*moduledata).funcName
0     0% 16.53%      0.03s  1.21%  runtime.(*pageAlloc).alloc
0     0% 16.53%     -0.01s   0.4%  runtime.(*pageAlloc).allocToCache
0     0% 16.53%      0.03s  1.21%  runtime.(*pageAlloc).find
0     0% 16.53%      0.01s   0.4%  runtime.(*pageAlloc).find.func1
0     0% 16.53%      0.01s   0.4%  runtime.(*pallocBits).find
0     0% 16.53%      0.01s   0.4%  runtime.(*pallocBits).findSmallN
0     0% 16.53%     -0.02s  0.81%  runtime.(*spanInlineMarkBits).tryAcquire
0     0% 16.53%     -0.01s   0.4%  runtime.(*spanQueue).steal
0     0% 16.53%      0.01s   0.4%  runtime.(*stackScanState).buildIndex (inline)
0     0% 16.53%      0.03s  1.21%  runtime.(*sweepLocked).sweep
0     0% 16.53%     -0.01s   0.4%  runtime.(*unwinder).init (inline)
0     0% 16.53%     -0.01s   0.4%  runtime.(*unwinder).initAt
0     0% 16.53%     -0.02s  0.81%  runtime.(*unwinder).next
0     0% 16.53%     -0.02s  0.81%  runtime.(*unwinder).resolveInternal
0     0% 16.53%     -0.01s   0.4%  runtime.acquirep
0     0% 16.53%     -0.01s   0.4%  runtime.acquirepNoTrace
0     0% 16.53%     -0.01s   0.4%  runtime.adjustframe
0     0% 16.53%     -0.07s  2.82%  runtime.bgsweep
0     0% 16.53%     -0.02s  0.81%  runtime.callers
0     0% 16.53%     -0.02s  0.81%  runtime.callers.func1
0     0% 16.53%      0.01s   0.4%  runtime.convT64
0     0% 16.53%     -0.01s   0.4%  runtime.copystack
0     0% 16.53%      0.11s  4.44%  runtime.deductAssistCredit
0     0% 16.53%     -0.02s  0.81%  runtime.deductSweepCredit
0     0% 16.53%      0.03s  1.21%  runtime.entersyscall
0     0% 16.53%      0.01s   0.4%  runtime.entersyscallWakeSysmon
0     0% 16.53%     -0.05s  2.02%  runtime.findRunnable
0     0% 16.53%      0.02s  0.81%  runtime.findnull
0     0% 16.53%      0.05s  2.02%  runtime.forEachP (inline)
0     0% 16.53%      0.05s  2.02%  runtime.forEachPInternal
0     0% 16.53%     -0.02s  0.81%  runtime.freeDeadSpanSPMCs
0     0% 16.53%      0.02s  0.81%  runtime.funcname (inline)
0     0% 16.53%     -0.02s  0.81%  runtime.funcspdelta (inline)
0     0% 16.53%      0.11s  4.44%  runtime.gcAssistAlloc
0     0% 16.53%      0.08s  3.23%  runtime.gcAssistAlloc.func2
0     0% 16.53%      0.08s  3.23%  runtime.gcAssistAlloc1
0     0% 16.53%     -0.08s  3.23%  runtime.gcBgMarkWorker.func2
0     0% 16.53%     -0.08s  3.23%  runtime.gcDrain
0     0% 16.53%     -0.03s  1.21%  runtime.gcDrainMarkWorkerDedicated (inline)
0     0% 16.53%     -0.05s  2.02%  runtime.gcDrainMarkWorkerIdle (inline)
0     0% 16.53%      0.06s  2.42%  runtime.gcDrainN
0     0% 16.53%      0.06s  2.42%  runtime.gcMarkDone
0     0% 16.53%     -0.01s   0.4%  runtime.gcMarkDone.forEachP.func5
0     0% 16.53%      0.07s  2.82%  runtime.gcMarkTermination
0     0% 16.53%      0.06s  2.42%  runtime.gcMarkTermination.forEachP.func7
0     0% 16.53%      0.06s  2.42%  runtime.gcMarkTermination.func4
0     0% 16.53%     -0.01s   0.4%  runtime.gcScanFinalizer
0     0% 16.53%      0.01s   0.4%  runtime.gcTrigger.test
0     0% 16.53%     -0.01s   0.4%  runtime.gcstopm
0     0% 16.53%     -0.01s   0.4%  runtime.goexit0
0     0% 16.53%     -0.03s  1.21%  runtime.gopreempt_m (inline)
0     0% 16.53%     -0.03s  1.21%  runtime.goschedImpl
0     0% 16.53%      0.02s  0.81%  runtime.gostringnocopy (inline)
0     0% 16.53%      0.04s  1.61%  runtime.growslice
0     0% 16.53%     -0.01s   0.4%  runtime.handoff
0     0% 16.53%      0.01s   0.4%  runtime.headTailIndex.split (inline)
0     0% 16.53%     -0.01s   0.4%  runtime.injectglist
0     0% 16.53%     -0.02s  0.81%  runtime.injectglist.func1
0     0% 16.53%      0.01s   0.4%  runtime.isSystemGoroutine
0     0% 16.53%     -0.04s  1.61%  runtime.lock (inline)
0     0% 16.53%     -0.03s  1.21%  runtime.lock2
0     0% 16.53%     -0.03s  1.21%  runtime.lockWithRank (inline)
0     0% 16.53%      0.01s   0.4%  runtime.mProf_Flush
0     0% 16.53%      0.01s   0.4%  runtime.mProf_FlushLocked
0     0% 16.53%     -0.01s   0.4%  runtime.mProf_Malloc
0     0% 16.53%      0.01s   0.4%  runtime.mProf_Malloc.func1
0     0% 16.53%      0.32s 12.90%  runtime.mallocgc
0     0% 16.53%      0.07s  2.82%  runtime.mallocgcLarge
0     0% 16.53%      0.01s   0.4%  runtime.mallocgcTiny
0     0% 16.53%     -0.01s   0.4%  runtime.mapaccess1_fast64
0     0% 16.53%     -0.01s   0.4%  runtime.mapaccess1_faststr
0     0% 16.53%     -0.01s   0.4%  runtime.mapassign
0     0% 16.53%      0.02s  0.81%  runtime.mapassign_faststr
0     0% 16.53%      0.01s   0.4%  runtime.mapdelete
0     0% 16.53%      0.01s   0.4%  runtime.markrootBlock
0     0% 16.53%      0.08s  3.23%  runtime.memclrNoHeapPointersChunked
0     0% 16.53%     -0.03s  1.21%  runtime.morestack
0     0% 16.53%      0.01s   0.4%  runtime.newMarkBits
0     0% 16.53%      0.01s   0.4%  runtime.newarray
0     0% 16.53%      0.01s   0.4%  runtime.newproc
0     0% 16.53%      0.01s   0.4%  runtime.newproc.func1
0     0% 16.53%      0.01s   0.4%  runtime.newproc1
0     0% 16.53%     -0.01s   0.4%  runtime.park_m
0     0% 16.53%      0.01s   0.4%  runtime.pcdatavalue
0     0% 16.53%      0.02s  0.81%  runtime.pollWork
0     0% 16.53%     -0.03s  1.21%  runtime.procyield (inline)
0     0% 16.53%     -0.01s   0.4%  runtime.profilealloc
0     0% 16.53%      0.03s  1.21%  runtime.reentersyscall
0     0% 16.53%      0.01s   0.4%  runtime.resetspinning
0     0% 16.53%      0.01s   0.4%  runtime.runqgrab
0     0% 16.53%      0.01s   0.4%  runtime.runqsteal
0     0% 16.53%      0.11s  4.44%  runtime.scanSpan
0     0% 16.53%     -0.01s   0.4%  runtime.scanframeworker
0     0% 16.53%     -0.01s   0.4%  runtime.scanstack
0     0% 16.53%     -0.04s  1.61%  runtime.schedule
0     0% 16.53%      0.01s   0.4%  runtime.setprofilebucket
0     0% 16.53%      0.01s   0.4%  runtime.shrinkstack
0     0% 16.53%     -0.01s   0.4%  runtime.slicebytetostring
0     0% 16.53%      0.01s   0.4%  runtime.stackcacherefill
0     0% 16.53%     -0.01s   0.4%  runtime.stackmapdata (inline)
0     0% 16.53%      0.01s   0.4%  runtime.stackpoolalloc
0     0% 16.53%      0.02s  0.81%  runtime.stopm
0     0% 16.53%      0.01s   0.4%  runtime.suspendG
0     0% 16.53%      0.09s  3.63%  runtime.systemstack
0     0% 16.53%     -0.02s  0.81%  runtime.tracebackPCs
0     0% 16.53%      0.01s   0.4%  runtime.typedmemclr
0     0% 16.53%      0.02s  0.81%  runtime.wakep
0     0% 16.53%     -0.03s  1.21%  runtime.wbBufFlush
0     0% 16.53%     -0.03s  1.21%  runtime.wbBufFlush.func1
0     0% 16.53%     -0.04s  1.61%  runtime.wbBufFlush1
0     0% 16.53%     -0.01s   0.4%  runtime/pprof.(*Profile).WriteTo
0     0% 16.53%     -0.01s   0.4%  runtime/pprof.(*profileBuilder).build
0     0% 16.53%     -0.01s   0.4%  runtime/pprof.writeHeap
0     0% 16.53%     -0.01s   0.4%  runtime/pprof.writeHeapInternal
0     0% 16.53%     -0.01s   0.4%  runtime/pprof.writeHeapProto
0     0% 16.53%      0.01s   0.4%  strconv.FormatFloat (inline)
0     0% 16.53%     -0.01s   0.4%  strconv.FormatInt (inline)
0     0% 16.53%      0.02s  0.81%  strings.Cut (inline)
0     0% 16.53%      0.01s   0.4%  sync.(*Cond).Broadcast
0     0% 16.53%      0.02s  0.81%  sync.(*Mutex).Lock (inline)
0     0% 16.53%      0.02s  0.81%  sync.(*Mutex).Unlock
0     0% 16.53%      0.01s   0.4%  sync.(*Once).Do (inline)
0     0% 16.53%      0.01s   0.4%  sync.(*Once).doSlow
0     0% 16.53%      0.01s   0.4%  sync.(*Pool).Get
0     0% 16.53%      0.02s  0.81%  sync.(*Pool).getSlow
0     0% 16.53%     -0.01s   0.4%  sync.(*Pool).pin
0     0% 16.53%     -0.01s   0.4%  sync.(*Pool).pinSlow
0     0% 16.53%      0.01s   0.4%  sync.runtime_notifyListNotifyAll
0     0% 16.53%     -0.01s   0.4%  sync.runtime_notifyListWait
0     0% 16.53%      0.01s   0.4%  sync/atomic.(*Value).Store
0     0% 16.53%      0.02s  0.81%  syscall.Accept4
0     0% 16.53%      0.04s  1.61%  syscall.RawSyscall6
0     0% 16.53%      0.09s  3.63%  syscall.Read (inline)
0     0% 16.53%      0.01s   0.4%  syscall.Stat (inline)
0     0% 16.53%      0.04s  1.61%  syscall.Syscall
0     0% 16.53%      0.03s  1.21%  syscall.Syscall6
0     0% 16.53%     -0.05s  2.02%  syscall.Write (inline)
0     0% 16.53%      0.02s  0.81%  syscall.accept4
0     0% 16.53%      0.01s   0.4%  syscall.fstatat
0     0% 16.53%      0.09s  3.63%  syscall.read
0     0% 16.53%     -0.05s  2.02%  syscall.write
0     0% 16.53%      0.01s   0.4%  time.Now
0     0% 16.53%     -0.01s   0.4%  time.Time.AppendFormat
```

## Diff heap

```bash
go tool pprof -top -diff_base=./profiles/base.heap.pprof ./profiles/result.heap.pprof
File: shortener
Build ID: 4baed019e76ecd5629a54e83bd1551c6a4d6c97c
Type: inuse_space
Time: 2026-06-07 14:58:10 MSK
Duration: 120.03s, Total samples = 7816.89kB 
Showing nodes accounting for -2247.12kB, 28.75% of 7816.89kB total
Dropped 6 nodes (cum <= 39.08kB)
      flat  flat%   sum%        cum   cum%
-1184.27kB 15.15% 15.15% -1184.27kB 15.15%  runtime/pprof.StartCPUProfile
 -548.84kB  7.02% 22.17%  -548.84kB  7.02%  compress/flate.(*compressor).initDeflate (inline)
    -514kB  6.58% 28.75%     -514kB  6.58%  bufio.NewWriterSize (inline)
         0     0% 28.75%  -548.84kB  7.02%  compress/flate.(*compressor).init
         0     0% 28.75%  -548.84kB  7.02%  compress/flate.NewWriter (inline)
         0     0% 28.75%  -548.84kB  7.02%  compress/gzip.(*Writer).Write
         0     0% 28.75%  -548.84kB  7.02%  encoding/json.(*Encoder).Encode
         0     0% 28.75% -1184.32kB 15.15%  github.com/go-chi/chi/v5.(*Mux).Mount.func1
         0     0% 28.75% -1733.16kB 22.17%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 28.75% -1733.16kB 22.17%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 28.75%  -548.84kB  7.02%  github.com/go-chi/chi/v5/middleware.(*basicWriter).Write
         0     0% 28.75% -1184.32kB 15.15%  github.com/go-chi/chi/v5/middleware.NoCache.func1
         0     0% 28.75% -1733.16kB 22.17%  github.com/go-chi/chi/v5/middleware.init.0.RequestLogger.func1.1
         0     0% 28.75%  -548.84kB  7.02%  github.com/liebeSonne/shortlink/internal/handler.(*shortLinkHandler).HandleGetUserUrls
         0     0% 28.75% -1733.16kB 22.17%  github.com/liebeSonne/shortlink/internal/handler/compress.NewCompressorMiddleware.NewEncodingMiddleware.NewGzipHandlerMiddleware.func1
         0     0% 28.75%  -548.84kB  7.02%  github.com/liebeSonne/shortlink/internal/handler/compress/gzip.(*gzipWriter).Write
         0     0% 28.75% -1733.16kB 22.17%  main.initRouter.LoggingMiddleware.func5
         0     0% 28.75% -1733.16kB 22.17%  main.initRouter.NewAuthCookieMiddleware.func4
         0     0% 28.75% -1733.16kB 22.17%  main.initRouter.NewAuthMiddleware.func3
         0     0% 28.75% -2247.17kB 28.75%  net/http.(*conn).serve
         0     0% 28.75% -1733.16kB 22.17%  net/http.HandlerFunc.ServeHTTP
         0     0% 28.75%     -514kB  6.58%  net/http.newBufioWriterSize
         0     0% 28.75% -1733.16kB 22.17%  net/http.serverHandler.ServeHTTP
         0     0% 28.75% -1184.27kB 15.15%  net/http/pprof.Profile
```
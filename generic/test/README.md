# Functions

- Amb(observable Observable, observables ...Observable) Observable
- CombineLatest(f FunctionN, observable Observable, observables ...Observable) Observable
Concat
Create
Defer
Empty
Error
From
FromChan
- FromEventSource(ch chan interface{}, opts ...options.Option) Observable
- FromIterable(it Iterable) Observable
- FromIterator(it Iterator) Observable
FromSlice
Froms
Interval
Just
Merge
Never
Range
Repeat
Start
Throw

# Operators

All
AsObservable
AutoConnect
Average
- BufferWithCount(count, skip int) Observable
- BufferWithTime(timespan, timeshift Duration) Observable
- BufferWithTimeOrCount(timespan Duration, count int) Observable
- Contains(equal Predicate) Single
Catch
Concat
ConcatAll
Count
Debounce
- DefaultIfEmpty(defaultValue interface{}) Observable
Delay
Distinct
- DistinctUntilChanged(apply Function) Observable
Do
DoOnComplete
DoOnError
ElementAt
Filter
Finally
First
- FirstOrDefault(defaultValue interface{}) Single
IgnoreCompletion
IgnoreElements
Last
- LastOrDefault(defaultValue interface{}) Single
Map
Max
Merge
MergeAll
MergeDelayError
MergeMap
Min
- OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable
- OnErrorReturn(resumeFunc ErrorFunction) Observable
- OnErrorReturnItem(item interface{}) Observable
ObserveOn
Only
Passthrough
Publish
PublishReplay
Reduce
RefCount
Repeat(count int) Observable
- Repeat(count int64, frequency Duration) Observable
Retry
Sample
Scan
- SequenceEqual(obs Observable) Single
Serialize
Single
Skip
SkipLast
- SkipWhile(apply Predicate) Observable
- StartWithItems(item interface{}, items ...interface{}) Observable
- StartWithIterable(iterable Iterable) Observable
- StartWithObservable(observable Observable) Observable
SubscribeOn
Sum
SwitchAll
SwitchMap
Take
TakeLast
TakeUntil
TakeWhile
Timeout
ToChan
ToSingle
ToSlice
- ZipFromObservable(publisher Observable, zipper Function2) Observable

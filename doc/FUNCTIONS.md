# Observables

The specification below is a template (generics) specification for golang.

Types are capitalized differently depending on where they are used in
the declaration. We use <T> to indicate the use of the capitalized type name 
and <t> for the actual real type name. e.g. <T> = Int and <t> = int.

We use the following type placeholders inside the function declarations.
<T> any type capitalized name.
<t> any type real name.
<U> second type capitalized name.
<u> second type real name.

A guard expresion can be prefixed to the function declaration and limits
for what types the function can be expanded. So <int> means only <t>=int
is allowed. Also <int,float> means that only <t>=int and <u>=float is
allowed. Using <*> means any type and <+> means any numerical type.
So e.g. <*,int> means function allowed for any type <t> but <u> has to
be int. The guard <*,*> indicates that any combination of types is allowed
for <t> and <u> in the function declaration.

## Creation
func Create<T>(f func(<T>Observer)) Observable<T>
func Emtpy<T>() Observable<T>
func Never<T>() Observable<T>
func Throw<T>(err error) Observable<T>
func From<T>Array(array []<t>) Observable<T>
func From<T>s(array ...<t>) Observable<T>
func From<T>Channel(data <-chan <t>) Observable<T>
func From<T>ChannelWithError(data <-chan <t>, errs <-chan error) Observable<T>
<int> func Interval(interval time.Duration) Observable<T>
func Just<T>(element <t>) Observable<T>
<int> func Range(start, count <t>) Observable<T>
func Repeat<T>(value <t>, count int) Observable<T>
func Start<T>(f func() (<t>, error)) Observable<T>
func Merge<T>(observables ...Observable<T>) Observable<T>
func Merge<T>DelayError(observables ...Observable<T>) Observable<T>

## Filters
func (o Observable<T>) Distinct() Observable<T>
func (o Observable<T>) ElementAt(n int) Observable<T>
func (o Observable<T>) Filter(f func(<t>) bool) Observable<T>
func (o Observable<T>) First() Observable<T>
func (o Observable<T>) Last() Observable<T>
func (o Observable<T>) Skip(n int) Observable<T>
func (o Observable<T>) SkipLast(n int) Observable<T>
func (o Observable<T>) Take(n int) Observable<T>
func (o Observable<T>) TakeLast(n int) Observable<T>
func (o Observable<T>) IgnoreElements() Observable<T>
func (o Observable<T>) IgnoreCompletion() Observable<T>
func (o Observable<T>) One() Observable<T>
func (o Observable<T>) Replay(size int, duration time.Duration) Observable<T>
func (o Observable<T>) Sample(duration time.Duration) Observable<T>
func (o Observable<T>) Debounce(duration time.Duration) Observable<T>

## Count
<*,int> func (o Observable<T>) Count() Observable<U>

## Miscelaneous
func (o Observable<T>) Concat(other ...Observable<T>) Observable<T>

func (o Observable<T>) Merge(other ...Observable<T>) Observable<T>
func (o Observable<T>) MergeDelayError(other ...Observable<T>) Observable<T>

func (o Observable<T>) Catch(catch Observable<T>) Observable<T>
func (o Observable<T>) Retry() Observable<T>

func (o Observable<T>) Do(f func(next <t>)) Observable<T>
func (o Observable<T>) DoOnError(f func(err error)) Observable<T>
func (o Observable<T>) DoOnComplete(f func()) Observable<T>

func (o Observable<T>) Finally(f func()) Observable<T>

func (o Observable<T>) Reduce(initial <t>, reducer func(<t>, <t>) <t>) Observable<T>
func (o Observable<T>) Scan(initial <t>, f func(<t>, <t>) <t>) Observable<T>
func (o Observable<T>) Timeout(timeout time.Duration) Observable<T>
func (o Observable<T>) Fork() Observable<T>

func (o Observable<T>) Publish() Connectable<T>
func (c Connectable<T>) Connect() Unsubscriber

## Mathematical
<+> func (o Observable<T>) Average() Observable<T>
<+> func (o Observable<T>) Sum() Observable<T>
<+> func (o Observable<T>) Min() Observable<T>
<+> func (o Observable<T>) Max() Observable<T>

## Map
<*,*> func (o Observable<T>) Map<U>(f func(<t>) <u>) Observable<U>
<*,*> func (o Observable<T>) FlatMap<U>(f func(<t>) Observable<U>) Observable<U>

## Scheduling
The filter operations Sample and Debounce use a goroutine internally.
Subscribe uses scheduler.Goroutines by default, which can be overridden using SubscribeOn()
Concat, Catch and Retry use scheduler.Immediate in their implementation.

func (o Observable<T>) SubscribeOn(scheduler Scheduler) Observable<T>
func (o Observalbe<T>) ObserveOn(scheduler Scheduler) Observable<T>

## Subscribing
func (o Observable<T>) Subscribe(observer <T>ObserverFunc) Unsubscriber
func (o Observable<T>) SubscribeNext(f func(v <t>)) Unsubscriber
func (o Observable<T>) ToOneWithError() (v <t>, e error)
func (o Observable<T>) ToOne() <t>
func (o Observable<T>) ToArrayWithError() (a []<t>, e error)
func (o Observable<T>) ToArray() []<t>
func (o Observable<T>) ToChannelWithError() (<-chan <t>, <-chan error)
func (o Observable<T>) ToChannel() <-chan <t>
func (o Observable<T>) Wait() error
